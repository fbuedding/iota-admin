package cookies

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

var (
	ErrValueTooLong = errors.New("cookie value too long")
	ErrInvalidValue = errors.New("invalid cookie value")
)

func New(name string, value string, expires time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Expires:  expires,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		Path:     "/",
		Value:    value,
	}
}

func Write(w http.ResponseWriter, cookie *http.Cookie) error {
	cookie.Value = base64.URLEncoding.EncodeToString([]byte(cookie.Value))

	if len(cookie.String()) > 4096 {
		return ErrValueTooLong
	}

	http.SetCookie(w, cookie)

	return nil
}

func Read(r *http.Request, name string) (string, error) {
	// Read the cookie as normal.
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	value, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return "", ErrInvalidValue
	}

	return string(value), nil
}

func Delete(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:   name,
		MaxAge: -1,
	})
}

func WriteSigned(w http.ResponseWriter, cookie *http.Cookie, secretKey []byte) error {
	mac := hmac.New(sha256.New, secretKey)
	_, err := mac.Write([]byte(cookie.Name))
	if err != nil {
		return err
	}

	_, err = mac.Write([]byte(cookie.Value))
	if err != nil {
		return err
	}

	signature := mac.Sum(nil)
	cookie.Value = string(signature) + cookie.Value

	return Write(w, cookie)
}

func ReadSigned(r *http.Request, name string, secretKey []byte) (string, error) {
	signedValue, err := Read(r, name)
	if err != nil {
		return "", err
	}

	if len(signedValue) < sha256.Size {
		log.Debug().Str("signed cookie value", signedValue).Int("Length", len(signedValue)).Send()
		return "", ErrInvalidValue
	}

	signature := signedValue[:sha256.Size]
	value := signedValue[sha256.Size:]

	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(name))
	mac.Write([]byte(value))
	expectedSignature := mac.Sum(nil)

	if !hmac.Equal([]byte(signature), expectedSignature) {
		log.Debug().Bytes("singnature", []byte(signature)).Bytes("expected signature", expectedSignature).Msg("Cookie has wrong signature")
		return "", ErrInvalidValue
	}
	return value, nil
}
