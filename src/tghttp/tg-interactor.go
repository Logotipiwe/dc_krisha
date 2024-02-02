package tghttp

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	config "github.com/logotipiwe/dc_go_config_lib"
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
	userHelp = `Бот работает следующим образом:
1. Вы отправляете фильтры, по которым ищете квартиру. Инструкция - /filterHelp
2. Бот присылает вам уведомления о новых квартирах
3. Вы можете писать /stop или /start, чтобы отключать или включать обратно уведомления
4. Вы можете отправить другой фильтр, чтобы бот перенастроился на него`
	ownerUnacceptedError = "unknown admin command"
	errorMessage         = "Произошла ошибка :( \r\n Попробуйте начать заново с /start"
	startAnswer          = `Привет! Это бот для получения уведомлений о новых квартирах по фильтрам. Для начала работы нужно установить нужный фильтр.
/help - общая информация и инструкция.`
	filterHelpAnswer = `Чтобы установить фильтр, нужно:
1. Зайти на https://krisha.kz/map/arenda/kvartiry/almaty/

2. Выбрать нужные фильтры в панели над картой (полезно бывает обвести область или поставить "от хозяев"). У вашего чата есть ограничение по количеству квартир в фильтре, поэтому постарайтесь оставить только те, которые вас правда интересуют

3. ВАЖНО - нажать синюю кнопку "показать результаты", чтобы фильтр отобразился в ссылке

4. Данные фильтра должны появиться в адресной строке браузера (ссылка на текущую страницу). Скопируйте текущую ссылку из адресной строки. Пример нужной ссылки:
https://krisha.kz/map/arenda/kvartiry/almaty/?areas=&das[live.rooms]=1&das[price][to]=234343&das[who]=1&das[who_match][4]=4&zoom=14&lat=43.23518&lon=76.93178

5. Отправьте ссылку сюда, бот сразу начнёт искать варианты по этому фильтру.

Чтобы ИЗМЕНИТЬ фильтр - просто отправьте новую ссылку с фильтром и бот перенастроится на него

/stop - остановить уведомления
/help - общая информация
По доп. вопросам обращайтесь к администратору :)`
	noAccessMessage            = "У вас нет доступа к боту"
	limitExceededMessageFormat = "Превышен лимит в %v квартир в вашем фильтре. Попробуйте другой фильтр"
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
		if update.Message.Chat.ID == pkg.GetOwnerChatID() {
			i.tgService.SendMessageToOwner(errorMessage + "\r\n" + err.Error())
		} else {
			i.tgService.SendMessage(update.Message.Chat.ID, errorMessage)
		}
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
		//keep silence if user has no access
		return nil
	}

	switch {
	case strings.HasPrefix(text, "/help"):
		return i.tgService.SendMessage(chatID, userHelp)
	case strings.HasPrefix(text, "/filterHelp"):
		return i.tgService.SendMessage(chatID, filterHelpAnswer)
	case strings.HasPrefix(text, "/start"):
		return i.handleUserStartCommand(chatID)
	case strings.Contains(text, "krisha.kz"):
		pair := strings.Split(text, "?")
		if len(pair) < 2 {
			return i.tgService.SendMessage(chatID, "Неверный формат ссылки с фильтрами. Ссылка должна начинаться на krisha.kz и содержать фильтры через знак ?")
		}
		settings, err := i.parserService.SetFilters(chatID, "?"+pair[1])
		if err != nil {
			return err
		}
		err, existed := i.parserService.StartParser(settings, true)
		if err != nil {
			return i.handleStartParserErr(settings, err)
		}
		if existed {
			return i.tgService.SendMessage(chatID, "Фильтр применен, перезапускаюсь...")
		} else {
			return i.tgService.SendMessage(chatID, "Фильтр успешно установлен и парсер запущен")
		}
	case strings.HasPrefix(text, "/stop"):
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

func (i *TgInteractor) handleUserStartCommand(chatID int64) error {
	settings, isErrNotFound, err := i.parserService.GetSettings(chatID)
	autoGrantLimit, intErr := config.GetConfigInt("AUTO_GRANT_LIMIT")

	if isErrNotFound {
		if intErr == nil && autoGrantLimit > 0 {
			err = i.parserService.CreateParserSettings(chatID, autoGrantLimit)
			if err != nil {
				return err
			}
		} else {
			return i.tgService.SendMessage(chatID, "У вашего чата есть доступ, но для него не найдено нужных настроек. Обратитесь к администратору")
		}
	} else {
		if err != nil {
			return err
		}
	}

	if settings.Filters == "" {
		return i.tgService.SendMessage(chatID, startAnswer)
	}
	err, existed := i.parserService.StartParser(settings, false)
	if err != nil {
		return i.handleStartParserErr(settings, err)
	}
	if existed {
		return i.tgService.SendMessage(chatID, "Парсер уже запущен")
	} else {
		return i.tgService.SendMessage(chatID, "Парсер запущен")
	}
}

func (i *TgInteractor) acceptAdminMessage(update tgbotapi.Update) error {
	text := update.Message.Text
	ownerChatMode := i.ownerChatMode
	i.ownerChatMode = defaultMode
	switch {
	case text == "/help":
		return i.tgService.SendMessageToOwner(ownerHelp + "\r\n\r\n" + userHelp)
	case text == "/grant":
		err := i.tgService.SendMessageToOwner("Какому чату выдать доступ? И через пробел - лимит")
		if err == nil {
			i.ownerChatMode = granting
		}
		return err
	case ownerChatMode == granting:
		return i.handleGrantCommand(text)
	case text == "/deny":
		err := i.tgService.SendMessageToOwner("Какому чату запретить доступ?")
		if err == nil {
			i.ownerChatMode = denying
		}
		return err
	case ownerChatMode == denying:
		return i.handleDenyCommand(text)
	}
	return errors.New(ownerUnacceptedError)
}

func (i *TgInteractor) handleDenyCommand(text string) error {
	denyingChat, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		i.ownerChatMode = denying
		return i.tgService.SendMessageToOwner("Кажется это не число, попробуй ещё раз")
	}
	stopParserErr := i.parserService.StopParser(denyingChat)
	if !i.permissionsService.HasAccess(denyingChat) {
		return i.tgService.SendMessageToOwner("У этого чата и так нет доступа. Спасибо")
	}
	err = i.permissionsService.DenyAccess(denyingChat)
	if err != nil {
		return err
	}
	if !errors.Is(stopParserErr, domain.ParserNotFoundErr) {
		i.tgService.SendMessage(denyingChat, "Парсер остановлен, обратитесь к администратору")
	}
	return i.tgService.SendMessageToOwner("Доступ запрещен чату " + text)
}

