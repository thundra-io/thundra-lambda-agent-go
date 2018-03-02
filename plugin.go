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

type ThundraError struct {
	Error string `json:"error"`// Error stack trace
	ErrorMessage string `json:"errorMessage"`
	ErrorType string `json:"errorType"`
}