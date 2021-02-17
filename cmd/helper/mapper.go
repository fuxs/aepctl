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
	"github.com/fuxs/aepctl/util"
)

// ChannelSToL maps short channel names to long ones
var ChannelSToL = util.Mapper{
	"email":   "https://ns.adobe.com/xdm/channel-types/email",
	"offline": "https://ns.adobe.com/xdm/channel-types/offline",
	"social":  "https://ns.adobe.com/xdm/channel-types/social",
	"web":     "https://ns.adobe.com/xdm/channel-types/web",
}

// ChannelLToS maps long channel names to short ones
var ChannelLToS = ChannelSToL.Invert()

// ContentSToL maps short content names to long ones
var ContentSToL = util.Mapper{
	"image": "https://ns.adobe.com/experience/offer-management/content-component-imagelink",
	"text":  "https://ns.adobe.com/experience/offer-management/content-component-text",
	"html":  "https://ns.adobe.com/experience/offer-management/content-component-html",
	"json":  "https://ns.adobe.com/experience/offer-management/content-component-json",
}

// ContentLToS maps long content names to short ones
var ContentLToS = ContentSToL.Invert()

// FilterTypeSToL maps short filter types to long ones
var FilterTypeSToL = util.Mapper{
	"all": "allTags",
	"any": "anyTags",
}

// NewNameToInstanceID resolves name to @id
func NewNameToInstanceID(auth *Configuration, schema string) *util.KVCache {
	return NewTemporaryCache(auth.AC, schema, []string{"_instance", "xdm:name"}, []string{"instanceId"})
}

// NewNameToID resolves name to @id
func NewNameToID(auth *Configuration, schema string) *util.KVCache {
	return NewTemporaryCache(auth.AC, schema, []string{"_instance", "xdm:name"}, []string{"_instance", "@id"})
}

// NameToID creates a name -> @id mapper from the passed ObjectCache
func NameToID(koc *util.ObjectCache) util.Mapper {
	return koc.Mapper([]string{"_instance", "xdm:name"}, []string{"_instance", "@id"})
}
