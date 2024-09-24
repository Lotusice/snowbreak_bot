package gatekeeper

import (
	"fmt"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	bot "snowbreak_bot/config"
	"snowbreak_bot/plugins/messagecleaner"
	"snowbreak_bot/utils"
	"strconv"
	"strings"
)

func CallBackData(callBack tgbotapi.Update) error {
	callbackQuery := callBack.CallbackQuery
	data := callBack.CallbackData()
	d := strings.Split(data, ",")

	if len(d) < 4 {
		return nil
	}

	userId, _ := strconv.ParseInt(d[1], 10, 64)
	chatId := callbackQuery.Message.Chat.ID
	joinMessageId, _ := strconv.Atoi(d[3])

	if d[2] == "PASS" || d[2] == "BAN" {

		if !bot.Snowbreak.IsAdminWithPermissions(chatId, callbackQuery.From.ID, tgbotapi.AdminCanRestrictMembers) {
			callbackQuery.Answer(true, "无使用权限！")
			return nil
		}
		if has, _ := verifySet.checkExistAndRemove(userId, chatId); has {
			if d[2] == "PASS" {
				err := pass(chatId, userId, callbackQuery, true)
				return err
			}

			if d[2] == "BAN" {
				ban(chatId, userId, callbackQuery, joinMessageId)
			}
		}
		return nil
	}

	if userId != callbackQuery.From.ID {
		callbackQuery.Answer(true, "这不是你的验证！")
		return nil
	}
	if has, correct := verifySet.checkExistAndRemove(userId, chatId); has {
		if d[2] != correct {
			callbackQuery.Answer(true, "验证未通过，请一分钟后再试！")
			ban(chatId, userId, callbackQuery, joinMessageId)
			go unban(chatId, userId)
			return nil
		}

		callbackQuery.Answer(true, "验证通过！")
		err := pass(chatId, userId, callbackQuery, false)
		return err
	}
	return nil
}

func pass(chatId int64, userId int64, callbackQuery *tgbotapi.CallbackQuery, adminPass bool) error {
	bot.Snowbreak.RestrictChatMember(chatId, userId, tgbotapi.AllPermissions)
	callbackQuery.Delete()
	if !adminPass {
		// 新人发送box提醒
		text := fmt.Sprintf("欢迎[%s](tg://user?id=%d)\n", tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, callbackQuery.From.FullName()), callbackQuery.From.ID)
		var joined utils.GroupJoined
		utils.GetJoinedByChatId(chatId).Scan(&joined)
		if joined.Reg != -1 {
			text += fmt.Sprintf("建议阅读群公约：[点击阅读](https://t.me/%s/%d)", callbackQuery.Message.Chat.UserName, joined.Reg)
		}
		sendMessage := tgbotapi.NewMessage(chatId, text)
		sendMessage.ParseMode = tgbotapi.ModeMarkdownV2
		msg, err := bot.Snowbreak.Send(sendMessage)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(chatId, msg.MessageID, 3600)

	}
	return nil
}

func ban(chatId int64, userId int64, callbackQuery *tgbotapi.CallbackQuery, joinMessageId int) {
	bot.Snowbreak.BanChatMember(chatId, userId)
	callbackQuery.Delete()
	delJoinMessage := tgbotapi.NewDeleteMessage(chatId, joinMessageId)
	bot.Snowbreak.Send(delJoinMessage)
}
