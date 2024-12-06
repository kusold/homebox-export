package cli

import (
	"fmt"
	"runtime"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func versionInfo() string {
	return fmt.Sprintf("homebox-export %s (%s) - built %s\n%s/%s",
		version,
		commit,
		date,
		runtime.GOOS,
		runtime.GOARCH,
	)
}
