package types

import (
	"github.com/zuchzub/Go/pkg/vc/ntgcalls"
)

type PendingConnection struct {
	MediaDescription ntgcalls.MediaDescription
	Payload          string
	Presentation     bool
}
