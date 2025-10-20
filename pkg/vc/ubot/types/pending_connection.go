package types

import (
	"github.com/AshokShau/TgMusicBot/pkg/vc/ntgcalls"
)

type PendingConnection struct {
	MediaDescription ntgcalls.MediaDescription
	Payload          string
	Presentation     bool
}
