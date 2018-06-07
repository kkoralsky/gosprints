package core

import (
	"log"
	"os"
)

var (
	// ErrorLogger -
	ErrorLogger = log.New(os.Stderr, "E ", log.Lshortfile)
	// InfoLogger -
	InfoLogger = log.New(os.Stdout, "I ", log.Lshortfile)
	// DebugLogger -
	DebugLogger = log.New(os.Stderr, "D ", log.Llongfile)
)
