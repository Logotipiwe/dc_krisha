package tghttp

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	config "github.com/logotipiwe/dc_go_config_lib"
	"krisha/src/internal"
	"krisha/src/internal/domain"
	"krisha/src/internal/service/admin"
	"krisha/src/internal/service/parser"
	"krisha/src/internal/service/tg"
	"krisha/src/pkg"
	"runtime/debug"
	"strconv"
	"strings"
)

type TgInteractor struct {
	tgService          tg.TgServicer
	parserService      *parser.Service
	permissionsService *internal.PermissionsService
	ownerChatMode      OwnerChatMode
	adminService       *admin.Service
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
/info - сводка
/grant - выдать доступ
/deny - забрать доступ
/grant with 0 - забрать доступ при включенном авто лимите`
	userHelp = `Бот работает следующим образом:
1. Вы отправляете фильтры, по которым ищете квартиру. Инструкция - /filterHelp
2. Бот присылает вам уведомления о новых квартирах
3. Вы можете писать /stop или /start, чтобы отключать или включать обратно уведомления
4. Вы можете отправить другой фильтр, чтобы бот перенастроился на него`
	ownerUnacceptedError = "unknown admin command"
	errorMessage         = "Произошла ошибка :( \r\n Попробуйте начать заново с /start"
	StartAnswer          = `Привет! Это бот для получения уведомлений о новых квартирах по фильтрам. Для начала работы нужно установить нужный фильтр.
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
	noAccessMessage            = "У вас нет доступа к боту. Обратитесь к администратору"
	limitExceededMessageFormat = "Превышен лимит в %v квартир в вашем фильтре. Попробуйте другой фильтр"
)

func NewTgInteractor(
	tgService tg.TgServicer,
	parserService *parser.Service,
	permissionsService *internal.PermissionsService,
	adminService *admin.Service,
) *TgInteractor {
	return &TgInteractor{
		tgService:          tgService,
		parserService:      parserService,
		permissionsService: permissionsService,
		ownerChatMode:      defaultMode,
		adminService:       adminService,
	}
}

func (i *TgInteractor) AcceptMessage(update tgbotapi.Update) error {
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

//TODO при авто гранте ни deny ни allow 0  не запрещают запуск
//TODO автостоп не работает
//TODO сделать /info в чате для получения chat_id ???

func (i *TgInteractor) acceptUserMessage(update tgbotapi.Update) error {
	text := update.Message.Text
	chatID := update.Message.Chat.ID

	err := i.createParserSettingsFromAutoGrantIfNeeded(chatID)
	if err != nil {
		return err
	}

	settings, err := i.parserService.GetSettings(chatID)
	if err != nil {
		return err
	}
	hasAccess := settings != nil && i.permissionsService.HasAccess(settings)
	isOwner := chatID == pkg.GetOwnerChatID()
	if !hasAccess && !isOwner {
		return i.tgService.SendMessage(chatID, noAccessMessage)
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
	settings, err := i.parserService.GetSettings(chatID)
	if err != nil {
		return err
	}
	if settings == nil || settings.Filters == "" {
		return i.tgService.SendMessage(chatID, StartAnswer)
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
	switch {
	case text == "/help":
		return i.tgService.SendMessageToOwner(ownerHelp + "\r\n\r\n" + userHelp)
	case text == "/info":
		info, err := i.adminService.GetGeneralInfo()
		if err != nil {
			fmt.Println(err)
			return i.tgService.SendLogMessageToOwner("Ошибка получения информации: " + err.Error())
		}
		message := formatAdminInfo(info)
		return i.tgService.SendMessageToOwner(message)
	case text == "/grant":
		err := i.tgService.SendMessageToOwner("Какому чату выдать доступ? И через пробел - лимит")
		if err == nil {
			i.ownerChatMode = granting
		}
		return err
	case ownerChatMode == granting:
		err := i.handleGrantCommand(text)
		if err == nil {
			i.ownerChatMode = defaultMode
		}
		return err
	case text == "/deny":
		err := i.tgService.SendMessageToOwner("Какому чату запретить доступ?")
		if err == nil {
			i.ownerChatMode = denying
		}
		return err
	case ownerChatMode == denying:
		err := i.handleDenyCommand(text)
		if err != nil {
			debug.PrintStack()
		} else {
			i.ownerChatMode = defaultMode
		}
		return err
	}
	return errors.New(ownerUnacceptedError)
}

func formatAdminInfo(info *domain.AdminInfo) string {
	ans := "Сводка о работе парсера:\n"
	ans += "Стандартный интервал: " + strconv.Itoa(info.DefaultInterval) + "сек\n"
	ans += "Авто лимит: " + strconv.Itoa(info.AutoGrantLimit) + "\n\n"
	if len(info.ActiveParsers) > 0 {
		ans += "Активные парсеры:\n"
		for _, settings := range info.ActiveParsers {
			ans += strconv.FormatInt(settings.ID, 10) +
				` - interval: ` + strconv.Itoa(settings.IntervalSec) +
				`, aps: ` + strconv.Itoa(0) + ` (TODO)` +
				`, explicit: ` + strconv.FormatBool(settings.IsGrantedExplicitly)
			if settings.IsGrantedExplicitly {
				ans += `, limit: ` + strconv.Itoa(settings.Limit)
			}
			ans += "\n"
		}
	} else {
		ans += "Нет активных парсеров"
	}
	return ans
}

func (i *TgInteractor) handleDenyCommand(text string) error {
	denyingChat, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		return i.tgService.SendMessageToOwner("Кажется это не число, попробуй ещё раз")
	}
	stopParserErr := i.parserService.StopParser(denyingChat)
	settings, err := i.parserService.GetSettings(denyingChat)
	if settings == nil && pkg.GetAutoGrantLimit() > 0 {
		fmt.Println("Creating parser settings from deny...")
		settings, err = i.parserService.CreateFromDenyCommand(denyingChat)
	} else {
		if !i.permissionsService.HasAccess(settings) {
			return i.tgService.SendMessageToOwner("У этого чата и так нет доступа. Спасибо")
		}
		err = i.parserService.UpdateFromDenyCommand(settings)
	}
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
		return i.tgService.SendMessageToOwner("Ошибка " + err.Error())
	}
	settings, err := i.parserService.GetSettings(grantingChat)
	if err != nil {
		return err
	}
	if settings == nil {
		err = i.parserService.CreateParserSettingsFromExplicitGrant(grantingChat, limit)
		if err != nil {
			return err
		}
	} else {
		err, stopped := i.parserService.UpdateLimitExplicitly(settings, limit)
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
	i.tgService.StartReceiveMessages(i.AcceptMessage)
}

func (i *TgInteractor) handleStartParserErr(settings *domain.ParserSettings, err error) error {
	if errors.Is(err, domain.LimitExceededErr) {
		var limit int
		if settings.IsGrantedExplicitly {
			limit = settings.Limit
		} else {
			limit = pkg.GetAutoGrantLimit()
		}
		return i.tgService.SendMessage(settings.ID, fmt.Sprintf(limitExceededMessageFormat, limit))
	}
	return nil
}

func (i *TgInteractor) createParserSettingsFromAutoGrantIfNeeded(chatID int64) error {
	settings, err := i.parserService.GetSettings(chatID)
	autoGrantLimit, _ := config.GetConfigInt("AUTO_GRANT_LIMIT")
	if err != nil {
		return err
	}

	if settings == nil && autoGrantLimit > 0 {
		err = i.parserService.CreateParserSettingsFromAutoGrant(chatID)
		if err != nil {
			return err
		}
	}
	return nil
}
