package base

type ModulationType int

//go:generate go run golang.org/x/tools/cmd/stringer@latest -type ModulationType

const (
	None ModulationType = iota
	Amplitude
	Frequency
	Phase
	DutyCycle
	Ratio
	Trigger
)
