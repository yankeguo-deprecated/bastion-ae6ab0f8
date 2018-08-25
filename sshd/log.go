package sshd

import (
	"encoding/base64"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh"
)

func cookEvent(cm ssh.ConnMetadata, e *zerolog.Event) *zerolog.Event {
	e = e.Str("conn", base64.URLEncoding.EncodeToString(cm.SessionID()))
	if conn, ok := cm.(*ssh.ServerConn); ok {
		if conn.Permissions.Extensions[extKeyStage] == stageLv1 {
			e = e.Str("stage", "lv1")
			e = e.Str("account", conn.Permissions.Extensions[extKeyAccount])
		} else if conn.Permissions.Extensions[extKeyStage] == stageLv2 {
			e = e.Str("stage", "lv2")
			e = e.Str("account", conn.Permissions.Extensions[extKeyAccount])
			e = e.Str("user", conn.Permissions.Extensions[extKeyUser])
			e = e.Str("address", conn.Permissions.Extensions[extKeyAddress])
			e = e.Str("hostname", conn.Permissions.Extensions[extKeyHostname])
		} else {
			e = e.Str("stage", "unknown")
		}
	} else {
		e = e.Str("stage", "pre")
	}
	return e
}

func DLog(conn ssh.ConnMetadata) *zerolog.Event {
	return cookEvent(conn, log.Debug())
}

func ILog(conn ssh.ConnMetadata) *zerolog.Event {
	return cookEvent(conn, log.Info())
}

func ELog(conn ssh.ConnMetadata) *zerolog.Event {
	return cookEvent(conn, log.Error())
}
