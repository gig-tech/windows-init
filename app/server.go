package app

import (
	"golang.org/x/sys/windows/svc/debug"
)

type server struct {
	winlog debug.Log
}
