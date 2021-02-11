/*
Package ui consists of console ui components

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
package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Tabable should be implemented by components supporting the tab key for
// navigation
type Tabable interface {
	SetNext(Tabable)
	SetPrev(Tabable)
	Primitive() tview.Primitive
}

// TabButton extends the tview.Button with a tab navigatioin
type TabButton struct {
	*tview.Button
	next tview.Primitive
	prev tview.Primitive
}

// ConnectTabable connects two Tabalbe components
func ConnectTabable(prev, next Tabable) {
	prev.SetNext(next)
	next.SetPrev(prev)
}

// NewTabButton creates an initialized TabButton object
func NewTabButton(label string) *TabButton {
	result := &TabButton{
		Button: tview.NewButton(label),
	}
	return result
}

// SetNext sets the next tview.Primitive in the navigation order
func (t *TabButton) SetNext(next Tabable) {
	t.next = next.Primitive()
}

// SetPrev sets the previous tview.Primitive in the navigation order
func (t *TabButton) SetPrev(prev Tabable) {
	t.prev = prev.Primitive()
}

// Primitive returns the this object with the tview.Primitive interface
func (t *TabButton) Primitive() tview.Primitive {
	return t
}

// InputHandler captures the tab events, other events are handled by the
// orignial InputHandler
func (t *TabButton) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return t.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		switch event.Key() {
		case tcell.KeyTab:
			setFocus(t.next)
		case tcell.KeyBacktab:
			setFocus(t.prev)
		default:
			t.Button.InputHandler()(event, setFocus)
		}
	})
}
