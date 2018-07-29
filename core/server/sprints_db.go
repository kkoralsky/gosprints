package server

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/kkoralsky/gosprints/core"
	pb "github.com/kkoralsky/gosprints/proto"
	"os"
)

const (
	dbFileFlags = os.O_EXCL | os.O_RDWR | os.O_SYNC
	dbFileMode  = 0666
)

type SprintsDb struct {
	file        *os.File
	tournaments *pb.Tournaments
}

func SetupSprintsDb(fileName string) (s *SprintsDb, err error) {
	var (
		b          []byte
		n          int
		dbFileInfo os.FileInfo
	)
	s = &SprintsDb{}
	dbFileInfo, err = os.Stat(fileName)
	if os.IsNotExist(err) {
		s.file, err = os.OpenFile(fileName, dbFileFlags|os.O_CREATE, dbFileMode)
	} else if err != nil {
		return
	} else {
		b = make([]byte, dbFileInfo.Size())
		s.file, err = os.OpenFile(fileName, dbFileFlags, dbFileMode)
	}
	if err != nil {
		return
	}
	n, err = s.file.Read(b)
	if err != nil {
		return
	}

	s.tournaments = &pb.Tournaments{}
	if n == 0 {
		core.InfoLogger.Printf("sprints db %s created or left empty\n", fileName)
	} else {
		err = proto.Unmarshal(b, s.tournaments)
		if err != nil {
			return
		}
		core.InfoLogger.Printf("loaded %d previous tournaments\n", len(s.tournaments.Tournament))
	}
	return s, nil
}

func (s *SprintsDb) SaveTournament(tournament *pb.Tournament) error {
	var (
		i            = s.getTournamentIndex(tournament.Name)
		err          error
		b            []byte
		bytesWritten int
		// fInfo        os.FileInfo
	)
	if i == -1 {
		s.tournaments.Tournament = append(s.tournaments.Tournament, tournament)
	} else {
		s.tournaments.Tournament[i] = tournament
	}

	// core.DebugLogger.Printf("len: %d i: %d, %v", len(s.tournaments.Tournament), i, s.tournaments.Tournament)

	b, err = proto.Marshal(s.tournaments)
	if err != nil {
		return err
	}
	err = s.file.Truncate(0)
	if err != nil {
		return err
	}
	bytesWritten, err = s.file.WriteAt(b, 0)
	if err != nil {
		return err
	}
	if bytesWritten == 0 {
		return errors.New("Nothing written to database")
	}
	err = s.file.Sync()
	if err != nil {
		return err
	}

	return nil
}

func (s *SprintsDb) getTournamentIndex(name string) int {
	for i, tournament := range s.tournaments.Tournament {
		if tournament.Name == name {
			return i
		}
	}
	return -1
}

func (s *SprintsDb) GetTournament(name string) (*pb.Tournament, error) {
	var i = s.getTournamentIndex(name)
	if i >= 0 {
		return s.tournaments.Tournament[i], nil
	}
	return nil, errors.New("not found")
}

func (s *SprintsDb) GetLastTournament() (*pb.Tournament, error) {
	var tournamentsLen = len(s.tournaments.Tournament)
	if tournamentsLen > 0 {
		return s.tournaments.Tournament[tournamentsLen-1], nil
	}
	return nil, errors.New("no tournaments")
}

func (s *SprintsDb) Close() error {
	return s.file.Close()
}
