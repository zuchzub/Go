package handlers

import (
	"https://github.com/iamnolimit/tggomusicbot/pkg/core/cache"
	"https://github.com/iamnolimit/tggomusicbot/pkg/core/db"
	"https://github.com/iamnolimit/tggomusicbot/pkg/lang"
	"https://github.com/iamnolimit/tggomusicbot/pkg/vc"

	"github.com/amarnathcjd/gogram/telegram"
)

// skipHandler handles the /skip command.
func skipHandler(m *telegram.NewMessage) error {
	chatID, _ := getPeerId(m.Client, m.ChatID())
	ctx, cancel := db.Ctx()
	defer cancel()
	langCode := db.Instance.GetLang(ctx, chatID)
	if !cache.ChatCache.IsActive(chatID) {
		_, _ = m.Reply(lang.GetString(langCode, "no_track_playing"))
		return nil
	}

	_ = vc.Calls.PlayNext(chatID)
	return nil
}
