package squid

import (
	"errors"
	"strings"
)

// format: %LOGIN %>{Proxy-Authorization} %SRC %SRCPORT %DST %PORT
type ExtAclRequest struct {
	Login   string
	Token   string
	Src     string
	SrcPort string
	Dst     string
	DstPort string
}

func NewExtAclRequest(line string) (*ExtAclRequest, error) {
	values := strings.Split(line, " ")
	if len(values) < 7 {
		return nil, errors.New("malformed request: insufficient format values")
	}

	token, err := PasswordFromProxyAuthorization(values[1])
	if err != nil {
		return nil, err
	}

	return &ExtAclRequest{
		Login:   values[0],
		Token:   token,
		Src:     values[2],
		SrcPort: values[3],
		Dst:     values[4],
		DstPort: values[5],
	}, nil
}
