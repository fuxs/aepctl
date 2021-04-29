/*
Package catalog consists of catalog functions.

Copyright 2021 Michael Bungenstock

Licensed under the Apache License, Version 2.0 (the "License"); you may not use
this file except in compliance with the License. You may obtain a copy of the
License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied. See the License for the
specific language governing permissions and limitations under the License.
*/
package api

import (
	"context"
	"net/http"
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
func CatalogCreateDataset(ctx context.Context, p *AuthenticationConfig, obj interface{}) (*http.Response, error) {
	return p.PostJSONRequestRaw(ctx, obj, "https://platform.adobe.io/data/foundation/catalog/dataSets")
}

func CatalogCreateProfileUnionDataset(ctx context.Context, p *AuthenticationConfig, name, format string) (*http.Response, error) {
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
