package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rm-Umar/email/internal/email"
)

func main() {
	fs := flag.NewFlagSet("send", flag.ExitOnError)
	nonInteractive := fs.Bool("non-interactive", false, "Disable interactive mode")
	to := fs.String("to", "", "Recipient email address")
	subject := fs.String("subject", "", "Email subject")
	body := fs.String("body", "", "Email body")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if *nonInteractive && (*to == "" && *subject == "" && *body == "") {
		fs.Usage()
		return
	}

	if err := email.Send(flag.Args(), !*nonInteractive); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}
}
