package app

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gig-tech/windows-init/windows"
	"github.com/mackerelio/go-osstat/uptime"
	"golang.org/x/sys/windows/svc"
)

// SvcName is the name you will use for the NET START command
const SvcName = "GIG Init"

// SvcNameLong is the name that will appear in the Services control panel
const SvcNameLong = "GIG Initializer"

// The wrapper of your app
func gigInit(s server) {
	defer stopSelf(s)

	// Get uptime
	up, err := uptime.Get()
	if err != nil {
		s.winlog.Error(1, "Could not get windows uptime!")
		return
	}
	bootTime := time.Now().Add(-up)
	start := time.Now()

	// Wait until cloudbase-init stops running and cloudbase-init.log has a recent file modification date since boot time.
	cbiLog := "C:\\Program Files\\Cloudbase Solutions\\Cloudbase-init\\log\\cloudbase-init.log"
	for {
		time.Sleep(time.Second * 5)
		if time.Since(start) > 10*time.Minute {
			s.winlog.Error(1, "Giving up after waiting for 10 minutes for cloudbase-init")
			return
		}
		info, err := os.Stat(cbiLog)
		if os.IsNotExist(err) {
			s.winlog.Info(1, "Cloudbase init log file does not exist yet")
			continue
		}
		running, err := windows.IsServiceRunning("cloudbase-init")
		if err != nil {
			s.winlog.Error(1, fmt.Sprintf("Could not determine the status of the Cloudbase-init service: %s", err))
			return
		}
		if info.ModTime().After(bootTime) && !running {
			s.winlog.Info(1, "Can proceed with running init scripts now")
			break
		}
	}

	if _, err := os.Stat("C:\\gig\\init"); os.IsNotExist(err) {
		s.winlog.Info(1, "No init scripts found.")
	} else {
		if files, err := ioutil.ReadDir("C:\\gig\\init"); err == nil {
			for _, f := range files {
				if strings.HasSuffix(f.Name(), ".ps1") {
					runOnce(s, "C:\\gig\\init\\"+f.Name())
				}
			}
		}

	}
}

func runOnce(s server, script string) {
	ran := script + ".executed"

	if _, err := os.Stat(ran); os.IsNotExist(err) {
		runPowershellScript(s, script)
		if file, err := os.Create(ran); err == nil {
			file.Close()
		}
	}
}

func runPowershellScript(s server, script string) {
	s.winlog.Info(1, "Executing "+script)
	ps, _ := exec.LookPath("powershell.exe")
	args := []string{"-NoProfile", "-NonInteractive", script}
	cmd := exec.Command(ps, args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Run()
	s.winlog.Info(1, stdout.String())
	s.winlog.Info(1, stderr.String())
}

func stopSelf(s server) {
	windows.ControlService(SvcName, svc.Stop, svc.Stopped)
}
