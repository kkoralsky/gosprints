package device

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestInitOnePlayer(t *testing.T) {
	var (
		s     = &ShmReader{}
		ports = []string{"1"}
	)
	s.Init(ports, 5, 4)
	if playerCount := s.GetPlayerCount(); playerCount != 1 {
		t.Errorf("player count should be 1, not %d", playerCount)
	}
}

func TestInitFileCreation(t *testing.T) {
	var (
		s     = &ShmReader{}
		ports = []string{"1", "2"}
		now   = time.Now()
	)
	err := s.Init(ports, 5, 4)
	if err != nil {
		t.Fatal(err)
	}
	if len(s.files) != 2 {
		t.Errorf("there should be 2 files opened, not %d", len(s.files))
	}

	for _, f := range s.files {
		fInfo, _ := f.Stat()
		if fInfo.Size() != 0 {
			t.Errorf("size of the SHM should be 0 at first; its: %d", fInfo.Size())
		}
		year, month, day := fInfo.ModTime().Date()
		if !(year == now.Year() && month == now.Month() && day == now.Day()) {
			t.Errorf("modification time of %s file is not today, its: %v",
				f.Name(), fInfo.ModTime())
		}
	}
}

func TestCounterCommandProcessStartup(t *testing.T) {
	var (
		s                  = &ShmReader{}
		ports              = []string{"1", "3"}
		counterExitHandled = make(chan struct{})
		err                error
	)
	err = s.Init(ports, 5, 6)
	if err != nil {
		t.Fatal(err)
	}
	err = s.Start()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("counter process PID: %d", s.counterProcess.Pid)

	go func() {
		counterProcess, _ := os.FindProcess(s.counterProcess.Pid)
		counterProcessState, _ := counterProcess.Wait()
		counterProcessWaitStatus := counterProcessState.Sys().(syscall.WaitStatus)
		if counterProcessWaitStatus.Exited() {
			t.Errorf("program exited by itself (it shouldnt); exit code: %d",
				counterProcessWaitStatus.ExitStatus())
		}
		counterExitHandled <- struct{}{}
	}()
	err = s.Close()
	if err != nil {
		t.Fatal(err)
	}

	select {
	case <-counterExitHandled:
		t.Log("counter exited successfully")
	case <-time.After(time.Second):
		t.Error("Process didnt catch SIGTERM signal after 1 sec.")
	}
}

func TestCloseErrors(t *testing.T) {
	var (
		s     = &ShmReader{}
		ports = []string{"1", "2"}
		err   error
	)
	err = s.Init(ports, 5, 6)
	if err != nil {
		t.Fatal(err)
	}
	err = s.Start()
	if err != nil {
		t.Fatal(err)
	}
	s.counterProcess.Kill()
	s.counterProcess.Wait()
	s.files[0].Close()

	err = s.Close()
	if err == nil {
		t.Fatal("there was no single error")
	}
	errs := strings.Split(err.Error(), "\n")
	if len(errs) != 3 {
		t.Errorf("there should be 3 lines of error message (2 errors). %d found",
			len(errs))
	}
}

func TestSimulation(t *testing.T) {
	var (
		s     = &ShmReader{}
		ports = []string{"1", "2"}
		err   error
	)
	if _, found := os.LookupEnv(shmSimulateEnv); !found {
		t.Fatalf("%s is not set", shmSimulateEnv)
	}
	err = s.Init(ports, 4, 5)
	if err != nil {
		t.Fatal(err)
	}
	err = s.Start()
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(200 * time.Millisecond) // wait for program to fill some data

	t.Run("test if dists are read", func(tt *testing.T) {
		for i, _ := range ports {
			if dist, err := s.GetDist(uint(i)); err != nil {
				tt.Error(err)
			} else if dist == 0 {
				tt.Errorf("distance for player #%d is eq. 0", i)
			} else {
				tt.Logf("distance for player #%d is %d", i, dist)
			}
		}
	})

	t.Run("test default reset", func(tt *testing.T) {
		s.Clean()
		time.Sleep(2800 * time.Millisecond)
		for i, _ := range ports {
			if dist, err := s.GetDist(uint(i)); err != nil {
				tt.Error(err)
			} else if dist != 0 {
				tt.Errorf("distance for player #%d is not 0, its %d", i, dist)
			}
		}
	})

	t.Run("test custom reset time", func(tt *testing.T) {
		s.counterCmd.Env = append(s.counterCmd.Env, fmt.Sprintf("%s=%d", shmWaitAfterResetEnv, 4))
		s.Clean()
		time.Sleep(3800 * time.Millisecond)
		for i, _ := range ports {
			if dist, err := s.GetDist(uint(i)); err != nil {
				tt.Error(err)
			} else if dist != 0 {
				tt.Errorf("distance for player #%d is not 0, its %d", i, dist)
			}
		}
	})

	err = s.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCheck(t *testing.T) {
	var (
		s           = &ShmReader{}
		ports       = []string{"1", "2", "3"}
		shmWFiles   = make([]*os.File, 3)
		openShmFile = func(index uint) *os.File {
			f, err := os.OpenFile(s.files[index].Name(), os.O_WRONLY, 0666)
			if err != nil {
				t.Fatal(err)
			}
			return f
		}
		writeShmFile = func(index uint, distance string) {
			shmWFiles[index].Seek(0, 0)
			shmWFiles[index].WriteString(distance)
			shmWFiles[index].Sync()
		}
		truncateShmFiles = func() {
			for _, f := range shmWFiles {
				f.Truncate(0)
			}
		}
	)

	err := s.Init(ports, 4, 5)
	if err != nil {
		t.Fatal(err)
	}

	shmWFiles[0] = openShmFile(0)
	shmWFiles[1] = openShmFile(1)
	shmWFiles[2] = openShmFile(2)

	t.Run("check first player false start", func(tt *testing.T) {
		truncateShmFiles()
		writeShmFile(0, "6")
		writeShmFile(1, "0")
		writeShmFile(2, "5")

		blame, err := s.Check()
		if err != nil {
			tt.Errorf("error on checking out: %s", err.Error())
		}
		if blame != 0 {
			tt.Errorf("first (0) player expected to blame; %d have", blame)
		}
	})
	t.Run("check no player false start", func(tt *testing.T) {
		truncateShmFiles()
		writeShmFile(0, "5")
		writeShmFile(1, "0")
		writeShmFile(2, "0")

		blame, err := s.Check()
		if err != nil {
			tt.Errorf("error on checking out: %s", err.Error())
		}
		if blame != -1 {
			tt.Errorf("no false start expected (-1); %d have", blame)
		}

	})
	t.Run("check player 2 and 3 false starting", func(tt *testing.T) {
		truncateShmFiles()
		writeShmFile(0, "0")
		writeShmFile(1, "6")
		writeShmFile(2, "7")

		blame, err := s.Check()
		if err != nil {
			tt.Errorf("error on checking out: %s", err.Error())
		}
		if !(blame == 1 || blame == 2) {
			tt.Errorf("expected second (1) or third (2) player to blame; got %d",
				blame)
		}
	})
	t.Run("check writing garbage to SHM file", func(tt *testing.T) {
		truncateShmFiles()
		writeShmFile(2, "xfg1")

		blame, err := s.Check()
		if err == nil {
			tt.Error("expected error, got nil")
		}
		if blame != -1 {
			tt.Errorf("we should not blame anyone if error occurs; got %d instead", blame)
		}
	})
}
