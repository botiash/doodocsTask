package handler

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/mail"
	"net/smtp"
	"strings"

	"github.com/gin-gonic/gin"

	"doodocsProg/internal/config"
)

func SendFileToEmails(file *multipart.FileHeader, emails []string) error {
	// Получить MIME-тип загруженного файла
	mimeType := file.Header.Get("Content-Type")

	// Проверить, является ли MIME-тип допустимым
	if !isValidMimeTypeForEmail(mimeType) {
		return fmt.Errorf("Invalid file format: %s", mimeType)
	}
	conf := config.Get()
	auth := smtp.PlainAuth("", conf.Username, conf.Password, conf.Smtp_server)

	// Create an email message with MIME format
	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)

	// Set up the email header
	fmt.Fprintf(buf,
		"To: %s\r\n"+
			"Subject: File Attachment\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: multipart/mixed; boundary=%s\r\n\r\n",
		strings.Join(emails, ","),
		writer.Boundary(),
	)

	// Add the message text part
	textPart, err := writer.CreatePart(nil)
	if err != nil {
		return err
	}
	fmt.Fprint(textPart, "Please find the attached file.\r\n\r\n")

	// Add the file attachment part
	filePart, err := writer.CreateFormFile("file", file.Filename)
	if err != nil {
		return err
	}

	// Open the uploaded file
	fileToAttach, err := file.Open()
	if err != nil {
		return err
	}
	defer fileToAttach.Close()

	// Copy the file into the email
	_, err = io.Copy(filePart, fileToAttach)
	if err != nil {
		return err
	}

	// Close the writer
	writer.Close()

	// Send the email
	err = smtp.SendMail(conf.Smtp_server, auth, conf.Username, emails, buf.Bytes())
	if err != nil {
		fmt.Println("Failed to send email:", err)
		return err
	}

	return nil
}

func parseEmails(emails string) ([]string, error) {
	emailList := strings.Split(emails, ",")

	validEmails := make([]string, 0)
	for _, email := range emailList {
		email = strings.TrimSpace(email)
		if isValidEmail(email) {
			validEmails = append(validEmails, email)
		} else {
			return nil, fmt.Errorf("Invalid email address: '%s'", email)
		}
	}

	if len(validEmails) == 0 {
		return nil, fmt.Errorf("No valid email addresses provided")
	}

	return validEmails, nil
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func isValidMimeTypeForEmail(mimeType string) bool {
	validMimeTypes := map[string]struct{}{
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": {},
		"application/pdf": {},
	}

	_, valid := validMimeTypes[mimeType]
	return valid
}

func SendFileToEmailsHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to retrieve the file"})
		return
	}

	emails, emailErr := parseEmails(c.PostForm("emails"))
	if emailErr != nil {
		c.JSON(400, gin.H{"error": emailErr.Error()})
		return
	}

	err1 := SendFileToEmails(file, emails)
	if err1 != nil {
		c.JSON(500, gin.H{"error": "Failed to send the file"})
		return
	}

	c.JSON(200, gin.H{"message": "File sent successfully"})
}
