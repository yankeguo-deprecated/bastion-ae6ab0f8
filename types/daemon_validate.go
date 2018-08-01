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

	GrantHostnamePatternPattern = regexp.MustCompile(`[0-9a-zA-Z_.*-]{4,64}`)
	GrantUserPattern            = UserAccountPattern

	errInvalidFingerprint = errInvalidField("fingerprint", "a valid ssh sha256 fingerprint of public key")
)

func errMissingField(key string) error {
	return status.Errorf(codes.InvalidArgument, "missing field '%s'", key)
}

func errInvalidField(key string, should string) error {
	return status.Errorf(codes.InvalidArgument, "invalid field '%s', should be %s", key, should)
}

func trimSpace(s *string) {
	if s != nil {
		*s = strings.TrimSpace(*s)
	}
}

type Validator interface {
	Validate() error
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

func (m *GetUserRequest) Validate() (err error) {
	trimSpace(&m.Account)
	if len(m.Account) == 0 {
		err = errMissingField("account")
		return
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

func (m *GetNodeRequest) Validate() (err error) {
	trimSpace(&m.Hostname)
	if len(m.Hostname) == 0 {
		err = errMissingField("hostname")
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
		err = errInvalidFingerprint
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
		err = errInvalidFingerprint
		return
	}
	return
}

func (m *GetKeyRequest) Validate() (err error) {
	trimSpace(&m.Fingerprint)
	if !KeyFingerprintPattern.MatchString(m.Fingerprint) {
		err = errInvalidFingerprint
		return
	}
	return
}

func (m *PutGrantRequest) Validate() (err error) {
	trimSpace(&m.Account)
	if len(m.Account) == 0 {
		err = errMissingField("account")
		return
	}
	trimSpace(&m.HostnamePattern)
	if !GrantHostnamePatternPattern.MatchString(m.HostnamePattern) {
		err = errInvalidField("hostname_pattern", "valid hostname pattern with options wildcards")
		return
	}
	trimSpace(&m.User)
	if !GrantUserPattern.MatchString(m.User) {
		err = errInvalidField("user", "a valid linux user")
		return
	}
	return
}

func (m *ListGrantsRequest) Validate() (err error) {
	trimSpace(&m.Account)
	if len(m.Account) == 0 {
		err = errMissingField("account")
		return
	}
	return
}

func (m *DeleteGrantRequest) Validate() (err error) {
	trimSpace(&m.Account)
	if len(m.Account) == 0 {
		err = errMissingField("account")
		return
	}
	trimSpace(&m.HostnamePattern)
	if !GrantHostnamePatternPattern.MatchString(m.HostnamePattern) {
		err = errInvalidField("hostname_pattern", "valid hostname pattern with wildcard support")
		return
	}
	trimSpace(&m.User)
	if !GrantUserPattern.MatchString(m.User) {
		err = errInvalidField("user", "valid linux user")
		return
	}
	return
}

func (m *CheckGrantRequest) Validate() (err error) {
	trimSpace(&m.Account)
	if len(m.Account) == 0 {
		err = errMissingField("account")
		return
	}
	trimSpace(&m.Hostname)
	if !NodeHostnamePattern.MatchString(m.Hostname) {
		err = errInvalidField("hostname", "valid hostname")
		return
	}
	trimSpace(&m.User)
	if !GrantUserPattern.MatchString(m.User) {
		err = errInvalidField("user", "valid linux user")
		return
	}
	return
}

func (m *ListGrantItemsRequest) Validate() (err error) {
	trimSpace(&m.Account)
	if len(m.Account) == 0 {
		err = errMissingField("account")
		return
	}
	return
}
