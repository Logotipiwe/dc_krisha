package service

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"krisha/src/internal/domain"
	"krisha/src/internal/service/apartments"
	"krisha/src/internal/service/api"
	"krisha/src/internal/service/tg"
	"log"
	"strconv"
	"time"
)

var Interval = 120 * time.Second
var Filters = "?das[_sys.hasphoto]=1&das[live.rooms][]=2&das[live.rooms][]=3&das[live.square][from]=30&das[live.square][to]=80&das[price][from]=200000&das[price][to]=330000&das[who]=1&lat=43.23814&lon=76.94297&zoom=13&precision=6&bounds=txwwjq%2Ctxwtz8&areas=p43.219849%2C76.932000%2C43.225373%2C76.925477%2C43.227256%2C76.916208%2C43.238928%2C76.916208%2C43.247588%2C76.914834%2C43.255493%2C76.921357%2C43.264338%2C76.932859%2C43.269167%2C76.940240%2C43.268352%2C76.961269%2C43.258629%2C76.963243%2C43.248215%2C76.967706%2C43.239305%2C76.966848%2C43.233281%2C76.968049%2C43.221732%2C76.971998%2C43.215706%2C76.942643%2C43.220351%2C76.930455%2C43.219849%2C76.932000"
var Enabled = false

const (
	mapDataUrl string = "https://krisha.kz/a/ajax-map/map/arenda/kvartiry/almaty/"
	url        string = "https://krisha.kz/a/ajax-map-list/map/arenda/kvartiry/almaty/"
)

type ParserService struct {
	KrishaClientService *api.KrishaClientService
	ApsCacheService     *apartments.ApsCacheService
	ApsLoggerService    *apartments.ApsLoggerService
	ApsTgSenderService  *apartments.ApsTgSenderService
	TgService           *tg.TgService
	db                  *gorm.DB
}

func NewParserService(
	krishaClientService *api.KrishaClientService,
	apsCacheService *apartments.ApsCacheService,
	apsTgSender *apartments.ApsTgSenderService,
	apsLoggerService *apartments.ApsLoggerService,
	tgService *tg.TgService,
	db *gorm.DB) *ParserService {
	return &ParserService{TgService: tgService, KrishaClientService: krishaClientService, ApsCacheService: apsCacheService,
		ApsTgSenderService: apsTgSender, ApsLoggerService: apsLoggerService, db: db}
}

func (s *ParserService) StartParse() {
	aps := make(map[string]*domain.Ap)
	first := true
	stopLogged := false
	filters := Filters
	for {
		if Enabled {
			stopLogged = false
			if filters != Filters {
				first = true
				aps = make(map[string]*domain.Ap)
				filters = Filters
			}
			data := s.KrishaClientService.RequestMapData(mapDataUrl + Filters + "&lat=43.23814&lon=76.94297&zoom=13&precision=6&bounds=txwwjn%2Ctxwtzb")
			_ = s.TgService.SendLogMessageInTg("Collecting " + strconv.Itoa(data.NbTotal) + " aps...")
			startTime := time.Now()
			newAps := s.CollectAllPages(url + Filters)
			elapsed := time.Since(startTime)
			log.Printf("collectAllPages took %s", elapsed)
			if !first {
				for id, ap := range newAps {
					has := s.ApsCacheService.IsInCache(id)
					if !has {
						s.ApsLoggerService.LogNewAp(ap)
						err := s.ApsTgSenderService.LogInTg(ap)
						if err != nil {
							log.Println(err)
						}

					}
				}
				for id, ap := range aps {
					_, has := newAps[id]
					if !has {
						s.ApsLoggerService.LogMissingAp(ap)
					}
				}
			} else {
				first = false
			}
			aps = newAps
			s.ApsCacheService.AddToCache(newAps)
			_ = s.TgService.SendLogMessageInTg(fmt.Sprintf(
				"Collected aps: %s in %s. Next fetch after %s",
				strconv.Itoa(len(aps)), elapsed.String(), Interval.String()))
			var sleeped float64 = 0
			for sleeped < Interval.Seconds() {
				time.Sleep(time.Second)
				sleeped += 1
			}
		} else {
			first = true
			if !stopLogged {
				log.Println("Parsing is disabled. Waiting for it to be enabled...")
				stopLogged = true
			}
			time.Sleep(time.Second * 2)
		}
	}
}

func (s *ParserService) CollectAllPages(url string) map[string]*domain.Ap {
	hasMore := true
	var aps = make(map[string]*domain.Ap)
	page := 1

	log.Println("Start collecting pages by url " + url)
	for hasMore {
		moreAps := s.KrishaClientService.RequestPage(url, page).Adverts
		if len(moreAps) > 0 {
			for s, i := range moreAps {
				if _, exists := aps[s]; exists {
					log.Println("WARINIG! Ap " + s + " already existed and rewritten")
				}
				aps[s] = i
			}
		} else {
			hasMore = false
		}
		page++
	}
	log.Println("Collected  " + strconv.Itoa(len(aps)) + " aps")
	return aps
}
