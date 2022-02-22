package boundary

import (
	"errors"
	"strings"
)

func TokenIdFromToken(token string) (string, error) {
	split := strings.Split(token, "_")
	if len(split) < 3 {
		return "", errors.New("unexpected stored token format")
	}

	return strings.Join(split[0:2], "_"), nil
}
