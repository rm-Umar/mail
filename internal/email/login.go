package email

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/rm-Umar/email/internal/config"
)

func Login(args []string) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to go-email setup!")
	fmt.Println("Please enter your email configuration:")

	cfg := config.DefaultConfig

	fmt.Print("IMAP Server [imap.gmail.com]: ")
	imapServer, _ := reader.ReadString('\n')
	imapServer = strings.TrimSpace(imapServer)
	if imapServer != "" {
		cfg.IMAPServer = imapServer
	}

	fmt.Print("IMAP Port [993]: ")
	imapPort, _ := reader.ReadString('\n')
	imapPort = strings.TrimSpace(imapPort)
	if imapPort != "" {
		fmt.Sscanf(imapPort, "%d", &cfg.IMAPPort)
	}

	fmt.Print("SMTP Server [smtp.gmail.com]: ")
	smtpServer, _ := reader.ReadString('\n')
	smtpServer = strings.TrimSpace(smtpServer)
	if smtpServer != "" {
		cfg.SMTPServer = smtpServer
	}

	fmt.Print("SMTP Port [587]: ")
	smtpPort, _ := reader.ReadString('\n')
	smtpPort = strings.TrimSpace(smtpPort)
	if smtpPort != "" {
		fmt.Sscanf(smtpPort, "%d", &cfg.SMTPPort)
	}

	fmt.Print("Email Address: ")
	username, _ := reader.ReadString('\n')
	cfg.Username = strings.TrimSpace(username)

	fmt.Print("Password/App Password: ")
	password, _ := reader.ReadString('\n')
	cfg.Password = strings.TrimSpace(password)

	fmt.Print("Use TLS? [Y/n]: ")
	useTLS, _ := reader.ReadString('\n')
	useTLS = strings.TrimSpace(strings.ToLower(useTLS))
	cfg.UseTLS = useTLS != "n"

	if err := config.SaveConfig(&cfg); err != nil {
		return fmt.Errorf("failed to save configuration: %v", err)
	}

	fmt.Println("\nConfiguration saved successfully!")
	fmt.Println("You can now use go-email commands.")
	return nil
}
