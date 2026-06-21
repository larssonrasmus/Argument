package auth

import (
	"unicode"

	"argument/internal/config"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var store = sessions.NewCookieStore([]byte(config.Conf.SessionKey))

func isNonspacing(r rune) bool {
	return unicode.Is(unicode.Mn, r)
}

// NormalizeUsername removes unicode to prevent similar looking usernames
func NormalizeUsername(username string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isNonspacing), norm.NFC)
	n, _, _ := transform.String(t, username)

	return n
}

// HashPassword generates a hash from a plaintext password
func HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)

	return string(bytes)
}

// CheckPassword compares a hash and plaintext password
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}
