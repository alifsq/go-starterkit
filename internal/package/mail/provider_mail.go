package mail

import "github.com/google/wire"

var MailSet = wire.NewSet(
	NewMailer,
)