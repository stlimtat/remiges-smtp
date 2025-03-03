package pmail

import (
	"github.com/mjl-/mox/smtp"
	"github.com/mjl-/mox/smtpclient"
)

type Mail struct {
	Body        []byte
	BodyHeaders map[string][]byte
	ContentType []byte
	DKIMHeaders []byte
	From        smtp.Address
	Metadata    map[string][]byte
	MsgID       []byte
	MsgPrefix   []byte
	Subject     []byte
	To          []smtp.Address
}

type Response struct {
	smtpclient.Response
}
