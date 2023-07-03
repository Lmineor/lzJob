package api

import (
	"github.com/Lmineor/lzJob/context"
	"github.com/emicklei/go-restful"
	"net/http"
)


var svc  TaskSvc
func RegisterTaskRoutes(ctx context.LZContext, ws *restful.WebService){
	idPathParam := ws.QueryParameter("id", "id")
	pageParam := ws.QueryParameter("page", "page")
	pageSizeParam := ws.QueryParameter("page_size", "page_size")
	ws.Route(
		ws.POST(Tasks).
			To(svc.RegisterTask(ctx)).
			Returns(http.StatusOK, http.StatusText(http.StatusOK), nil),
	)
	ws.Route(
		ws.GET(Task).
			To(svc.GetTask(ctx)).
			Param(idPathParam).
			Returns(http.StatusOK, http.StatusText(http.StatusOK), nil),
		)
	ws.Route(
		ws.GET(TaskResult).
			To(svc.GetTaskResult(ctx)).
			Param(idPathParam).
			Param(pageSizeParam).
			Param(pageParam).
			Returns(http.StatusOK, http.StatusText(http.StatusOK), nil),
	)
	ws.Route(
		ws.PUT(Task).
			To(svc.UpdateTask(ctx)).
			Param(idPathParam).
			Returns(http.StatusOK, http.StatusText(http.StatusOK), nil),
	)
	ws.Route(
		ws.DELETE(Task).
			To(svc.DeleteTask(ctx)).
			Param(idPathParam).
			Returns(http.StatusOK, http.StatusText(http.StatusOK), nil),
	)
	ws.Route(
		ws.GET(Tasks).
			To(svc.GetTasks(ctx)).
			Param(pageSizeParam).
			Param(pageParam).
			Returns(http.StatusOK, http.StatusText(http.StatusOK), nil),
		)
	ws.Route(
		ws.GET(TasksResult).
			To(svc.GetTasksResult(ctx)).
			Param(pageSizeParam).
			Param(pageParam).
			Returns(http.StatusOK, http.StatusText(http.StatusOK), nil),
	)
}