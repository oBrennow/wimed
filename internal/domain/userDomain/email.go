package userDomain

import (
	"strings"
)

type Email struct {
	value string
}

func NewEmail(raw string) (Email, error) {
	v := strings.TrimSpace(strings.ToLower(raw))
	if v == "" || !looksLikeEmail(v) {
		return Email{}, ErrInvalidEmail
	}
	return Email{value: v}, nil
}

func (e Email) Value() string {
	return e.value
}

func looksLikeEmail(s string) bool {
	at := strings.Index(s, "@")
	dot := strings.LastIndex(s, ".")

	return at > 0 &&
		dot > at+1 &&
		dot < len(s)-1
}
