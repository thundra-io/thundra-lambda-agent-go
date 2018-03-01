package thundra

import (
	"context"
	"sync"
	"encoding/json"
)

type Plugin interface {
	BeforeExecution(ctx context.Context, request interface{}, wg *sync.WaitGroup)
	AfterExecution(ctx context.Context, request interface{}, response interface{}, error interface{}, wg *sync.WaitGroup)
	OnPanic(ctx context.Context, request json.RawMessage, panic *ThundraPanic, wg *sync.WaitGroup)
}

type ThundraPanic struct {
	ErrInfo    error
	StackTrace string
	ErrType    string
}

const timeFormat = "2006-01-02 15:04:05.000 -0700"
