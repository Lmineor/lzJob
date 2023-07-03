package context

import (
	"context"
	"github.com/Lmineor/lzJob/config"
	"github.com/Lmineor/lzJob/pqueue"
	"gorm.io/gorm"
)


type LZContext struct {
	Ctx context.Context
	DB  *gorm.DB
	PQ  *pqueue.PriorityQueue
	Cfg *config.Config
}
