package water_heater

type DeviceInterface interface {
	Close() (err error)
	GetRTDTemperature() (temp float64, err error)
	GetRTDOhms()
	SetDACVoltage(rvOut float64) (err error)
	GetMinDACVoltage()
	GetMaxDACVoltage()
}
