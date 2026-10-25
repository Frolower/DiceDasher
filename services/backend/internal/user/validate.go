package user

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/google/uuid"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+$`)

func validateUserUpdateRequest(input UserUpdateRecord) error {
	var errs []error

	if input.ID == uuid.Nil {
		errs = append(errs, errors.New("missing ID"))
	}

	if len(input.Username) != 0 {
		if n := utf8.RuneCountInString(input.Username); n < 5 || n > 64 {
			errs = append(errs, errors.New("username must be between 5 and 64 characters"))
		}

		if strings.ContainsFunc(input.Username, unicode.IsSpace) {
			errs = append(errs, errors.New("username must not contain whitespace"))
		}

		if strings.ContainsFunc(input.Username, unicode.IsControl) {
			errs = append(errs, errors.New("username must not contain control characters"))
		}
	}

	if strings.ContainsFunc(input.Username, func(r rune) bool { return r > unicode.MaxASCII }) {
		errs = append(errs, errors.New("username must use only English letters, digits and ASCII symbols"))
	}

	if len(input.Email) != 0 {
		if len(input.Email) < 5 || len(input.Email) > 254 {
			errs = append(errs, errors.New("email must be between 5 and 254 characters"))
		} else if !emailRegex.MatchString(input.Email) {
			errs = append(errs, errors.New("email is incorrect"))
		}
	}

	if input.Username == "" && input.Email == "" {
		errs = append(errs, errors.New("missing arguments to update"))
	}

	return errors.Join(errs...)
}

func validatePassword(input string) error {
	var errs []error

	var (
		hasLower   bool
		hasUpper   bool
		hasDigit   bool
		hasSpecial bool
	)

	if utf8.RuneCountInString(input) < 8 {
		errs = append(errs, errors.New("password must be at least 8 characters"))
	}

	if len(input) > 72 {
		errs = append(errs, errors.New("password must be at most 72 characters"))
	}

	if strings.ContainsFunc(input, unicode.IsSpace) {
		errs = append(errs, errors.New("password must not contain whitespace"))
	}

	validASCII := true
	for _, r := range input {
		switch {
		case r >= 'a' && r <= 'z':
			hasLower = true
		case r >= 'A' && r <= 'Z':
			hasUpper = true
		case r >= '0' && r <= '9':
			hasDigit = true
		case r >= '!' && r <= '~':
			hasSpecial = true
		default:
			validASCII = false
		}
	}

	if !validASCII {
		errs = append(errs, errors.New("password must contain only English letters, digits and ASCII special characters, with no spaces"))
	}

	if !hasLower || !hasUpper || !hasDigit || !hasSpecial {
		errs = append(errs, errors.New("password must contain lower case letters, upper case letters, special symbols and digits"))
	}

	return errors.Join(errs...)
}

func normalizeUserUpdateInput(input UserUpdateRecord) UserUpdateRecord {
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	input.Username = strings.TrimSpace(input.Username)
	return input
}
