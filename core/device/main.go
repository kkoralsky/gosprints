package device

import (
	"github.com/pkg/errors"
	"strings"
)

// InputDevice -
type InputDevice interface {
	Init(players []string, threshold uint, falseStart uint) error // initialization of a device
	Start() error                                                 // start the device
	GetDist(playerID uint) (uint, error)                          // returns actual distance ridden by the given player
	GetPlayerCount() uint                                         // returns player count initialized
	Clean() error                                                 // resets players distance to 0
	Check() (int, error)                                          // checks if any player has exceeded the falseStart distance
	Close() error                                                 // performs cleanups: closes all devices, files etc.
}

// SetupDevice parses device configuration string and returns proper InputDevice interface
// implementation already initiaited
func SetupDevice(deviceConf string, samplingRate uint, failstartThreshold uint) (device InputDevice, err error) {
	deviceConfTuple := strings.Split(deviceConf, ":")
	switch deviceConfTuple[0] {
	case string("SHM"):
		device = &ShmReader{}
	default:
		device = &ShmReader{}
	}

	err = device.Init(strings.Split(deviceConfTuple[1], ","), samplingRate, failstartThreshold)
	if err != nil {
		return nil, errors.Wrap(err, "device initialization")
	}

	return device, nil
}
