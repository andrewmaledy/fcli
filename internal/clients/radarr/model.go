package radarr

type Movie struct {
	Id         int    `json:"id"`
	Title      string `json:"title"`
	SizeOnDisk int64  `json:"sizeOnDisk"`
	TmdbId     int    `json:"tmdbId"`
}
