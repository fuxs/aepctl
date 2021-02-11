/*
Package configure contains the configuration command

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
package configure

import (
	"github.com/fuxs/aepctl/api"
	"github.com/fuxs/aepctl/ui"
	"github.com/fuxs/aepctl/util"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

var (
	// TODO long description
	configureLong = util.LongDesc(`
	Open configuration.

	This command supports the initial configuration. It opens a terminal based
	user interface with mouse and keyboard. Just copy the required credentials
	and paste the values into the right fields.

	The following keys can be used:
	
	* TAB go to next item. Alternatives are
	  * RETURN
	  * DOWN
	  * RIGHT
	* SHIFT+TAB go to the previous item
	  * UP
	  * LEFT
	* F1 show and hide help
	* F2 save current settings
	* F5 test the connection. It is recommended to test your settings, otherwise
	you might encounter problems with your settings
	* ESC exit the dialog

	aepctl will create a
 	`)
	configureExample = util.Example(`
	# Start initial configuration
	aepctl configure
	`)
)

var helpTextShort = "[white]Use Mouse and Keyboard, press [green]F1[white] for help"
var helpText = util.Form(`
					[white]Mouse and Keyboard are supported
					[yellow]Navigate                  Paste  Test Connection  Save  Hide Help  Exit
					[green]TAB, SHIFT+TAB, UP, DOWN  CMD+V  F5               F2    F1         ESC
`)

func newApp(cfg *util.ConfigFile) {
	pages := ui.NewFocusPages()
	form := ui.NewFocusFlex()
	showHelp := false
	contextHelp := tview.NewTextView().SetDynamicColors(true)
	contextHelpText := ""
	message := tview.NewTextView().SetDynamicColors(true).SetText("[green]Loaded [white]" + cfg.Path)
	help := tview.NewTextView().SetDynamicColors(true).SetText(helpTextShort)
	help.SetBorderPadding(0, 0, 1, 1) //.SetBorder(false)

	clientID := ui.NewPasteField("CLIENT ID            ").SetText(cfg.ClientID())
	clientID.SetFocusedFunc(form, func() {
		contextHelpText = "[yellow]CLIENT ID: [white]paste value with [green]CMD+V[white], go to next with [green]RETURN"
		if showHelp {
			contextHelp.SetText(contextHelpText)
		}
	})
	clientSecret := ui.NewPasteField("CLIENT SECRET        ").SetText(cfg.ClientSecret()).After(clientID)
	clientSecret.SetFocusedFunc(form, func() {
		contextHelpText = "[yellow]CLIENT SECRET: [white]paste value with [green]CMD+V[white], go to next with [green]RETURN"
		if showHelp {
			contextHelp.SetText(contextHelpText)
		}
	})
	techAccount := ui.NewPasteField("TECHNICAL ACCOUNT ID ").SetText(cfg.TechAccount()).After(clientSecret)
	techAccount.SetFocusedFunc(form, func() {
		contextHelpText = "[yellow]TECHNICAL ACCOUNT ID: [white]paste value with [green]CMD+V[white], go to next with [green]RETURN"
		if showHelp {
			contextHelp.SetText(contextHelpText)
		}
	})
	organization := ui.NewPasteField("ORGANIZATION ID      ").SetText(cfg.Organization()).After(techAccount)
	organization.SetFocusedFunc(form, func() {
		contextHelpText = "[yellow]ORGANIZATION ID: [white]paste value with [green]CMD+V[white], go to next with [green]RETURN"
		if showHelp {
			contextHelp.SetText(contextHelpText)
		}
	})
	sandbox := ui.NewPasteField("SANDBOX              ").SetText(util.StringOr(cfg.Sandbox(), "prod")).After(organization)
	sandbox.SetFocusedFunc(form, func() {
		contextHelpText = "[yellow]SANDBOX: [white]paste value with [green]CMD+V[white], go to next with [green]RETURN"
		if showHelp {
			contextHelp.SetText(contextHelpText)
		}
	})
	key := ui.NewFileField("PRIVATE KEY FILE     ", pages).SetText(cfg.Key()).After(sandbox)
	key.SetFocusedFunc(form,
		func() {
			contextHelpText = "[yellow]PRIVATE KEY FILE: [white]paste value with [green]CMD+V[white], go to next with [green]RETURN"
			if showHelp {
				contextHelp.SetText(contextHelpText)
			}
		},
		func() {
			contextHelpText = "[yellow]PRIVATE KEY FILE: [white]open file selector with [green]RETURN[white], [green]LEFT[white] key to edit"
			if showHelp {
				contextHelp.SetText(contextHelpText)
			}
		})

	helpAction := func() {
		showHelp = !showHelp
		if showHelp {
			help.SetText(helpText)
			contextHelp.SetText(contextHelpText)
		} else {
			help.SetText(helpTextShort)
			contextHelp.SetText("")
		}
	}

	testAction := func() {
		auth := &api.AuthenticationConfig{
			ClientID:         clientID.GetText(),
			ClientSecret:     clientSecret.GetText(),
			TechnicalAccount: techAccount.GetText(),
			Organization:     organization.GetText(),
			Sandbox:          sandbox.GetText(),
			Key:              key.GetText(),
		}
		token, err := auth.GetToken()
		if err != nil {
			ui.ErrorDialog(pages, err)
			return
		}
		ui.InfoDialog(pages, "Success: Retrieved token starting with "+token.Token[:16]+"...")
	}
	title := "Editing: " + cfg.Path
	changed := false

	saveAction := func() {
		cfg.SetClientID(clientID.GetText())
		cfg.SetClientSecret(clientSecret.GetText())
		cfg.SetTechAccount(techAccount.GetText())
		cfg.SetOrganization(organization.GetText())
		cfg.SetSandbox(sandbox.GetText())
		cfg.SetKey(key.GetText())
		err := cfg.Save()
		if err != nil {
			ui.ErrorDialog(pages, err)
			return
		}
		form.SetTitle(title)
		changed = false
		message.SetText("[green]Saved [white]" + cfg.Path)
	}

	app := tview.NewApplication()

	exitAction := func() {
		if !changed {
			app.Stop()
		}
		pages.AddAndSwitchToPage("exitDialog",
			tview.NewModal().SetText("Unsaved changes").AddButtons([]string{"Return", "Exit Anyway"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					if buttonIndex == 1 {
						app.Stop()
					}
					pages.RemovePage("exitDialog")
				}), true)
	}

	testButton := ui.NewFormButton("Test Connection").After(key)
	testButton.SetFocusedFunc(form, func() {
		contextHelpText = "[yellow]Test Connection: [white]press [green]ENTER[white] to test current configuration"
		if showHelp {
			contextHelp.SetText(contextHelpText)
		}
	}).SetSelectedFunc(testAction)
	saveButton := ui.NewFormButton("Save").After(testButton)
	saveButton.SetFocusedFunc(form, func() {
		contextHelpText = "[yellow]Save: [white]press [green]ENTER[white] to save current configuration"
		if showHelp {
			contextHelp.SetText(contextHelpText)
		}
	}).SetSelectedFunc(saveAction)
	exitButton := ui.NewFormButton("Exit").After(saveButton)
	exitButton.SetFocusedFunc(form, func() {
		contextHelpText = "[yellow]Exit: [white]press [green]ENTER[white] to exit"
		if showHelp {
			contextHelp.SetText(contextHelpText)
		}
	}).SetSelectedFunc(exitAction)
	clientID.After(exitButton)

	form.SetDirection(tview.FlexRow).
		AddItem(clientID, 2, 1, true).
		AddItem(clientSecret, 2, 1, false).
		AddItem(techAccount, 2, 1, false).
		AddItem(organization, 2, 1, false).
		AddItem(sandbox, 2, 1, false).
		AddItem(key, 2, 1, false).
		AddItem(tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(testButton, 17, 1, false).
			AddItem(nil, 1, 1, false).
			AddItem(saveButton, 6, 1, false).
			AddItem(nil, 1, 1, false).
			AddItem(exitButton, 6, 1, false).
			AddItem(nil, 0, 1, false), 1, 1, false).
		AddItem(nil, 0, 1, false).
		AddItem(contextHelp, 2, 1, false).
		AddItem(message, 1, 1, false)
	form.SetBorder(true).SetBorderPadding(1, 1, 1, 1).SetTitle(title)
	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			exitAction()
		case tcell.KeyF1:
			helpAction()
		case tcell.KeyF2:
			saveAction()
		case tcell.KeyF5:
			testAction()
		}

		return event
	})

	cf := func() {
		if !changed {
			changed = true
			form.SetTitle(title + "*")
			message.SetText("[yellow]Unsaved changes for [white]" + cfg.Path)
		}
	}

	clientID.Input.SetChangedFunc(func(text string) { cf() })
	clientSecret.Input.SetChangedFunc(func(text string) { cf() })
	techAccount.Input.SetChangedFunc(func(text string) { cf() })
	organization.Input.SetChangedFunc(func(text string) { cf() })
	sandbox.Input.SetChangedFunc(func(text string) { cf() })
	key.Input.SetChangedFunc(func(text string) { cf() })

	all := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(form, 0, 1, true).
		AddItem(help, 4, 1, false)

	pages.AddAndSwitchToPage("form", all, true)

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

// NewConfigureCommand creates an initialized command object
func NewConfigureCommand(gcfg *util.GlobalConfig) *cobra.Command {

	cmd := &cobra.Command{
		Use:                   "configure",
		Short:                 "Open configuration",
		Long:                  configureLong,
		Example:               configureExample,
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := util.LoadConfigFile(gcfg.Config)
			if err != nil {
				return err
			}
			newApp(cfg)
			return nil
		},
	}
	return cmd
}
