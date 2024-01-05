package service

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

var Interval = 120 * time.Second
var Filters = "?das[price][from]=200000&das[price][to]=330000&areas=p43.220505%2C76.932670%2C43.228163%2C76.922885%2C43.237325%2C76.916018%2C43.242973%2C76.915847%2C43.254266%2C76.923056%2C43.257027%2C76.928721%2C43.259787%2C76.939021%2C43.259661%2C76.952411%2C43.256776%2C76.968890%2C43.253514%2C76.970092%2C43.248369%2C76.964599%2C43.243977%2C76.962710%2C43.228163%2C76.964084%2C43.222263%2C76.966144%2C43.217743%2C76.962710%2C43.216864%2C76.960135%2C43.216488%2C76.948977%2C43.218245%2C76.936618%2C43.221007%2C76.931468%2C43.220505%2C76.932670"
var Enabled = false

//var cache = make(map[string]interface{})

const (
	mapDataUrl string = "https://krisha.kz/a/ajax-map/map/arenda/kvartiry/almaty/"
	url        string = "https://krisha.kz/a/ajax-map-list/map/arenda/kvartiry/almaty/"
)

func StartParse(db *gorm.DB) {
	aps := make(map[string]interface{})
	first := true
	filters := Filters
	for {
		if Enabled {
			if filters != Filters {
				first = true
				aps = make(map[string]interface{})
				filters = Filters
			}
			startTime := time.Now()
			newAps := collectAllPages(url + Filters)
			elapsed := time.Since(startTime)
			log.Printf("collectAllPages took %s", elapsed)
			if !first {
				for id, apData := range newAps {
					_, has := aps[id]
					if !has {
						logNewAp(apData.(map[string]interface{}))
					}
				}
			} else {
				first = false
			}
			aps = newAps
			SendMessageInTg(fmt.Sprintf(
				"Collected aps: %s in %s. Next fetch after %s",
				strconv.Itoa(len(aps)), elapsed.String(), Interval.String()))
			time.Sleep(Interval)
		} else {
			log.Println("Parsing is disabled. Waiting for it to be enabled...")
			time.Sleep(time.Second * 2)
		}
	}
}

func logNewAp(data map[string]interface{}) {
	log.Println("=======================================================================")
	log.Println("NEW AP FOUND")
	log.Println("ID")
	log.Println(getId(data))
	link := "Link: https://krisha.kz/a/show/" + getId(data)
	log.Println(link)
	log.Println("=======================================================================")
	SendMessageInTg(link)
}

func getId(data map[string]interface{}) string {
	return strconv.FormatFloat(data["id"].(float64), 'f', -1, 64)
}

func collectAllPages(url string) map[string]interface{} {
	hasMore := true
	var aps = make(map[string]interface{})
	page := 1

	log.Println("Start collecting pages by url " + url)
	for hasMore {
		moreAps := requestPage(url, page)
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

func requestPage(url string, page int) map[string]interface{} {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, _ := http.NewRequest("GET", url+"&page="+strconv.Itoa(page), nil)
	req.Header.Add("x-requested-with", "XMLHttpRequest")

	log.Println("Requesting page " + strconv.Itoa(page) + "...")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	log.Println(resp.Status)

	result := make(map[string]interface{})

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		panic(err)
	}
	aps := result["adverts"]
	apsMap, ok := aps.(map[string]interface{})
	if !ok {
		_, empty := aps.([]interface{})
		if empty {
			return make(map[string]interface{})
		}
	}
	log.Println("Found " + strconv.Itoa(len(apsMap)) + " aps")
	return apsMap
}

func requestMapData(url string) map[string]interface{} {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("x-requested-with", "XMLHttpRequest")

	log.Println("Requesting map data...")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	log.Println(resp.Status)

	result := make(map[string]interface{})

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
