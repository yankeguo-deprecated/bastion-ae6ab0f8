package types

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"regexp"
	"strings"
)

const (
	NodeSourceManual = "manual"
	NodeSourceConsul = "consul"

	NodeUserRoot = "root"
)

var (
	UserAccountPattern    = regexp.MustCompile(`^[a-zA-Z][0-9a-zA-Z_.-]{3,24}$`)
	UserNicknameMaxLength = 10
	UserPasswordMinLength = 6

	NodeHostnamePattern = regexp.MustCompile(`[0-9a-zA-Z_.-]{4,64}`)
	NodeUserPattern     = UserAccountPattern
	NodeAddressPattern  = regexp.MustCompile(`^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}(:\d+)?$`)

	KeyFingerprintPattern = regexp.MustCompile(`SHA256:[0-9a-zA-Z+/]{43}`)
)

func errMissingField(key string) error {
	return status.Errorf(codes.InvalidArgument, "missing field '%s'", key)
}

func errInvalidField(key string, should string) error {
	return status.Errorf(codes.InvalidArgument, "invalid field '%s', should be %s", key, should)
}

func trimSpace(s *string) {
	*s = strings.TrimSpace(*s)
}

type Validator interface {
	Validate() (error)
}

func (m *CreateUserRequest) Validate() (err error) {
	trimSpace(&m.Account)
	if !UserAccountPattern.MatchString(m.Account) {
		err = errInvalidField("account", "valid linux user name")
		return
	}
	if len(m.Password) < UserPasswordMinLength {
		err = errInvalidField("password", fmt.Sprintf("longer than %d characteristics", UserPasswordMinLength))
		return
	}
	trimSpace(&m.Nickname)
	if len(m.Nickname) > UserNicknameMaxLength {
		err = errInvalidField("nickname", fmt.Sprintf("shorter than %d characterstics", UserNicknameMaxLength))
		return
	} else if len(m.Nickname) == 0 {
		m.Nickname = m.Account
	}
	return
}

func (m *UpdateUserRequest) Validate() (err error) {
	trimSpace(&m.Account)
	if m.UpdateNickname {
		trimSpace(&m.Nickname)
		if len(m.Nickname) > UserNicknameMaxLength {
			err = errInvalidField("nickname", fmt.Sprintf("shorter than %d characterstics", UserNicknameMaxLength))
			return
		} else if len(m.Nickname) == 0 {
			m.UpdateNickname = false
		}
	}
	if m.UpdatePassword {
		if len(m.Password) < UserPasswordMinLength {
			err = errInvalidField("password", fmt.Sprintf("longer than %d characteristics", UserPasswordMinLength))
			return
		}
	}
	return
}

func (m *PutNodeRequest) Validate() (err error) {
	trimSpace(&m.Hostname)
	if !NodeHostnamePattern.MatchString(m.Hostname) {
		err = errInvalidField("hostname", "valid hostname")
		return
	}
	trimSpace(&m.User)
	if len(m.User) == 0 {
		m.User = NodeUserRoot
	} else {
		if !NodeUserPattern.MatchString(m.User) {
			err = errInvalidField("user", "valid linux user name")
			return
		}
	}
	trimSpace(&m.Address)
	if !NodeAddressPattern.MatchString(m.Address) {
		err = errInvalidField("address", "IPv4 address with an optional port")
		return
	}
	trimSpace(&m.Source)
	if len(m.Source) == 0 {
		m.Source = NodeSourceManual
	} else if m.Source != NodeSourceManual && m.Source != NodeSourceConsul {
		err = errInvalidField("source", "one of 'manual' or 'consul'")
		return
	}
	return
}

func (m *CreateKeyRequest) Validate() (err error) {
	trimSpace(&m.Account)
	if len(m.Account) == 0 {
		err = errMissingField("account")
		return
	}
	trimSpace(&m.Fingerprint)
	if !KeyFingerprintPattern.MatchString(m.Fingerprint) {
		err = errInvalidField("fingerprint", "a valid ssh fingerprint in sha256 digest")
		return
	}
	trimSpace(&m.Name)
	if len(m.Name) == 0 {
		m.Name = "no name"
	}
	return
}

func (m *DeleteKeyRequest) Validate() (err error) {
	trimSpace(&m.Fingerprint)
	if !KeyFingerprintPattern.MatchString(m.Fingerprint) {
		err = errInvalidField("fingerprint", "a valid ssh fingerprint in sha256 digest")
		return
	}
	return
}
