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
package cache

import (
	"time"

	"github.com/fuxs/aepctl/util"
)

// ODNameToID returns a transformation xdm:name -> @id
func ODNameToID() *TransformMap {
	return NewTransformMap("_embedded", "results").
		K("_instance", "xdm:name").
		V("_instance", "@id")
}

// ODNameToID returns a transformation xdm:name -> instanceId
func ODNameToInstanceID() *TransformMap {
	return NewTransformMap("_embedded", "results").
		K("_instance", "xdm:name").
		V("instanceId")
}

// NewODNameToID creates an initialzed MapFileCache object with a xdm:name -> @id
func NewODNameToID(ac *AutoContainer, name, schema string, pp util.PathProvider) *MapFileCache {
	call := NewODCall(ac, schema)
	t := ODNameToID()
	return NewMapFileCache(call, t, time.Minute, util.NewLazyPath(pp, name+"_n2id.json"))
}

// NewODNameToInstanceID creates an initialzed MapFileCache object with a xdm:name -> instanceId
func NewODNameToInstanceID(ac *AutoContainer, name, schema string, pp util.PathProvider) *MapFileCache {
	call := NewODCall(ac, schema)
	t := ODNameToInstanceID()
	return NewMapFileCache(call, t, time.Minute, util.NewLazyPath(pp, name+"_n2iid.json"))
}

// NewODNameToIDMem creates an initialzed MemCache object with a xdm:name -> @id
func NewODNameToIDMem(ac *AutoContainer, schema string) *MapMemCache {
	return NewMapMemCache(NewODCall(ac, schema), ODNameToID())
}

// NewODNameToIDMem creates an initialzed MemCache object with a xdm:name -> instanceId
func NewODNameToInstanceIDMem(ac *AutoContainer, schema string) *MapMemCache {
	call := NewODCall(ac, schema)
	t := ODNameToInstanceID()
	return NewMapMemCache(call, t)
}

// NewODTrans returns a transformation with the right path for offer
// decisionsing HAL responses
func NewODTrans() *TransformMap {
	return NewTransformMap("_embedded", "results")
}
