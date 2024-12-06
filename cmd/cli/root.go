package cli

import (
	"fmt"
	"io"
	"os"
)

const Version = "0.0.1"

type App struct {
	out io.Writer
}

func New() *App {
	return &App{
		out: os.Stdout,
	}
}

func (a *App) Execute(args []string) error {
	if len(args) == 0 {
		a.printHelp()
		return nil
	}

	switch args[0] {
	case "help", "-h", "--help":
		a.printHelp()
		return nil
	case "version", "-v", "--version":
		fmt.Fprintf(a.out, "%s", versionInfo())
		return nil
	case "export":
		return a.handleExport(args[1:])
	default:
		return fmt.Errorf("unknown command %q\nRun 'homebox-export help' for usage", args[0])
	}
}

func (a *App) printHelp() {
	help := `Usage: homebox-export <command> [options]

Commands:
  export        Download all items and their attachments
  help          Show this help message
  version       Show version information

Export Options:
  -server       Homebox server URL
  -user         Username for authentication
  -pass         Password for authentication
  -output       Output directory (default: ./downloads)
  -pagesize     Number of items per page (default: 100)

Environment Variables:
  HOMEBOX_SERVER   Server URL
  HOMEBOX_USER     Username
  HOMEBOX_PASS     Password
  HOMEBOX_OUTPUT   Output directory
  HOMEBOX_PAGESIZE Number of items per page

Examples:
  homebox-export export -server http://homebox.local -user admin -pass secret
  homebox-export export -output ./my-backup

For more information, visit: https://github.com/kusold/homebox-export`

	fmt.Fprintln(a.out, help)
}
