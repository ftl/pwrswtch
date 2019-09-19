package core

type Configuration struct {
	Levels      []PowerLevel
	TuningValue string
	TRXHost     string
}

type PowerLevel struct {
	Watts int
	Value string
}
