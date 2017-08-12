package main

import (
	"errors"
)

const (
	hlMagic = "HLB"
)

// Error messages
var (
	ErrNotValidHLB = errors.New("No valid HLB identifier")
	ErrUnsupported = errors.New("Unsupported HLB version")
)
