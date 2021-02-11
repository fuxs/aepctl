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
	"os"
	"strings"

	"github.com/fuxs/aepctl/util"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type tableComp struct {
	*tview.Table
	dir          *util.Dir
	hidden       bool
	separated    bool
	byDate       bool
	byName       bool
	bySize       bool
	ascending    bool
	pages        *FocusPages
	selectedFunc func(string)
	next         tview.Primitive
	prev         tview.Primitive
}

func (t *tableComp) SetNext(next Tabable) {
	t.next = next.Primitive()
}

func (t *tableComp) SetPrev(prev Tabable) {
	t.prev = prev.Primitive()
}

func (t *tableComp) Primitive() tview.Primitive {
	return t
}

func (t *tableComp) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return t.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		switch event.Key() {
		case tcell.KeyTab:
			setFocus(t.next)
		case tcell.KeyBacktab:
			setFocus(t.prev)
		default:
			t.Table.InputHandler()(event, setFocus)
		}
	})
}

func (t *tableComp) selectAction() bool {
	row, _ := t.Table.GetSelection()
	hp := t.dir.HasParent()

	if row == 1 && hp {
		parent, err := t.dir.Parent()
		if err != nil {
			ErrorDialog(t.pages, err)
			return true
		}
		t.dir = parent
		t.reloadView()
		return true
	}
	row--
	if hp {
		row--
	}
	if t.dir.IsDir(row) {
		child, err := t.dir.ChildI(row)
		if err != nil {
			ErrorDialog(t.pages, err)
			return true
		}
		t.dir = child
		t.reloadView()
	} else {
		// "returns" the value by passing the value to the provided function
		if t.selectedFunc != nil {
			t.selectedFunc(t.dir.PathI(row))
		}
		t.pages.RemovePage("fileDialog")
	}
	return true
}

func (t *tableComp) exitAction() {
	t.pages.RemovePage("fileDialog")
}

func newTableComp(pages *FocusPages, selectedFunc func(string)) (*tableComp, error) {
	t := &tableComp{
		hidden:       false,
		separated:    true,
		byName:       true,
		ascending:    true,
		pages:        pages,
		selectedFunc: selectedFunc,
	}

	table := tview.NewTable()

	table.SetSelectable(true, false).Select(1, 0).SetFixed(1, 0).
		SetSelectionChangedFunc(func(row, column int) {
			if row == 0 {
				table.Select(1, 0)
				return
			}
			if row == -1 {
				table.Select(table.GetRowCount()-1, 0)
			}
		}).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyRune {
				switch event.Rune() {
				case 'd':
					if t.byDate {
						t.ascending = !t.ascending
					} else {
						t.byDate, t.byName, t.bySize, t.ascending = true, false, false, true
					}
					t.updateView()
					return nil
				case 'n':
					if t.byName {
						t.ascending = !t.ascending
					} else {
						t.byDate, t.byName, t.bySize, t.ascending = false, true, false, true
					}
					t.updateView()
					return nil
				case 's':
					if t.bySize {
						t.ascending = !t.ascending
					} else {
						t.byDate, t.byName, t.bySize, t.ascending = false, false, true, true
					}
					t.updateView()
					return nil
				case 'h':
					t.hidden = !t.hidden
					if !t.hidden {
						t.Table.Clear()
					}
					t.updateView()
					return nil
				case 'f':
					t.separated = !t.separated
					t.updateView()
					return nil
				}
			} else if event.Key() == tcell.KeyEnter {
				if t.selectAction() {
					return nil
				}
			}
			return event
		}).SetBorder(true)
	table.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		if action == tview.MouseLeftDoubleClick {
			_, y := event.Position()
			if y > 1 && y < table.GetRowCount()+2 {
				if t.selectAction() {
					return tview.MouseLeftClick, event
				}
			}
		}
		return action, event
	})

	t.Table = table
	dir, err := util.NewDir()
	if err != nil {
		return nil, err
	}
	t.dir = dir

	/*table.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEsc {
			t.exitAction()
		}
	})*/

	t.updateView()
	return t, nil
}

func (t *tableComp) updateRow(i int, f os.FileInfo) {
	var name string
	var color tcell.Color
	if f.IsDir() {
		name = "/" + f.Name()
		color = tcell.ColorBlue
	} else {
		name = f.Name()
		color = tcell.ColorLightGray
	}
	t.Table.SetCell(i, 0, tview.NewTableCell(name).SetTextColor(color))
	t.Table.SetCell(i, 1, tview.NewTableCell(f.ModTime().Format("02.01.06 15:04")).SetTextColor(color))
	t.Table.SetCell(i, 2, tview.NewTableCell(util.ByteCountSI(f.Size())).SetTextColor(color))
}

