package device

import (
	"os"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestInitOnePlayer(t *testing.T) {
	var (
		s       = &ShmReader{}
		players = []string{"one"}
	)
	s.Init(players, 5, 4)
	if playerCount := s.GetPlayerCount(); playerCount != 1 {
		t.Errorf("player count should be 1, not %d", playerCount)
	}
}

func TestInitFileCreation(t *testing.T) {
	var (
		s       = &ShmReader{}
		players = []string{"red", "blue"}
		now     = time.Now()
	)
	err := s.Init(players, 5, 4)
	if err != nil {
		t.Fatal(err)
	}
	if len(s.files) != 2 {
		t.Errorf("there should be 2 files opened, not %d", len(s.files))
	}

	for _, f := range s.files {
		f_info, _ := f.Stat()
		if f_info.Size() != 0 {
			t.Errorf("size of the SHM should be 0 at first; its: %d", f_info.Size())
		}
		year, month, day := f_info.ModTime().Date()
		if !(year == now.Year() && month == now.Month() && day == now.Day()) {
			t.Errorf("modification time of %s file is not today, its: %v",
				f.Name(), f_info.ModTime())
		}
	}
}

func TestCounterCommandProcessStartup(t *testing.T) {
	var (
		s                  = &ShmReader{}
		players            = []string{"1", "3"}
		counterExitHandled = make(chan struct{})
	)
	err := s.Init(players, 5, 6)
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
		s       = &ShmReader{}
		players = []string{"1", "2"}
	)
	err := s.Init(players, 5, 6)
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
