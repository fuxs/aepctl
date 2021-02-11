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

// FocusPages fixes a problem with greedy MouseHandler implementations. If you
// want to switch a page in a SelectFunc called by the MouseHandler, then this
// is not possible. Table for example sets back the focus to itself after the
// execution and that is not what we want.
type FocusPages struct {
	*tview.Pages
}

// NewFocusPages returns an initialized FocusPages object
func NewFocusPages() *FocusPages {
	return &FocusPages{
		Pages: tview.NewPages(),
	}
}

// MouseHandler implements the fix. If the numer of pages changes then the focus
// will be set to the current page.
func (f *FocusPages) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		i := f.Pages.GetPageCount()
		consumed, capture = f.Pages.MouseHandler()(action, event, setFocus)
		j := f.Pages.GetPageCount()
		if i != j {
			setFocus(f.Pages)
		}
		return
	}
}

// FocusFlex stores the Last primitive with focus and returns the focus to it if
// method Focus have been called
type FocusFlex struct {
	*tview.Flex
	Last tview.Primitive
}

// NewFocusFlex creates an initialized FocusFlex object
func NewFocusFlex() *FocusFlex {
	return &FocusFlex{
		Flex: tview.NewFlex(),
	}
}

// Focus overrules the original Flex Focus method
func (f *FocusFlex) Focus(delegate func(p tview.Primitive)) {
	if f.Last == nil {
		f.Flex.Focus(delegate)
		return
	}
	delegate(f.Last)
}
