package ui

import "github.com/itohio/gortag/pkg/signal"

type Signal interface {
	UpdateGenerator(signal.Generator)
}
