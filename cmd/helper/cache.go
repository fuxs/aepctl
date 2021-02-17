/*
Package helper consists of helping functions.

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
package helper

import (
	"context"
	"strings"

	"github.com/fuxs/aepctl/api/od"
	"github.com/fuxs/aepctl/api/sandbox"
	"github.com/fuxs/aepctl/util"
)

// NewContainerCache creates a new container cache
func NewContainerCache(auth *Configuration) *util.KVCache {
	return util.NewKVCache(
		auth.ReadCache,
		auth.WriteCache,
		func() (map[string]string, error) {
			containers, err := od.ListContainer(context.Background(), auth.NoDryRun())
			if err != nil {
				return nil, err
			}
			query := util.NewQuery(containers)
			m := make(map[string]string)
			query.Path("_embedded", "https://ns.adobe.com/experience/xcore/container").Range(func(q *util.Query) {
				m[q.Str("_instance", "parentName")] = q.Str("instanceId")
			})
			return m, nil
		},
		func() []string { return auth.UniquePath("container.json") })
}

// NewSandboxCache creates a new sandbox cache
func NewSandboxCache(auth *Configuration) *util.StringCache {
	return util.NewStringCache(
		auth.ReadCache,
		auth.WriteCache,
		func() ([]string, error) {
			sandboxes, err := sandbox.List(context.Background(), auth.NoDryRun())
			if err != nil {
				return nil, err
			}
			query := util.NewQuery(sandboxes)
			list := make([]string, 0, 16)
			query.Path("sandboxes").Range(func(q *util.Query) {
				list = append(list, q.Str("name"))
			})
			return list, nil
		},
		func() []string { return auth.UniquePath("sandbox.json") })
}

// NewActivityCache creats a new name -> instanceId cache
func NewActivityCache(ac *AutoContainer) *util.KVCache {
	return newCache(ac, "activities.json", od.ActivitySchema)
}

// NewActivityIDCache creats a new name -> @id cache
func NewActivityIDCache(ac *AutoContainer) *util.KVCache {
	return newIDCache(ac, "activity_ids.json", od.ActivitySchema)
}

// NewTagCache creats a new name -> instanceId cache
func NewTagCache(ac *AutoContainer) *util.KVCache {
	return newCache(ac, "tags.json", od.TagSchema)
}

// NewTagIDCache creats a new name -> @id cache
func NewTagIDCache(ac *AutoContainer) *util.KVCache {
	return newIDCache(ac, "tag_ids.json", od.TagSchema)
}

// NewPlacementCache creats a new name -> instanceId cache
func NewPlacementCache(ac *AutoContainer) *util.KVCache {
	return newCache(ac, "placements.json", od.PlacementSchema)
}

// NewPlacementIDCache creats a new name -> @id cache
func NewPlacementIDCache(ac *AutoContainer) *util.KVCache {
	return newIDCache(ac, "placement_ids.json", od.PlacementSchema)
}

// NewOfferCache creats a new name -> instanceId cache
func NewOfferCache(ac *AutoContainer) *util.KVCache {
	return newCache(ac, "offers.json", od.OfferSchema)
}

// NewOfferIDCache creats a new name -> @id cache
func NewOfferIDCache(ac *AutoContainer) *util.KVCache {
	return newIDCache(ac, "offer_ids.json", od.OfferSchema)
}

// NewFallbackCache creats a new name -> instanceId cache
func NewFallbackCache(ac *AutoContainer) *util.KVCache {
	return newCache(ac, "fallbacks.json", od.FallbackSchema)
}

// NewFallbackIDCache creats a new name -> @id cache
func NewFallbackIDCache(ac *AutoContainer) *util.KVCache {
	return newIDCache(ac, "fallback_ids.json", od.FallbackSchema)
}

// NewCollectionCache creats a new name -> instanceId cache
func NewCollectionCache(ac *AutoContainer) *util.KVCache {
	return newCache(ac, "collections.json", od.CollectionSchema)
}

// NewCollectionIDCache creats a new name -> @id cache
func NewCollectionIDCache(ac *AutoContainer) *util.KVCache {
	return newIDCache(ac, "collection_ids.json", od.CollectionSchema)
}

// NewRuleCache creats a new name -> instanceId cache
func NewRuleCache(ac *AutoContainer) *util.KVCache {
	return newCache(ac, "rules.json", od.RuleSchema)
}

// NewRuleIDCache creats a new name -> @id cache
func NewRuleIDCache(ac *AutoContainer) *util.KVCache {
	return newIDCache(ac, "rule_ids.json", od.RuleSchema)
}

func newCache(ac *AutoContainer, file, schema string) *util.KVCache {
	return NewCache(ac, file, schema, []string{"_instance", "xdm:name"}, []string{"instanceId"})
}

func newIDCache(ac *AutoContainer, file, schema string) *util.KVCache {
	return NewCache(ac, file, schema, []string{"_instance", "xdm:name"}, []string{"_instance", "@id"})
}

// NewCache creates a new persistent cache for JSON objects
func NewCache(ac *AutoContainer, file, schema string, key []string, value []string) *util.KVCache {
	auth := ac.Auth
	return util.NewKVCache(
		auth.ReadCache,
		auth.WriteCache,
		func() (map[string]string, error) {
			if err := ac.AutoFillContainer(); err != nil {
				return nil, err
			}
			objs, err := od.List(context.Background(), auth.NoDryRun(), ac.ContainerID, schema)
			if err != nil {
				return nil, err
			}
			query := util.NewQuery(objs)
			m := make(map[string]string)
			query.Path("_embedded", "results").Range(func(q *util.Query) {
				m[q.Str(key...)] = q.Str(value...)
			})
			return m, nil
		},
		func() []string { return ac.UniquePath(file) })
}

// NewTemporaryCache creates  a new temporary cache for JSON objects
func NewTemporaryCache(ac *AutoContainer, schema string, key []string, value []string) *util.KVCache {
	auth := ac.Auth
	return util.NewKVCache(
		auth.ReadCache,
		auth.WriteCache,
		func() (map[string]string, error) {
			if err := ac.AutoFillContainer(); err != nil {
				return nil, err
			}
			objs, err := od.List(context.Background(), auth.NoDryRun(), ac.ContainerID, schema)
			if err != nil {
				return nil, err
			}
			query := util.NewQuery(objs)
			m := make(map[string]string)
			query.Path("_embedded", "results").Range(func(q *util.Query) {
				m[strings.TrimSpace(q.Str(key...))] = q.Str(value...)
			})
			return m, nil
		}, nil)
}

// NewPlacementsCache creates a cache for placement objects
func NewPlacementsCache(ac *AutoContainer) *util.ObjectCache {
	return NewObjectCache(ac, od.PlacementSchema)
}

// NewObjectCache creates a new object cache
func NewObjectCache(ac *AutoContainer, schema string) *util.ObjectCache {
	auth := ac.Auth
	return util.NewObjectCache(func() ([]interface{}, error) {
		if err := ac.AutoFillContainer(); err != nil {
			return nil, err
		}
		objs, err := od.List(context.Background(), auth.NoDryRun(), ac.ContainerID, schema)
		if err != nil {
			return nil, err
		}
		query := util.NewQuery(objs).Path("_embedded", "results")
		m := make([]interface{}, query.Length())
		query.RangeI(func(i int, q *util.Query) {
			m[i] = q.Interface()
		})
		return m, nil
	})
}
