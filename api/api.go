package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/fuxs/aepctl/util"
)

type Func func(context.Context, *AuthenticationConfig, util.Params) (*http.Response, error)
type FuncTable map[string]Func

var All = map[string]Func{
	"ODQuery":            ODQuery,
	"FlowGetConnections": FlowGetConnections,
}

func (ft FuncTable) Call(name string, ctx context.Context, auth *AuthenticationConfig, params util.Params) (*http.Response, error) {
	f := ft[name]
	if f == nil {
		return nil, fmt.Errorf("function %s not found", name)
	}
	return f(ctx, auth, params)
}
