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

// StackInt is a simple int stack (lifo)
type StackInt []int

// Push pushes the integer i to the stack
func (s *StackInt) Push(i int) {
	*s = append(*s, i)
}

// Pop removes the last element from the stack and returns it
func (s *StackInt) Pop() (int, bool) {
	l := len(*s)
	if l == 0 {
		return 0, false
	}
	result := (*s)[l-1]
	*s = (*s)[:l-1]
	return result, true
}
