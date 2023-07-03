package api

import (
	"fmt"
	"github.com/Lmineor/lzJob/context"
	"github.com/emicklei/go-restful"
	"k8s.io/klog/v2"
	"net/http"
)

func NewServer(ctx context.LZContext)*http.Server{
	container := restful.NewContainer()
	container.Router(restful.CurlyRouter{})
	ws := new(restful.WebService)
	ws.Path(RootPath).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	RegisterTaskRoutes(ctx, ws)
	container.Add(ws)
	cors := restful.CrossOriginResourceSharing{
		AllowedHeaders: []string{"Content-Type", "Accept", "X-Auth-Token"},
		CookiesAllowed: false,
		Container:      container}
	container.Filter(cors.Filter)


	ad := fmt.Sprintf("%s:%d", ctx.Cfg.Server.ServerAddr, ctx.Cfg.Server.ListenPort)
	server := &http.Server{Addr: ad, Handler: container}
	klog.Infof("start listen on %s", ad)
	return server
}