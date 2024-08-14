package overseer

// apiResponse models the structure of the API's JSON response.
type apiResponse struct {
	PageInfo pageInfo    `json:"pageInfo"`
	Results  []MediaItem `json:"results"`
}

// pageInfo contains pagination details from the API response.
type pageInfo struct {
	Page    int `json:"page"`
	Pages   int `json:"pages"`
	Results int `json:"results"`
}

// mediaItem represents each media entry returned by the API.
type MediaItem struct {
	Id     int    `json:"id"`
	TmdbId int    `json:"tmdbId"`
	Title  string `json:"title"`
	Size   int64  `json:"size"`
	// Other fields are omitted for brevity.
}
