package server

import (
	"github.com/google/wire"
)

// ProviderSet is user providers.
var ProviderSet = wire.NewSet(NewHTTPServer, NewGRPCServer, NewRocketMqConsumerServer)
