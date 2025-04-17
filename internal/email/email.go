package email

import (
	"flag"
	"fmt"
	"os"
	"sort"

	"slices"

	"github.com/emersion/go-imap"
	"github.com/rivo/tview"
)

func List(args []string, interactive bool) error {
	if !interactive {
		// Non-interactive mode
		fs := flag.NewFlagSet("list", flag.ExitOnError)
		fs.Usage = func() {
			fmt.Fprintf(os.Stderr, "Usage: list [options]\n\n")
			fmt.Fprintf(os.Stderr, "Options:\n")
			fs.PrintDefaults()
		}

		if err := fs.Parse(args); err != nil {
			return err
		}

		client, err := NewClient()
		if err != nil {
			return err
		}
		defer client.Close()

		messages, err := client.ListMessages()
		if err != nil {
			return err
		}

		sort.Slice(messages, func(i, j int) bool {
			return messages[i].Envelope.Date.After(messages[j].Envelope.Date)
		})

		for _, msg := range messages {
			envelope := msg.Envelope
			fmt.Printf("%d. From: %s, Subject: %s, Date: %s\n",
				msg.SeqNum,
				envelope.From[0].Address(),
				envelope.Subject,
				envelope.Date.Format("2006-01-02 15:04:05"))
		}

		return nil
	}

	// Interactive mode
	client, err := NewClient()
	if err != nil {
		return err
	}
	defer client.Close()

	messages, err := client.ListMessages()
	if err != nil {
		return err
	}

	sort.Slice(messages, func(i, j int) bool {
		return messages[i].Envelope.Date.After(messages[j].Envelope.Date)
	})

	emailList := make([]*EmailList, len(messages))
	for i, msg := range messages {
		envelope := msg.Envelope
		emailList[i] = &EmailList{
			UID:     msg.SeqNum,
			From:    envelope.From[0].Address(),
			Subject: envelope.Subject,
			Date:    envelope.Date,
			Unread:  !contains(msg.Flags, imap.SeenFlag),
		}
	}

	app := tview.NewApplication()
	var list *tview.List
	list = createListView(emailList, func(email *EmailList) {
		body, err := client.GetMessage(email.UID)
		if err != nil {
			modal := tview.NewModal().
				SetText(fmt.Sprintf("Error: %v", err)).
				AddButtons([]string{"Close"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					app.SetRoot(list, true)
				})
			app.SetRoot(modal, true)
			return
		}

		emailView := createEmailView(email, body)
		emailView.GetItem(1).(*tview.Button).SetSelectedFunc(func() {
			app.SetRoot(list, true)
		})

		app.SetRoot(emailView, true)
	})

	if err := app.SetRoot(list, true).Run(); err != nil {
		return fmt.Errorf("failed to run application: %v", err)
	}

	return nil
}

func Send(args []string, interactive bool) error {
	if !interactive {
		// Non-interactive mode
		fs := flag.NewFlagSet("send", flag.ExitOnError)
		to := fs.String("to", "", "Recipient email address")
		subject := fs.String("subject", "", "Email subject")
		body := fs.String("body", "", "Email body")

		fs.Usage = func() {
			fmt.Fprintf(os.Stderr, "Usage: go-email send -to <recipient> -subject <subject> -body <body>\n\n")
			fmt.Fprintf(os.Stderr, "Options:\n")
			fs.PrintDefaults()
		}

		if err := fs.Parse(args); err != nil {
			return err
		}

		if *to == "" || *subject == "" || *body == "" {
			fs.Usage()
			return fmt.Errorf("to, subject, and body are required")
		}

		client, err := NewClient()
		if err != nil {
			return err
		}
		defer client.Close()

		return client.SendMessage(*to, *subject, *body)
	}

	// Interactive mode
	client, err := NewClient()
	if err != nil {
		return err
	}
	defer client.Close()

	app := tview.NewApplication()
	var form *tview.Form
	form = createSendView(func(to, subject, body string) {
		if err := client.SendMessage(to, subject, body); err != nil {
			modal := tview.NewModal().
				SetText(fmt.Sprintf("Error sending email: %v", err)).
				AddButtons([]string{"Close"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					app.SetRoot(form, true)
				})
			app.SetRoot(modal, true)
			return
		}

		modal := tview.NewModal().
			SetText("Email sent successfully!").
			AddButtons([]string{"Close"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				app.Stop()
			})
		app.SetRoot(modal, true)
	})

	if err := app.SetRoot(form, true).Run(); err != nil {
		return fmt.Errorf("failed to run application: %v", err)
	}

	return nil
}

func contains(slice []string, str string) bool {
	return slices.Contains(slice, str)
}
