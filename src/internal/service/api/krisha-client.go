package api

import (
	"encoding/json"
	"github.com/Logotipiwe/krisha_model/model"
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

var (
	TargetDomain              string
	TargetMapDataPath         string
	TargetPath                string
	PageSize                  int
	TargetMapDataFilterParams string
)

func initVariables() {
	TargetDomain = config.GetConfigOr("TARGET_HOST", "https://krisha.kz")
	TargetMapDataPath = config.GetConfigOr("TARGET_MAPDATA_PATH", "/a/ajax-map")
	TargetPath = config.GetConfigOr("TARGET_PATH", "/a/ajax-map-list")
	PageSize = 20
	TargetMapDataFilterParams = config.GetConfigOr("TARGET_MAPDATA_FILTER_PARAMS", "&lat=43.23814&lon=76.94297&zoom=13&precision=6&bounds=txwwjn%2Ctxwtzb")
}

type KrishaClientService struct {
	client    *http.Client
	tgService tg.TgServicer
}

func NewKrishaClientService(tgService tg.TgServicer) *KrishaClientService {
	s := &KrishaClientService{
		tgService: tgService,
		client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
	initVariables()
	return s
}

func (s *KrishaClientService) CollectAllPages(settings *domain.ParserSettings, stopped *bool) map[string]*model.Ap {
	//initVariables() //TODO Rewrites in every call, but in var or init blocks CS is not loaded - think where to put it
	data := s.RequestMapData(settings.Filters)
	requestUrl := TargetDomain + TargetPath + settings.Filters

	aps := make(map[string]*model.Ap)
	if data.NbTotal <= 0 {
		return aps
	}
	requestsCount := int(math.Ceil(float64(data.NbTotal) / float64(PageSize)))

	log.Println("Start collecting pages by url " + requestUrl)
	jobs := make([]func() map[string]*model.Ap, 0)
	for i := 0; i < requestsCount; i++ {
		num := i + 1
		jobs = append(jobs, func() map[string]*model.Ap {
			log.Printf("[%v. ] Requesting page %v...\n", strconv.FormatInt(settings.ID, 10), strconv.Itoa(num))
			return s.requestPage(requestUrl, num).Adverts
		})
	}

	workersNum, err := config.GetConfigInt("REQUEST_WORKERS_NUM")
	if err != nil {
		log.Println(err)
		workersNum = 1
	}
	pages := parallel.DoJobs(jobs, workersNum, stopped)
	if *stopped {
		return nil
	}
	for pageIndex, mapp := range pages {
		for id, ap := range mapp {
			if _, exists := aps[id]; exists {
				log.Printf("WARINIG! Ap %v from page %v already existed and rewritten", id, pageIndex+1)
			}
			aps[id] = ap
		}
	}
	log.Println("Collected  " + strconv.Itoa(len(aps)) + " aps")
	return aps
}

// TODO return err from jobs if occured
func (s *KrishaClientService) requestPage(url string, page int) model.ApsResult {

	req, _ := http.NewRequest("GET", url+"&page="+strconv.Itoa(page), nil)
	req.Header.Add("x-requested-with", "XMLHttpRequest")

	resp, err := s.client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	log.Println(resp.Status)

	var resultRaw model.ApsResultRaw

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &resultRaw)
	if err != nil {
		panic(err)
	}

	var aps map[string]*model.Ap
	if string(resultRaw.Adverts) != "[]" {
		err = json.Unmarshal(resultRaw.Adverts, &aps)
		if err != nil {
			panic(err)
		}
	} else {
		aps = make(map[string]*model.Ap)
	}

	log.Println("Found " + strconv.Itoa(len(aps)) + " aps")
	return resultRaw.ToResult(aps)
}

func (s *KrishaClientService) RequestMapData(filters string) *model.MapData {
	//initVariables() //TODO check todo above
	req, _ := http.NewRequest("GET", TargetDomain+TargetMapDataPath+filters+TargetMapDataFilterParams, nil)
	req.Header.Add("x-requested-with", "XMLHttpRequest")

	log.Println("Requesting map data.json...")
	resp, err := s.client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	log.Println(resp.Status)

	var result model.MapData

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		panic(err)
	}
	return &result
}
