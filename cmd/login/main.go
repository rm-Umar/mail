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
	fs := flag.NewFlagSet("login", flag.ExitOnError)
	versionFlag := fs.Bool("version", false, "Show version information")
	fs.Parse(os.Args[1:])

	if *versionFlag {
		fmt.Printf("login version %s\n", version)
		os.Exit(0)
	}

	if err := email.Login(fs.Args()); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
