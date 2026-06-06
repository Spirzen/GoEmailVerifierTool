package main

import (
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"time"
)

// verifySMTP подключается к почтовому серверу и проверяет адрес через RCPT TO.
func verifySMTP(mxHost, email string, timeout time.Duration) (exists bool, detail string, err error) {
	host := strings.TrimSuffix(mxHost, ".")
	addr := net.JoinHostPort(host, "25")

	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return false, "", fmt.Errorf("TCP-соединение с %s: %w", addr, err)
	}
	defer conn.Close()

	if err := conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		return false, "", err
	}

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return false, "", fmt.Errorf("SMTP-клиент: %w", err)
	}
	defer client.Close()

	if err := client.Hello("email-verifier.local"); err != nil {
		return false, "", fmt.Errorf("EHLO/HELO: %w", err)
	}

	// Отправитель для проверки — типичный null-sender.
	if err := client.Mail(""); err != nil {
		return false, "", fmt.Errorf("MAIL FROM: %w", err)
	}

	if err := client.Rcpt(email); err != nil {
		return false, parseSMTPError(err), nil
	}

	_ = client.Quit()
	return true, "сервер принял RCPT TO", nil
}

func parseSMTPError(err error) string {
	msg := err.Error()
	// Код 550 — mailbox unavailable (чаще всего «не существует»).
	if strings.Contains(msg, "550") {
		return "адрес не существует (550)"
	}
	if strings.Contains(msg, "551") {
		return "адрес не локален для сервера (551)"
	}
	if strings.Contains(msg, "553") {
		return "некорректный адрес (553)"
	}
	return msg
}
