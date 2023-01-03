package raspberry

import (
	"fmt"
	_ "github.com/gostonefire/water-heater/common"
	"net"
	"strconv"
)

const (
	minVout float64 = 0.0
	maxVout float64 = 5.0
)

type Raspberry struct {
	address       string
	dacBias       float64
	biasSlope     float64
	biasIntercept float64
}

// NewRaspberry - Initializes and returns a pointer to a new Raspberry struct
// The IP address needs to be a valid IP address and the port in the range 0 to 65535.
// The dacBias is used to set a min output voltage from the DAC, not affecting the max output voltage, i.e. it
// re-slopes the default curve from 0.0-1.0 as minVout -> maxVout to 0.0-1.0 as minVout+dacBias -> maxVout.
// Before returning, it sets the DAC output to dacBias.
//
// Note, the DAC does remember last set output even after powering down, hence if critical to always start at dacBias
// either explicitly call the Raspberry.SetDACVoltage method or call the Raspberry.Close method which does it for you.
func NewRaspberry(ipAddress string, ipPort uint16, dacBias float64) (device *Raspberry, err error) {
	// Some basic error checking on given IP address and DAC bias
	if net.ParseIP(ipAddress) == nil {
		err = fmt.Errorf("invalid IP address: %s", ipAddress)
		return
	}
	if dacBias < 0.0 || dacBias > 4.0 {
		err = fmt.Errorf("dacBias range is between 0.0 and 4.0, but %f was given", dacBias)
		return
	}

	device = &Raspberry{
		address:       fmt.Sprintf("%s:%d", ipAddress, ipPort),
		dacBias:       dacBias,
		biasIntercept: dacBias / maxVout,
		biasSlope:     1.0 - dacBias/maxVout,
	}

	return
}

// Close - Sets the output voltage to dacBias, hence it doesn't really close anything but here for convenience
func (R *Raspberry) Close() (err error) {
	err = R.SetDACVoltage(0.0)
	return
}

// GetRTDTemperature - Reads temperature from RTD device and returns as a float value
func (R *Raspberry) GetRTDTemperature() (temp float64, err error) {
	temp, err = R.sendCmd("temp:0000000000")
	return
}

// GetRTDOhms - Reads resistance from RTD device and returns as a float value
func (R *Raspberry) GetRTDOhms() (ohms float64, err error) {
	ohms, err = R.sendCmd("ohms:0000000000")
	return
}

// SetDACVoltage - Sets the output voltage from the DAC
// - rvOut is the relative value (0.0-1.0) of the output voltage to set between DAC bias and 5 volts
func (R *Raspberry) SetDACVoltage(rvOut float64) (err error) {
	sloped := rvOut*R.biasSlope + R.biasIntercept
	cmd := fmt.Sprintf("vout:%010.3f", sloped)

	_, err = R.sendCmd(cmd)
	return
}

// GetMinDACVoltage - Returns min DAC output voltage, not taking bias into account
func (R *Raspberry) GetMinDACVoltage() (vout float64) {
	return minVout
}

// GetMaxDACVoltage - Returns max DAC output voltage
func (R *Raspberry) GetMaxDACVoltage() (vout float64) {
	return maxVout
}

// sendCmd - Sends a command with its value to the device
func (R *Raspberry) sendCmd(cmd string) (value float64, err error) {
	conn, err := net.Dial("tcp", R.address)
	if err != nil {
		fmt.Printf("error while connecting to server: %s\n", err)
		return
	}
	defer func(conn net.Conn) { _ = conn.Close() }(conn)

	_, err = conn.Write([]byte(cmd))
	if err != nil {
		fmt.Printf("error while writing to server: %s\n", err)
		return
	}

	buf := make([]byte, 1024)
	_, err = conn.Read(buf)
	if err != nil {
		fmt.Printf("error while reading from server: %s\n", err)
		return
	}

	// Convert to float64, first 5 bytes are command + colon, e.g. "temp:"
	value, err = strconv.ParseFloat(string(buf[5:]), 64)

	return
}
