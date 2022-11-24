package signal

import (
	"sort"

	"github.com/itohio/gortag/pkg/signal/file"
	"github.com/itohio/gortag/pkg/signal/message"
	"github.com/itohio/gortag/pkg/signal/noise"
	"github.com/itohio/gortag/pkg/signal/normalize"
	"github.com/itohio/gortag/pkg/signal/rectangle"
	"github.com/itohio/gortag/pkg/signal/sine"
	"github.com/itohio/gortag/pkg/signal/sinesweep"
	"github.com/itohio/gortag/pkg/signal/triangle"
)

// factory helper
var (
	signals = map[string]func() Generator{
		sine.Name: func() Generator {
			return sine.New()
		},
		rectangle.Name: func() Generator {
			return rectangle.New()
		},
		triangle.Name: func() Generator {
			return triangle.New()
		},
		noise.Name: func() Generator {
			return noise.New()
		},
		sinesweep.Name: func() Generator {
			return sinesweep.New()
		},
		file.Name: func() Generator {
			return file.New()
		},
		message.Name: func() Generator {
			return message.New()
		},
		normalize.Name: func() Generator {
			return normalize.New()
		},
	}
	signalNames []string
)

func init() {
	signalNames = make([]string, 0, len(signals))
	for name := range signals {
		signalNames = append(signalNames, name)
	}
	sort.Strings(signalNames)
}
