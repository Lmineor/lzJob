package api

import (
	"github.com/Lmineor/lzJob/context"
	"github.com/Lmineor/lzJob/store"
	"github.com/Lmineor/lzJob/task"
	"github.com/emicklei/go-restful"
	"k8s.io/klog/v2"
)

type TaskSvc struct{}

func (t TaskSvc) GetTask(ctx context.LZContext) func(request *restful.Request, response *restful.Response) {
	return func(request *restful.Request, response *restful.Response) {
		taskId := request.PathParameter("id")

		klog.Infof("get task %s", taskId)
		ts, err := task.GetTask(ctx, taskId)
		if err != nil {
			BadRequestResp(response, err)
			return
		} else {
			SuccessRespWithData(response, ts)
		}
	}
}

func (t TaskSvc) GetTaskResult(ctx context.LZContext) func(request *restful.Request, response *restful.Response) {
	return func(request *restful.Request, response *restful.Response) {
		taskId := request.PathParameter("id")
		pageSizeP := request.QueryParameter("page_size")
		pageP := request.QueryParameter("page")
		page, pageSize, err := ParsePagination(pageP, pageSizeP)
		if err != nil {
			BadRequestResp(response, err)
			return
		}
		klog.Infof("get task result %s", taskId)
		ts, total, err := task.GetTaskResult(ctx, taskId, pageSize, page)
		if err != nil {
			BadRequestResp(response, err)
			return
		} else {
			res := PageResult{
				List:     ts,
				Total:    total,
				Page:     page,
				PageSize: pageSize,
			}
			SuccessRespWithData(response, res)
		}
	}
}

func (t TaskSvc) GetTasksResult(ctx context.LZContext) func(request *restful.Request, response *restful.Response) {
	return func(request *restful.Request, response *restful.Response) {
		pageSizeP := request.QueryParameter("page_size")
		pageP := request.QueryParameter("page")
		page, pageSize, err := ParsePagination(pageP, pageSizeP)
		if err != nil {
			BadRequestResp(response, err)
			return
		}
		klog.Infof("get tasks result")
		tasksResult, total, err := task.GetTasksResult(ctx, pageSize, page)
		if err != nil {
			BadRequestResp(response, err)
			return
		} else {
			res := PageResult{
				List:     tasksResult,
				Total:    total,
				Page:     page,
				PageSize: pageSize,
			}
			SuccessRespWithData(response, res)
		}
	}
}

func (t TaskSvc) GetTasks(ctx context.LZContext) func(request *restful.Request, response *restful.Response) {
	return func(request *restful.Request, response *restful.Response) {
		pageSizeP := request.QueryParameter("page_size")
		pageP := request.QueryParameter("page")
		page, pageSize, err := ParsePagination(pageP, pageSizeP)
		if err != nil {
			BadRequestResp(response, err)
			return
		}
		klog.Infof("get tasks")
		tasks, total, err := task.GetTasks(ctx, pageSize, page)
		if err != nil {
			BadRequestResp(response, err)
			return
		} else {
			res := PageResult{
				List:     tasks,
				Total:    total,
				Page:     page,
				PageSize: pageSize,
			}
			SuccessRespWithData(response, res)
		}
	}
}

func (t TaskSvc) RegisterTask(ctx context.LZContext) func(request *restful.Request, response *restful.Response) {
	return func(request *restful.Request, response *restful.Response) {
		var t TaskReqBody
		request.ReadEntity(&t)

		klog.Infof("register task %s", t)
		taskObj := toTaskObj(t)
		err := task.AddTask(ctx, taskObj)
		if err != nil {
			BadRequestResp(response, err)
			return
		} else {
			SuccessResp(response)
		}
	}
}

func (t TaskSvc) UpdateTask(ctx context.LZContext) func(request *restful.Request, response *restful.Response) {
	return func(request *restful.Request, response *restful.Response) {
		var t TaskReqBody

		taskId := request.PathParameter("id")
		request.ReadEntity(&t)

		klog.Infof("update task %s", t)
		taskObj := toTaskObj(t)
		err := task.UpdateTask(ctx, taskId, taskObj)
		if err != nil {
			BadRequestResp(response, err)
			return
		} else {
			SuccessResp(response)
		}
	}
}

func (t TaskSvc) DeleteTask(ctx context.LZContext) func(request *restful.Request, response *restful.Response) {
	return func(request *restful.Request, response *restful.Response) {
		taskId := request.PathParameter("id")

		klog.Infof("delete task %s", taskId)
		err := task.DeleteTask(ctx, taskId)
		if err != nil {
			BadRequestResp(response, err)
			return
		} else {
			SuccessResp(response)
		}
	}
}

func toTaskObj(req TaskReqBody) store.Task {
	return store.Task{
		Type:          req.Type,
		Status:        req.Status,
		Cron:          req.Cron,
		ExecTime:      req.ExecTime,
		ParamBody:     req.ParamBody,
		ExtInfo:       req.ExtInfo,
		TriggerSpi:    req.TriggerSpi,
		TriggerMethod: req.TriggerMethod,
	}
}
