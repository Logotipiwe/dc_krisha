package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	url := "https://krisha.kz/a/ajax-map-list/map/arenda/kvartiry/almaty/?das[live.rooms][0]=1&das[live.rooms][1]=2&das[price][from]=200000&das[price][to]=300000&das[_sys.hasphoto]=1&das[who]=1&zoom=11&lat=43.28590&lon=76.91290&bounds=43.43820082112147%2C76.69214346923823%2C43.12315556729036%2C77.1336565307617&page=100"
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, _ := http.NewRequest("GET", url, nil)
	//req.Header.Add("authority", "krisha.kz")
	//req.Header.Add("accept", "application/json, text/javascript, */*; q=0.01")
	//req.Header.Add("accept-language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7,und;q=0.6")
	//req.Header.Add("cache-control", "no-cache")
	//req.Header.Add("cookie", "krishauid=64e81b43bf82804956354f4e3bc96ac789c3323f; _ym_uid=168948411715301724; _ym_d=1689484117; ssaid=d19a9720-2396-11ee-91da-5f487a928ae4; _gcl_au=1.1.1819964744.1703768801; _tt_enable_cookie=1; _ttp=5B74kymHQOHb7YLkgALPAYL2m60; __gsas=ID=cded80d0439c45a4:T=1703768831:RT=1703768831:S=ALNI_MYZQcDb-vQLIr6P4aW0y4ez2HtqZQ; tutorial=%7B%22add-note%22%3A%22viewed%22%2C%22advPage%22%3A%22viewed%22%7D; _gid=GA1.2.2115420768.1704284700; re2Qqt9pLYYRJ8Br=1; _ym_isad=1; krssid=96jlcr9hmhi5np4sce86decon6; kr_cdn_host=//alakt-kz.kcdn.online; _ym_visorc=w; ksq_region=2; hist_region=2; __tld__=null; _ga=GA1.2.237972167.1703768799; _ga_6YZLS7YDS7=GS1.1.1704435475.18.1.1704437759.60.0.0")
	//req.Header.Add("pragma", "no-cache")
	//req.Header.Add("referer", "https://krisha.kz/map/arenda/kvartiry/almaty/?das[live.rooms][0]=1&das[live.rooms][1]=2&das[price][from]=200000&das[price][to]=300000&das[_sys.hasphoto]=1&das[who]=1&lat=43.28590&lon=76.91290&zoom=11&sidebarPage=5")
	//req.Header.Add("sec-ch-ua", `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`)
	//req.Header.Add("sec-ch-ua-mobile", "?0")
	//req.Header.Add("sec-ch-ua-platform", "Windows")
	//req.Header.Add("sec-fetch-dest", "empty")
	//req.Header.Add("sec-fetch-mode", "cors")
	//req.Header.Add("sec-fetch-site", "same-origin")
	//req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Add("x-requested-with", "XMLHttpRequest")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println(resp.Status)

	result := make(map[string]interface{})

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		panic(err)
	}
	apsMap := result["adverts"]
	i, ok := apsMap.(map[string]interface{})
	if !ok {
		_, empty := apsMap.([]interface{})
		if !empty {
			panic("Invalid data type")
		}
	}

	var aps []interface{}
	for _, value := range i {
		aps = append(aps, value)
	}

	fmt.Println(len(aps))
}
