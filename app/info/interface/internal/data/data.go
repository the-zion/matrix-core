package data

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewData, NewUserRepo)

type Data struct {
}

func NewData(logger log.Logger) (*Data, func(), error) {
	//l := log.NewHelper(log.With(logger, "module", "user/data/new-data"))
	d := &Data{}
	return d, func() {}, nil
}
