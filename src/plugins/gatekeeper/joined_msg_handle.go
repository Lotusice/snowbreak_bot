package gatekeeper

import (
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"github.com/spf13/viper"
	"snowbreak_bot/utils"
)

func JoinedMsgHandle(update tgbotapi.Update) error {
	message := update.Message
	message.Delete()
	for _, member := range message.NewChatMembers {
		if member.ID == message.From.ID { // 自己加入群组
			continue
		}
		// 机器人被邀请加群
		if member.UserName == viper.GetString("bot.name") {
			utils.SaveJoined(message)
			continue
		}
		// 邀请加入群组，无需进行验证
		utils.SaveInvite(message, &member)
	}
	return nil
}
