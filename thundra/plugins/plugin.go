package plugins

import (
	"context"
	"sync"
	"fmt"
	"time"
)

type Plugin interface {
	BeforeExecution(ctx context.Context, request interface{}, wg *sync.WaitGroup)
	AfterExecution(ctx context.Context, request interface{}, response interface{}, error interface{}, wg *sync.WaitGroup)
}

func ParseTimeFormat(t time.Time) time.Time {
	date, err := time.Parse("2006-01-02 15:04:05.000 -07:00", t.Format("2006-01-02 15:04:05.000 -07:00"))
	if err!= nil{
		fmt.Println("parseTimeForat Error",err)
	}
	return date
}

const timeFormat = "2006-01-02 15:04:05.000 -07:00"