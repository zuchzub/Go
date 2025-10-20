package types

import (
	"https://github.com/iamnolimit/tggomusicbot/pkg/vc/ntgcalls"
)

type PendingConnection struct {
	MediaDescription ntgcalls.MediaDescription
	Payload          string
	Presentation     bool
}
