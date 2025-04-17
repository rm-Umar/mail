package email

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/smtp"
	"strings"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/rm-Umar/email/internal/config"
)

type Client struct {
	imapClient *client.Client
	config     *config.Config
}

func NewClient() (*Client, error) {
	config, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	addr := fmt.Sprintf("%s:%d", config.IMAPServer, config.IMAPPort)
	var c *client.Client
	var err2 error

	if config.UseTLS {
		c, err2 = client.DialTLS(addr, &tls.Config{InsecureSkipVerify: true})
	} else {
		c, err2 = client.Dial(addr)
	}

	if err2 != nil {
		return nil, fmt.Errorf("failed to connect to server: %v", err2)
	}

	// Login
	if err := c.Login(config.Username, config.Password); err != nil {
		return nil, fmt.Errorf("failed to login: %v", err)
	}

	return &Client{
		imapClient: c,
		config:     config,
	}, nil
}

func (c *Client) Close() error {
	return c.imapClient.Logout()
}

func (c *Client) ListMessages() ([]*imap.Message, error) {
	mbox, err := c.imapClient.Select(c.config.InboxFolder, false)
	if err != nil {
		return nil, fmt.Errorf("failed to select mailbox: %v", err)
	}

	seqset := new(imap.SeqSet)
	seqset.AddRange(1, mbox.Messages)

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.imapClient.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope, imap.FetchFlags}, messages)
	}()

	var result []*imap.Message
	for msg := range messages {
		result = append(result, msg)
	}

	if err := <-done; err != nil {
		return nil, fmt.Errorf("failed to fetch messages: %v", err)
	}

	return result, nil
}

func (c *Client) GetMessage(seqNum uint32) (string, error) {
	_, err := c.imapClient.Select(c.config.InboxFolder, false)
	if err != nil {
		return "", fmt.Errorf("failed to select INBOX: %v", err)
	}

	seqset := new(imap.SeqSet)
	seqset.AddNum(seqNum)

	section := &imap.BodySectionName{
		BodyPartName: imap.BodyPartName{
			Specifier: imap.TextSpecifier,
		},
		Peek: false, // set to true if you don't want to mark it as Seen
	}

	// Create a channel for the message
	messages := make(chan *imap.Message, 1)
	done := make(chan error, 1)

	// Fetch the message with body and envelope
	go func() {
		done <- c.imapClient.Fetch(seqset, []imap.FetchItem{section.FetchItem()}, messages)
	}()

	// Wait for the message
	msg := <-messages
	if msg == nil {
		return "", fmt.Errorf("message not found")
	}

	if err := <-done; err != nil {
		return "", fmt.Errorf("failed to fetch message: %v", err)
	}

	// Get the body content
	var body strings.Builder
	for _, literal := range msg.Body {
		_, err := io.Copy(&body, literal)
		if err != nil {
			return "", fmt.Errorf("failed to read message body: %v", err)
		}
	}

	return body.String(), nil
}

func (c *Client) SendMessage(to, subject, body string) error {
	host := c.config.SMTPServer
	port := c.config.SMTPPort
	addr := fmt.Sprintf("%s:%d", host, port)

	// 1) Connect over TLS (implicit SSL)
	tlsConn, err := tls.Dial("tcp", addr, &tls.Config{
		ServerName: host,
	})
	if err != nil {
		return fmt.Errorf("TLS dial to %s failed: %w", addr, err)
	}
	defer tlsConn.Close()

	// 2) Create new SMTP client over that connection
	client, err := smtp.NewClient(tlsConn, host)
	if err != nil {
		return fmt.Errorf("creating SMTP client: %w", err)
	}
	defer client.Close()

	// 3) Authenticate
	auth := smtp.PlainAuth("", c.config.Username, c.config.Password, host)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("auth failed: %w", err)
	}

	// 4) Set the sender and recipient
	if err := client.Mail(c.config.Username); err != nil {
		return fmt.Errorf("setting MAIL FROM: %w", err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("setting RCPT TO: %w", err)
	}

	// 5) Send the data
	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("getting Data writer: %w", err)
	}
	defer wc.Close()

	// 6) Write headers + body
	msg := []byte(
		fmt.Sprintf("From: %s\r\n", c.config.Username) +
			fmt.Sprintf("To: %s\r\n", to) +
			fmt.Sprintf("Subject: %s\r\n", subject) +
			"\r\n" +
			body + "\r\n",
	)
	if _, err := wc.Write(msg); err != nil {
		return fmt.Errorf("writing message: %w", err)
	}

	// 7) Close the DATA command
	if err := wc.Close(); err != nil {
		return fmt.Errorf("closing DATA: %w", err)
	}

	// 8) Quit politely
	if err := client.Quit(); err != nil {
		return fmt.Errorf("QUIT failed: %w", err)
	}

	return nil
}
