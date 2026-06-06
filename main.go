package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 15*time.Second, "таймаут DNS/SMTP операций")
	mxOnly := flag.Bool("mx", false, "только показать MX-записи домена")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		printUsage()
		os.Exit(1)
	}

	target := args[0]

	if *mxOnly {
		domain, err := ParseDomain(target)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ошибка: %v\n", err)
			os.Exit(1)
		}
		if err := PrintMX(domain); err != nil {
			fmt.Fprintf(os.Stderr, "ошибка: %v\n", err)
			os.Exit(1)
		}
		return
	}

	v := NewVerifier(*timeout)
	result := v.Verify(target)

	printResult(result)

	if result.Error != nil || !result.Valid || (result.Exists != nil && !*result.Exists) {
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `Email Verifier Tool — проверка существования email через DNS (MX) и SMTP

Использование:
  email-verifier <email>           проверить адрес
  email-verifier -mx <domain>      показать MX-записи домена
  email-verifier -timeout 30s <email>

Примеры:
  email-verifier user@gmail.com
  email-verifier -mx google.com
  email-verifier -mx user@gmail.com

`)
}

func printResult(r Result) {
	fmt.Printf("Email:  %s\n", r.Email)
	fmt.Printf("Формат: %s\n", yesNo(r.Valid))

	if r.Domain != "" {
		fmt.Printf("Дomain: %s\n", r.Domain)
	}
	if r.MXHost != "" {
		fmt.Printf("MX:     %s\n", r.MXHost)
	}

	switch {
	case r.Exists == nil:
		fmt.Printf("SMTP:   не проверен (%s)\n", r.Detail)
	case *r.Exists:
		fmt.Printf("SMTP:   вероятно существует (%s)\n", r.Detail)
	default:
		fmt.Printf("SMTP:   не существует (%s)\n", r.Detail)
	}

	if r.Error != nil {
		fmt.Printf("Ошибка: %v\n", r.Error)
	}
}

func yesNo(v bool) string {
	if v {
		return "корректный"
	}
	return "некорректный"
}
