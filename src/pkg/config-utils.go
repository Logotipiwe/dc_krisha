package pkg

import (
	"errors"
	_ "github.com/lib/pq"
	. "github.com/logotipiwe/dc_go_config_lib"
	"log"
	"strconv"
)

func GetTgFilterHelpMessage() string {
	return GetConfigOr("TG_MESSAGE_FILTER_HELP", "Ошибка получения сообщения с инструкцией. Обратитесь к администратору")
}

func GetTgStartMessage() string {
	return GetConfigOr("TG_MESSAGE_START", "Ошибка получения приветственного сообщения. Обратитесь к администратору")
}

func GetTgFaqMessage() string {
	return GetConfigOr("TG_MESSAGE_FAQ", "Часто задаваемые вопросы ещё не заполнены. Обратитесь к администратору")
}

func GetOwnerChatID() int64 {
	ownerChatIdStr := GetConfig("OWNER_TG_CHAT_ID")
	if ownerChatIdStr == "" {
		log.Println(errors.New("empty owner chat, unable to send message"))
		return 0
	}
	ownerChatID, _ := strconv.ParseInt(ownerChatIdStr, 10, 64)
	return ownerChatID
}

func IsTesting() bool {
	is, _ := GetConfigBool("IS_TESTING")
	return is
}

func GetAutoGrantLimit() int {
	configInt, err := GetConfigInt("AUTO_GRANT_LIMIT")
	if err != nil {
		return 0
	}
	return configInt
}

func GetAutoStopHours() int {
	return GetConfigIntOr("AUTO_STOP_HOURS", 0)
}

func GetTgHelpMessage() string {
	return GetConfigOr("TG_MESSAGE_HELP", "Ошибка получения сообщения с инструкцией. Обратитесь к администратору")
}