func (i *TgInteractor) handleGrantCommand(text string) error {
	grantingChat, limit, err := parseGrantMessage(text)
	if err != nil {
		i.ownerChatMode = granting
		return i.tgService.SendMessageToOwner("Ошибка " + err.Error())
	}
	err = i.permissionsService.GrantAccess(grantingChat)
	if err != nil {
		return err
	}
	settings, isErrNotFound, err := i.parserService.GetSettings(grantingChat)
	if err != nil && !isErrNotFound {
		return err
	}
	if isErrNotFound {
		err = i.parserService.CreateParserSettings(grantingChat, limit)
		if err != nil {
			return err
		}
	} else {
		err, stopped := i.parserService.UpdateLimit(settings, limit)
		if err != nil {
			return err
		}
		if stopped {
			return i.tgService.SendMessage(grantingChat, fmt.Sprintf(limitExceededMessageFormat, limit))
		}
	}
	return i.tgService.SendMessageToOwner(
		fmt.Sprintf("Доступ выдан для чата %v с лимитом %v", grantingChat, limit))
}

func parseGrantMessage(text string) (int64, int, error) {
	args := strings.Split(text, " ")
	if len(args) != 2 {
		return 0, 0, errors.New("expected 2 args instead of " + strconv.Itoa(len(args)))
	}
	grantingChatStr := args[0]
	limitStr := args[1]
	grantingChat, err := strconv.ParseInt(grantingChatStr, 10, 64)
	if err != nil {
		return 0, 0, err
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return 0, 0, err
	}
	return grantingChat, limit, nil
}

func (i *TgInteractor) StartHandleTgMessages() {
	i.tgService.StartReceiveMessages(i.acceptMessage)
}

func (i *TgInteractor) handleStartParserErr(settings *domain.ParserSettings, err error) error {
	if errors.Is(err, domain.LimitExceededErr) {
		return i.tgService.SendMessage(settings.ID, fmt.Sprintf(limitExceededMessageFormat, settings.Limit))
	}
	return nil
}
