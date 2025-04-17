package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/rm-Umar/email/internal/email"
)

func main() {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	nonInteractive := fs.Bool("non-interactive", false, "Disable interactive mode")
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s list [--non-interactive] [read <seq>]\n", os.Args[0])
		fs.PrintDefaults()
	}

	if err := fs.Parse(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	args := fs.Args()

	if len(args) >= 2 && args[0] == "read" {
		seqStr := args[1]
		n, err := strconv.Atoi(seqStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid message number %q: must be an integer\n", seqStr)
			os.Exit(1)
		}

		// Initialize client and fetch body
		client, err := email.NewClient()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Connection error: %v\n", err)
			os.Exit(1)
		}
		defer client.Close()

		body, err := client.GetMessage(uint32(n))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fetch error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print("Body of text", body)
		return
	}

	// 2) Otherwise, just list messages (interactive or not)
	if err := email.List(args, !*nonInteractive); err != nil {
		fmt.Fprintf(os.Stderr, "Error listing messages: %v\n", err)
		os.Exit(1)
	}
}
