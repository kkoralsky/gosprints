package device

import (
	// "encoding/binary"
	"fmt"
	"github.com/hidez8891/shm"
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
	shmMapSize                  = 20
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
	playersCount   uint
	files          []*os.File
	segments       []*shm.Memory
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
		counterArgs = append(counterArgs, "-E", counterExecPath)
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

	if counterArgs[0] == "-E" {
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
	s.playersCount = uint(len(ports))

	for i := 0; i < len(ports); i++ {
		filename := filepath.Join(shmDevice, fmt.Sprintf("%s%d", shmPrefix, i))
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDONLY, 0666)
		file.Truncate(0)
		if err != nil {
			return err
		}
		s.files = append(s.files, file)
	}
	for i := 0; i < len(ports); i++ {
		mem, err := shm.Open(fmt.Sprintf("%s%d", shmPrefix, i), shmMapSize)
		if err != nil {
			return err
		}
		s.segments = append(s.segments, mem)
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

func (s *ShmReader) readSegment(i uint) (uint, error) {
	var (
		b   = make([]byte, shmMapSize)
		res uint64
		n   int
		err error
	)
	if int(i) > len(s.segments) {
		return 0, errors.New("segment not initialized")
	}

	_, err = s.segments[i].Seek(0, 0)
	if err != nil {
		return 0, err
	}
	_, err = s.segments[i].Read(b)
	if err != nil {
		return 0, err
	}
	for n = 0; b[n] != '\x00' && n < shmMapSize; n++ {
	}

	if n == shmMapSize {
		return 0, fmt.Errorf("map size: %d too small", shmMapSize)
	} else if n == 0 {
		return 0, errors.New("nothing has been read")
	}

	res, err = strconv.ParseUint(string(b[:n]), 10, 64)
	if err != nil {
		return 0, err
	}

	return uint(res), nil
}

// GetDist reads current distance of a player from a SHM file
func (s *ShmReader) GetDist(playerID uint) (uint, error) {
	var errCtx = fmt.Sprintf("reading for %d player failed", playerID)

	res, err := s.readSegment(playerID)
	// log.DebugLogger.Printf("player %d, distance: %d", playerID, res)
	return res, errors.Wrap(err, errCtx)
}

// GetPlayerCount returns number of players that were defined
func (s *ShmReader) GetPlayerCount() uint {
	return s.playersCount
}

// Clean triggers distance reset to 0 in the measuring program
func (s *ShmReader) Clean() error {
	return s.counterProcess.Signal(shmResetSignal)
}

// Check checks whether in any of the input SHM files distance of the allowed
// falseStart was exceeded
func (s *ShmReader) Check() (int, error) {
	for i := 0; i < int(s.playersCount); i++ {
		errCtx := fmt.Sprintf("check failed for %d", i)

		res, err := s.readSegment(uint(i))
		if err != nil {
			return -1, errors.Wrap(err, errCtx)
		}

		log.DebugLogger.Printf("#%d distance: %d", i, res)
		if res > s.falseStart {
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
