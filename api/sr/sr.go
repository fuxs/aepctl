package sr

import (
	"context"

	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/util"
)

// GetStats returns schema registry informations
func GetStats(ctx context.Context, p *api.AuthenticationConfig) (util.JSONResponse, error) {
	res, err := p.GetRequestRaw(ctx, "https://platform.adobe.io/data/foundation/schemaregistry/stats")
	if err != nil {
		return nil, err
	}
	return util.NewJSONIterator(res.Body)
}
