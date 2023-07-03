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

func New(cfg *config.Config, db *gorm.DB, pq *pqueue.PriorityQueue) LZContext {
	ctx := context.Background()
	return LZContext{
		Ctx: ctx,
		DB:  db,
		PQ:  pq,
		Cfg: cfg,
	}
}
