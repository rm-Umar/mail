# Go Email Client

A command-line email client written in Go that supports both interactive and non-interactive modes. It provides a simple interface for managing emails through IMAP and SMTP protocols.

## Features

- **Interactive TUI Mode**: Beautiful terminal user interface for managing emails
- **Command-line Mode**: Scriptable commands for automation
- **Email Operations**:
  - List emails with sorting and unread indicators
  - Read specific emails by ID
  - Send new emails
  - Configure email settings
- **Security**: Supports TLS for both IMAP and SMTP connections
- **Configuration**: Easy setup with interactive configuration

## Modes

### Interactive TUI Mode
The interactive mode provides a beautiful terminal user interface for managing your emails. It's perfect for:
- Daily email management
- Reading and composing emails
- Visual email organization
- Quick navigation between emails

### Non-interactive Mode
The non-interactive mode is designed for automation and scripting
- Automated email notifications
- Scripted email operations
- Integration with other tools
- Scheduled tasks
- CI/CD pipelines
- System monitoring alerts

Example automation use cases:
```bash
# Send automated system status emails
email send --non-interactive -to "admin@example.com" -subject "System Status" -body "$(systemctl status)"

# Monitor log files and send alerts
tail -f /var/log/app.log | while read line; do
  if echo "$line" | grep -q "ERROR"; then
    email send --non-interactive -to "devops@example.com" -subject "Error Alert" -body "$line"
  fi
done

# Backup notification
if tar -czf backup.tar.gz /data; then
  email send --non-interactive -to "backup@example.com" -subject "Backup Complete" -body "Backup completed successfully"
else
  email send --non-interactive -to "backup@example.com" -subject "Backup Failed" -body "Backup failed. Please check logs."
fi
```

## Installation

1. Clone the repository:
```bash
git clone https://github.com/rm-Umar/email.git
cd email
```

2. Build the project:
```bash
./build.sh
```

This will create the following binaries in the `bin/` directory:
- `email`: Main binary with all commands
- `list`: Standalone list command
- `send`: Standalone send command
- `login`: Standalone login command

## Configuration

Before using the client, you need to configure your email settings:

```bash
# Using the main binary
email login

# Or using the standalone command
login
```

The configuration will be saved in `~/.email/config.json` with the following settings:
- IMAP server and port
- SMTP server and port
- Email username and password
- TLS settings

## Usage

### Main Binary

The main binary provides all commands in one package:

```bash
# Show help
email --help

# List emails (interactive mode)
email list

# List emails (non-interactive mode)
email list --non-interactive

# Read a specific email
email list read <email-id>

# Send an email (interactive mode)
email send

# Send an email (non-interactive mode)
email send --non-interactive
```

### Standalone Commands

Each command is also available as a standalone binary:

```bash
# List emails
list [--non-interactive]

# Send an email
send [--non-interactive]

# Configure settings
login
```


## Interactive Mode Features

### Email List View
- Shows email ID, sender, subject, date, and unread status
- Emails are sorted by date (newest first)
- Unread emails are marked with `[~]`

### Email View
- Displays full email content
- Shows sender, subject, and date
- Includes a back button to return to the list

### Send View
- Form for entering recipient, subject, and body
- Validation for required fields
- Success/error messages

## Non-interactive Mode

### Listing Emails
```bash
email list --non-interactive
```
Output format:
```
ID. From: <sender>, Subject: <subject>, Date: <date>
```

### Sending Emails
```bash
email send --non-interactive -to <recipient> -subject <subject> -body <body>
```

### Reading Emails
```bash
email list read <email-id>
```

## Configuration File

The configuration is stored in `~/.go-email/config.json` with the following structure:
```json
{
  "imap_server": "imap.gmail.com",
  "imap_port": 993,
  "inbox_folder": "INBOX",
  "smtp_server": "smtp.gmail.com",
  "smtp_port": 465,
  "username": "your-email@gmail.com",
  "password": "your-password",
  "use_tls": true
}
```

## Dependencies

- [go-imap](https://github.com/emersion/go-imap): IMAP client library
- [tview](https://github.com/rivo/tview): Terminal UI library

