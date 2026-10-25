package auth

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+$`)

func validateCreateInput(input CreateInput) error {
	var errs []error
	var (
		hasLower   bool
		hasUpper   bool
		hasDigit   bool
		hasSpecial bool
	)

	if n := utf8.RuneCountInString(input.Username); n < 5 || n > 64 {
		errs = append(errs, errors.New("username must be between 5 and 64 characters"))
	}

	// Check all Unicode whitespace, including tabs and non-breaking spaces.
	if strings.ContainsFunc(input.Username, unicode.IsSpace) {
		errs = append(errs, errors.New("username must not contain whitespace"))
	}

	if strings.ContainsFunc(input.Username, unicode.IsControl) {
		errs = append(errs, errors.New("username must not contain control characters"))
	}

	if strings.ContainsFunc(input.Username, func(r rune) bool { return r > unicode.MaxASCII }) {
		errs = append(errs, errors.New("username must use only English letters, digits and ASCII symbols"))
	}

	if len(input.Email) < 5 || len(input.Email) > 254 {
		errs = append(errs, errors.New("email must be between 5 and 254 characters"))
	} else if !emailRegex.MatchString(input.Email) {
		errs = append(errs, errors.New("email is incorrect"))
	}

	if utf8.RuneCountInString(input.Password) < 8 {
		errs = append(errs, errors.New("password must be at least 8 characters"))
	}

	if len(input.Password) > 72 {
		errs = append(errs, errors.New("password must be at most 72 characters"))
	}

	if strings.ContainsFunc(input.Password, unicode.IsSpace) {
		errs = append(errs, errors.New("password must not contain whitespace"))
	}

	validASCII := true
	for _, r := range input.Password {
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

func normalizeCreateInput(input CreateInput) CreateInput {
	input.Email = strings.ToLower(input.Email)
	return input
}

// Login validates only the request format. New password complexity rules
// must not prevent login with a previously set password.
func validateLoginInput(input LoginInput) error {
	if input.Username == "" || input.Password == "" {
		return errors.New("username and password are required")
	}
	invalidCharacter := func(r rune) bool { return r < '!' || r > '~' }
	if strings.ContainsFunc(input.Username, invalidCharacter) || strings.ContainsFunc(input.Password, invalidCharacter) {
		return errors.New("username and password must use only English letters, digits and ASCII symbols, with no spaces")
	}
	if utf8.RuneCountInString(input.Username) > 64 || len(input.Password) > 72 {
		return errors.New("username or password is too long")
	}
	return nil
}
