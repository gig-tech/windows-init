package main

import (
	"git.gig.tech/openvcloud/windows-init/app"
	"github.com/pkg/errors"
)

// This is assigned the full SHA1 hash from GIT
var sha1ver string

func svcLauncher() error {

	err := app.Run(elog, app.SvcName, sha1ver)
	if err != nil {
		return errors.Wrap(err, "app.run")
	}

	return nil
}
