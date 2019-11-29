package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gig-tech/windows-init/app"
	"github.com/gig-tech/windows-init/windows"
	"golang.org/x/sys/windows/svc"
)

func usage(errmsg string) {
	fmt.Fprintf(os.Stderr,
		"%s\n\n"+
			"usage: %s <command>\n"+
			"       where <command> is one of\n"+
			"       install, remove, debug, start, stop, pause or continue.\n",
		errmsg, os.Args[0])
	os.Exit(2)
}

func main() {
	isIntSess, err := svc.IsAnInteractiveSession()
	if err != nil {
		log.Fatalf("failed to determine if we are running in an interactive session: %v", err)
	}
	if !isIntSess {
		runService(app.SvcName, false)
		return
	}

	if len(os.Args) < 2 {
		usage("no command specified")
	}

	cmd := strings.ToLower(os.Args[1])
	switch cmd {
	case "debug":
		runService(app.SvcName, true)
		return
	case "install":
		err = installService(app.SvcName, app.SvcNameLong)
	case "remove":
		err = removeService(app.SvcName)
	case "start":
		err = windows.StartService(app.SvcName)
	case "stop":
		err = windows.ControlService(app.SvcName, svc.Stop, svc.Stopped)
	case "pause":
		err = windows.ControlService(app.SvcName, svc.Pause, svc.Paused)
	case "continue":
		err = windows.ControlService(app.SvcName, svc.Continue, svc.Running)
	default:
		usage(fmt.Sprintf("invalid command %s", cmd))
	}
	if err != nil {
		log.Fatalf("failed to %s %s: %v", cmd, app.SvcName, err)
	}
	return
}
