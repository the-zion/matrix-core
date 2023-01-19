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
		"SetArticleImageIrregularToCache": `
					local key = KEYS[1]
					local value = ARGV[1]
					local exist = redis.call("EXISTS", key)
					if exist == 1 then
						redis.call("LPUSH", key, value)
					end
					return 0
	`,
		"SetArticleContentIrregularToCache": `
					local key = KEYS[1]
					local value = ARGV[1]
					local exist = redis.call("EXISTS", key)
					if exist == 1 then
						redis.call("LPUSH", key, value)
					end
					return 0
	`,
		"CreateArticleCache": `
					local articleStatistic = KEYS[1]
					local article = KEYS[2]
					local articleHot = KEYS[3]
					local leaderboard = KEYS[4]
					local userArticleList = KEYS[5]
					local userArticleListVisitor = KEYS[6]
					local creationUser = KEYS[7]
					local creationUserVisitor = KEYS[8]

					local uuid = ARGV[1]
					local auth = ARGV[2]
					local id = ARGV[3]
					local member = ARGV[4]
					local mode = ARGV[5]
					local member2 = ARGV[6]

					local userArticleListExist = redis.call("EXISTS", userArticleList)
					local creationUserExist = redis.call("EXISTS", creationUser)
					local articleExist = redis.call("EXISTS", article)
					local articleHotExist = redis.call("EXISTS", articleHot)
					local leaderboardExist = redis.call("EXISTS", leaderboard)
					local userArticleListVisitorExist = redis.call("EXISTS", userArticleListVisitor)
					local creationUserVisitorExist = redis.call("EXISTS", creationUserVisitor)

					redis.call("HSETNX", articleStatistic, "uuid", uuid)
					redis.call("HSETNX", articleStatistic, "agree", 0)
					redis.call("HSETNX", articleStatistic, "collect", 0)
					redis.call("HSETNX", articleStatistic, "view", 0)
					redis.call("HSETNX", articleStatistic, "comment", 0)
					redis.call("HSETNX", articleStatistic, "auth", auth)
					redis.call("EXPIRE", articleStatistic, 1800)

					if userArticleListExist == 1 then
						redis.call("ZADD", userArticleList, id, member)
					end

					if (creationUserExist == 1) and (mode == "create") then
						redis.call("HINCRBY", creationUser, "article", 1)
					end

					if auth == "2" then
						return 0
					end

					if articleExist == 1 then
						redis.call("ZADD", article, id, member)
					end

					if articleHotExist == 1 then
						redis.call("ZADD", articleHot, 0, member)
					end

					if leaderboardExist == 1 then
						redis.call("ZADD", leaderboard, 0, member2)
					end

					if userArticleListVisitorExist == 1 then
						redis.call("ZADD", userArticleListVisitor, id, member)
					end

					if (creationUserVisitorExist == 1) and (mode == "create") then
						redis.call("HINCRBY", creationUserVisitor, "article", 1)
					end
					return 0
	`,
		"DeleteArticleCache": `
					local key1 = KEYS[1]
					local key2 = KEYS[2]
					local key3 = KEYS[3]
					local key4 = KEYS[4]
					local key5 = KEYS[5]
					local key6 = KEYS[6]
					local key7 = KEYS[7]
					local key8 = KEYS[8]
					local key9 = KEYS[9]

                    local member = ARGV[1]
					local leadMember = ARGV[2]
					local auth = ARGV[3]

                    redis.call("ZREM", key1, member)
					redis.call("ZREM", key2, member)
					redis.call("ZREM", key3, leadMember)
					redis.call("DEL", key4)
					redis.call("DEL", key5)

                    redis.call("ZREM", key6, member)

                    local exist = redis.call("EXISTS", key8)
					if exist == 1 then
						local number = tonumber(redis.call("HGET", key8, "article"))
						if number > 0 then
  							redis.call("HINCRBY", key8, "article", -1)
						end
					end

                    if auth == 2 then
						return 0
					end

					redis.call("ZREM", key7, member)

					local exist = redis.call("EXISTS", key9)
					if exist == 1 then
						local number = tonumber(redis.call("HGET", key9, "article"))
						if number > 0 then
  							redis.call("HINCRBY", key9, "article", -1)
						end
					end

					return 0
	`,
		"AddArticleCommentToCache": `
					local key = KEYS[1]
					local keyExist = redis.call("EXISTS", key)
					if keyExist == 1 then
						redis.call("HINCRBY", key, "comment", 1)
					end
					return 0
	`,
		"ReduceArticleCommentToCache": `
					local key = KEYS[1]
					local value = redis.call("EXISTS", key)
					if value == 1 then
						local number = tonumber(redis.call("HGET", key, "comment"))
						if number > 0 then
  							redis.call("HINCRBY", key, "comment", -1)
						end
					end
					return 0
	`,
		"SetArticleAgreeToCache": `
					local hotKey = KEYS[1]
					local statisticKey = KEYS[2]
					local boardKey = KEYS[3]
					local userKey = KEYS[4]

					local member1 = ARGV[1]
					local member2 = ARGV[2]
					local id = ARGV[3]

					local hotKeyExist = redis.call("EXISTS", hotKey)
					local statisticKeyExist = redis.call("EXISTS", statisticKey)
					local boardKeyExist = redis.call("EXISTS", boardKey)
					local userKeyExist = redis.call("EXISTS", userKey)

					if hotKeyExist == 1 then
						redis.call("ZINCRBY", hotKey, 1, member1)
					end

					if statisticKeyExist == 1 then
						redis.call("HINCRBY", statisticKey, "agree", 1)
					end

					if boardKeyExist == 1 then
						redis.call("ZINCRBY", boardKey, 1, member2)
					end

					if userKeyExist == 1 then
						redis.call("SADD", userKey, id)
					end
					return 0
	`,
		"SetArticleViewToCache": `
					local key = KEYS[1]
					local value = ARGV[1]
					local exist = redis.call("EXISTS", key)
					if exist == 1 then
						redis.call("HINCRBY", key, "view", value)
					end
					return 0
	`,
		"SetArticleCollectToCache": `
					local statisticKey = KEYS[1]
					local collectKey = KEYS[2]
					local collectionsKey = KEYS[3]
					local creationKey = KEYS[4]
					local userKey = KEYS[5]

					local member = ARGV[1]

					local statisticKeyExist = redis.call("EXISTS", statisticKey)
					local collectKeyExist = redis.call("EXISTS", collectKey)
					local collectionsKeyExist = redis.call("EXISTS", collectionsKey)
					local creationKeyExist = redis.call("EXISTS", creationKey)
					local userKeyExist = redis.call("EXISTS", userKey)
					
					if statisticKeyExist == 1 then
						redis.call("HINCRBY", statisticKey, "collect", 1)
					end

					if collectKeyExist == 1 then
						redis.call("ZADD", collectKey, 0, member)
					end

					if collectionsKeyExist == 1 then
						redis.call("HINCRBY", collectionsKey, "article", 1)
					end

					if creationKeyExist == 1 then
						redis.call("HINCRBY", creationKey, "collect", 1)
					end

					if userKeyExist == 1 then
						redis.call("SADD", userKey, id)
					end
					return 0
	`,
		"SetUserArticleAgreeToCache": `
					local key = KEYS[1]
                    local change = ARGV[1]
					local value = redis.call("EXISTS", key)
					if value == 1 then
  						redis.call("SADD", key, change)
					end
					return 0
	`,
		"SetUserArticleCollectToCache": `
					local key = KEYS[1]
                    local change = ARGV[1]
					local value = redis.call("EXISTS", key)
					if value == 1 then
  						redis.call("SADD", key, change)
					end
					return 0
	`,
		"CancelArticleAgreeFromCache": `
					local hotKey = KEYS[1]
                    local member = ARGV[1]
					local hotKeyExist = redis.call("EXISTS", hotKey)
					if hotKeyExist == 1 then
						local score = tonumber(redis.call("ZSCORE", hotKey, member))
						if score > 0 then
  							redis.call("ZINCRBY", hotKey, -1, member)
						end
					end

					local boardKey = KEYS[2]
                    local member = ARGV[2]
					local boardKeyExist = redis.call("EXISTS", boardKey)
					if boardKeyExist == 1 then
						local score = tonumber(redis.call("ZSCORE", boardKey, member))
						if score > 0 then
  							redis.call("ZINCRBY", boardKey, -1, member)
						end
					end

					local statisticKey = KEYS[3]
					local statisticKeyExist = redis.call("EXISTS", statisticKey)
					if statisticKeyExist == 1 then
						local number = tonumber(redis.call("HGET", statisticKey, "agree"))
						if number > 0 then
  							redis.call("HINCRBY", statisticKey, "agree", -1)
						end
					end

					local userKey = KEYS[4]
					local commentId = ARGV[3]
					redis.call("SREM", userKey, commentId)
					return 0
	`,
		"CancelArticleCollectFromCache": `
					local statisticKey = KEYS[1]
					local statisticKeyExist = redis.call("EXISTS", statisticKey)
					if statisticKeyExist == 1 then
						local number = tonumber(redis.call("HGET", statisticKey, "collect"))
						if number > 0 then
  							redis.call("HINCRBY", statisticKey, "collect", -1)
						end
					end

					local collectKey = KEYS[2]
					local articleMember = ARGV[1]
					redis.call("ZREM", collectKey, articleMember)

					local collectionsKey = KEYS[3]
					local collectionsKeyExist = redis.call("EXISTS", collectionsKey)
					if collectionsKeyExist == 1 then
						local number = tonumber(redis.call("HGET", collectionsKey, "article"))
						if number > 0 then
  							redis.call("HINCRBY", collectionsKey, "article", -1)
						end
					end

					local creationKey = KEYS[4]
					local creationKeyExist = redis.call("EXISTS", creationKey)
					if creationKeyExist == 1 then
						local number = tonumber(redis.call("HGET", creationKey, "collect"))
						if number > 0 then
  							redis.call("HINCRBY", creationKey, "collect", -1)
						end
					end

					local userKey = KEYS[5]
					local articleId = ARGV[2]
					redis.call("SREM", userKey, articleId)

					return 0
	`,
		"CancelUserArticleAgreeFromCache": `
					local key = KEYS[1]
                    local change = ARGV[1]
					local value = redis.call("EXISTS", key)
					if value == 1 then
						redis.call("SREM", key, change)
					end
					return 0
	`,
		"CancelUserArticleCollectFromCache": `
					local key = KEYS[1]
                    local change = ARGV[1]
					local value = redis.call("EXISTS", key)
					if value == 1 then
  						redis.call("SREM", key, change)
					end
					return 0
	`,
		"AddColumnIncludesToCache": `
					local key = KEYS[1]
					local member = KEYS[2]
                    local change = ARGV[1]
					local value = redis.call("EXISTS", key)
					if value == 1 then
  						redis.call("ZADD", key, change, member)
					end
					return 0
	`,
		"CreateColumnCache": `
					local columnStatistic = KEYS[1]
					local column = KEYS[2]
					local columnHot = KEYS[3]
					local leaderboard = KEYS[4]
					local userColumnList = KEYS[5]
					local userColumnListVisitor = KEYS[6]
					local creationUser = KEYS[7]
					local creationUserVisitor = KEYS[8]

					local uuid = ARGV[1]
					local auth = ARGV[2]
					local id = ARGV[3]
					local member = ARGV[4]
					local mode = ARGV[5]
					local member2 = ARGV[6]

					local userColumnListExist = redis.call("EXISTS", userColumnList)
					local creationUserExist = redis.call("EXISTS", creationUser)
					local columnExist = redis.call("EXISTS", column)
					local columnHotExist = redis.call("EXISTS", columnHot)
					local leaderboardExist = redis.call("EXISTS", leaderboard)
					local userColumnListVisitorExist = redis.call("EXISTS", userColumnListVisitor)
					local creationUserVisitorExist = redis.call("EXISTS", creationUserVisitor)

					redis.call("HSETNX", columnStatistic, "uuid", uuid)
					redis.call("HSETNX", columnStatistic, "agree", 0)
					redis.call("HSETNX", columnStatistic, "collect", 0)
					redis.call("HSETNX", columnStatistic, "view", 0)
					redis.call("HSETNX", columnStatistic, "comment", 0)
					redis.call("HSETNX", columnStatistic, "auth", auth)
					redis.call("EXPIRE", columnStatistic, 1800)

					if userColumnListExist == 1 then
						redis.call("ZADD", userColumnList, id, member)
					end

					if (creationUserExist == 1) and (mode == "create") then
						redis.call("HINCRBY", creationUser, "column", 1)
					end

					if auth == "2" then
						return 0
					end

					if columnExist == 1 then
						redis.call("ZADD", column, id, member)
					end

					if columnHotExist == 1 then
						redis.call("ZADD", columnHot, 0, member)
					end

					if leaderboardExist == 1 then
						redis.call("ZADD", leaderboard, 0, member2)
					end

					if userColumnListVisitorExist == 1 then
						redis.call("ZADD", userColumnListVisitor, id, member)
					end

					if (creationUserVisitorExist == 1) and (mode == "create") then
						redis.call("HINCRBY", creationUserVisitor, "column", 1)
					end
					return 0
	`,
		"DeleteColumnCache": `
					local key1 = KEYS[1]
					local key2 = KEYS[2]
					local key3 = KEYS[3]
					local key4 = KEYS[4]
					local key5 = KEYS[5]
					local key6 = KEYS[6]
					local key7 = KEYS[7]
					local key8 = KEYS[8]
					local key9 = KEYS[9]

                    local member = ARGV[1]
					local leadMember = ARGV[2]
					local auth = ARGV[3]

                    redis.call("ZREM", key1, member)
					redis.call("ZREM", key2, member)
					redis.call("ZREM", key3, leadMember)
					redis.call("DEL", key4)
					redis.call("DEL", key5)

                    redis.call("ZREM", key6, member)

                    local exist = redis.call("EXISTS", key8)
					if exist == 1 then
						local number = tonumber(redis.call("HGET", key8, "column"))
						if number > 0 then
  							redis.call("HINCRBY", key8, "column", -1)
						end
					end

                    if auth == "2" then
						return 0
					end

					redis.call("ZREM", key7, member)

					local exist = redis.call("EXISTS", key9)
					if exist == 1 then
						local number = tonumber(redis.call("HGET", key9, "column"))
						if number > 0 then
  							redis.call("HINCRBY", key9, "column", -1)
						end
					end

					return 0
	`,
		"SetColumnViewToCache": `
					local key = KEYS[1]
					local value = ARGV[1]
					local exist = redis.call("EXISTS", key)
					if exist == 1 then
						redis.call("HINCRBY", key, "view", value)
					end
					return 0
	`,
		"SetColumnAgreeToCache": `
					local hotKey = KEYS[1]
					local statisticKey = KEYS[2]
					local boardKey = KEYS[3]
					local userKey = KEYS[4]

					local member1 = ARGV[1]
					local member2 = ARGV[2]
					local id = ARGV[3]

					local hotKeyExist = redis.call("EXISTS", hotKey)
					local statisticKeyExist = redis.call("EXISTS", statisticKey)
					local boardKeyExist = redis.call("EXISTS", boardKey)
					local userKeyExist = redis.call("EXISTS", userKey)

					if hotKeyExist == 1 then
						redis.call("ZINCRBY", hotKey, 1, member1)
					end

					if statisticKeyExist == 1 then
						redis.call("HINCRBY", statisticKey, "agree", 1)
					end

					if boardKeyExist == 1 then
						redis.call("ZINCRBY", boardKey, 1, member2)
					end

					if userKeyExist == 1 then
						redis.call("SADD", userKey, id)
					end
					return 0
	`,
		"SetColumnCollectToCache": `
					local statisticKey = KEYS[1]
					local collectKey = KEYS[2]
					local collectionsKey = KEYS[3]
					local creationKey = KEYS[4]
					local userKey = KEYS[5]

					local member = ARGV[1]

					local statisticKeyExist = redis.call("EXISTS", statisticKey)
					local collectKeyExist = redis.call("EXISTS", collectKey)
					local collectionsKeyExist = redis.call("EXISTS", collectionsKey)
					local creationKeyExist = redis.call("EXISTS", creationKey)
					local userKeyExist = redis.call("EXISTS", userKey)
					
					if statisticKeyExist == 1 then
						redis.call("HINCRBY", statisticKey, "collect", 1)
					end

					if collectKeyExist == 1 then
						redis.call("ZADD", collectKey, 0, member)
					end

					if collectionsKeyExist == 1 then
						redis.call("HINCRBY", collectionsKey, "column", 1)
					end

					if creationKeyExist == 1 then
						redis.call("HINCRBY", creationKey, "collect", 1)
					end

					if userKeyExist == 1 then
						redis.call("SADD", userKey, id)
					end
					return 0
	`,
		"SetUserColumnAgreeToCache": `
					local key = KEYS[1]
                    local change = ARGV[1]
					local value = redis.call("EXISTS", key)
					if value == 1 then
  						redis.call("SADD", key, change)
					end
					return 0
	`,
		"SetUserColumnCollectToCache": `
					local key = KEYS[1]
                    local change = ARGV[1]
					local value = redis.call("EXISTS", key)
					if value == 1 then
  						redis.call("SADD", key, change)
					end
					return 0
	`,
		"CancelColumnAgreeFromCache": `
					local hotKey = KEYS[1]
                    local member = ARGV[1]
					local hotKeyExist = redis.call("EXISTS", hotKey)
					if hotKeyExist == 1 then
						local score = tonumber(redis.call("ZSCORE", hotKey, member))
						if score > 0 then
  							redis.call("ZINCRBY", hotKey, -1, member)
						end
					end

					local boardKey = KEYS[2]
                    local member = ARGV[2]
					local boardKeyExist = redis.call("EXISTS", boardKey)
					if boardKeyExist == 1 then
						local score = tonumber(redis.call("ZSCORE", boardKey, member))
						if score > 0 then
  							redis.call("ZINCRBY", boardKey, -1, member)
						end
					end

					local statisticKey = KEYS[3]
					local statisticKeyExist = redis.call("EXISTS", statisticKey)
					if statisticKeyExist == 1 then
						local number = tonumber(redis.call("HGET", statisticKey, "agree"))
						if number > 0 then
  							redis.call("HINCRBY", statisticKey, "agree", -1)
						end
					end

					local userKey = KEYS[4]
					local commentId = ARGV[3]
					redis.call("SREM", userKey, commentId)
					return 0
	`,
		"CancelColumnCollectFromCache": `
					local statisticKey = KEYS[1]
					local statisticKeyExist = redis.call("EXISTS", statisticKey)
					if statisticKeyExist == 1 then
						local number = tonumber(redis.call("HGET", statisticKey, "collect"))
						if number > 0 then
  							redis.call("HINCRBY", statisticKey, "collect", -1)
						end
					end

					local collectKey = KEYS[2]
					local columnMember = ARGV[1]
					redis.call("ZREM", collectKey, columnMember)

					local collectionsKey = KEYS[3]
					local collectionsKeyExist = redis.call("EXISTS", collectionsKey)
					if collectionsKeyExist == 1 then
						local number = tonumber(redis.call("HGET", collectionsKey, "column"))
						if number > 0 then
  							redis.call("HINCRBY", collectionsKey, "column", -1)
						end
					end

					local creationKey = KEYS[4]
					local creationKeyExist = redis.call("EXISTS", creationKey)
					if creationKeyExist == 1 then
						local number = tonumber(redis.call("HGET", creationKey, "collect"))
						if number > 0 then
  							redis.call("HINCRBY", creationKey, "collect", -1)
						end
					end

					local userKey = KEYS[5]
					local columnId = ARGV[2]
					redis.call("SREM", userKey, columnId)

					return 0
	`,
		"SetUserColumnSubscribeToCache": `
					local key = KEYS[1]
                    local change = ARGV[1]
					local value = redis.call("EXISTS", key)
					if value == 1 then
  						redis.call("SADD", key, change)
					end
					return 0
	`,
		"SetColumnSubscribeToCache": `
					local key = KEYS[1]
                    local change = ARGV[1]
					local value = redis.call("EXISTS", key)
					if value == 1 then
  						redis.call("SADD", key, change)
					end

					local key2 = KEYS[2]
					local member = KEYS[3]
					local value = redis.call("EXISTS", key2)
					if value == 1 then
  						redis.call("ZADD", key2, change, member)
					end

					local key3 = KEYS[4]
					local value = redis.call("EXISTS", key3)
					if value == 1 then
  						redis.call("HINCRBY", key3, "subscribe", 1)
					end

					return 0
	`,
		"SetColumnImageIrregularToCache": `
					local key = KEYS[1]
					local value = ARGV[1]
					local exist = redis.call("EXISTS", key)
					if exist == 1 then
						redis.call("LPUSH", key, value)
					end
					return 0
	`,
		"SetColumnContentIrregularToCache": `
					local key = KEYS[1]
					local value = ARGV[1]
					local exist = redis.call("EXISTS", key)
					if exist == 1 then
						redis.call("LPUSH", key, value)
					end
					return 0
	`,
		"CancelUserColumnAgreeFromCache": `
					local key = KEYS[1]
                    local change = ARGV[1]
					local value = redis.call("EXISTS", key)
					if value == 1 then
  						redis.call("SREM", key, change)
					end
					return 0
	`,
		"CancelUserColumnCollectFromCache": `
					local key = KEYS[1]
                    local change = ARGV[1]
					local value = redis.call("EXISTS", key)
					if value == 1 then
  						redis.call("SREM", key, change)
					end
					return 0
	`,
		"CancelUserColumnSubscribeFromCache": `
					local key = KEYS[1]
                    local change = ARGV[1]
					redis.call("SREM", key, change)

					return 0
	`,
		"CancelColumnSubscribeFromCache": `
					local key = KEYS[1]
                    local change = ARGV[1]
					redis.call("SREM", key, change)

					local key2 = KEYS[2]
					local member = KEYS[3]
					redis.call("ZREM", key2, member)

					local key3 = KEYS[4]
					local key3Exist = redis.call("EXISTS", key3)
					if key3Exist == 1 then
						local number = tonumber(redis.call("HGET", key3, "subscribe"))
						if number > 0 then
  							redis.call("HINCRBY", key3, "subscribe", -1)
						end
					end
					return 0
	`,
		"SetCollectionsContentIrregularToCache": `
					local key = KEYS[1]
					local value = ARGV[1]
					local exist = redis.call("EXISTS", key)
					if exist == 1 then
						redis.call("LPUSH", key, value)
					end
					return 0
	`,
		"CreateCollectionsCache": `
					local collectionsStatistic = KEYS[1]
					local userCollectionsList = KEYS[2]
					local userCollectionsListVisitor = KEYS[3]
					local creationUser = KEYS[4]
					local creationUserVisitor = KEYS[5]
					local userCollectionsListAll = KEYS[6]

					local uuid = ARGV[1]
					local auth = ARGV[2]
					local id = ARGV[3]
					local ids = ARGV[4]
					local mode = ARGV[5]

					local userCollectionsListExist = redis.call("EXISTS", userCollectionsList)
					local userCollectionsListVisitorExist = redis.call("EXISTS", userCollectionsListVisitor)
					local creationUserExist = redis.call("EXISTS", creationUser)
					local creationUserVisitorExist = redis.call("EXISTS", creationUserVisitor)
					local userCollectionsListAllExist = redis.call("EXISTS", userCollectionsListAll)

					redis.call("HSETNX", collectionsStatistic, "uuid", uuid)
					redis.call("HSETNX", collectionsStatistic, "auth", auth)
					redis.call("HSETNX", collectionsStatistic, "article", 0)
					redis.call("HSETNX", collectionsStatistic, "column", 0)
					redis.call("HSETNX", collectionsStatistic, "talk", 0)
					redis.call("EXPIRE", collectionsStatistic, 1800)

					if userCollectionsListExist == 1 then
						redis.call("ZADD", userCollectionsList, id, ids)
					end

					if userCollectionsListAllExist == 1 then
						redis.call("ZADD", userCollectionsListAll, id, ids)
					end

					if (creationUserExist == 1) and (mode == "create") then
						redis.call("HINCRBY", creationUser, "collections", 1)
					end

					if auth == "2" then
						return 0
					end

					if userCollectionsListVisitorExist == 1 then
						redis.call("ZADD", userCollectionsListVisitor, id, ids)
					end

					if (creationUserVisitorExist == 1) and (mode == "create") then
						redis.call("HINCRBY", creationUserVisitor, "collections", 1)
					end
					return 0
	`,
		"CreateTimeLineCache": `
					local userTimeLineList = KEYS[1]
					local id = ARGV[1]
					local member = ARGV[2]

					local userTimeLineListExist = redis.call("EXISTS", userTimeLineList)

					if userTimeLineListExist == 1 then
						redis.call("ZADD", userTimeLineList, id, member)
					end
					return 0
	`,
		"DeleteCreationCache": `
					local key1 = KEYS[1]
					local key2 = KEYS[2]
					local key3 = KEYS[3]
					local key4 = KEYS[4]
					local key5 = KEYS[5]
					local key6 = KEYS[6]

					local member = ARGV[1]
					local auth = ARGV[2]

                    redis.call("ZREM", key1, member)
					redis.call("ZREM", key5, member)
					redis.call("DEL", key6)

                    local exist = redis.call("EXISTS", key3)
					if exist == 1 then
						local number = tonumber(redis.call("HGET", key3, "collections"))
						if number > 0 then
  							redis.call("HINCRBY", key3, "collections", -1)
						end
					end

                    if auth == "2" then
						return 0
					end

					redis.call("ZREM", key2, member)

					local exist = redis.call("EXISTS", key4)
					if exist == 1 then
						local number = tonumber(redis.call("HGET", key4, "collections"))
						if number > 0 then
  							redis.call("HINCRBY", key4, "collections", -1)
						end
					end
					return 0
	`,
		"DeleteTalkCache": `
					local key1 = KEYS[1]
					local key2 = KEYS[2]
					local key3 = KEYS[3]
					local key4 = KEYS[4]
					local key5 = KEYS[5]
					local key6 = KEYS[6]
					local key7 = KEYS[7]
					local key8 = KEYS[8]
					local key9 = KEYS[9]

                    local member = ARGV[1]
					local leadMember = ARGV[2]
					local auth = ARGV[3]

                    redis.call("ZREM", key1, member)
					redis.call("ZREM", key2, member)
					redis.call("ZREM", key3, leadMember)
					redis.call("DEL", key4)
					redis.call("DEL", key5)

                    redis.call("ZREM", key6, member)

                    local exist = redis.call("EXISTS", key8)
					if exist == 1 then
						local number = tonumber(redis.call("HGET", key8, "talk"))
						if number > 0 then
  							redis.call("HINCRBY", key8, "talk", -1)
						end
					end

                    if auth == "2" then
						return 0
					end

					redis.call("ZREM", key7, member)

					local exist = redis.call("EXISTS", key9)
					if exist == 1 then
						local number = tonumber(redis.call("HGET", key9, "talk"))
						if number > 0 then
  							redis.call("HINCRBY", key9, "talk", -1)
						end
					end

					return 0
	`,
		"CreateTalkCache": `
					local talkStatistic = KEYS[1]
					local talk = KEYS[2]
					local talkHot = KEYS[3]
					local leaderboard = KEYS[4]
					local userTalkList = KEYS[5]
					local userTalkListVisitor = KEYS[6]
					local creationUser = KEYS[7]
					local creationUserVisitor = KEYS[8]

					local uuid = ARGV[1]
					local auth = ARGV[2]
					local id = ARGV[3]
					local member = ARGV[4]
					local mode = ARGV[5]
					local member2 = ARGV[6]

					local userTalkListExist = redis.call("EXISTS", userTalkList)
					local creationUserExist = redis.call("EXISTS", creationUser)
					local talkExist = redis.call("EXISTS", talk)
					local talkHotExist = redis.call("EXISTS", talkHot)
					local leaderboardExist = redis.call("EXISTS", leaderboard)
					local userTalkListVisitorExist = redis.call("EXISTS", userTalkListVisitor)
					local creationUserVisitorExist = redis.call("EXISTS", creationUserVisitor)

					redis.call("HSETNX", talkStatistic, "uuid", uuid)
					redis.call("HSETNX", talkStatistic, "agree", 0)
					redis.call("HSETNX", talkStatistic, "collect", 0)
					redis.call("HSETNX", talkStatistic, "view", 0)
					redis.call("HSETNX", talkStatistic, "comment", 0)
					redis.call("HSETNX", talkStatistic, "auth", auth)
					redis.call("EXPIRE", talkStatistic, 1800)

					if userTalkListExist == 1 then
						redis.call("ZADD", userTalkList, id, member)
					end

					if (creationUserExist == 1) and (mode == "create") then
						redis.call("HINCRBY", creationUser, "talk", 1)
					end

					if auth == "2" then
						return 0
					end

					if talkExist == 1 then
						redis.call("ZADD", talk, id, member)
					end

					if talkHotExist == 1 then
						redis.call("ZADD", talkHot, 0, member)
					end

					if leaderboardExist == 1 then
						redis.call("ZADD", leaderboard, 0, member2)
					end

					if userTalkListVisitorExist == 1 then
						redis.call("ZADD", userTalkListVisitor, id, member)
					end

					if (creationUserVisitorExist == 1) and (mode == "create") then
						redis.call("HINCRBY", creationUserVisitor, "talk", 1)
					end
					return 0
	`,
		"SetTalkAgreeToCache": `
					local hotKey = KEYS[1]
					local statisticKey = KEYS[2]
					local boardKey = KEYS[3]
					local userKey = KEYS[4]

					local member1 = ARGV[1]
					local member2 = ARGV[2]
					local id = ARGV[3]

					local hotKeyExist = redis.call("EXISTS", hotKey)
					local statisticKeyExist = redis.call("EXISTS", statisticKey)
					local boardKeyExist = redis.call("EXISTS", boardKey)
					local userKeyExist = redis.call("EXISTS", userKey)

					if hotKeyExist == 1 then
						redis.call("ZINCRBY", hotKey, 1, member1)
					end

					if statisticKeyExist == 1 then
						redis.call("HINCRBY", statisticKey, "agree", 1)
					end

					if boardKeyExist == 1 then
						redis.call("ZINCRBY", boardKey, 1, member2)
					end

					if userKeyExist == 1 then
						redis.call("SADD", userKey, id)
					end
					return 0
	`,
		"CancelTalkAgreeFromCache": `
					local hotKey = KEYS[1]
                    local member = ARGV[1]
					local hotKeyExist = redis.call("EXISTS", hotKey)
					if hotKeyExist == 1 then
						local score = tonumber(redis.call("ZSCORE", hotKey, member))
						if score > 0 then
  							redis.call("ZINCRBY", hotKey, -1, member)
						end
					end

					local boardKey = KEYS[2]
                    local member = ARGV[2]
					local boardKeyExist = redis.call("EXISTS", boardKey)
					if boardKeyExist == 1 then
						local score = tonumber(redis.call("ZSCORE", boardKey, member))
						if score > 0 then
  							redis.call("ZINCRBY", boardKey, -1, member)
						end
					end

					local statisticKey = KEYS[3]
					local statisticKeyExist = redis.call("EXISTS", statisticKey)
					if statisticKeyExist == 1 then
						local number = tonumber(redis.call("HGET", statisticKey, "agree"))
						if number > 0 then
  							redis.call("HINCRBY", statisticKey, "agree", -1)
						end
					end

					local userKey = KEYS[4]
					local commentId = ARGV[3]
					redis.call("SREM", userKey, commentId)
					return 0
	`,
		"SetTalkViewToCache": `
					local key = KEYS[1]
					local value = ARGV[1]
					local exist = redis.call("EXISTS", key)
					if exist == 1 then
						redis.call("HINCRBY", key, "view", value)
					end
					return 0
	`,
		"SetTalkCollectToCache": `
					local statisticKey = KEYS[1]
					local collectKey = KEYS[2]
					local collectionsKey = KEYS[3]
					local creationKey = KEYS[4]
					local userKey = KEYS[5]

					local member = ARGV[1]

					local statisticKeyExist = redis.call("EXISTS", statisticKey)
					local collectKeyExist = redis.call("EXISTS", collectKey)
					local collectionsKeyExist = redis.call("EXISTS", collectionsKey)
					local creationKeyExist = redis.call("EXISTS", creationKey)
					local userKeyExist = redis.call("EXISTS", userKey)
					
					if statisticKeyExist == 1 then
						redis.call("HINCRBY", statisticKey, "collect", 1)
					end

					if collectKeyExist == 1 then
						redis.call("ZADD", collectKey, 0, member)
					end

					if collectionsKeyExist == 1 then
						redis.call("HINCRBY", collectionsKey, "talk", 1)
					end

					if creationKeyExist == 1 then
						redis.call("HINCRBY", creationKey, "collect", 1)
					end

					if userKeyExist == 1 then
						redis.call("SADD", userKey, id)
					end
					return 0
	`,
		"SetUserTalkAgreeToCache": `
					local key = KEYS[1]
                    local change = ARGV[1]
					local value = redis.call("EXISTS", key)
					if value == 1 then
						redis.call("SADD", key, change)
					end
					return 0
	`,
		"SetUserTalkCollectToCache": `
					local key = KEYS[1]
                    local change = ARGV[1]
					local value = redis.call("EXISTS", key)
					if value == 1 then
  						redis.call("SADD", key, change)
					end
					return 0
	`,
		"SetTalkImageIrregularToCache": `
					local key = KEYS[1]
					local value = ARGV[1]
					local exist = redis.call("EXISTS", key)
					if exist == 1 then
						redis.call("LPUSH", key, value)
					end
					return 0
	`,
		"SetTalkContentIrregularToCache": `
					local key = KEYS[1]
					local value = ARGV[1]
					local exist = redis.call("EXISTS", key)
					if exist == 1 then
						redis.call("LPUSH", key, value)
					end
					return 0
	`,
		"CancelTalkCollectFromCache": `
					local statisticKey = KEYS[1]
					local statisticKeyExist = redis.call("EXISTS", statisticKey)
					if statisticKeyExist == 1 then
						local number = tonumber(redis.call("HGET", statisticKey, "collect"))
						if number > 0 then
  							redis.call("HINCRBY", statisticKey, "collect", -1)
						end
					end

					local collectKey = KEYS[2]
					local talkMember = ARGV[1]
					redis.call("ZREM", collectKey, talkMember)

					local collectionsKey = KEYS[3]
					local collectionsKeyExist = redis.call("EXISTS", collectionsKey)
					if collectionsKeyExist == 1 then
						local number = tonumber(redis.call("HGET", collectionsKey, "talk"))
						if number > 0 then
  							redis.call("HINCRBY", collectionsKey, "talk", -1)
						end
					end

					local creationKey = KEYS[4]
					local creationKeyExist = redis.call("EXISTS", creationKey)
					if creationKeyExist == 1 then
						local number = tonumber(redis.call("HGET", creationKey, "collect"))
						if number > 0 then
  							redis.call("HINCRBY", creationKey, "collect", -1)
						end
					end

					local userKey = KEYS[5]
					local talkId = ARGV[2]
					redis.call("SREM", userKey, talkId)

					return 0
	`,
		"CancelUserTalkAgreeFromCache": `
					local key = KEYS[1]
                    local change = ARGV[1]
					local value = redis.call("EXISTS", key)
					if value == 1 then
  						redis.call("SREM", key, change)
					end
					return 0
	`,
		"CancelUserTalkCollectFromCache": `
					local key = KEYS[1]
                    local change = ARGV[1]
					local value = redis.call("EXISTS", key)
					if value == 1 then
						redis.call("SREM", key, change)
					end
					return 0
	`,
		"AddTalkCommentToCache": `
					local key = KEYS[1]
					local keyExist = redis.call("EXISTS", key)
					if keyExist == 1 then
						redis.call("HINCRBY", key, "comment", 1)
					end
					return 0
	`,
		"ReduceTalkCommentToCache": `
					local key = KEYS[1]
					local value = redis.call("EXISTS", key)
					if value == 1 then
						local number = tonumber(redis.call("HGET", key, "comment"))
						if number > 0 then
  							redis.call("HINCRBY", key, "comment", -1)
						end
					end
					return 0
	`,
	}
)

func NewRedis(logger log.Logger) redis.Cmdable {
	l := log.NewHelper(log.With(logger, "lua", "redis"))

	client := redis.NewClient(&redis.Options{
		Addr:        addr,
		DB:          1,
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