func (t *tableComp) reloadView() {
	t.Table.Clear()
	t.updateView()
	t.Table.Select(1, 0).ScrollToBeginning()
}

func buildHeader(text string, set, asc bool) string {
	if !set {
		return text
	}
	var b strings.Builder
	b.WriteString(text)
	if asc {
		b.WriteString(" ▲")
	} else {
		b.WriteString(" ▼")
	}
	return b.String()
}

func (t *tableComp) updateView() {

	t.Table.SetTitle(t.dir.Path).SetTitleAlign(tview.AlignLeft)
	t.Table.SetCell(0, 0, tview.NewTableCell(buildHeader("Name", t.byName, t.ascending)).SetTextColor(tcell.ColorYellow)).
		SetCell(0, 1, tview.NewTableCell(buildHeader("Mod. Date", t.byDate, t.ascending)).SetTextColor(tcell.ColorYellow)).
		SetCell(0, 2, tview.NewTableCell(buildHeader("Size", t.bySize, t.ascending)).SetTextColor(tcell.ColorYellow))
	offset := 1
	if t.dir.HasParent() {
		t.Table.SetCell(1, 0, tview.NewTableCell("/..").SetTextColor(tcell.ColorBlue)).
			SetCell(1, 1, tview.NewTableCell("").SetTextColor(tcell.ColorBlue)).
			SetCell(1, 2, tview.NewTableCell("").SetTextColor(tcell.ColorBlue))
		offset = 2
	}
	if t.byName {
		for i, f := range t.dir.SortedByName(t.hidden, t.separated, t.ascending) {
			t.updateRow(i+offset, f)
		}
	} else if t.byDate {
		for i, f := range t.dir.SortedByModTime(t.hidden, t.separated, t.ascending) {
			t.updateRow(i+offset, f)
		}
	} else if t.bySize {
		for i, f := range t.dir.SortedBySize(t.hidden, t.separated, t.ascending) {
			t.updateRow(i+offset, f)
		}
	}
}

var helpTextShort = util.Form(`
	[white]Press [green]F1
	[white]for help`)

var helpText = util.Form(`
	[yellow]Navigate
	[green]UP, DOWN
	  
	[yellow]Select
	[green]ENTER, DBLCLK
	  
	[yellow]Cancel
	[green]ESC

	[yellow]Sort
	[green]N [lightgray]by name
	[green]D [lightgray]by date
	[green]S [lightgray]by size
					
	Press repeated
	to toggle

	[yellow]Other
	[green]H [lightgray]hidden files
	[green]F [lightgray]dirs first

	[yellow]Help
	[green]F1 [lightgray]hide
					
`)

// OpenFileDialog opens a new file dialog
func OpenFileDialog(pages *FocusPages, selectedFunc func(string)) {
	comp, err := newTableComp(pages, selectedFunc)
	if err != nil {
		ErrorDialog(pages, err)
		return
	}

	help := tview.NewTextView().SetDynamicColors(true).SetText(helpTextShort)
	help.SetBorderPadding(1, 1, 2, 1).SetBorder(false)

	selectButton := NewTabButton("Select")
	selectButton.SetSelectedFunc(func() {
		comp.selectAction()
	})
	cancelButton := NewTabButton("Cancel")
	cancelButton.SetSelectedFunc(func() {
		comp.exitAction()
	})
	ConnectTabable(comp, selectButton)
	ConnectTabable(selectButton, cancelButton)
	ConnectTabable(cancelButton, comp)

	layout := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(comp, 0, 1, true).
			AddItem(tview.NewFlex().
				AddItem(nil, 0, 1, false).
				AddItem(selectButton, 8, 1, false).
				AddItem(nil, 1, 1, false).
				AddItem(cancelButton, 8, 1, false).
				AddItem(nil, 0, 1, false), 1, 1, false).
			AddItem(nil, 1, 1, false),
			0, 1, true).
		AddItem(help, 18, 1, false)
	showHelp := false
	layout.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			comp.exitAction()
			return nil
		case tcell.KeyF1:
			showHelp = !showHelp
			if showHelp {
				help.SetText(helpText)
			} else {
				help.SetText(helpTextShort)
			}
		}
		return event
	})

	pages.AddAndSwitchToPage("fileDialog", layout, true)
}

// NewPadded returns a new
func NewPadded(p tview.Primitive, left, top int) tview.Primitive {
	return tview.NewGrid().
		SetColumns(left, 0, left).
		SetRows(top, 0, top).
		AddItem(p, 1, 1, 1, 1, 0, 0, true)
}
