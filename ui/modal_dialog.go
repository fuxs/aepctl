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

// ErrorDialog shows a modal dialog with a message
func ErrorDialog(pages *FocusPages, err error) {
	modal := tview.NewModal().
		SetText(err.Error()).
		AddButtons([]string{"Ok"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			pages.RemovePage("errorDialog")
		}).SetBackgroundColor(tcell.ColorRed)
	pages.AddAndSwitchToPage("errorDialog", modal, true)
}

// InfoDialog shows a modal dialog with a message
func InfoDialog(pages *FocusPages, message string) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"Ok"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			pages.RemovePage("infoDialog")
		}).SetBackgroundColor(tcell.ColorGreen)
	pages.AddAndSwitchToPage("infoDialog", modal, true)
}
