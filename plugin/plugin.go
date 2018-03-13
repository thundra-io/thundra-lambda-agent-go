package plugin

import (
	"context"
	"sync"
	"encoding/json"
)

type Plugin interface {
	BeforeExecution(ctx context.Context, request interface{}, wg *sync.WaitGroup)
	AfterExecution(ctx context.Context, request interface{}, response interface{}, err interface{}, wg *sync.WaitGroup) (interface{}, string)
	OnPanic(ctx context.Context, request json.RawMessage, err interface{}, stackTrace []byte, wg *sync.WaitGroup) (interface{}, string)
}