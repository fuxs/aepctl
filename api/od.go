package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/fuxs/aepctl/util"
)

type ODQueryParames struct {
	ContainerID string
	Schema      string
	Query       string
	QOP         string
	Field       string
	OrderBy     string
	Limit       int
}

func (p *ODQueryParames) Params() util.Params {
	var limit string
	if p.Limit > 0 {
		limit = strconv.FormatInt(int64(p.Limit), 10)
	}
	return util.NewParams(
		"containerID", p.ContainerID,
		"schema", p.Schema,
		"q", p.Query,
		"qop", p.QOP,
		"field", p.Field,
		"oderBy", p.OrderBy,
		"limit", limit,
	)
}

func ODQuery(ctx context.Context, p *AuthenticationConfig, params util.Params) (*http.Response, error) {
	containerID := params.Get("containerID")
	if containerID == "" {
		return nil, errors.New("container-id is empty")
	}
	return p.GetRequestRaw(ctx,
		"https://platform.adobe.io/data/core/xcore/%s/queries/core/search%s",
		containerID,
		params.EncodeWithout("containerID"),
	)
}
