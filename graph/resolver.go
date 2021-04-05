package graph

import (
	"sync"

	"github.com/amenabe22/chachata_backend/graph/chans"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	mu         sync.Mutex
	AdminChans map[string]*chans.CoreAdminChannel
}
