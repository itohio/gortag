package app

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type device struct {
	q          string
	id         string
	name       string
	channel    []string
	sampleRate []string
	format     []string
}

func (d device) String() string {
	if d.id == "wav" {
		return fmt.Sprintf("wav?name=%s&ch=%s&sr=%s&fmt=%s", url.QueryEscape(d.name), d.channel[0], d.sampleRate[0], d.format[0])
	} else {
		return fmt.Sprintf("writer?id=%s&name=%s&ch=%s&sr=%s&fmt=%s", d.id, url.QueryEscape(d.name), d.channel[0], d.sampleRate[0], d.format[0])
	}
}

func (a *App) initPreferences() {
	p := a.app.Preferences()

	_ = p
}

func deduplicate(arr []string) (res []string) {
	m := make(map[string]struct{}, len(arr))
	for _, s := range arr {
		m[s] = struct{}{}
	}
	res = make([]string, 0, len(m))
	for k := range m {
		res = append(res, k)
	}
	return
}

func deviceNames(dev []device) []string {
	arr := make([]string, len(dev))
	for i, d := range dev {
		arr[i] = d.name
	}
	return arr
}

func (a *App) devices() (devices []device, err error) {
	devs := a.engine.Devices()
	for i, dev := range devs {
		if dev == "wav" {
			devices = append(devices, device{
				q:          dev,
				id:         "wav",
				name:       "wav",
				channel:    []string{"1", "2", "3", "4", "5", "6"},
				sampleRate: []string{"8000", "22000", "44100", "48000", "96000"},
				format:     []string{"S16", "S24", "S32", "F32"},
			})
			continue
		}

		v, err := url.ParseQuery(dev)
		if err != nil {
			return nil, err
		}
		devices = append(devices, device{
			q:          dev,
			id:         v.Get("id"),
			name:       v.Get("name"),
			channel:    deduplicate(v["ch"]),
			sampleRate: deduplicate(v["sr"]),
			format:     deduplicate(v["fmt"]),
		})

		_ = i
		// dflt := " "
		// if v.Has("default") {
		// 	dflt = "*"
		// }
		// if v.Get("name") != "" {
		// 	fmt.Printf(" %d%s: %s (%v ch, %vsps) id=%s\n", i, dflt, v.Get("name"), v["ch"], v["sr"], v.Get("id"))
		// } else {
		// 	fmt.Printf(" %d%s: %s\n", i, dflt, dev)
		// }
	}
	return
}

func (a *App) device() (d device, err error) {
	q, err := a.engine.Device().Get()
	if err != nil {
		return
	}
	fmt.Println(q)
	p, err := url.ParseQuery(q)
	if err != nil {
		return
	}

	name := p.Get("name")
	ch := p.Get("ch")
	sr := p.Get("sr")
	fmt := p.Get("fmt")
	id := p.Get("id")

	if strings.HasPrefix("wav", q) {
		id = "wav"
	}

	return device{
		q:          q,
		id:         id,
		channel:    []string{ch},
		name:       name,
		sampleRate: []string{sr},
		format:     []string{fmt},
	}, nil
}

func StrIndex(arr []device, name string) int {
	for i, s := range arr {
		if s.name == name {
			return i
		}
	}
	return -1
}

func (a *App) ShowPreferences() {
	devs, err := a.devices()
	if err != nil {
		dialog.NewError(err, a.window).Show()
		return
	}
	currentDev, err := a.device()
	if err != nil {
		dialog.NewError(err, a.window).Show()
		return
	}
	originalName := currentDev.name
	originalCh := currentDev.channel[0]
	originalSr := currentDev.sampleRate[0]
	originalFormat := currentDev.format[0]

	eName := widget.NewEntry()
	eName.Validator = func(s string) error {
		if s == "" {
			return errors.New("must not be empty")
		}
		if currentDev.id == "wav" && !strings.HasSuffix(s, ".wav") {
			return errors.New("must be wav file")
		}
		return nil
	}
	bBrowse := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {
		dlg := dialog.NewFileSave(func(uc fyne.URIWriteCloser, err error) {
			currentDev.name = uc.URI().String()
			eName.SetText(currentDev.name)
		}, a.window)
		dlg.SetFileName(currentDev.name)
		dlg.SetFilter(storage.NewExtensionFileFilter([]string{".wav"}))
		dlg.Show()
	})
	cName := container.NewBorder(nil, nil, nil, bBrowse, eName)

	sRate := widget.NewSelect(nil, func(s string) { currentDev.sampleRate[0] = s })
	sFormat := widget.NewSelect(nil, func(s string) { currentDev.format[0] = s })
	sChannels := widget.NewSelect(nil, func(s string) { currentDev.channel[0] = s })

	names := deviceNames(devs)
	sDevice := widget.NewSelect(
		names,
		func(s string) {
			idx := StrIndex(devs, s)
			if devs[idx].id == "wav" {
				eName.SetText("")
				eName.Enable()
				currentDev.id = "wav"
				bBrowse.Show()
				cName.Refresh()
				sChannels.Options = devs[idx].channel
				sRate.Options = devs[idx].sampleRate
				sFormat.Options = devs[idx].format
				sChannels.SetSelectedIndex(0)
				sRate.SetSelectedIndex(2)
				sFormat.SetSelectedIndex(0)
				return
			} else {
				eName.SetText(s)
				eName.Disable()
				bBrowse.Hide()
				cName.Refresh()
			}

			currentDev = devs[idx]
			sRate.Options = currentDev.sampleRate
			sChannels.Options = currentDev.channel
			sFormat.Options = currentDev.format
			if currentDev.name == originalName {
				sRate.SetSelected(originalSr)
				sChannels.SetSelected(originalCh)
				sFormat.SetSelected(originalFormat)
				return
			}
			sRate.SetSelectedIndex(0)
			sChannels.SetSelectedIndex(0)
			sFormat.SetSelectedIndex(0)
		},
	)
	if currentDev.id == "wav" {
		sDevice.SetSelected("wav")
	} else {
		sDevice.SetSelected(currentDev.name)
		sChannels.SetSelected(currentDev.channel[0])
		sRate.SetSelected(currentDev.sampleRate[0])
		sFormat.SetSelected(currentDev.format[0])
	}

	dlg := dialog.NewForm(
		"Settings",
		"OK", "Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("Output Device", sDevice),
			widget.NewFormItem("Name", cName),
			widget.NewFormItem("Sample Rate", sRate),
			widget.NewFormItem("Channels", sChannels),
			widget.NewFormItem("Sample Format", sFormat),
		},
		func(b bool) {
			if !b {
				return
			}

			if currentDev.id == "" {
				return
			}
			// fmt.Println(currentDev.String())
			a.engine.Device().Set(currentDev.String())
		},
		a.window,
	)

	dlg.Resize(a.window.Canvas().Size())
	dlg.Show()
}
