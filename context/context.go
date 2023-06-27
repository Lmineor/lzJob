package context

import "context"

import (
	"github.com/Lmineor/lzJob/pqueue"
	"gorm.io/gorm"
)

type LZContext struct {
	Ctx context.Context
	DB  *gorm.DB
	PQ  *pqueue.PriorityQueue
}
