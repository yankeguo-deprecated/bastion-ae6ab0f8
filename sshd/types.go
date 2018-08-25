package sshd

// this package contains SSH protocol specified constants and structures
// see https://www.iana.org/assignments/ssh-parameters/ssh-parameters.xhtml
// see https://tools.ietf.org/html/rfc4254

const (
	ChannelTypeDirectTCPIP = "direct-tcpip"
	ChannelTypeSession     = "session"

	RequestTypePtyReq       = "pty-req"
	RequestTypeX11Req       = "x11-req"
	RequestTypeEnv          = "env"
	RequestTypeShell        = "shell"
	RequestTypeExec         = "exec"
	RequestTypeSubsystem    = "subsystem"
	RequestTypeWindowChange = "window-change"
	RequestTypeExitStatus   = "exit-status"
)

const (
	extKeyStage    = "bastion-mode"
	extKeyAccount  = "bastion-account"
	extKeyHostname = "bastion-hostname"
	extKeyUser     = "bastion-user"
	extKeyAddress  = "bastion-address"

	stagePre = "pre"
	stageLv1 = "lv1"
	stageLv2 = "lv2"
)

type DirectTCPIPExtraData struct {
	Host           string
	Port           uint32
	OriginatorIP   string
	OriginatorPort uint32
}

type PtyRequestPayload struct {
	Term   string
	Cols   uint32
	Rows   uint32
	Width  uint32
	Height uint32
	Modes  string
}

type EnvRequestPayload struct {
	Name  string
	Value string
}

type ExecRequestPayload struct {
	Command string
}

type WindowChangeRequestPayload struct {
	Cols   uint32
	Rows   uint32
	Width  uint32
	Height uint32
}

type ExitStatusRequestPayload struct {
	Code uint32
}
