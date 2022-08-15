package biz

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewUserUseCase, NewCreationUseCase, NewAchievementUseCase, NewCommentUseCase)

type Jwt interface {
	JwtCheck(token string) (string, error)
}
