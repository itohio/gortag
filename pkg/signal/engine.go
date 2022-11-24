package signal

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"fyne.io/fyne/v2/data/binding"
	"github.com/itohio/gocoustics/pkg/audio"
	"github.com/itohio/gocoustics/pkg/audio/system"
	"github.com/itohio/gocoustics/pkg/dsp"
)

type Engine interface {
	Running() binding.Bool
	Generators() binding.UntypedList

	Devices() []string
	DefaultDevice() string
	Device() binding.String
	Destroy()

	ValidateDevice(string) error
}

type engine struct {
	sync.Mutex
	channels  atomic.Int32
	isRunning atomic.Bool

	ctx        context.Context
	cancel     context.CancelFunc
	sys        system.AudioSystem[float64]
	generators binding.UntypedList

	frames  int
	running binding.Bool
	device  binding.String
}

func NewEngine(ctx context.Context) *engine {
	ret := &engine{
		ctx:        ctx,
		sys:        audio.New[float64](),
		generators: binding.NewUntypedList(),
		frames:     2048,
		running:    binding.NewBool(),
		device:     binding.NewString(),
	}

	ret.device.AddListener(binding.NewDataListener(ret.updateDevice))
	ret.running.AddListener(binding.NewDataListener(ret.updateRunning))
	ret.device.Set(ret.DefaultDevice())

	return ret
}

func (e *engine) Channels() int {
	return int(e.channels.Load())
}

func (e *engine) Destroy() {
	if e.sys != nil {
		e.sys.Close()
		e.sys = nil
	}
}

func (e *engine) updateDevice() {
	if running, _ := e.running.Get(); running {
		e.stop()
	}
	device, err := e.device.Get()
	if err != nil {
		panic(err)
	}
	if device == "" {
		return
	}
	v, err := url.ParseQuery(device)
	if err != nil {
		panic(err)
	}

	channels, err := strconv.Atoi(v.Get("ch"))
	if err != nil {
		panic(err)
	}
	e.channels.Store(int32(channels))
}

func (e *engine) updateRunning() {
	e.Lock()
	defer e.Unlock()

	running, err := e.running.Get()
	if err != nil {
		panic(err)
	}
	if running {
		e.start()
	} else {
		e.stop()
	}
}

func (e *engine) Devices() []string {
	return e.sys.Writers()
}

func (e *engine) Device() binding.String {
	return e.device
}

func (e *engine) Generators() binding.UntypedList {
	return e.generators
}

func (e *engine) Running() binding.Bool {
	return e.running
}

func (e *engine) DefaultDevice() string {
	for _, w := range e.Devices() {
		v, err := url.ParseQuery(w)
		if err != nil {
			panic(err)
		}

		if v.Has("default") {
			return w
		}
	}
	return ""
}

func (e *engine) ValidateDevice(string) error {
	return nil // TODO
}

func (e *engine) start() {
	if e.isRunning.Load() {
		return
	}

	device, err := e.device.Get()
	if err != nil {
		panic(err)
	}

	writer := e.sys.Writer(device)
	err = writer.Open()
	if err != nil {
		fmt.Println(err.Error())
		e.running.Set(false)
		return
	}
	e.channels.Store(int32(writer.Channels()))

	ctx, cancel := context.WithCancel(e.ctx)
	e.cancel = cancel
	e.isRunning.Store(true)

	go e.player(ctx, writer)
}

func (e *engine) stop() {
	if !e.isRunning.Load() {
		return
	}

	if e.cancel == nil {
		return
	}

	e.cancel()
}

func (e *engine) onStop() {
	generators, err := e.generators.Get()
	if err != nil {
		return
	}
	for _, genI := range generators {
		gen, ok := genI.(Generator)
		if !ok {
			return
		}
		gen.Reset()
	}
	e.isRunning.Store(false)
	e.running.Set(false)
}

func (e *engine) player(ctx context.Context, dev system.WriterCloserDevice[float64]) {
	defer e.onStop()
	defer dev.Close()
	buf := dev.NewSampleBuffer(e.frames, -1)
	ticker := time.NewTicker(buf.Frame2Duration(e.frames / 2))
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
		health := dev.Health()
		if health[audio.AH_OutQueue] > health[audio.AH_MaxOutQueue]*2/3 {
			continue
		}

		generators, err := e.generators.Get()
		if err != nil {
			return
		}

		dsp.FillSilence(buf.Buffer(), -1, buf.Channels())

		for _, genI := range generators {
			gen, ok := genI.(Generator)
			if !ok {
				return
			}
			if a, _ := gen.Active().Get(); !a {
				continue
			}

			gen.FillBuffer(buf, nil)
		}

		dev.Write(buf.Buffer())
	}
}
