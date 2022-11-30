package app

import (
	"time"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/itohio/gortag/pkg/signal"
)

func (a *App) makeMenu() *widget.Toolbar {

	lastClick := time.Now()
	playAction := widget.NewToolbarAction(theme.MediaPlayIcon(), func() {
		if time.Since(lastClick) < time.Millisecond*500 {
			return
		}
		lastClick = time.Now()
		running, _ := a.engine.Running().Get()
		a.engine.Running().Set(!running)
	})

	a.engine.Running().AddListener(
		binding.NewDataListener(
			func() {
				running, _ := a.engine.Running().Get()
				if running {
					playAction.SetIcon(theme.MediaStopIcon())
				} else {
					playAction.SetIcon(theme.MediaPlayIcon())
				}
			},
		),
	)

	options := signal.Signals()

	addAction := widget.NewToolbarAction(theme.ContentAddIcon(), func() {
		sel := widget.NewSelect(options, nil)
		dlg := dialog.NewForm(
			"Add Signal",
			"Add",
			"Cancel",
			[]*widget.FormItem{
				widget.NewFormItem("Name:", sel),
			},
			func(b bool) {
				if !b {
					return
				}

				idx := sel.SelectedIndex()
				if idx < 0 {
					return
				}
				gen, err := signal.New(options[idx])
				if err != nil {
					return
				}

				a.engine.Generators().Append(gen)
			},
			a.window,
		)
		dlg.Show()
	})

	settingsAction := widget.NewToolbarAction(theme.SettingsIcon(), func() {
		a.ShowPreferences()
	})

	var actions []widget.ToolbarItem

	actions = append(
		actions,
		[]widget.ToolbarItem{
			widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
			}),
			widget.NewToolbarAction(theme.DownloadIcon(), func() {
			}),
		}...,
	)

	actions = append(
		actions,
		[]widget.ToolbarItem{
			addAction,
			widget.NewToolbarSpacer(),
			playAction,
			settingsAction,
			widget.NewToolbarSeparator(),
		}...,
	)

	return widget.NewToolbar(actions...)
}
