package dto

import (
	"errors"
	"regexp"
)

const MaxMessageLength = 100

// for TR phone numbers
var phoneRegex = regexp.MustCompile(`^\+905\d{9}$`)

func (m *MessageRequest) Validate() error {
	if len(m.To) == 0 {
		return errors.New("phone number is required")
	}
	if !phoneRegex.MatchString(m.To) {
		return errors.New("invalid phone number format, expected +905xxxxxxxxx")
	}

	if len(m.Content) == 0 {
		return errors.New("content cannot be empty")
	}
	if len(m.Content) > MaxMessageLength {
		return errors.New("content exceeds maximum length of 500 characters")
	}

	return nil
}
