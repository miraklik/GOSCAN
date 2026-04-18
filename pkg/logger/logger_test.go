package logger

import (
	"testing"
)

func TestInit(t *testing.T) {
	Init(false)
	Init(true)
}

func TestSetLevel(t *testing.T) {
	SetLevel(DEBUG)
	SetLevel(INFO)
	SetLevel(WARN)
	SetLevel(ERROR)
}