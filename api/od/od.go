/*
Package od contains offer decisiong related functions.

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
package od

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/util"
)

const (
	// ActivitySchema references the activity schema
	ActivitySchema = "https://ns.adobe.com/experience/offer-management/offer-activity" //;version=0.5
	// TagSchema references the tag schema
	TagSchema = "https://ns.adobe.com/experience/offer-management/tag" // ;version=0.1
	// CollectionSchema references the collection schema
	CollectionSchema = "https://ns.adobe.com/experience/offer-management/offer-filter" //;version=0.2
	// PlacementSchema references the placement schema
	PlacementSchema = "https://ns.adobe.com/experience/offer-management/offer-placement" //;version=0.4
	// OfferSchema references the offer schema
	OfferSchema = "https://ns.adobe.com/experience/offer-management/personalized-offer" //;version=0.6
	// FallbackSchema references the fallback schema
	FallbackSchema = "https://ns.adobe.com/experience/offer-management/fallback-offer" //;version=0.1
	// RuleSchema references the rule schema
	RuleSchema = "https://ns.adobe.com/experience/offer-management/eligibility-rule" // ;version=0.3
)

// Update is used for file based updates
type Update struct {
	IDs   []string           `json:"ids" yaml:"ids"`
	Apply []*UpdateOperation `json:"apply" yaml:"apply"`
}

// UpdateOperation defines the update operation
type UpdateOperation struct {
	Operation string `json:"op" yaml:"op"`
	Path      string `json:"path" yaml:"path"`
	Value     string `json:"value" yaml:"value"`
}

// ListContainer returns a list of container
func ListContainer(ctx context.Context, p *api.AuthenticationConfig) (interface{}, error) {
	return p.GetRequest(ctx, "https://platform.adobe.io/data/core/xcore/?product=acp&property=_instance.containerType==decisioning")
}

// Get returns a collection by name (wild cards are supported) or id (exact match)
func Get(ctx context.Context, p *api.AuthenticationConfig, containerID, schema, id string) (interface{}, error) {
	if containerID == "" {
		return nil, errors.New("container-id is empty")
	}
	if schema == "" {
		return nil, errors.New("schema is empty")
	}
	return p.GetRequest(ctx,
		"https://platform.adobe.io/data/core/xcore/%s/instances%s",
		containerID,
		util.Par("schema", schema, "id", id),
	)
}

// Query queries all objects
func Query(ctx context.Context, p *api.AuthenticationConfig, containerID, schema, q, qop, field, orderBy, limit string) (interface{}, error) {
	if containerID == "" {
		return nil, errors.New("container-id is empty")
	}
	if schema == "" {
		return nil, errors.New("schema is empty")
	}
	return p.GetRequest(ctx,
		"https://platform.adobe.io/data/core/xcore/%s/queries/core/search%s",
		containerID,
		util.Par("schema", schema, "q", q, "qop", qop, "field", field, "orderBy", orderBy, "limit", limit),
	)
}

// List lists all objects
func List(ctx context.Context, p *api.AuthenticationConfig, containerID, schema string) (interface{}, error) {
	if containerID == "" {
		return nil, errors.New("container-id is empty")
	}
	if schema == "" {
		return nil, errors.New("schema is empty")
	}
	return p.GetRequest(ctx,
		"https://platform.adobe.io/data/core/xcore/%s/queries/core/search%s",
		containerID,
		util.Par("schema", schema),
	)
}

// Delete deletes the object with the passed id
func Delete(ctx context.Context, p *api.AuthenticationConfig, containerID, id string) error {
	if containerID == "" {
		return errors.New("container-id is empty")
	}
	_, err := p.DeleteRequest(ctx,
		"https://platform.adobe.io/data/core/xcore/%s/instances/%s",
		containerID,
		id,
	)
	return err
}

// Patch patches the object with the passed id
func Patch(ctx context.Context, p *api.AuthenticationConfig, containerID, id, schema string, ops ...*UpdateOperation) (interface{}, error) {
	if containerID == "" {
		return nil, errors.New("container-id is empty")
	}
	if schema == "" {
		return nil, errors.New("schema is empty")
	}
	data, err := json.Marshal(ops)
	if err != nil {
		return nil, err
	}
	ct := fmt.Sprintf(`application/vnd.adobe.platform.xcore.patch.hal+json; version=1; schema="%s"`, schema)
	return p.PatchRequest(ctx, map[string]string{"Content-Type": ct}, data,
		"https://platform.adobe.io/data/core/xcore/%s/instances/%s",
		containerID,
		id,
	)
}

// Create creates a new object
func Create(ctx context.Context, p *api.AuthenticationConfig, containerID, schema string, obj interface{}) (interface{}, error) {
	if containerID == "" {
		return nil, errors.New("container-id is empty")
	}
	if schema == "" {
		return nil, errors.New("schema is empty")
	}
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	ct := fmt.Sprintf(`application/schema-instance+json; version=1; schema="%s"`, schema)
	return p.PostRequest(ctx, map[string]string{"Content-Type": ct}, data,
		"https://platform.adobe.io/data/core/xcore/%s/instances",
		containerID)
}
