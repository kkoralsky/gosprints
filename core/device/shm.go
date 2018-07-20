package device

import (
	"encoding/binary"
	"fmt"
	log "github.com/kkoralsky/gosprints/core"
	"github.com/pkg/errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

const (
	shmResetSignal              = syscall.SIGABRT
	shmCloseSignal              = syscall.SIGTERM
	shmMapSize                  = 4
	shmPrefix                   = "/gosprints"
	shmDevice                   = "/dev/shm"
	defaultShmCounterExecutable = "raspio/goldio"
	shmExecutableEnv            = "GOSPRINTS_SHM_EXEC"
	shmSleepEnv                 = "GOSPRINTS_SHM_SLEEP"
	shmSimulateEnv              = "GOSPRINTS_SHM_SIMULATE"
	shmPullUpEnv                = "GOSPRINTS_SHM_PULLUP"
	shmSudoEnv                  = "GOSPRINTS_SHM_SUDO"
	shmWaitAfterResetEnv        = "GOSPRINTS_GOLDIO_WAIT"
)

// ShmReader represents SHM connection to read players distance; implements InputDevice
type ShmReader struct {
	playersNum     uint
	files          []*os.File
	falseStart     uint
	threshold      uint
	counterProcess *os.Process
	counterCmd     *exec.Cmd
}

// Init creates SHM "sockets" where input device data will be written
func (s *ShmReader) Init(ports []string, samplingRate uint, falseStart uint) error {
	var (
		counterExecPath        string
		counterArgs            []string
		found                  bool
		counterSleep           string
		counterSleepAfterReset string
	)
	counterExecPath, found = os.LookupEnv(shmExecutableEnv)
	if !found {
		counterExecPath = defaultShmCounterExecutable
	}
	if _, found = os.LookupEnv(shmSudoEnv); found {
		counterArgs[0] = counterExecPath
	}
	if _, found = os.LookupEnv(shmPullUpEnv); found {
		counterArgs = append(counterArgs, "-p")
	}
	if _, found = os.LookupEnv(shmSimulateEnv); found {
		counterArgs = append(counterArgs, "-s")
	}
	if counterSleep, found = os.LookupEnv(shmSleepEnv); found {
		counterArgs = append(counterArgs, fmt.Sprintf("-w %s", counterSleep))
	}

	counterArgs = append(counterArgs, fmt.Sprintf("-t %d", samplingRate), strings.Join(ports, ","))

	if counterExecPath == counterArgs[0] {
		s.counterCmd = exec.Command("sudo", counterArgs...)
	} else {
		s.counterCmd = exec.Command(counterExecPath, counterArgs...)
	}

	s.counterCmd.Env = os.Environ()
	if counterSleepAfterReset, found = os.LookupEnv(shmWaitAfterResetEnv); found {
		s.counterCmd.Env = append(s.counterCmd.Env, shmWaitAfterResetEnv+"="+counterSleepAfterReset)
	}

	s.threshold = samplingRate
	s.falseStart = falseStart
	s.playersNum = uint(len(ports))

	for i := 0; i < len(ports); i++ {
		filename := filepath.Join(shmDevice, fmt.Sprintf("%s%d", shmPrefix, i))
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDONLY, 0666)
		file.Truncate(0)
		if err != nil {
			return err
		}
		s.files = append(s.files, file)
	}

	return nil
}

// Start starts a counter program that will write distances into SHM files
func (s *ShmReader) Start() error {
	var (
		err                                      error
		counterStdOutReader, counterStdErrReader io.ReadCloser
	)
	counterStdOutReader, err = s.counterCmd.StdoutPipe()
	if err != nil {
		return err
	}
	counterStdErrReader, err = s.counterCmd.StderrPipe()
	if err != nil {
		return err
	}

	if err = s.counterCmd.Start(); err != nil {
		return err
	}
	counterReader := io.MultiReader(counterStdErrReader, counterStdOutReader)

	go func() {
		_, err := io.Copy(os.Stdout, counterReader)
		if err != nil {
			log.ErrorLogger.Printf("error on copying counter output: %v", err)
		}
	}()

	s.counterProcess = s.counterCmd.Process
	return nil
}

// GetDist reads current distance of a player from a SHM file
func (s *ShmReader) GetDist(playerID uint) (uint, error) {
	var (
		b      = make([]byte, shmMapSize)
		res    uint64
		n      int
		err    error
		errCtx = fmt.Sprintf("read from %s failed", s.files[playerID].Name())
	)
	s.files[playerID].Seek(0, 0)
	if _, err := s.files[playerID].Read(b); err != nil {
		return 0, errors.Wrap(err, errCtx)
	}

	res, n = binary.Uvarint(b)

	if n < 0 {
		return 0, errors.Wrap(err, errCtx)
	}
	return uint(res), nil
}

// GetPlayerCount returns number of players that were defined
func (s *ShmReader) GetPlayerCount() uint {
	return s.playersNum
}

// Clean triggers distance reset to 0 in the measuring program
func (s *ShmReader) Clean() error {
	return s.counterProcess.Signal(shmResetSignal)
}

// Check checks whether in any of the input SHM files distance of the allowed
// falseStart was exceeded
func (s *ShmReader) Check() (int, error) {
	var buf = make([]byte, 1024)
	for i, f := range s.files {
		_, err := f.Seek(0, 0)
		if err != nil {
			return -1, err
		}
		err = f.Sync()
		if err != nil {
			return -1, err
		}
		if _, err := f.Read(buf); err != nil {
			return -1, err
		}
		rotations, err := strconv.ParseUint(strings.TrimRight(string(buf), "\x00"), 10, 32)
		if err != nil {
			return -1, err
		}
		if rotations > uint64(s.falseStart) {
			return i, nil
		}
	}
	return -1, nil
}

// Close terminates counter companion program and closes all SHM files
func (s *ShmReader) Close() error {
	var errs []string

	if err := s.counterProcess.Signal(shmCloseSignal); err != nil {
		errs = append(errs, err.Error())
	}
	for _, f := range s.files {
		if err := f.Close(); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) != 0 {
		return errors.New("Closing failed:\n" + strings.Join(errs, "\n"))
	}
	return nil
}
