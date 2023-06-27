package lzJob

import (
	gocontext "context"
	"github.com/Lmineor/lzJob/context"
	"github.com/Lmineor/lzJob/pqueue"
)

func buildContext(cfg *Config) context.LZContext {
	db := initMysql(&cfg.Mysql)
	ctx := gocontext.Background()
	pq := pqueue.NewPriorityQueue()
	return context.LZContext{
		Ctx: ctx,
		DB:  db,
		PQ:  pq,
	}
}

func main() {
	cfg := InitConfig("/abc/def")
	ctx := buildContext(cfg)
	Trigger(ctx)
}
