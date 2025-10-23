package ubot

import (
	"https://github.com/zuchzub/Go/pkg/vc/ntgcalls"

	tg "github.com/amarnathcjd/gogram/telegram"
)

func (ctx *Context) Calls() map[int64]*ntgcalls.CallInfo {
	return ctx.binding.Calls()
}

func (ctx *Context) InputGroupCall(chatId int64) tg.InputGroupCall {
	return ctx.inputGroupCalls[chatId]
}
