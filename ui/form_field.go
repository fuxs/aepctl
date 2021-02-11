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
	"github.com/fuxs/aepctl/util"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type focusedInputField struct {
	*tview.InputField
	Focused func()
	Parent  *FocusFlex
}

func newFocusedInputField() *focusedInputField {
	return &focusedInputField{
		InputField: tview.NewInputField(),
	}
}

func (f *focusedInputField) SetFocusedFunc(ff *FocusFlex, focused func()) *focusedInputField {
	f.Focused = focused
	f.Parent = ff
	return f
}

func (f *focusedInputField) Focus(delegate func(p tview.Primitive)) {
	f.InputField.Focus(delegate)
	f.Parent.Last = f
	if f.Focused != nil {
		f.Focused()
	}
}

type focusedButton struct {
	*tview.Button
	Focused func()
	Parent  *FocusFlex
}

func newFocusedButton(label string) *focusedButton {
	return &focusedButton{
		Button: tview.NewButton(label),
	}
}

func (f *focusedButton) SetFocusedFunc(ff *FocusFlex, focused func()) *focusedButton {
	f.Focused = focused
	f.Parent = ff
	return f
}

func (f *focusedButton) Focus(delegate func(p tview.Primitive)) {
	f.Button.Focus(delegate)
	f.Parent.Last = f
	if f.Focused != nil {
		f.Focused()
	}
}

// FormField consists of an input fiels and a button.
type FormField struct {
	*tview.Box
	Input       *focusedInputField
	Button      *focusedButton
	ButtonFirst bool
	next        tview.Primitive
	prev        tview.Primitive
}

// SetNext sets the next tview.Primitive in the navigation order
func (f *FormField) SetNext(next Tabable) {
	f.next = next.Primitive()
}

// SetPrev sets the previous tview.Primitive in the navigation order
func (f *FormField) SetPrev(prev Tabable) {
	f.prev = prev.Primitive()
}

// Primitive returns the this object with the tview.Primitive interface
func (f *FormField) Primitive() tview.Primitive {
	return f
}

// SetFocusedFunc sets the Focused functions for the InpurtForm and the Button
func (f *FormField) SetFocusedFunc(ff *FocusFlex, focused ...func()) *FormField {
	l := len(focused)
	if l == 1 {
		f.Input.SetFocusedFunc(ff, focused[0])
		f.Button.SetFocusedFunc(ff, focused[0])
	} else if l > 1 {
		f.Input.SetFocusedFunc(ff, focused[0])
		f.Button.SetFocusedFunc(ff, focused[1])
	}
	return f
}

// Focus sets the focus to the InputForm or Button
func (f *FormField) Focus(delegate func(p tview.Primitive)) {
	if f.ButtonFirst {
		delegate(f.Button)
	} else {
		delegate(f.Input)
	}
}

// Blur forwards the call to the Inputform and Button
func (f *FormField) Blur() {
	f.Input.Blur()
	f.Button.Blur()
}

// After connects this button with the previous Tabable
func (f *FormField) After(previous Tabable) *FormField {
	ConnectTabable(previous, f)
	return f
}

// HasFocus returns true if this primitve has the focus
func (f *FormField) HasFocus() bool {
	return f.Button.HasFocus() || f.Input.HasFocus()
}

// SetText sets the text
func (f *FormField) SetText(text string) *FormField {
	f.Input.SetText(text)
	return f
}

// GetText returns the text of the InputField
func (f *FormField) GetText() string {
	return f.Input.GetText()
}

// Draw draws this primitive
func (f *FormField) Draw(screen tcell.Screen) {
	f.Box.DrawForSubclass(screen, f)
	x, y, width, height := f.GetInnerRect()
	l := len(f.Button.GetLabel()) + 4
	f.Input.SetRect(x, y, width-(l+1), height)
	f.Input.Draw(screen)
	f.Button.SetRect(x+width-l, y, l, 1)
	f.Button.Draw(screen)
}

// InputHandler handles the keyboard input
func (f *FormField) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return f.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		key := event.Key()
		switch key {
		case tcell.KeyUp, tcell.KeyBacktab:
			setFocus(f.prev)
			return
		case tcell.KeyDown, tcell.KeyTab:
			setFocus(f.next)
			return
		default:

		}
		if f.Button.HasFocus() {
			switch key {
			case tcell.KeyRight:
				setFocus(f.next)
			case tcell.KeyLeft:
				setFocus(f.Input)
			case tcell.KeyEnter:
				f.Button.InputHandler()(event, setFocus)
			}
		} else {
			if key == tcell.KeyEnter {
				setFocus(f.next)
			}
			f.Input.InputHandler()(event, setFocus)
		}
	})
}

// MouseHandler handles mouse input
func (f *FormField) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return f.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		if !f.InRect(event.Position()) {
			return false, nil
		}
		consumed, capture = f.Input.MouseHandler()(action, event, setFocus)
		if consumed {
			setFocus(f.Input)
			return
		}
		if action == tview.MouseLeftClick {
			setFocus(f.Button)
		}

		consumed, capture = f.Button.MouseHandler()(action, event, setFocus)
		return
	})
}

/*
// GetFieldWidth returns the width
func (f *FormField) GetFieldWidth() int {
	return 0
}

// GetLabel returns the label
func (f *FormField) GetLabel() string {
	return f.Input.GetLabel()
}

// SetFormAttributes set form specific attributes
func (f *FormField) SetFormAttributes(labelWidth int, labelColor, bgColor, fieldTextColor, fieldBgColor tcell.Color) tview.FormItem {
	f.Input.SetFormAttributes(labelWidth, labelColor, bgColor, fieldTextColor, fieldBgColor)
	f.Button.SetBackgroundColor(fieldBgColor)
	f.Button.SetLabelColor(fieldTextColor)
	return f
}

func (f *FormField) SetFinishedFunc(handler func(key tcell.Key)) tview.FormItem {
	f.Input.SetFinishedFunc(handler)
	return f
}*/

// NewPasteField returns a textfield with a past button
func NewPasteField(label string) *FormField {
	field := newFocusedInputField()
	field.SetLabel(label)
	button := newFocusedButton("Paste")
	button.SetSelectedFunc(func() {
		field.SetText(util.Paste())
	})
	result := &FormField{
		Box:    tview.NewBox(),
		Input:  field,
		Button: button,
	}
	return result
}

// NewFileField returns a textfield with a file
func NewFileField(label string, pages *FocusPages) *FormField {
	field := newFocusedInputField()
	field.SetLabel(label)
	button := newFocusedButton("Select File")
	button.SetSelectedFunc(func() {
		OpenFileDialog(pages, func(s string) { field.SetText(s) })
	})
	result := &FormField{
		Box:         tview.NewBox(),
		ButtonFirst: true,
		Input:       field,
		Button:      button,
	}
	return result
}
