package sonarr

// Series represents the data structure for a series in Sonarr.
type Series struct {
	ID           int      `json:"id"`
	Title        string   `json:"title"`
	SeasonCount  int      `json:"seasonCount"`
	SizeOnDisk   int64    `json:"sizeOnDisk"`
	Monitored    bool     `json:"monitored"`
	ProfileID    int      `json:"profileId"`
	Path         string   `json:"path"`
	SeriesType   string   `json:"seriesType"`
	Network      string   `json:"network"`
	Status       string   `json:"status"`
	Year         int      `json:"year"`
	SeasonFolder bool     `json:"seasonFolder"`
	LastInfoSync string   `json:"lastInfoSync"`
	Runtime      int      `json:"runtime"`
	Images       []Image  `json:"images"`
	Genres       []string `json:"genres"`
}

// Season represents the data structure for a season of a series in Sonarr.
type Season struct {
	ID           int        `json:"id"`
	SeriesID     int        `json:"seriesId"`
	SeasonNumber int        `json:"seasonNumber"`
	Monitored    bool       `json:"monitored"`
	Statistics   Statistics `json:"statistics"`
}

// Image represents the data structure for an image in Sonarr.
type Image struct {
	CoverType string `json:"coverType"`
	URL       string `json:"url"`
}

// Statistics represents the data structure for season statistics in Sonarr.
type Statistics struct {
	PreviousAiring    string `json:"previousAiring"`
	EpisodeFileCount  int    `json:"episodeFileCount"`
	EpisodeCount      int    `json:"episodeCount"`
	TotalEpisodeCount int    `json:"totalEpisodeCount"`
	SizeOnDisk        int64  `json:"sizeOnDisk"`
}
