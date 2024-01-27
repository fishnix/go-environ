package ds18b20

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gobot.io/x/gobot"
)

var (
	sysBusDevices = "/sys/bus/w1/devices"
)

// ThermalProbeDriver is the Gobot driver for the DS18B20 temperature sensor
type ThermalProbeDriver struct {
	DeviceID   string
	name       string
	connection gobot.Connection
}

func NewThermalProbeDriver(name string) *ThermalProbeDriver {
	return &ThermalProbeDriver{
		name: name,
	}
}

// Name returns the name of the device.
func (d *ThermalProbeDriver) Name() string {
	return d.name
}

// SetName sets the name of the device.
func (d *ThermalProbeDriver) SetName(s string) {
	d.name = s
}

// Start initializes the device.
func (d *ThermalProbeDriver) Start() error {
	exec.Command("modprobe", "w1-gpio")
	exec.Command("modprobe", "w1-therm")

	time.Sleep(time.Second)

	if err := d.findDeviceID(); err != nil {
		return errors.Wrap(err, "DS18B20 device not found")
	}
	return nil
}

// Halt stops the device, a noop.
func (d *ThermalProbeDriver) Halt() error {
	return nil
}

// Connection returns the connection of the device.
func (d *ThermalProbeDriver) Connection() gobot.Connection {
	return d.connection
}

// findDeviceID finds the device ID of the DS18B20 device.
func (d *ThermalProbeDriver) findDeviceID() error {
	files, err := os.ReadDir(sysBusDevices)
	if err != nil {
		return errors.Wrap(err, "Failed to read device ID")
	}

	rp := regexp.MustCompile(`28-[0-9a-zA-Z]*`)

	for _, f := range files {
		// DS18B20 devices IDs start with 28
		if rp.MatchString(f.Name()) {
			d.DeviceID = f.Name()
			return nil
		}
	}

	if d.DeviceID == "" {
		return errors.New("DS18B20 device not found. Failed to read device ID.")
	}

	return nil
}

// ReadTempC reads the temperature in Celsius from the device.
func (d *ThermalProbeDriver) ReadTempC() (float32, error) {
	if d.DeviceID == "" {
		return float32(0.0), fmt.Errorf("Device ID not found or device %s is not ready!", d.name)
	}

	reading, err := os.ReadFile(fmt.Sprintf("%s/%s/temperature", sysBusDevices, d.DeviceID))
	if err != nil {
		return float32(0.0), fmt.Errorf("Device ID not found or device %s is not ready! %s", d.name, err)
	}

	temp, err := strconv.Atoi(strings.TrimSpace(string(reading)))
	if err != nil {
		return float32(0.0), err
	}

	return float32(temp) / 1000, nil
}
