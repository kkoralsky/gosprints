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
	RollerCircum       float64
	DestValue          uint
	MovingUnit         uint
	VisName            string
	SamplingRate       uint
	FailstartThreshold uint
	VisPort            uint
	CmdPort            uint
	RaceMode           rune
	raceMode           string
	InputDevice        string
}

// ClientConfig is client configuration struct
type ClientConfig struct {
	ServerAddr string
	MovingUnit uint
	VisName    string
}

var (
	defaultServerConfig = ServerConfig{
		RollerCircum:       .00025,
		DestValue:          400,
		MovingUnit:         1,
		SamplingRate:       5,
		FailstartThreshold: 5,
		VisPort:            9998,
		CmdPort:            9999,
		RaceMode:           't',
		InputDevice:        "SHM:5,6",
	}
	defaultClientConfig = ClientConfig{
		ServerAddr: "goldie1:9998",
	}
)

// Setup maps command line options into ServerConfig struct
func (s *ServerConfig) Setup() *flag.FlagSet {
	cfg := flag.NewFlagSet("server", flag.ExitOnError)
	cfg.Usage = func() {
		fmt.Printf("\nserver configuration\n")
		cfg.PrintDefaults()
	}
	cfg.Float64Var(&s.RollerCircum, "roller_circum", defaultServerConfig.RollerCircum,
		"roller circum in km")
	cfg.UintVar(&s.DestValue, "dest_value", defaultServerConfig.DestValue,
		"destination value to reach during a race")
	cfg.UintVar(&s.MovingUnit, "moving_unit", defaultServerConfig.MovingUnit,
		"how many pixels to move in animation on one -sampling_rate")
	cfg.UintVar(&s.SamplingRate, "sampling_rate", defaultServerConfig.SamplingRate,
		"how many wheel turnovers causes animation to move")
	cfg.UintVar(&s.FailstartThreshold, "failstart_threshold",
		defaultServerConfig.FailstartThreshold,
		"how many wheel turnovers is acceptable during countdown")
	cfg.UintVar(&s.VisPort, "vis_port", defaultServerConfig.VisPort,
		"UDP port for remote visualizations")
	cfg.UintVar(&s.CmdPort, "cmd_port", defaultServerConfig.CmdPort,
		"TCP port for remote race control")
	cfg.StringVar(&s.raceMode, "race_mode", string(defaultServerConfig.RaceMode),
		"race mode: either t for time constrained race or d for distance constrained")
	cfg.StringVar(&s.InputDevice, "input_device", defaultServerConfig.InputDevice,
		"in the form: <type>:<device1_spec>,<device2_spec>,...")

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

// Setup maps command line options into ClientConfig struct
func (c *ClientConfig) Setup() *flag.FlagSet {
	cfg := flag.NewFlagSet("server", flag.ExitOnError)
	cfg.Usage = func() {
		fmt.Printf("\nclient configuration\n")
		cfg.PrintDefaults()
	}
	cfg.StringVar(&c.ServerAddr, "server_addr", defaultClientConfig.ServerAddr,
		"server address in host:port format")

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
	if len(validationErrs) != 0 || len(args) == 0 {
		flagset.Usage()
		os.Exit(1)
	}
}
