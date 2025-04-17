package email

import (
	"fmt"
	"strings"
	"time"

	"github.com/rivo/tview"
)

type EmailList struct {
	UID     uint32
	From    string
	Subject string
	Date    time.Time
	Unread  bool
}

func formatEmailList(email *EmailList) string {
	unread := "[ ]"
	if email.Unread {
		unread = "[~]"
	}
	return fmt.Sprintf("[white]%-6d | %-4s %-30s │ %-50s │ %-20s",
		email.UID,
		unread,
		truncateString(email.From, 30),
		truncateString(email.Subject, 50),
		email.Date.Format("2006-01-02 15:04"))
}

func truncateString(str string, length int) string {
	if len(str) <= length {
		return str
	}
	return str[:length-3] + "..."
}

func createListView(emails []*EmailList, onSelect func(*EmailList)) *tview.List {
	list := tview.NewList().
		ShowSecondaryText(false).
		SetHighlightFullLine(true)

	// Add header
	header := fmt.Sprintf("[white]%-6s | %-30s │ %-50s │ %-20s", "ID", "From", "Subject", "Date")
	list.AddItem(header, "", 0, nil)

	divider := fmt.Sprintf("[white]%s", strings.Repeat("─", 100))
	list.AddItem(divider, "", 0, nil)
	for _, email := range emails {
		msg := email
		list.AddItem(formatEmailList(msg), "", 0, func() {
			onSelect(msg)
		})
	}

	return list
}

func createEmailView(email *EmailList, body string) *tview.Flex {
	text := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetWrap(true)

	content := fmt.Sprintf("[white]From: [yellow]%s\n[white]Subject: [yellow]%s\n[white]Date: [yellow]%s\n\n[white]%s",
		email.From,
		email.Subject,
		email.Date.Format("2006-01-02 15:04:05"),
		body)

	text.SetText(content)
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(text, 0, 1, false).
		AddItem(tview.NewButton("Back").SetSelectedFunc(func() {
			// The back action will be handled by the caller
		}), 1, 0, true)

	return flex
}

func createSendView(onSend func(to, subject, body string)) *tview.Form {
	var form *tview.Form
	form = tview.NewForm().
		AddInputField("To", "", 50, nil, nil).
		AddInputField("Subject", "", 50, nil, nil).
		AddTextArea("Body", "", 50, 10, 0, nil).
		AddButton("Send", func() {
			to := form.GetFormItem(0).(*tview.InputField).GetText()
			subject := form.GetFormItem(1).(*tview.InputField).GetText()
			body := form.GetFormItem(2).(*tview.TextArea).GetText()

			if to == "" || subject == "" || body == "" {
				return
			}

			onSend(to, subject, body)
		})

	form.SetBorder(true).SetTitle("Compose Email")
	return form
}
