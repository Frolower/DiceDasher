package auth

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)

func validateCreateInput(input CreateInput) error {
	var errs []error
	var (
		hasLower bool
		hasUpper bool
		hasDigit bool
	)

	if n := utf8.RuneCountInString(input.Username); n < 5 || n > 64 {
		errs = append(errs, errors.New("username must be between 5 and 64 characters"))
	}

	// Проверяем все Unicode-пробелы, включая табуляцию и неразрывный пробел.
	if strings.ContainsFunc(input.Username, unicode.IsSpace) {
		errs = append(errs, errors.New("username must not contain whitespace"))
	}

	if strings.ContainsFunc(input.Username, unicode.IsControl) {
		errs = append(errs, errors.New("username must not contain control characters"))
	}

	if len(input.Email) < 3 || len(input.Email) > 254 {
		errs = append(errs, errors.New("email must be between 3 and 254 characters"))
	} else if !emailRegex.MatchString(input.Email) {
		errs = append(errs, errors.New("email is incorrect"))
	}

	if utf8.RuneCountInString(input.Password) < 8 {
		errs = append(errs, errors.New("password must be at least 8 characters"))
	}

	if len(input.Password) > 72 {
		errs = append(errs, errors.New("password must be at most 72 bytes"))
	}

	if strings.ContainsFunc(input.Password, unicode.IsSpace) {
		errs = append(errs, errors.New("password must not contain whitespace"))
	}

	for _, r := range input.Password {
		switch {
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsDigit(r):
			hasDigit = true
		}
		if hasLower && hasUpper && hasDigit {
			break
		}
	}

	if !hasLower || !hasUpper || !hasDigit {
		errs = append(errs, errors.New("password must contain lower case letters, upper case letters and digits"))
	}

	return errors.Join(errs...)
}

func normalizeCreateInput(input CreateInput) CreateInput {
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	return input
}
