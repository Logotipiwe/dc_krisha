package service

type MapData struct {
	IsTooManyAdverts bool      `json:"isTooManyAdverts"`
	ListURL          *string   `json:"listUrl"`
	MetaData         *MetaData `json:"metaData"`
	NbTotal          *int      `json:"nbTotal"`
	Results          *[]Result `json:"results"`
}

type MetaData struct {
	CanonicalURL *string `json:"canonicalUrl"`
	Description  *string `json:"description"`
	Header       *string `json:"header"`
	Keywords     *string `json:"keywords"`
	Title        *string `json:"title"`
}

type Result struct {
	Geometry   *Geometry   `json:"geometry"`
	ID         *string     `json:"id"`
	Number     *int        `json:"number"`
	Properties *Properties `json:"properties"`
	Type       *string     `json:"type"`
}

type Geometry struct {
	Coordinates *[]float64 `json:"coordinates"`
	Type        *string    `json:"type"`
}

type Properties struct {
	Agents       *int    `json:"agents"`
	CounterValue *int    `json:"counterValue"`
	Geohash      *string `json:"geohash"`
	IconContent  *any    `json:"iconContent"`
	AdvertId     *int64  `json:"advertId"`
}

type ApsResult struct {
	HTML         *string        `json:"html"`
	PriceHistory *[]interface{} `json:"priceHistory"`
	Adverts      *map[string]Ap `json:"adverts"`
	Pager        *string        `json:"pager"`
	Page         *int           `json:"page"`
}

type Ap struct {
	ID                      *int64    `json:"id"`
	Storage                 *string   `json:"storage"`
	CommentsType            *string   `json:"commentsType"`
	IsCommentable           *bool     `json:"isCommentable"`
	IsCommentableByEveryone *bool     `json:"isCommentableByEveryone"`
	IsOnMap                 *bool     `json:"isOnMap"`
	HasPrice                *bool     `json:"hasPrice"`
	Price                   *int64    `json:"price"`
	Photos                  *[]Photos `json:"photos"`
	HasPackages             *bool     `json:"hasPackages"`
	Title                   *string   `json:"title"`
	Addresstitle            *string   `json:"addresstitle"`
	UserType                *string   `json:"userType"`
	Square                  *float64  `json:"square"`
	Rooms                   *int      `json:"rooms"`
	OwnerName               *string   `json:"ownerName"`
	Status                  *string   `json:"status"`
	Map                     *Map      `json:"map"`
}

type Photos struct {
	Src   *string `json:"src"`
	W     *int    `json:"w"`
	H     *int    `json:"h"`
	Title *string `json:"title"`
	Alt   *string `json:"alt"`
}

type Map struct {
	Lat  *float64 `json:"lat"`
	Lon  *float64 `json:"lon"`
	Zoom *int     `json:"zoom"`
	Type *string  `json:"type,omitempty"`
}
