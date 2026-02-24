package members

import (
	"fmt"
	"net/mail"
	"strings"
	"unicode"
)

func normalizeAndValidateEmail(email string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(email))
	if normalized == "" {
		return "", fmt.Errorf("email is required")
	}

	parsed, err := mail.ParseAddress(normalized)
	if err != nil || parsed == nil || !strings.EqualFold(parsed.Address, normalized) {
		return "", fmt.Errorf("email format is invalid")
	}

	return normalized, nil
}

func normalizeAndValidatePhone(phone string) (string, error) {
	normalized := strings.TrimSpace(phone)
	if normalized == "" {
		return "", fmt.Errorf("phone is required")
	}

	builder := strings.Builder{}
	for _, char := range normalized {
		if unicode.IsDigit(char) {
			builder.WriteRune(char)
			continue
		}

		if char == ' ' || char == '-' {
			continue
		}

		return "", fmt.Errorf("phone format is invalid")
	}

	digits := builder.String()
	if len(digits) != 10 {
		return "", fmt.Errorf("phone must be exactly 10 digits")
	}

	return digits, nil
}
