package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"os"
	"time"
)

var (
	addr      string
	password  string
	scriptBox = map[string]string{
		"SetUserCommentAgreeToCache": `
					local key = KEYS[1]
                    local change = ARGV[1]
					local value = redis.call("EXISTS", key)
					if value == 1 then
  						redis.call("SADD", key, change)
					end
					return 0
	`,
		"SetCommentAgreeToCache": `
					local hotKey = KEYS[1]
                    local member = ARGV[1]
					local hotKeyExist = redis.call("EXISTS", hotKey)
					if hotKeyExist == 1 then
						redis.call("ZINCRBY", hotKey, 1, member)
					end

					local statisticKey = KEYS[2]
					local statisticKeyExist = redis.call("EXISTS", statisticKey)
					if statisticKeyExist == 1 then
						redis.call("HINCRBY", statisticKey, "agree", 1)
					end

					local userKey = KEYS[3]
					local commentId = ARGV[2]
					redis.call("SADD", userKey, commentId)
					return 0
	`,
		"SetSubCommentAgreeToCache": `
					local statisticKey = KEYS[1]
					local statisticKeyExist = redis.call("EXISTS", statisticKey)
					if statisticKeyExist == 1 then
						redis.call("HINCRBY", statisticKey, "agree", 1)
					end

					local userKey = KEYS[2]
					local commentId = ARGV[1]
					redis.call("SADD", userKey, commentId)
					return 0
	`,
		"SetCommentContentIrregularToCache": `
					local key = KEYS[1]
					local value = ARGV[1]
					local exist = redis.call("EXISTS", key)
					if exist == 1 then
						redis.call("LPUSH", key, value)
					end
					return 0
	`,
		"CancelCommentAgreeFromCache": `
					local hotKey = KEYS[1]
                    local member = ARGV[1]
					local hotKeyExist = redis.call("EXISTS", hotKey)
					if hotKeyExist == 1 then
						local score = tonumber(redis.call("ZSCORE", hotKey, member))
						if score > 0 then
  							redis.call("ZINCRBY", hotKey, -1, member)
						end
					end

					local statisticKey = KEYS[2]
					local statisticKeyExist = redis.call("EXISTS", statisticKey)
					if statisticKeyExist == 1 then
						local number = tonumber(redis.call("HGET", statisticKey, "agree"))
						if number > 0 then
  							redis.call("HINCRBY", statisticKey, "agree", -1)
						end
					end

					local userKey = KEYS[3]
					local commentId = ARGV[2]
					redis.call("SREM", userKey, commentId)
					return 0
	`,
		"CancelSubCommentAgreeFromCache": `
					local statisticKey = KEYS[1]
					local statisticKeyExist = redis.call("EXISTS", statisticKey)
					if statisticKeyExist == 1 then
						local number = tonumber(redis.call("HGET", statisticKey, "agree"))
						if number > 0 then
  							redis.call("HINCRBY", statisticKey, "agree", -1)
						end
					end

					local userKey = KEYS[2]
					local commentId = ARGV[1]
					redis.call("SREM", userKey, commentId)
					return 0
	`,
		"CreateCommentCache": `
					local newKey = KEYS[1]
					local newMember = ARGV[2]
					local hotKey = KEYS[2]
					local hotMember = ARGV[2]
					local commentKey = KEYS[3]
					local commentUserKey1 = KEYS[4]
					local commentUserKey2 = KEYS[5]
					local commentUserReplyField = KEYS[6]
					local commentUserRepliedField = KEYS[7]
					local userCommentCreationReplyList = KEYS[8]
					local userCommentCreationRepliedList = KEYS[9]
					local userCommentCreationMessageRepliedList = KEYS[10]
					local memberReply = ARGV[3]
					local memberReplied = ARGV[4]
					local memberMessageReplied = ARGV[5]

					local newKeyExist = redis.call("EXISTS", newKey)
					local hotKeyExist = redis.call("EXISTS", hotKey)
					local commentUserKey1Exist = redis.call("EXISTS", commentUserKey1)
					local commentUserKey2Exist = redis.call("EXISTS", commentUserKey2)
					local userCommentCreationReplyListExist = redis.call("EXISTS", userCommentCreationReplyList)
					local userCommentCreationRepliedListExist = redis.call("EXISTS", userCommentCreationRepliedList)
					local userCommentCreationMessageRepliedListExist = redis.call("EXISTS", userCommentCreationMessageRepliedList)

					redis.call("HSETNX", commentKey, "agree", 0)
					redis.call("HSETNX", commentKey, "comment", 0)
					redis.call("EXPIRE", commentKey, 1800)

					local time = ARGV[1]

					if newKeyExist == 1 then
						redis.call("ZADD", newKey, time, newMember)
						redis.call("EXPIRE", newKey, 1800)
					end

					if hotKeyExist == 1 then
						redis.call("ZADD", hotKey, 0, hotMember)
						redis.call("EXPIRE", hotKey, 1800)
					end

					if commentUserKey1Exist == 1 then
						redis.call("HINCRBY", commentUserKey1, "comment", 1)
						redis.call("HINCRBY", commentUserKey1, commentUserReplyField, 1)
						redis.call("EXPIRE", commentUserKey1, 1800)
					end

					if commentUserKey2Exist == 1 then
						redis.call("HINCRBY", commentUserKey2, commentUserRepliedField, 1)
					end

					if userCommentCreationReplyListExist == 1 then
						redis.call("ZADD", userCommentCreationReplyList, time, memberReply)
						redis.call("EXPIRE", userCommentCreationReplyList, 1800)
					end

					if userCommentCreationRepliedListExist == 1 then
						redis.call("ZADD", userCommentCreationRepliedList, time, memberReplied)
						redis.call("EXPIRE", userCommentCreationRepliedList, 1800)
					end

					if userCommentCreationMessageRepliedListExist == 1 then
						redis.call("ZADD", userCommentCreationMessageRepliedList, time, memberMessageReplied)
						redis.call("EXPIRE", userCommentCreationMessageRepliedList, 1800)
					end
					return 0
	`,
		"CreateSubCommentCache": `
					local comment = KEYS[1]
					local commentRoot = KEYS[2]
					local subCommentList = KEYS[3]
					local commentUserKey1 = KEYS[4]
					local commentUserKey2 = KEYS[5]
					local commentUserKey3 = KEYS[6]
					local commentUserReplyField = KEYS[7]
					local commentUserRepliedField = KEYS[8]
					local userSubCommentCreationReplyList = KEYS[9]
					local userSubCommentCreationRepliedListForRoot = KEYS[10]
					local userSubCommentCreationRepliedListForParent = KEYS[11]
					local userSubCommentCreationMessageRepliedListForRoot = KEYS[12]
					local userSubCommentCreationMessageRepliedListForParent = KEYS[13]

					local time = ARGV[1]
					local subCommentListMember = ARGV[2]
					local memberReply = ARGV[3]
					local memberReplied = ARGV[4]
					local memberMessageReplied = ARGV[5]
					local parentId = ARGV[6]
					local rootUser = ARGV[7]
					local reply = ARGV[8]

					local commentRootExist = redis.call("EXISTS", commentRoot)
					local subCommentListExist = redis.call("EXISTS", subCommentList)
					local commentUserKey1Exist = redis.call("EXISTS", commentUserKey1)
					local commentUserKey2Exist = redis.call("EXISTS", commentUserKey2)
					local commentUserKey3Exist = redis.call("EXISTS", commentUserKey3)
					local userSubCommentCreationReplyListExist = redis.call("EXISTS", userSubCommentCreationReplyList)
					local userSubCommentCreationRepliedListForRootExist = redis.call("EXISTS", userSubCommentCreationRepliedListForRoot)
					local userSubCommentCreationRepliedListForParentExist = redis.call("EXISTS", userSubCommentCreationRepliedListForParent)
					local userSubCommentCreationMessageRepliedListForRootExist = redis.call("EXISTS", userSubCommentCreationMessageRepliedListForRoot)
					local userSubCommentCreationMessageRepliedListForParentExist = redis.call("EXISTS", userSubCommentCreationMessageRepliedListForParent)

					redis.call("HSETNX", comment, "agree", 0)
					redis.call("EXPIRE", comment, 1800)

					if commentRootExist == 1 then
						redis.call("HINCRBY", commentRoot, "comment", 1)
						redis.call("EXPIRE", commentRoot, 1800)
					end

					if subCommentListExist == 1 then
						redis.call("ZADD", subCommentList, time, subCommentListMember)
						redis.call("EXPIRE", subCommentList, 1800)
					end

					if commentUserKey1Exist == 1 then
						redis.call("HINCRBY", commentUserKey1, "comment", 1)
						redis.call("HINCRBY", commentUserKey1, commentUserReplyField, 1)
						redis.call("EXPIRE", commentUserKey1, 1800)
					end

					if commentUserKey2Exist == 1 then
						redis.call("HINCRBY", commentUserKey2, commentUserRepliedField, 1)
					end

					if (parentId ~= 0) and (rootUser ~= reply) and (commentUserKey3Exist == 1) then
						redis.call("HINCRBY", commentUserKey3, commentUserRepliedField, 1)
					end

					if userSubCommentCreationReplyListExist == 1 then
						redis.call("ZADD", userSubCommentCreationReplyList, time, memberReply)
						redis.call("EXPIRE", userSubCommentCreationReplyList, 1800)
					end

					if userSubCommentCreationRepliedListForRootExist == 1 then
						redis.call("ZADD", userSubCommentCreationRepliedListForRoot, time, memberReplied)
						redis.call("EXPIRE", userSubCommentCreationRepliedListForRoot, 1800)
					end

					if userSubCommentCreationRepliedListForParentExist == 1 then
						redis.call("ZADD", userSubCommentCreationRepliedListForParent, time, memberReplied)
						redis.call("EXPIRE", userSubCommentCreationRepliedListForParent, 1800)
					end

					if userSubCommentCreationMessageRepliedListForRootExist == 1 then
						redis.call("ZADD", userSubCommentCreationMessageRepliedListForRoot, time, memberMessageReplied)
						redis.call("EXPIRE", userSubCommentCreationMessageRepliedListForRoot, 1800)
					end

					if userSubCommentCreationMessageRepliedListForParentExist == 1 then
						redis.call("ZADD", userSubCommentCreationMessageRepliedListForParent, time, memberMessageReplied)
						redis.call("EXPIRE", userSubCommentCreationMessageRepliedListForParent, 1800)
					end
					return 0
	`,
		"RemoveCommentCache": `
					local newKey = KEYS[1]
					local newMember = ARGV[2]
					local hotKey = KEYS[2]
					local hotMember = ARGV[2]
					local commentKey = KEYS[3]
					local commentUserKey1 = KEYS[4]
					local commentUserKey2 = KEYS[5]
					local commentUserReplyField = KEYS[6]
					local commentUserRepliedField = KEYS[7]
					local userCommentCreationReplyList = KEYS[8]
					local userCommentCreationRepliedList = KEYS[9]
					local userCommentCreationMessageRepliedList = KEYS[10]
					local memberReply = ARGV[3]
					local memberReplied = ARGV[4]
					local memberMessageReplied = ARGV[5]


					local commentUserKey1Exist = redis.call("EXISTS", commentUserKey1)
					local commentUserKey2Exist = redis.call("EXISTS", commentUserKey2)

					redis.call("DEL", commentKey)

					redis.call("ZREM", newKey, newMember)
					redis.call("ZREM", hotKey, hotMember)

					if commentUserKey1Exist == 1 then
						local number = tonumber(redis.call("HGET", commentUserKey1, "comment"))
						if number > 0 then
  							redis.call("HINCRBY", commentUserKey1, "comment", -1)
						end

						local number = tonumber(redis.call("HGET", commentUserKey1, commentUserReplyField))
						if number > 0 then
  							redis.call("HINCRBY", commentUserKey1, commentUserReplyField, -1)
						end
					end

					if commentUserKey2Exist == 1 then
						local number = tonumber(redis.call("HGET", commentUserKey2, commentUserRepliedField))
						if number > 0 then
  							redis.call("HINCRBY", commentUserKey2, commentUserRepliedField, -1)
						end
					end

					redis.call("ZREM", userCommentCreationReplyList, memberReply)
					redis.call("ZREM", userCommentCreationRepliedList, memberReplied)
					redis.call("ZREM", userCommentCreationMessageRepliedList, memberMessageReplied)
					return 0
	`,
		"RemoveSubCommentCache": `
					local comment = KEYS[1]
					local commentRoot = KEYS[2]
					local subCommentList = KEYS[3]
					local commentUserKey1 = KEYS[4]
					local commentUserKey2 = KEYS[5]
					local commentUserKey3 = KEYS[6]
					local commentUserReplyField = KEYS[7]
					local commentUserRepliedField = KEYS[8]
					local userSubCommentCreationReplyList = KEYS[9]
					local userSubCommentCreationRepliedListForRoot = KEYS[10]
					local userSubCommentCreationRepliedListForParent = KEYS[11]
					local userSubCommentCreationMessageRepliedListForRoot = KEYS[12]
					local userSubCommentCreationMessageRepliedListForParent = KEYS[13]

					local commentId = ARGV[1]
					local subCommentListMember = ARGV[2]
					local memberReply = ARGV[3]
					local memberReplied = ARGV[4]
					local memberMessageReplied = ARGV[5]
					local parentId = ARGV[6]
					local rootUser = ARGV[7]
					local reply = ARGV[8]

					local commentRootExist = redis.call("EXISTS", commentRoot)
					local commentUserKey1Exist = redis.call("EXISTS", commentUserKey1)
					local commentUserKey2Exist = redis.call("EXISTS", commentUserKey2)
					local commentUserKey3Exist = redis.call("EXISTS", commentUserKey3)

					redis.call("DEL", comment)

					if commentRootExist == 1 then
						local number = tonumber(redis.call("HGET", commentRoot, "comment"))
						if number > 0 then
  							redis.call("HINCRBY", commentRoot, "comment", -1)
						end
					end

					redis.call("ZREM", subCommentList, subCommentListMember)

					if commentUserKey1Exist == 1 then
						local number = tonumber(redis.call("HGET", commentUserKey1, "comment"))
						if number > 0 then
  							redis.call("HINCRBY", commentUserKey1, "comment", -1)
						end

						local number = tonumber(redis.call("HGET", commentUserKey1, commentUserReplyField))
						if number > 0 then
  							redis.call("HINCRBY", commentUserKey1, commentUserReplyField, -1)
						end
					end

					if commentUserKey2Exist == 1 then
						local number = tonumber(redis.call("HGET", commentUserKey2, commentUserRepliedField))
						if number > 0 then
  							redis.call("HINCRBY", commentUserKey2, commentUserRepliedField, -1)
						end
					end

					if (parentId ~= 0) and (rootUser ~= reply) and (commentUserKey3Exist == 1) then
						local number = tonumber(redis.call("HGET", commentUserKey3, commentUserRepliedField))
						if number > 0 then
  							redis.call("HINCRBY", commentUserKey3, commentUserRepliedField, -1)
						end
					end

					redis.call("ZREM", userSubCommentCreationReplyList, memberReply)
					redis.call("ZREM", userSubCommentCreationRepliedListForRoot, memberReplied)
					redis.call("ZREM", userSubCommentCreationRepliedListForParent, memberReplied)
					redis.call("ZREM", userSubCommentCreationMessageRepliedListForRoot, memberMessageReplied)
					redis.call("ZREM", userSubCommentCreationMessageRepliedListForParent, memberMessageReplied)
					return 0
	`,
	}
)

func NewRedis(logger log.Logger) redis.Cmdable {
	l := log.NewHelper(log.With(logger, "lua", "redis"))

	client := redis.NewClient(&redis.Options{
		Addr:        addr,
		DB:          3,
		DialTimeout: time.Second * 2,
		PoolSize:    10,
		Password:    password,
	})
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*2)
	defer cancelFunc()
	err := client.Ping(timeout).Err()
	if err != nil {
		l.Fatalf("redis connect error: %v", err)
	}
	return client
}

func main() {
	flag.Parse()
	logger := log.NewStdLogger(os.Stdout)
	r := NewRedis(logger)
	ScriptLoad(r)
}

func ScriptLoad(r redis.Cmdable) {
	for key, value := range scriptBox {
		result, err := r.ScriptLoad(context.Background(), value).Result()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(key, result)
	}
}

func init() {
	flag.StringVar(&addr, "addr",
		"ip:port",
		"redis addr, eg: -addr ip:port")
	flag.StringVar(&password, "password",
		"abc",
		"redis password, eg: -password abc")
}
