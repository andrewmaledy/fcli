// Model.go

package overseer

import "time"

// apiResponse models the structure of the API's JSON response for media.
type GetMediaResponse struct {
	PageInfo pageInfo `json:"pageInfo"`
	Results  []Media  `json:"results"`
}

// apiRequestResponse models the structure of the API's JSON response for requests.
type getRequestsResponse struct {
	PageInfo pageInfo  `json:"pageInfo"`
	Results  []Request `json:"results"`
}

// pageInfo contains pagination details from the API response.
type pageInfo struct {
	Page    int `json:"page"`
	Pages   int `json:"pages"`
	Results int `json:"results"`
}

// mediaItem represents each media entry returned by the API.
type Media struct {
	Id     int    `json:"id"`
	TmdbId int    `json:"tmdbId"`
	TvdbId int    `json:"tvdbId"`
	Title  string `json:"title"`
	Size   int64  `json:"size"`
	// Other fields are omitted for brevity.
}
type Season struct {
	Id           int `json:"id"`
	SeasonNumber int `json:"seasonNumber"`
}

// RequestItem represents each request entry returned by the API.
type Request struct {
	Id        int       `json:"id"`
	Is4k      bool      `json:"is4k"`
	CreatedAt time.Time `json:"createdAt"`
	Media     Media     `json:"media"`
	Season    []Season  `json:"seasons"`
}
