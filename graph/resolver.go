package graph

import (
	"sync"

	"github.com/amenabe22/chachata_backend/graph/chans"
	"github.com/amenabe22/chachata_backend/graph/model"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	RedisClient *redis.Client
	mu          sync.Mutex
	AdminChans  map[string]*chans.CoreAdminChannel
	Coredb      *gorm.DB
	Rooms       map[string]*model.Chatroom
}
