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

type Request struct {
	ID          int         `json:"id"`
	Status      int         `json:"status"`
	Media       string      `json:"media"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
	RequestedBy RequestedBy `json:"requestedBy"`
	ModifiedBy  ModifiedBy  `json:"modifiedBy"`
	Is4K        bool        `json:"is4k"`
	ServerID    int         `json:"serverId"`
	ProfileID   int         `json:"profileId"`
	RootFolder  string      `json:"rootFolder"`
}
type RequestedBy struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	PlexToken    string    `json:"plexToken"`
	PlexUsername string    `json:"plexUsername"`
	UserType     int       `json:"userType"`
	Permissions  int       `json:"permissions"`
	Avatar       string    `json:"avatar"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	RequestCount int       `json:"requestCount"`
}
type ModifiedBy struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	PlexToken    string    `json:"plexToken"`
	PlexUsername string    `json:"plexUsername"`
	UserType     int       `json:"userType"`
	Permissions  int       `json:"permissions"`
	Avatar       string    `json:"avatar"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	RequestCount int       `json:"requestCount"`
}
