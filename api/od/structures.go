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

// ActivityCriteria contains the criteria of an activity
type ActivityCriteria struct {
	Placements []string `json:"xdm:placements" yaml:"placements"`
	Selection  struct {
		Filter string `json:"xdm:filter" yaml:"filter"`
	} `json:"xdm:optionSelection" yaml:"selection"`
}

// Activity contains the logic that informs the selection of an offer
type Activity struct {
	Name      string              `json:"xdm:name" yaml:"name"`
	StartDate string              `json:"xdm:startDate" yaml:"startDate"`
	EndDate   string              `json:"xdm:endDate" yaml:"endDate"`
	Status    string              `json:"xdm:status" yaml:"status"`
	Criteria  []*ActivityCriteria `json:"xdm:criteria" yaml:"criteria"`
	Fallback  string              `json:"xdm:fallback" yaml:"fallback"`
}

// Collection filters a set of offers
type Collection struct {
	Name   string   `json:"xdm:name" yaml:"name"`
	Filter string   `json:"xdm:filterType" yaml:"filter"`
	IDs    []string `json:"xdm:ids" yaml:"ids"`
}

// Fallback represents a default offer
type Fallback struct {
	Name            string                 `json:"xdm:name" yaml:"name"`
	Status          string                 `json:"xdm:status" yaml:"status"`
	Representations []*OfferRepresentation `json:"xdm:representations" yaml:"representations"`
	Tags            []string               `json:"xdm:tags,omitempty" yaml:"tags,omitempty"`
	Attributes      map[string]string      `json:"xdm:characteristics,omitempty" yaml:"attributes,omitempty"`
}

// OfferComponent contains the content of the offer
type OfferComponent struct {
	Name      string   `json:"repo:name" yaml:"name"`
	Type      string   `json:"@type" yaml:"type"`
	Format    string   `json:"dc:format" yaml:"format"`
	Content   string   `json:"xdm:content,omitempty" yaml:"content,omitempty"`
	URL       string   `json:"xdm:deliveryURL,omitempty" yaml:"url,omitempty"`
	Link      string   `json:"xdm:linkURL,omitempty" yaml:"link,omitempty"`
	Languages []string `json:"dc:language" yaml:"language"`
}

// OfferRepresentation connects content and placement
type OfferRepresentation struct {
	Components []*OfferComponent `json:"xdm:components,omitempty" yaml:"components,omitempty"`
	Placement  string            `json:"xdm:placement" yaml:"placement"`
	Channel    string            `json:"xdm:channel,omitempty" yaml:"channel,omitempty"`
}

// OfferConstraint contains the offer rule
type OfferConstraint struct {
	StartDate string `json:"xdm:startDate" yaml:"startDate"`
	EndDate   string `json:"xdm:endDate" yaml:"endDate"`
	Rule      string `json:"xdm:eligibilityRule" yaml:"rule"`
}

// Offer is a customizable marketing message based on eligibility rules and constraints
type Offer struct {
	Name            string                 `json:"xdm:name" yaml:"name"`
	Status          string                 `json:"xdm:status" yaml:"status"`
	Representations []*OfferRepresentation `json:"xdm:representations" yaml:"representations"`
	Constraint      *OfferConstraint       `json:"xdm:selectionConstraint" yaml:"constraint"`
	Rank            struct {
		Priority int `json:"xdm:priority" yaml:"priority"`
	} `json:"xdm:rank" yaml:"rank"`
	Capping *struct {
		Global int `json:"xdm:globalCap,omitempty" yaml:"global,omitempty"`
	} `json:"xdm:cappingConstraint,omitempty" yaml:"capping,omitempty"`
	Tags       []string          `json:"xdm:tags,omitempty" yaml:"tags,omitempty"`
	Attributes map[string]string `json:"xdm:characteristics,omitempty" yaml:"attributes,omitempty"`
}

// Placement is a container that is used to showcase your offers
type Placement struct {
	Name        string `json:"xdm:name" yaml:"name"`
	Content     string `json:"xdm:componentType" yaml:"content"`
	Channel     string `json:"xdm:channel" yaml:"channel"`
	Description string `json:"xdm:description" yaml:"description"`
}

// RuleCondition contains the definition of the rule
type RuleCondition struct {
	Type   string `json:"type" yaml:"type"`
	Format string `json:"format" yaml:"format"`
	Value  string `json:"value" yaml:"value"`
}

// RuleVersionedSchema ...
type RuleVersionedSchema struct {
	Ref     string `json:"$ref" yaml:"ref"`
	Version string `json:"version" yaml:"version"`
}

// RuleProfile ...
type RuleProfile struct {
	Schema *RuleVersionedSchema `json:"xdm:schema" yaml:"schema"`
	Paths  []string             `json:"xdm:referencePaths" yaml:"paths"`
}

// RuleDefinedOn ...
type RuleDefinedOn struct {
	Profile *RuleProfile `json:"profile" yaml:"profile"`
}

// Rule is used to determine eligibility
type Rule struct {
	Name        string         `json:"xdm:name" yaml:"name"`
	Description string         `json:"xdm:description" yaml:"description"`
	Condition   *RuleCondition `json:"condition" yaml:"condition"`
	On          *RuleDefinedOn `json:"xdm:definedOn" yaml:"on"`
}

// Tag allow you to better organize and sort through your offers
type Tag struct {
	Name string `json:"xdm:name" yaml:"name"`
}
