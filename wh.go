package waterheater

type Device interface {
	Close() (err error)
	GetRTDTemperature() (temp float64, err error)
	GetRTDOhms() (ohms float64, err error)
	SetDACVoltage(rvOut float64) (vout float64, err error)
	GetMinDACVoltage() (vout float64)
	GetMaxDACVoltage() (vout float64)
}

type WaterHeater struct {
	device Device
}

func NewWaterHeater(device Device) (waterHeater *WaterHeater) {
	waterHeater = &WaterHeater{device: device}

	return
}

func (W *WaterHeater) SetHeat(heat float64) (vout float64, err error) {
	vout, err = W.device.SetDACVoltage(heat)
	return
}

func (W *WaterHeater) GetTemp() (temp float64, err error) {
	temp, err = W.device.GetRTDTemperature()
	return
}
