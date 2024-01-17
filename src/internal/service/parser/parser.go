package parser

import (
	"fmt"
	"krisha/src/internal/domain"
	"krisha/src/pkg"
	"log"
	"strconv"
	"time"
)

//var Interval = 120 * time.Second
//var Filters = "?das[_sys.hasphoto]=1&das[live.rooms][]=2&das[live.rooms][]=3&das[live.square][from]=30&das[live.square][to]=80&das[price][from]=200000&das[price][to]=330000&das[who]=1&lat=43.23814&lon=76.94297&zoom=13&precision=6&bounds=txwwjq%2Ctxwtz8&areas=p43.219849%2C76.932000%2C43.225373%2C76.925477%2C43.227256%2C76.916208%2C43.238928%2C76.916208%2C43.247588%2C76.914834%2C43.255493%2C76.921357%2C43.264338%2C76.932859%2C43.269167%2C76.940240%2C43.268352%2C76.961269%2C43.258629%2C76.963243%2C43.248215%2C76.967706%2C43.239305%2C76.966848%2C43.233281%2C76.968049%2C43.221732%2C76.971998%2C43.215706%2C76.942643%2C43.220351%2C76.930455%2C43.219849%2C76.932000"
//var Enabled = false

type Parser struct {
	factory                    *Factory
	settings                   *domain.ParserSettings
	enabled                    bool
	areAllCurrentApsCollected  bool
	areCollectApsTriesExceeded bool
	collectedAps               map[string]*domain.Ap
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
		collectedAps:               make(map[string]*domain.Ap),
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
			p.doParseWithNotification()
			time.Sleep(time.Duration(p.settings.IntervalSec) * time.Second)
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
	aps := p.factory.krishaClient.CollectAllPages(p.settings.Filters, &p.stopped)
	if !p.enabled {
		return
	}
	for id, ap := range aps {
		_, has := p.collectedAps[id]
		if !has {
			photosUrls := pkg.Map(ap.Photos, func(p *domain.Photo) string {
				return p.Src
			})
			p.factory.tgService.SendImgMessage(p.settings.ID, "Новая квартира: "+ap.GetLink(), photosUrls[0:pkg.Min(len(photosUrls), 10)])
		}
		p.collectedAps[id] = ap
	}
}

func (p *Parser) doParseForCollectAps() {
	p.factory.tgService.SendMessage(p.settings.ID, "Начинаю собирать существующие квартиры, это займет немного времени...")
	attempts := 0
	for !p.areAllCurrentApsCollected && !p.areCollectApsTriesExceeded {
		aps := p.factory.krishaClient.CollectAllPages(p.settings.Filters, &p.stopped)
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

func (p *Parser) StartParse(filters string) {
	return
	//p.settings.Filters = filters
	//p.settings.Enabled = true
	//p.parserSettingsRepo.Update(&p.settings)

	//first := true
	//go func() {
	//	for p.settings.Enabled {
	//		if first {
	//p.doFirstParse() //TODO пробовать начинать писать с верхнего уровня абстракции
	//first = false
	//continue
	//} else {
	//	if p.areAllCurrentApsCollected {
	//		p.doParseWithNotification()
	//	} else {
	//		p.doCollectAps()
	//	}
	//}
	//time.Sleep(time.Duration(p.settings.IntervalSec) * time.Second)
	//}
	//}()

	//return
	//
	//aps := make(map[string]*domain.Ap)
	//for p.settings.Enabled {
	//	startTime := time.Now()
	//	newAps := p.KrishaClientService.CollectAllPages(p.settings.Filters)
	//	elapsed := time.Since(startTime)
	//	log.Printf("collectAllPages took %s", elapsed)
	//	if p.areAllCurrentApsCollected {
	//		for id, ap := range newAps {
	//			has := p.ApsCacheService.IsInCache(id)
	//			if !has {
	//				p.ApsLoggerService.LogNewAp(ap)
	//				err := p.ApsTgSenderService.LogInTg(ap)
	//				if err != nil {
	//					log.Println(err)
	//				}
	//			}
	//		}
	//		for id, ap := range aps {
	//			_, has := newAps[id]
	//			if !has {
	//				p.ApsLoggerService.LogMissingAp(ap)
	//			}
	//		}
	//	}
	//aps = newAps
	//p.ApsCacheService.AddToCache(newAps)
	//_ = p.tgService.SendLogMessageToOwner(fmt.Sprintf(
	//	"Collected aps: %p in %p. Next fetch after %p",
	//	strconv.Itoa(len(aps)), elapsed.String(), p.interval.String()))
	//var sleeped float64 = 0
	//for sleeped < p.interval.Seconds() {
	//	time.Sleep(time.Second)
	//	sleeped += 1
	//}
	//}
	//first = true
	//log.Println("Parsing is disabled. Waiting for it to be enabled...")
}

//func (p *Parser) resetFilters(newFilters string) {
//
//}
