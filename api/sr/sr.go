package sr

import (
	"context"
	"net/http"

	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/util"
)

// GetStats returns schema registry informations
func GetStatsRaw(ctx context.Context, p *api.AuthenticationConfig) (*http.Response, error) {
	return p.GetRequestRaw(ctx, "https://platform.adobe.io/data/foundation/schemaregistry/stats")
}

// GetStats returns schema registry informations
func GetStats(ctx context.Context, p *api.AuthenticationConfig) (util.JSONResponse, error) {
	return api.NewJSONIterator(p.GetRequestRaw(ctx, "https://platform.adobe.io/data/foundation/schemaregistry/stats"))
}
