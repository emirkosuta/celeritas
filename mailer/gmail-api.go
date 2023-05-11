package mailer

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/hashicorp/go-multierror"
	"google.golang.org/api/gmail/v1"
)

const boundaryText = "boundaryText"

func (m *Mail) SendEmailUsingGmailApi(msg Message) error {
	formattedMessage, err := m.buildHTMLMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to build HTML message: %w", err)
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "From: %s\r\n", msg.From)
	fmt.Fprintf(&sb, "To: %s\r\n", msg.To)

	if len(msg.Cc) > 0 {
		fmt.Fprintf(&sb, "Cc: %s\r\n", strings.Join(msg.Cc, ","))
	}
	if len(msg.Bcc) > 0 {
		fmt.Fprintf(&sb, "Bcc: %s\r\n", strings.Join(msg.Bcc, ","))
	}

	fmt.Fprintf(&sb, "Subject: %s\r\n", msg.Subject)
	fmt.Fprintf(&sb, "Content-Type: multipart/mixed; boundary=\"%s\"\r\n\r\n", boundaryText)
	fmt.Fprintf(&sb, "--%s\r\n", boundaryText)
	fmt.Fprintf(&sb, "Content-Type: text/html; charset=\"UTF-8\"\r\n")
	fmt.Fprintf(&sb, "%s\r\n", formattedMessage)

	var result error

	for _, x := range msg.Attachments {
		fileBytes, err := os.ReadFile(x)
		if err != nil {
			result = multierror.Append(result, fmt.Errorf("failed to read attachment file %s: %w", x, err))
			continue
		}

		fileMIMEType := mimetype.Detect(fileBytes)
		fileName := filepath.Base(x)

		fmt.Fprintf(&sb, "--%s\r\n", boundaryText)
		fmt.Fprintf(&sb, "Content-Type: %s\r\n", fileMIMEType)
		fmt.Fprintf(&sb, "Content-Transfer-Encoding: Base64\r\n")
		fmt.Fprintf(&sb, "Content-Disposition: attachment; filename=\"%s\"\r\n\r\n", fileName)

		encoder := base64.NewEncoder(base64.StdEncoding, &sb)
		_, err = encoder.Write(fileBytes)
		if err != nil {
			result = multierror.Append(result, fmt.Errorf("failed to write encoded file data: %w", err))
			continue
		}

		encoder.Close()
		sb.WriteString("\r\n")
	}

	fmt.Fprintf(&sb, "--%s--", boundaryText)

	if result != nil {
		return result
	}

	gMessage := &gmail.Message{Raw: base64.URLEncoding.EncodeToString([]byte(sb.String()))}
	_, err = GmailSrv.Users.Messages.Send("me", gMessage).Do()
	if err != nil {
		return fmt.Errorf("failed to send mail: %w", err)
	}

	return nil
}
