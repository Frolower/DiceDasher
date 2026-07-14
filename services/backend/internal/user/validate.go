package user

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)

func validateCreateInput(input CreateInput) error {
	var errs []error
	var (
		hasLower bool
		hasUpper bool
		hasDigit bool
	)

	if len(input.Username) < 3 {
		errs = append(errs, errors.New("username must be at least 3 characters"))
	}

	if len(input.Email) < 3 || len(input.Email) > 254 {
		errs = append(errs, errors.New("email must be between 3 and 254 characters"))
	} else if !emailRegex.MatchString(input.Email) {
		errs = append(errs, errors.New("email is incorrect"))
	}

	if len(input.Password) < 8 {
		errs = append(errs, errors.New("password must be at least 8 characters"))
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
	input.Username = strings.TrimSpace(input.Username)
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	return input
}
