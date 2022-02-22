package squid

import (
	"encoding/base64"
	"errors"
	"net/url"
	"strings"
)

func PasswordFromProxyAuthorization(value string) (string, error) {
	unescape, err := url.QueryUnescape(value)
	if err != nil {
		return "", err
	}

	parts := strings.Split(unescape, " ")
	if len(parts) < 2 {
		return "", errors.New("bad header value")
	}

	raw, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", err
	}

	parts = strings.Split(string(raw), ":")
	if len(parts) < 2 {
		return "", errors.New("bad header value")
	}

	return parts[1], nil
}
