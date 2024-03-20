package parser

import (
	"fmt"
	"github.com/Logotipiwe/krisha_model/model"
	"krisha/src/internal/domain"
	"krisha/src/pkg"
	"log"
	"strconv"
	"time"
)

type Parser struct {
	factory                    *Factory
	settings                   *domain.ParserSettings
	enabled                    bool
	areAllCurrentApsCollected  bool
	areCollectApsTriesExceeded bool
	collectedAps               map[string]*model.Ap
	stopped                    bool
	initialApsCountInFilter    int
	startTime                  time.Time
}

func (p *Parser) startParsing(shouldNotifyWhenStart bool) error {
	p.enabled = true
	fmt.Println("Started parser for chat " + strconv.FormatInt(p.settings.ID, 10))
	go func() {
		p.initParsing(shouldNotifyWhenStart)
		p.doParseForCollectAps(shouldNotifyWhenStart)
		for p.enabled {
			p.sleepForInterval() //TODO cover with tests
			p.doParseWithNotification()
			p.checkForExpiredAutoStopTime() //TODO cover with tests
		}
	}()
	return nil
}

func (p *Parser) initParsing(shouldNotify bool) {
	log.Println("Parse for chat " + strconv.FormatInt(p.settings.ID, 10))
	data := p.getMapData()
	if shouldNotify {
		p.factory.tgService.SendMessage(p.settings.ID, "Квартир: "+strconv.Itoa(data.NbTotal))
	}
	p.factory.tgService.SendLogMessageToOwner(fmt.Sprintf(
		"Parser started for chat %v. filter %v. Interval: %v", p.settings.ID, p.settings.Filters, p.settings.IntervalSec))
}

func (p *Parser) doParseWithNotification() {
	if !p.enabled {
		return
	}
	aps := p.factory.krishaClient.CollectAllPages(p.settings, &p.stopped)
	for id, ap := range aps {
		_, has := p.collectedAps[id]
		if !has {
			photosUrls := pkg.Map(ap.Photos, func(p *model.Photo) string {
				return p.Src
			})
			p.factory.tgService.SendImgMessage(p.settings.ID, "Новое объявление: "+ap.GetLink(), photosUrls[0:pkg.Min(len(photosUrls), 10)])
			go p.factory.tgService.SendLogMessageToOwner(fmt.Sprintf("У чата %v квартира %v", p.settings.ID, ap.GetLink()))
		}
		p.collectedAps[id] = ap
	}
	apsCount := len(aps)
	p.updateApsCount(apsCount)
}

func (p *Parser) doParseForCollectAps(shouldNotiy bool) {
	if shouldNotiy {
		p.factory.tgService.SendMessage(p.settings.ID, "Начинаю собирать существующие квартиры, это займет немного времени...")
	}
	attempts := 0
	for !p.areAllCurrentApsCollected && !p.areCollectApsTriesExceeded {
		aps := p.factory.krishaClient.CollectAllPages(p.settings, &p.stopped)
		if p.stopped {
			return
		}
		for id, ap := range aps {
			p.collectedAps[id] = ap
		}
		attempts++

		if len(p.collectedAps) >= p.initialApsCountInFilter {
			p.areAllCurrentApsCollected = true
			if shouldNotiy {
				p.factory.tgService.SendMessage(p.settings.ID, "Существующие квартиры собраны, начинаю присылать уведомления о новых...")
			}
		}
		if attempts > 5 {
			p.areCollectApsTriesExceeded = true
			if shouldNotiy {
				p.factory.tgService.SendMessage(p.settings.ID, "Существующие квартиры собраны, но из-за большого их количества в фильтре - могут иногда присылаться уведомления не по новым квартирам, а по уже существующим")
			}
		}
	}
}

func (p *Parser) getMapData() *model.MapData {
	return p.factory.krishaClient.RequestMapData(p.settings.Filters)
}

func (p *Parser) disable() {
	p.enabled = false
	p.stopped = true
}

func (p *Parser) updateApsCount(apsCount int) {
	p.settings.ApsCount = apsCount
	err := p.factory.parserSettingsRepo.Update(p.settings)
	if err != nil {
		fmt.Println("[ERROR] err updating curr aps count " + err.Error())
	}
}

func (p *Parser) sleepForInterval() {
	pkg.SleepWithInterruption(func() time.Duration {
		return pkg.GetParserSleepingInterval(p.settings)
	}, 10*time.Second)
}

func (p *Parser) checkForExpiredAutoStopTime() {
	stopHours := pkg.GetAutoStopHours()
	if stopHours == 0 {
		return
	}
	sinceStart := time.Since(p.startTime)
	if sinceStart > (time.Duration(stopHours) * time.Hour) {
		err := p.factory.parserService.StopParser(p.settings.ID)
		if err == nil {
			p.factory.tgService.SendMessage(p.settings.ID, "Парсер автоматически остановлен по истечению таймера. Если вы успешно нашли жилье - просим отправить фидбек администрации :) \n\n Если вы ещё не нашли нужный вариант - запустите парсер заново командой /start")
		}
	}
}
