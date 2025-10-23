package handlers

import (
	"fmt"
	"github.com/zuchzub/Go/pkg/core/cache"
	"github.com/zuchzub/Go/pkg/core/db"
	"github.com/zuchzub/Go/pkg/lang"
	"github.com/zuchzub/Go/pkg/vc"

	"github.com/amarnathcjd/gogram/telegram"
)

// stopHandler handles the /stop command.
func stopHandler(m *telegram.NewMessage) error {
	chatID, _ := getPeerId(m.Client, m.ChatID())
	ctx, cancel := db.Ctx()
	defer cancel()
	langCode := db.Instance.GetLang(ctx, chatID)
	if !cache.ChatCache.IsActive(chatID) {
		_, _ = m.Reply(lang.GetString(langCode, "no_track_playing"))
		return nil
	}

	if err := vc.Calls.Stop(chatID); err != nil {
		_, _ = m.Reply(fmt.Sprintf(lang.GetString(langCode, "stop_error"), err.Error()))
		return err
	}

	_, _ = m.Reply(fmt.Sprintf(lang.GetString(langCode, "stop_success"), m.Sender.FirstName))
	return nil
}
