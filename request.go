package lzJob

import (
	"fmt"
	"strconv"
	"time"
)

type PageInfo struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type TaskReqBody struct {
	Type          string    `json:"type"`
	Status        string    `json:"status"`
	Cron          string    `json:"cron"`
	ExecTime      time.Time `json:"exec_time"`
	ParamBody     string    `json:"param_body"`
	ExtInfo       string    `json:"ext_info"`
	TriggerMethod string    `json:"trigger_method"`
	TriggerSpi    string    `json:"trigger_spi"`
}

func ParsePagination(pageParam, pageSizeParam string) (page int, pageSize int, err error) {
	strPage := pageParam
	strPageSize := pageSizeParam
	page, err = strconv.Atoi(strPage)
	if err != nil {
		return 0, 0, fmt.Errorf("page 和pageSize值不合法")

	}
	pageSize, err = strconv.Atoi(strPageSize)
	if err != nil {
		return 0, 0, fmt.Errorf("page 和pageSize值不合法")
	}
	if page == 0 || pageSize == 0 {
		return 0, 0, fmt.Errorf("page 和pageSize值不合法")
	}
	return page, pageSize, nil
}
