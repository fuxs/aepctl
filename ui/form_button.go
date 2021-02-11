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

// FormButton implements a button for home-made forms
type FormButton struct {
	*tview.Button
	Focused func()
	Parent  *FocusFlex
	next    tview.Primitive
	prev    tview.Primitive
}

// NewFormButton creates an initialzed FormButton object
func NewFormButton(label string) *FormButton {
	button := tview.NewButton(label)
	return &FormButton{
		Button: button,
	}
}

// SetNext sets the next tview.Primitive in the navigation order
func (f *FormButton) SetNext(next Tabable) {
	f.next = next.Primitive()
}

// SetPrev sets the previous tview.Primitive in the navigation order
func (f *FormButton) SetPrev(prev Tabable) {
	f.prev = prev.Primitive()
}

// Primitive returns the this object with the tview.Primitive interface
func (f *FormButton) Primitive() tview.Primitive {
	return f
}

// SetFocusedFunc sets the function which will be called when this FormButton
// receives the focus
func (f *FormButton) SetFocusedFunc(ff *FocusFlex, focused func()) *FormButton {
	f.Focused = focused
	f.Parent = ff
	return f
}

// After connects this button with the previous Tabable
func (f *FormButton) After(previous Tabable) *FormButton {
	ConnectTabable(previous, f)
	return f
}

// Focus calls the original Button Focus function, stores this FormButton as last
// focused object and calls the Focused function if set
func (f *FormButton) Focus(delegate func(p tview.Primitive)) {
	f.Button.Focus(delegate)
	f.Parent.Last = f
	if f.Focused != nil {
		f.Focused()
	}
}

// InputHandler handles the keyboard input
func (f *FormButton) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return f.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		key := event.Key()
		switch key {
		case tcell.KeyUp, tcell.KeyLeft, tcell.KeyBacktab:
			setFocus(f.prev)
		case tcell.KeyDown, tcell.KeyRight, tcell.KeyTab:
			setFocus(f.next)
		default:
			f.Button.InputHandler()(event, setFocus)
		}
	})
}
