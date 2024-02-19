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
}

func newParser(settings *domain.ParserSettings, apsInFilter int, factory *Factory) *Parser {
	return &Parser{
		factory:                    factory,
		settings:                   settings,
		areAllCurrentApsCollected:  false,
		areCollectApsTriesExceeded: false,
		enabled:                    true,
		collectedAps:               make(map[string]*model.Ap),
		stopped:                    false,
		initialApsCountInFilter:    apsInFilter,
	}
}

func (p *Parser) startParsing() error {
	p.enabled = true
	go func() {
		p.initParsing()
		p.doParseForCollectAps()
		for p.enabled {
			time.Sleep(time.Duration(p.settings.IntervalSec) * time.Second)
			p.doParseWithNotification()
		}
	}()
	return nil
}

func (p *Parser) initParsing() {
	log.Println("Parse for chat " + strconv.FormatInt(p.settings.ID, 10))
	data := p.factory.krishaClient.RequestMapData(p.settings.Filters)
	p.factory.tgService.SendMessage(p.settings.ID, "Квартир: "+strconv.Itoa(data.NbTotal))
	p.factory.tgService.SendLogMessageToOwner(fmt.Sprintf(
		"Parser started for chat %v. filter %v. Interval: %v", p.settings.ID, p.settings.Filters, p.settings.IntervalSec))
}

func (p *Parser) doParseWithNotification() {
	aps := p.factory.krishaClient.CollectAllPages(p.settings.Filters, p.settings.ID, &p.stopped)
	if !p.enabled {
		return
	}
	for id, ap := range aps {
		_, has := p.collectedAps[id]
		if !has {
			photosUrls := pkg.Map(ap.Photos, func(p *model.Photo) string {
				return p.Src
			})
			p.factory.tgService.SendImgMessage(p.settings.ID, "Новая квартира: "+ap.GetLink(), photosUrls[0:pkg.Min(len(photosUrls), 10)])
			go p.factory.tgService.SendLogMessageToOwner(fmt.Sprintf("У чата %v квартира %v", p.settings.ID, ap.GetLink()))
		}
		p.collectedAps[id] = ap
	}
}

func (p *Parser) doParseForCollectAps() {
	p.factory.tgService.SendMessage(p.settings.ID, "Начинаю собирать существующие квартиры, это займет немного времени...")
	attempts := 0
	for !p.areAllCurrentApsCollected && !p.areCollectApsTriesExceeded {
		aps := p.factory.krishaClient.CollectAllPages(p.settings.Filters, p.settings.ID, &p.stopped)
		if p.stopped {
			return
		}
		for id, ap := range aps {
			p.collectedAps[id] = ap
		}
		attempts++

		if len(p.collectedAps) >= p.initialApsCountInFilter {
			p.areAllCurrentApsCollected = true
			p.factory.tgService.SendMessage(p.settings.ID, "Существующие квартиры собраны, начинаю присылать уведомления о новых...")
		}
		if attempts > 5 {
			p.areCollectApsTriesExceeded = true
			p.factory.tgService.SendMessage(p.settings.ID, "Существующие квартиры собраны, но из-за большого их количества в фильтре - могут иногда присылаться уведомления не по новым квартирам, а по уже существующим")
		}
	}
}

func (p *Parser) disable() {
	p.enabled = false
	p.stopped = true
}
