package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rm-Umar/email/internal/email"
)

var (
	version = "dev"
)

func main() {
	// Parse command line flags
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <command> [args...]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Commands:\n")
		fmt.Fprintf(os.Stderr, "  login    Configure email settings\n")
		fmt.Fprintf(os.Stderr, "  list     List email messages (interactive by default)\n")
		fmt.Fprintf(os.Stderr, "  send     Send an email message (interactive by default)\n")
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		fmt.Fprintf(os.Stderr, "  --non-interactive    Disable interactive mode\n")
		fmt.Fprintf(os.Stderr, "  --version           Show version information\n")
	}

	versionFlag := flag.Bool("version", false, "Show version information")
	nonInteractive := flag.Bool("non-interactive", false, "Disable interactive mode")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("go-email version %s\n", version)
		return
	}

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		return
	}

	cmd := args[0]
	cmdArgs := args[1:]

	var err error
	switch cmd {
	case "login":
		err = email.Login(cmdArgs)
	case "list":
		err = email.List(cmdArgs, !*nonInteractive)
	case "send":
		err = email.Send(cmdArgs, !*nonInteractive)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		flag.Usage()
		return
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}
}
