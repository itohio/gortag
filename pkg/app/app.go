package app

import (
	"context"
	"fmt"
	"net/url"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"github.com/itohio/gortag/pkg/signal"
	"github.com/itohio/gortag/pkg/ui"
)

type App struct {
	sync.Mutex
	app    fyne.App
	window fyne.Window
	cancel context.CancelFunc

	engine signal.Engine
}

var (
	appFQDN = "itohio.acoustics.rtag"
)

func New(name string) *App {
	a := app.NewWithID(appFQDN)
	w := a.NewWindow(name)
	w.Resize(fyne.NewSize(350, 600))

	ctx, cancel := context.WithCancel(context.Background())

	ret := &App{
		app:    a,
		window: w,
		cancel: cancel,
		engine: signal.NewEngine(ctx),
	}

	ret.initPreferences()

	for i, w := range ret.engine.Devices() {
		v, err := url.ParseQuery(w)
		if err != nil {
			panic(err)
		}
		dflt := " "
		if v.Has("default") {
			dflt = "*"
		}
		if v.Get("name") != "" {
			fmt.Printf(" %d%s: %s (%s ch, %ssps) id=%s\n", i, dflt, v.Get("name"), v.Get("ch"), v.Get("sr"), v.Get("id"))
		} else {
			fmt.Printf(" %d%s: %s\n", i, dflt, w)
		}
	}

	return ret
}

func (a *App) Run() {
	defer a.cancel()
	a.makeContent()
	a.window.ShowAndRun()
}

func (a *App) makeContent() {
	a.window.SetContent(
		container.NewBorder(
			a.makeMenu(), nil, // Top, Bottom
			nil, nil, // Left, Right
			container.NewVScroll(ui.NewSignalList(a.engine)),
			// widget.NewListWithData(
			// 	a.engine.Generators(),
			// 	func() fyne.CanvasObject {
			// 		return ui.NewSignal()
			// 	},
			// 	func(di binding.DataItem, co fyne.CanvasObject) {
			// 		untyped, ok := di.(binding.Untyped)
			// 		if !ok {
			// 			fmt.Println("data not ok", di)
			// 			return
			// 		}
			// 		val, err := untyped.Get()
			// 		if err != nil {
			// 			fmt.Println("failed get", err)
			// 			return
			// 		}
			// 		gen, ok := val.(signal.Generator)
			// 		if !ok {
			// 			fmt.Println("data not gen", untyped)
			// 			return
			// 		}
			// 		co.(ui.Signal).UpdateGenerator(gen)
			// 	},
			// ),
		),
	)
}
