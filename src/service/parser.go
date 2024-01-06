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
var Filters = "?das[_sys.hasphoto]=1&das[live.rooms][]=2&das[live.rooms][]=3&das[live.square][from]=30&das[live.square][to]=80&das[price][from]=200000&das[price][to]=330000&das[who]=1&lat=43.23814&lon=76.94297&zoom=13&precision=6&bounds=txwwjq%2Ctxwtz8&areas=p43.219849%2C76.932000%2C43.225373%2C76.925477%2C43.227256%2C76.916208%2C43.238928%2C76.916208%2C43.247588%2C76.914834%2C43.255493%2C76.921357%2C43.264338%2C76.932859%2C43.269167%2C76.940240%2C43.268352%2C76.961269%2C43.258629%2C76.963243%2C43.248215%2C76.967706%2C43.239305%2C76.966848%2C43.233281%2C76.968049%2C43.221732%2C76.971998%2C43.215706%2C76.942643%2C43.220351%2C76.930455%2C43.219849%2C76.932000"
var Enabled = false

var cache = make(map[string]*CacheItem)

type CacheItem struct {
	Count int
	Data  interface{}
}

func NewCacheItem() *CacheItem {
	return &CacheItem{0, make(map[string]interface{})}
}

const (
	mapDataUrl string = "https://krisha.kz/a/ajax-map/map/arenda/kvartiry/almaty/"
	url        string = "https://krisha.kz/a/ajax-map-list/map/arenda/kvartiry/almaty/"
)

func StartParse(db *gorm.DB) {
	aps := make(map[string]*Ap)
	first := true
	stopLogged := false
	filters := Filters
	for {
		if Enabled {
			stopLogged = false
			if filters != Filters {
				first = true
				aps = make(map[string]*Ap)
				filters = Filters
			}
			data := requestMapData(mapDataUrl + Filters + "&lat=43.23814&lon=76.94297&zoom=13&precision=6&bounds=txwwjn%2Ctxwtzb")
			SendLogInTg("Collecting " + strconv.Itoa(data.NbTotal) + " aps...")
			startTime := time.Now()
			newAps := collectAllPages(url + Filters)
			elapsed := time.Since(startTime)
			log.Printf("collectAllPages took %s", elapsed)
			if !first {
				for id, ap := range newAps {
					_, has := cache[id]
					if !has {
						logNewAp(ap)
					}
				}
				for id, ap := range aps {
					_, has := newAps[id]
					if !has {
						logMissingAp(ap)
					}
				}
			} else {
				first = false
			}
			aps = newAps
			addToCache(newAps)
			SendLogInTg(fmt.Sprintf(
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

func logMissingAp(m *Ap) {
	log.Println(fmt.Sprintf("Missing ap %s", strconv.FormatInt(m.ID, 10)))
	log.Println(m)
}

func addToCache(aps map[string]*Ap) {
	for id, val := range aps {
		cacheItem := cache[id]
		if cacheItem == nil {
			cacheItem = NewCacheItem()
		}
		cacheItem.Count++
		cacheItem.Data = val
		cache[id] = cacheItem
	}
}

func logNewAp(data *Ap) {
	log.Println("=======================================================================")
	log.Println("NEW AP FOUND")
	log.Println("ID")
	log.Println(data.ID)
	link := "https://krisha.kz/a/show/" + strconv.FormatInt(data.ID, 10)
	log.Println(link)
	log.Println("=======================================================================")
	var imagesUrls = make([]string, 0)
	photos := data.Photos
	message := fmt.Sprintf("Link: %s\r\n", link)
	if len(photos) > 0 {
		for _, photo := range photos {
			imagesUrls = append(imagesUrls, photo.Src)
		}
		if len(imagesUrls) > 10 {
			imagesUrls = imagesUrls[:10]
		}
		SendMessageInTgWithImages(message, imagesUrls)
	}
	SendMessageInTg(message)
}

func getId(data map[string]interface{}) string {
	return strconv.FormatFloat(data["id"].(float64), 'f', -1, 64)
}

func collectAllPages(url string) map[string]*Ap {
	hasMore := true
	var aps = make(map[string]*Ap)
	page := 1

	log.Println("Start collecting pages by url " + url)
	for hasMore {
		moreAps := requestPage(url, page).Adverts
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

func requestPage(url string, page int) ApsResult {
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

	var resultRaw ApsResultRaw

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &resultRaw)
	if err != nil {
		panic(err)
	}

	var aps map[string]*Ap
	if string(resultRaw.Adverts) != "[]" {
		err = json.Unmarshal(resultRaw.Adverts, &aps)
		if err != nil {
			panic(err)
		}
	} else {
		aps = make(map[string]*Ap)
	}

	log.Println("Found " + strconv.Itoa(len(aps)) + " aps")
	return resultRaw.toResult(aps)
}

func requestMapData(url string) MapData {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("x-requested-with", "XMLHttpRequest")

	log.Println("Requesting map data.json...")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	log.Println(resp.Status)

	var result MapData

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
