package main

import (
	"errors"
)

// Error messages
var (
	ErrNotValidHLB = errors.New("No valid HLB identifier")
	ErrUnsupported = errors.New("Unsupported HLB version")
)

var (
	ErrBadOpCode = errors.New("Bad op code")
)
