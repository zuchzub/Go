package pkg

import (
	"https://github.com/iamnolimit/tggomusicbot/pkg/config"
	"https://github.com/iamnolimit/tggomusicbot/pkg/handlers"
	"https://github.com/iamnolimit/tggomusicbot/pkg/vc"

	tg "github.com/amarnathcjd/gogram/telegram"
)

func Init(client *tg.Client) error {
	for _, session := range config.Conf.SessionStrings {
		_, err := vc.Calls.StartClient(config.Conf.ApiId, config.Conf.ApiHash, session)
		if err != nil {
			return err
		}
	}

	vc.Calls.RegisterHandlers(client)
	handlers.LoadModules(client)
	return nil
}
