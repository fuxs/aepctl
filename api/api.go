package api

import (
	"context"
	"net/http"

	"github.com/fuxs/aepctl/util"
)

type Func func(context.Context, *AuthenticationConfig, util.Params) (*http.Response, error)
