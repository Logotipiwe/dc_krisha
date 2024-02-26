package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jinzhu/gorm"
	"krisha/src/internal/service/parser"
	"krisha/src/internal/service/tg"
	"krisha/src/pkg"
	"krisha/src/tghttp"
	"os"
	"strconv"
)

type Controller struct {
	TgInteractor *tghttp.TgInteractor
	Router       *gin.Engine
}

type MockIncomeMessage struct {
	ChatID int64
	Text   string
}

func NewController(tgInteractor *tghttp.TgInteractor, db *gorm.DB,
	parserService *parser.Service) *Controller {
	router := gin.Default()
	controller := &Controller{tgInteractor, router}

	testGroup := router.Group("/tests")
	testGroup.GET("/enabled", func(c *gin.Context) {
		if pkg.IsTesting() {
			c.String(200, "true")
		} else {
			c.String(200, "false")
		}
	})
	if pkg.IsTesting() {
		testGroup.POST("/sendMessage", func(c *gin.Context) {
			var msg MockIncomeMessage
			if err := c.BindJSON(&msg); err != nil {
				c.Status(400)
				return
			}
			fmt.Printf("Got mock message %v\n", msg)
			err := tgInteractor.AcceptMessage(createMockUpdate(msg))
			if err != nil {
				c.Status(500)
				return
			}
			c.Status(201)
			return
		})
		testGroup.GET("/getAnswerMessages", func(c *gin.Context) {
			c.JSON(200, tg.GetSentMessages())
		})
		testGroup.POST("/reset", func(c *gin.Context) {
			err := parserService.StopAllParsersOnlyInGoroutines()
			if err != nil {
				c.AbortWithError(500, err)
				return
			}
			fmt.Println("Clearing all db data!")
			err = db.Begin().
				Exec("TRUNCATE TABLE krisha.parsers_settings").Commit().Error
			if err != nil {
				c.AbortWithError(500, err)
				return
			}
			tg.GetSentMessages() //Perform clearing
			c.AbortWithStatus(200)
			return
		})
		testGroup.POST("/set-auto-grant-limit", func(c *gin.Context) {
			nStr := c.Query("n")
			os.Setenv("AUTO_GRANT_LIMIT", nStr)
			c.Status(200)
		})
		testGroup.POST("/set-auto-stop-interval", func(c *gin.Context) {
			nStr := c.Query("n")
			os.Setenv("AUTO_STOP_INTERVAL_SEC", nStr)
			c.Status(200)
		})
	}

	return controller
}

func createMockUpdate(msg MockIncomeMessage) tgbotapi.Update {
	return tgbotapi.Update{
		UpdateID: 0,
		Message: &tgbotapi.Message{
			MessageID: 0,
			From: &tgbotapi.User{
				ID:                      msg.ChatID,
				IsBot:                   false,
				FirstName:               "mock name",
				LastName:                "mock lastname",
				UserName:                "mock username",
				LanguageCode:            "",
				CanJoinGroups:           false,
				CanReadAllGroupMessages: false,
				SupportsInlineQueries:   false,
			},
			SenderChat: nil,
			Date:       0,
			Chat: &tgbotapi.Chat{
				ID:                    msg.ChatID,
				Type:                  "",
				Title:                 "Title " + strconv.FormatInt(msg.ChatID, 10),
				UserName:              "",
				FirstName:             "",
				LastName:              "",
				Photo:                 nil,
				Bio:                   "",
				HasPrivateForwards:    false,
				Description:           "",
				InviteLink:            "",
				PinnedMessage:         nil,
				Permissions:           nil,
				SlowModeDelay:         0,
				MessageAutoDeleteTime: 0,
				HasProtectedContent:   false,
				StickerSetName:        "",
				CanSetStickerSet:      false,
				LinkedChatID:          0,
				Location:              nil,
			},
			ForwardFrom:                   nil,
			ForwardFromChat:               nil,
			ForwardFromMessageID:          0,
			ForwardSignature:              "",
			ForwardSenderName:             "",
			ForwardDate:                   0,
			IsAutomaticForward:            false,
			ReplyToMessage:                nil,
			ViaBot:                        nil,
			EditDate:                      0,
			HasProtectedContent:           false,
			MediaGroupID:                  "",
			AuthorSignature:               "",
			Text:                          msg.Text,
			Entities:                      nil,
			Animation:                     nil,
			Audio:                         nil,
			Document:                      nil,
			Photo:                         nil,
			Sticker:                       nil,
			Video:                         nil,
			VideoNote:                     nil,
			Voice:                         nil,
			Caption:                       "",
			CaptionEntities:               nil,
			Contact:                       nil,
			Dice:                          nil,
			Game:                          nil,
			Poll:                          nil,
			Venue:                         nil,
			Location:                      nil,
			NewChatMembers:                nil,
			LeftChatMember:                nil,
			NewChatTitle:                  "",
			NewChatPhoto:                  nil,
			DeleteChatPhoto:               false,
			GroupChatCreated:              false,
			SuperGroupChatCreated:         false,
			ChannelChatCreated:            false,
			MessageAutoDeleteTimerChanged: nil,
			MigrateToChatID:               0,
			MigrateFromChatID:             0,
			PinnedMessage:                 nil,
			Invoice:                       nil,
			SuccessfulPayment:             nil,
			ConnectedWebsite:              "",
			PassportData:                  nil,
			ProximityAlertTriggered:       nil,
			VoiceChatScheduled:            nil,
			VoiceChatStarted:              nil,
			VoiceChatEnded:                nil,
			VoiceChatParticipantsInvited:  nil,
			ReplyMarkup:                   nil,
		},
		EditedMessage:      nil,
		ChannelPost:        nil,
		EditedChannelPost:  nil,
		InlineQuery:        nil,
		ChosenInlineResult: nil,
		CallbackQuery:      nil,
		ShippingQuery:      nil,
		PreCheckoutQuery:   nil,
		Poll:               nil,
		PollAnswer:         nil,
		MyChatMember:       nil,
		ChatMember:         nil,
		ChatJoinRequest:    nil,
	}
}

func (c *Controller) Start() {
	c.Router.Run(":8083")
}
