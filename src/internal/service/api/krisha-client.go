package api

import (
	"encoding/json"
	config "github.com/logotipiwe/dc_go_config_lib"
	"io"
	"krisha/src/internal/domain"
	"krisha/src/internal/service/parallel"
	"krisha/src/internal/service/tg"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

const (
	mapDataUrl string = "https://krisha.kz/a/ajax-map/map/arenda/kvartiry/almaty/"
	url        string = "https://krisha.kz/a/ajax-map-list/map/arenda/kvartiry/almaty/"
	pageSize   int    = 20
)

type KrishaClientService struct {
	client    *http.Client
	tgService *tg.TgService
}

func NewKrishaClientService(tgService *tg.TgService) *KrishaClientService {
	return &KrishaClientService{
		tgService: tgService,
		client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (s *KrishaClientService) CollectAllPages(filters string) map[string]*domain.Ap {
	data := s.requestMapData(mapDataUrl + filters + "&lat=43.23814&lon=76.94297&zoom=13&precision=6&bounds=txwwjn%2Ctxwtzb")
	_ = s.tgService.SendLogMessageInTg("Collecting " + strconv.Itoa(data.NbTotal) + " aps...")
	requestUrl := url + filters

	var aps = make(map[string]*domain.Ap)
	if data.NbTotal <= 0 {
		return aps
	}
	requestsCount := int(math.Ceil(float64(data.NbTotal) / float64(pageSize)))

	log.Println("Start collecting pages by url " + requestUrl)
	jobs := make([]func() map[string]*domain.Ap, 0)
	for i := 0; i < requestsCount; i++ {
		num := i
		jobs = append(jobs, func() map[string]*domain.Ap {
			println("!!!!!!!!!!!!!! REQUESTING PAGE !!!!!!!!!!!!!!!!!!")
			return s.requestPage(requestUrl, num).Adverts
			//println(num)
			//return map[string]*domain.Ap{}
		})
	}

	workersNum, err := config.GetConfigInt("REQUEST_WORKERS_NUM")
	if err != nil {
		log.Println(err)
		workersNum = 1
	}
	pages := parallel.DoJobs(jobs, workersNum)

	for _, mapp := range pages {
		for id, ap := range mapp {
			if _, exists := aps[id]; exists {
				log.Println("WARINIG! Ap " + id + " already existed and rewritten")
			}
			aps[id] = ap
		}
	}
	/*page := 1

	log.Println("Start collecting pages by url " + requestUrl)
	for hasMore {
		moreAps := s.requestPage(requestUrl, page).Adverts
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
	}*/
	log.Println("Collected  " + strconv.Itoa(len(aps)) + " aps")
	return aps
}

func (s *KrishaClientService) requestPage(url string, page int) domain.ApsResult {

	req, _ := http.NewRequest("GET", url+"&page="+strconv.Itoa(page), nil)
	req.Header.Add("x-requested-with", "XMLHttpRequest")

	log.Println("Requesting page " + strconv.Itoa(page) + "...")
	resp, err := s.client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	log.Println(resp.Status)

	var resultRaw domain.ApsResultRaw

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &resultRaw)
	if err != nil {
		panic(err)
	}

	var aps map[string]*domain.Ap
	if string(resultRaw.Adverts) != "[]" {
		err = json.Unmarshal(resultRaw.Adverts, &aps)
		if err != nil {
			panic(err)
		}
	} else {
		aps = make(map[string]*domain.Ap)
	}

	log.Println("Found " + strconv.Itoa(len(aps)) + " aps")
	return resultRaw.ToResult(aps)
}

func (s *KrishaClientService) requestMapData(url string) domain.MapData {

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("x-requested-with", "XMLHttpRequest")

	log.Println("Requesting map data.json...")
	resp, err := s.client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	log.Println(resp.Status)

	var result domain.MapData

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		panic(err)
	}
	return result
}
