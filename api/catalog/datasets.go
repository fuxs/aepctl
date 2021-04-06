package catalog

import (
	"context"
	"net/http"

	"github.com/fuxs/aepctl/api"
)

type FileDescription struct {
	Format          string `json:"format,omitempty" yaml:"format,omitempty"`
	ContainerFormat string `json:"containerFormat,omitempty" yaml:"containerFormat,omitempty"`
	Persisted       bool   `json:"persisted,omitempty" yaml:"persisted,omitempty"`
}

type SchemaRef struct {
	ID          string `json:"id,omitempty" yaml:"id,omitempty"`
	ContentType string `json:"contentType,omitempty" yaml:"contentType,omitempty"`
}

type DataSetRequest struct {
	ConncectorID    string           `json:"connectorId,omitempty" yaml:"connectorId,omitempty"`
	ConncectoionID  string           `json:"connectionId,omitempty" yaml:"connectionId,omitempty"`
	Name            string           `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace       string           `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	FileDescription *FileDescription `json:"fileDescription,omitempty" yaml:"fileDescription,omitempty"`
	SchemaRef       *SchemaRef       `json:"schemaRef,omitempty" yaml:"schemaRef,omitempty"`
}

// Create creates a new object
func CreateDataset(ctx context.Context, p *api.AuthenticationConfig, obj interface{}) (*http.Response, error) {
	return p.PostJSONRequestRaw(ctx, obj, "https://platform.adobe.io/data/foundation/catalog/dataSets")
}

func CreateProfileUnionDataset(ctx context.Context, p *api.AuthenticationConfig, name, format string) (*http.Response, error) {
	obj := &DataSetRequest{
		Name: name,
		SchemaRef: &SchemaRef{
			ID:          "https://ns.adobe.com/xdm/context/profile__union",
			ContentType: "application/vnd.adobe.xed+json;version=1",
		},
		FileDescription: &FileDescription{
			Persisted: true,
			//ContainerFormat: format,
			Format: format,
		},
	}
	return p.PostJSONRequestRaw(ctx, obj, "https://platform.adobe.io/data/foundation/catalog/dataSets")
}
