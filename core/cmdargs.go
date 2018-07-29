package core

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

type cliSetup interface {
	Setup() *flag.FlagSet
	Validate() []error
}

// ServerConfig is server configuration struct
type ServerConfig struct {
	DestValue          uint
	SamplingRate       uint
	FailstartThreshold uint
	Port               uint
	RaceMode           rune
	raceMode           string
	DbPath             string
	InputDevice        string
	OutputVisuals      string
	GrpcDebug          bool
	Fullscreen         bool
}

// VisualConfig is visual configuration struct
type VisualConfig struct {
	DistFactor       uint
	Port             uint
	HostName         string
	MovingUnit       uint
	VisName          string
	Fullscreen       bool
	ResolutionWidth  uint
	ResolutionHeight uint
	GrpcDebug        bool
}

var (
	defaultServerConfig = ServerConfig{
		DestValue:          400,
		SamplingRate:       5,
		FailstartThreshold: 5,
		Port:               9999,
		RaceMode:           't',
		InputDevice:        "SHM:5,6",
		OutputVisuals:      "localhost:9998",
		DbPath:             "sprints.pb",
		GrpcDebug:          false,
	}
	defaultVisConfig = VisualConfig{
		DistFactor:       25 * 5, // 25cm * 5
		Port:             9998,
		HostName:         "vision",
		VisName:          "bar",
		MovingUnit:       1,
		Fullscreen:       false,
		ResolutionWidth:  640,
		ResolutionHeight: 480,
		GrpcDebug:        false,
	}
)

// Setup maps command line options into ServerConfig struct
func (s *ServerConfig) Setup() *flag.FlagSet {
	cfg := flag.NewFlagSet("server", flag.ExitOnError)
	cfg.Usage = func() {
		fmt.Printf("\nserver configuration\n")
		cfg.PrintDefaults()
	}
	cfg.UintVar(&s.DestValue, "dest_value", defaultServerConfig.DestValue,
		"destination value to reach during a race")
	cfg.UintVar(&s.SamplingRate, "sampling_rate", defaultServerConfig.SamplingRate,
		"how many wheel turnovers causes animation to move")
	cfg.UintVar(&s.FailstartThreshold, "failstart_threshold",
		defaultServerConfig.FailstartThreshold,
		"how many wheel turnovers is acceptable during countdown")
	cfg.UintVar(&s.Port, "port", defaultServerConfig.Port,
		"TCP port for remote race control")
	cfg.StringVar(&s.raceMode, "race_mode", string(defaultServerConfig.RaceMode),
		"race mode: either t for time constrained race or d for distance constrained")
	cfg.StringVar(&s.InputDevice, "input_device", defaultServerConfig.InputDevice,
		"in the form: <type>:<device1_spec>,<device2_spec>,...")
	cfg.StringVar(&s.DbPath, "db_path", defaultServerConfig.DbPath,
		"file where to load & save tournament data")
	cfg.StringVar(&s.OutputVisuals, "visuals", defaultServerConfig.OutputVisuals,
		"comma seperated output visual addresses ie. ip:port,hostname:port etc.")
	cfg.BoolVar(&s.GrpcDebug, "grpc_debug", defaultServerConfig.GrpcDebug,
		"run GRPC server in debug mode")
	cfg.BoolVar(&s.Fullscreen, "fullscreen", defaultServerConfig.Fullscreen,
		"run visualization in fullscreen mode")

	return cfg
}

// Validate validates whether server configuration is correct
func (s *ServerConfig) Validate() (errs []error) {
	err := errors.New("race mode should be either d or t")

	if len(s.raceMode) != 1 {
		errs = append(errs, err)
		ErrorLogger.Println(err)
	}

	s.RaceMode = rune(s.raceMode[0])

	if s.RaceMode != 'd' && s.RaceMode != 't' {
		errs = append(errs, err)
		ErrorLogger.Println(err)
	}

	return
}

// Setup maps command line options into VisualConfig struct
func (c *VisualConfig) Setup() *flag.FlagSet {
	hostName, err := os.Hostname()
	if err != nil {
		hostName = defaultVisConfig.HostName
	}
	cfg := flag.NewFlagSet("visual", flag.ExitOnError)
	cfg.Usage = func() {
		fmt.Printf("\nvisual configuration\n")
		cfg.PrintDefaults()
	}
	cfg.UintVar(&c.DistFactor, "dist_factor", defaultVisConfig.DistFactor,
		"roller circum in cm * sampling rate (as in server)")
	cfg.UintVar(&c.Port, "port", defaultVisConfig.Port,
		"TCP port for GRPC communication w/ \"server\"")
	cfg.StringVar(&c.HostName, "name", hostName,
		"vision instance identification")
	cfg.StringVar(&c.VisName, "vis_name", defaultVisConfig.VisName,
		"visual name")
	cfg.UintVar(&c.MovingUnit, "moving_unit", defaultVisConfig.MovingUnit,
		"how many pixels to move in animation on one -sampling_rate")
	cfg.BoolVar(&c.Fullscreen, "fullscreen", defaultVisConfig.Fullscreen,
		"run visualization in fullscreen mode")
	cfg.UintVar(&c.ResolutionWidth, "width", defaultVisConfig.ResolutionWidth,
		"visualisation window/screen width in pixels")
	cfg.UintVar(&c.ResolutionHeight, "height", defaultVisConfig.ResolutionHeight,
		"visualization window/screen height in pixels")
	cfg.BoolVar(&c.GrpcDebug, "grpc_debug", defaultVisConfig.GrpcDebug,
		"run GRPC server in debug mode")

	return cfg
}

// FlagsetParse parses flags and prints usage if no options are given
func FlagsetParse(flagset *flag.FlagSet, args []string, argsValidation func() []error) {
	flagset.Parse(args)
	//flagset.Args()

	var validationErrs []error
	if argsValidation != nil {
		validationErrs = argsValidation()
	}
	if len(validationErrs) != 0 {
		flagset.Usage()
		os.Exit(1)
	}
}
