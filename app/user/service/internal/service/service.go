package service

import (
	"context"
	v1 "github.com/Cube-v2/matrix-core/api/user/service/v1"
	"github.com/Cube-v2/matrix-core/app/user/service/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	v4 "github.com/golang-jwt/jwt/v4"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewUserService)

type UserService struct {
	v1.UnimplementedUserServer
	uc   *biz.UserUseCase
	ac   *biz.AuthUseCase
	pc   *biz.ProfileUseCase
	achc *biz.AchievementUseCase
	log  *log.Helper
}

func GetUserIdFromContext(ctx context.Context) int64 {
	claims, _ := jwt.FromContext(ctx)
	mapClaims := *(claims.(*v4.MapClaims))
	return int64(mapClaims["user_id"].(float64))
}

func NewUserService(uc *biz.UserUseCase, ac *biz.AuthUseCase, pc *biz.ProfileUseCase, logger log.Logger) *UserService {
	return &UserService{
		log: log.NewHelper(log.With(logger, "module", "user/service")),
		uc:  uc,
		ac:  ac,
		pc:  pc,
	}
}
