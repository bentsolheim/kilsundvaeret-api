package service

type Debug struct {
	SignalStrength   string
	TimeSpent        int32
	Iteration        int32
	Errors           int32
	MillisSinceStart int64
	Battery          BatteryLevel
}

type BatteryLevel struct {
	AnalogReading int32
	Voltage       float32
	Level         int32
}

type DebugResponse struct {
	Items Debug
}

type DataLoggerService struct {
}
