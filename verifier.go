package main

import (
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// Result — результат проверки одного адреса.
type Result struct {
	Email   string
	Valid   bool   // синтаксис корректен
	Exists  *bool  // nil — не удалось проверить через SMTP
	Domain  string
	MXHost  string
	Detail  string
	Error   error
}

// Verifier проверяет email через DNS (MX) и SMTP (RCPT TO).
type Verifier struct {
	Timeout time.Duration
}

func NewVerifier(timeout time.Duration) *Verifier {
	return &Verifier{Timeout: timeout}
}

func (v *Verifier) Verify(email string) Result {
	email = strings.TrimSpace(strings.ToLower(email))
	result := Result{Email: email}

	if !emailRegex.MatchString(email) {
		result.Detail = "некорректный формат email"
		return result
	}
	result.Valid = true

	parts := strings.SplitN(email, "@", 2)
	result.Domain = parts[1]

	mxRecords, err := lookupMX(result.Domain)
	if err != nil {
		result.Error = err
		result.Detail = err.Error()
		return result
	}

	var lastErr error
	for _, mx := range mxRecords {
		host := strings.TrimSuffix(mx.Host, ".")
		exists, detail, err := verifySMTP(host, email, v.Timeout)
		if err != nil {
			lastErr = err
			continue
		}

		result.MXHost = host
		result.Exists = &exists
		result.Detail = detail
		return result
	}

	result.Error = lastErr
	if lastErr != nil {
		result.Detail = fmt.Sprintf("не удалось связаться ни с одним MX: %v", lastErr)
	} else {
		result.Detail = "нет доступных MX-серверов"
	}
	return result
}

// PrintMX выводит MX-записи домена (для демонстрации работы DNS).
func PrintMX(domain string) error {
	records, err := lookupMX(domain)
	if err != nil {
		return err
	}
	fmt.Printf("MX-записи для %s:\n", domain)
	for _, mx := range records {
		fmt.Printf("  приоритет %d → %s\n", mx.Pref, strings.TrimSuffix(mx.Host, "."))
	}
	return nil
}

// ParseDomain извлекает домен из email или возвращает строку как есть.
func ParseDomain(input string) (string, error) {
	input = strings.TrimSpace(input)
	if strings.Contains(input, "@") {
		parts := strings.SplitN(input, "@", 2)
		if len(parts[1]) == 0 {
			return "", fmt.Errorf("домен не указан")
		}
		return parts[1], nil
	}
	if net.ParseIP(input) != nil {
		return "", fmt.Errorf("ожидается домен или email, не IP")
	}
	return input, nil
}
