/*
Package util util consists of general utility functions and structures.

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
package util

// NVPTransformer is a generic transformer to Name, Value and Path
type NVPTransformer struct{}

func (*NVPTransformer) Header(wide bool) []string {
	if wide {
		return []string{"NAME", "VALUE", "PATH"}
	}
	return []string{"PATH", "VALUE"}
}

// Preprocess has nothing to do
func (*NVPTransformer) Preprocess(i JSONResponse) error {
	return nil
}

// WriteRow writes name, value and path
// TODO make Truncate configurable
func (t *NVPTransformer) WriteRow(q *Query, w *RowWriter, wide bool) error {
	if wide {
		return w.Write(
			q.JSONName(),
			q.String(),
			q.JSONPath(),
		)
	}
	return w.Write(
		q.JSONFullPath(),
		q.String(),
	)
}

// Iterator returns a JSONValueIterator returning all values
func (*NVPTransformer) Iterator(c *JSONCursor) (JSONResponse, error) {
	return NewJSONValueIterator(c, nil), nil
}
