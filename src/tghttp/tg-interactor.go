package tghttp

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"krisha/src/internal"
	"krisha/src/internal/domain"
	"krisha/src/internal/service/parser"
	"krisha/src/internal/service/tg"
	"krisha/src/pkg"
	"strconv"
	"strings"
)

type TgInteractor struct {
	tgService          *tg.TgService
	parserService      *parser.Service
	permissionsService *internal.PermissionsService
	ownerChatMode      OwnerChatMode
}

type OwnerChatMode int

const (
	defaultMode OwnerChatMode = iota
	granting
	denying
)

const (
	unknownMessage = "Не понял команду. Попробуйте /help, чтобы получить список команд"
	ownerHelp      = `Здарова админ
/grant - выдать доступ
/deny - забрать доступ`
	userHelp = `/start - запустить бота
/stop - остановить остановить бота

После запуска отправьте нужный фильтр (/filter - инструкция по получению фильтра)
После отправки фильтра вам начнут приходить уведомления по этому фильтру. Боту потребуется некоторое время, чтобы настроиться и обработать объявления, чтобы начать присылать вам сообщения. Это может занять около 10 минут, в зависимости от размера фильтра. `
	ownerUnacceptedError = "unknown admin command"
	errorMessage         = "Произошла ошибка :( \r\n Попробуйте начать заново с /start"
	startAnswer          = "Отправьте фильтр с krisha.kz. Фильтр должен начинаться с знака '?'"
	noAccessMessage      = "У вас нет доступа к боту"
)

func NewTgInteractor(
	tgService *tg.TgService,
	parserService *parser.Service,
	permissionsService *internal.PermissionsService,
) *TgInteractor {
	return &TgInteractor{
		tgService:          tgService,
		parserService:      parserService,
		permissionsService: permissionsService,
		ownerChatMode:      defaultMode,
	}
}

func (i *TgInteractor) acceptMessage(update tgbotapi.Update) error {
	err := i.acceptMessageUnsafe(update)
	if err != nil {
		i.tgService.SendMessage(update.Message.Chat.ID, errorMessage)
	}
	return err
}

func (i *TgInteractor) acceptMessageUnsafe(update tgbotapi.Update) error {
	ownerChatID := pkg.GetOwnerChatID()
	if ownerChatID == update.Message.Chat.ID {
		err := i.acceptAdminMessage(update)
		if err != nil && err.Error() == ownerUnacceptedError {
			return i.acceptUserMessage(update)
		}
		return err
	} else {
		return i.acceptUserMessage(update)
	}
}

func (i *TgInteractor) acceptUserMessage(update tgbotapi.Update) error {
	text := update.Message.Text
	chatID := update.Message.Chat.ID

	hasAccess := i.permissionsService.HasAccess(chatID)
	isOwner := chatID == pkg.GetOwnerChatID()
	if !hasAccess && !isOwner {
		return i.tgService.SendMessage(chatID, noAccessMessage)
	}

	switch {
	case text == "/help":
		return i.tgService.SendMessage(chatID, userHelp)
	case text == "/start":
		err := i.parserService.InitParserSettings(chatID)
		if err != nil {
			return err
		}
		return i.tgService.SendMessage(chatID, startAnswer)
	case strings.HasPrefix(text, "?"):
		err, existed := i.parserService.SetFiltersAndStartParser(chatID, text)
		if err != nil {
			return err
		}
		if existed {
			return i.tgService.SendMessage(chatID, "Фильтр применен")
		} else {
			return i.tgService.SendMessage(chatID, "Фильтр успешно установлен и парсер запущен")
		}
	case text == "/stop":
		err := i.parserService.StopParser(chatID)
		if err != nil {
			if errors.Is(err, domain.ParserNotFoundErr) {
				return i.tgService.SendMessage(chatID, "Парсер уже остановлен")
			}
			return err
		}
		return i.tgService.SendMessage(chatID, "Парсер остановлен")
	}
	return i.tgService.SendMessage(chatID, unknownMessage)
}

func (i *TgInteractor) acceptAdminMessage(update tgbotapi.Update) error {
	text := update.Message.Text
	ownerChatMode := i.ownerChatMode
	i.ownerChatMode = defaultMode
	switch {
	case text == "/help":
		return i.tgService.SendMessageToOwner(ownerHelp + "\r\n\r\n" + userHelp)
	case text == "/grant":
		err := i.tgService.SendMessageToOwner("Какому чату выдать доступ?")
		if err == nil {
			i.ownerChatMode = granting
		}
		return err
	case ownerChatMode == granting:
		grantingChat, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			i.ownerChatMode = granting
			return i.tgService.SendMessageToOwner("Кажется это не число, попробуй ещё раз")
		}
		if i.permissionsService.HasAccess(grantingChat) {
			return i.tgService.SendMessageToOwner("У этого чата и так есть доступ. Спасибо")
		}
		err = i.permissionsService.GrantAccess(grantingChat)
		if err != nil {
			return err
		}
		return i.tgService.SendMessageToOwner("Доступ выдан для чата " + text)
	case text == "/deny":
		err := i.tgService.SendMessageToOwner("Какому чату запретить доступ?")
		if err == nil {
			i.ownerChatMode = denying
		}
		return err
	case ownerChatMode == denying:
		denyingChat, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			i.ownerChatMode = denying
			return i.tgService.SendMessageToOwner("Кажется это не число, попробуй ещё раз")
		}
		if !i.permissionsService.HasAccess(denyingChat) {
			return i.tgService.SendMessageToOwner("У этого чата и так нет доступа. Спасибо")
		}
		//TODO stop chat parser
		err = i.permissionsService.DenyAccess(denyingChat)
		if err != nil {
			return err
		}
		return i.tgService.SendMessageToOwner("Доступ запрещен чату " + text)
	}
	return errors.New(ownerUnacceptedError)
}

func (i *TgInteractor) StartHandleTgMessages() {
	i.tgService.StartReceiveMessages(i.acceptMessage)
}
