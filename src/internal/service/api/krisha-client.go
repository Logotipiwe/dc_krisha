package api

import (
	"encoding/json"
	"io"
	"krisha/src/internal/domain"
	"log"
	"net/http"
	"strconv"
	"time"
)

type KrishaClientService struct {
	client *http.Client
}

func NewKrishaClientService() *KrishaClientService {
	return &KrishaClientService{
		client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (s *KrishaClientService) RequestPage(url string, page int) domain.ApsResult {

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

func (s *KrishaClientService) RequestMapData(url string) domain.MapData {

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
