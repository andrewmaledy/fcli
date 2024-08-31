package sonarr

type Series struct {
	Title             string           `json:"title"`
	AlternateTitles   []AlternateTitle `json:"alternateTitles,omitempty"`
	SortTitle         string           `json:"sortTitle,omitempty"`
	Status            string           `json:"status"`
	Ended             bool             `json:"ended"`
	Overview          string           `json:"overview,omitempty"`
	PreviousAiring    string           `json:"previousAiring,omitempty"`
	Network           string           `json:"network"`
	AirTime           string           `json:"airTime,omitempty"`
	Images            []Image          `json:"images"`
	RemotePoster      string           `json:"remotePoster,omitempty"`
	Seasons           []Season         `json:"seasons"`
	Year              int              `json:"year"`
	Path              string           `json:"path"`
	QualityProfileID  int              `json:"qualityProfileId"`
	SeasonFolder      bool             `json:"seasonFolder"`
	Monitored         bool             `json:"monitored"`
	MonitorNewItems   string           `json:"monitorNewItems,omitempty"`
	UseSceneNumbering bool             `json:"useSceneNumbering"`
	Runtime           int              `json:"runtime"`
	TvdbId            int              `json:"tvdbId"`
	TvrageID          int              `json:"tvRageId,omitempty"`
	TVMazeID          int              `json:"tvMazeId,omitempty"`
	TmdbID            int              `json:"tmdbId"`
	FirstAired        string           `json:"firstAired,omitempty"`
	LastAired         string           `json:"lastAired,omitempty"`
	SeriesType        string           `json:"seriesType"`
	CleanTitle        string           `json:"cleanTitle,omitempty"`
	ImdbID            string           `json:"imdbId,omitempty"`
	TitleSlug         string           `json:"titleSlug,omitempty"`
	RootFolderPath    string           `json:"rootFolderPath,omitempty"`
	Folder            string           `json:"folder,omitempty"`
	Certification     string           `json:"certification,omitempty"`
	Genres            []string         `json:"genres"`
	Tags              []int            `json:"tags,omitempty"`
	Added             string           `json:"added,omitempty"`
	Statistics        *Statistics      `json:"statistics,omitempty"`
	ID                int              `json:"id"`
}

type AlternateTitle struct {
	Title             string `json:"title"`
	SeasonNumber      int    `json:"seasonNumber"`
	SceneSeasonNumber int    `json:"sceneSeasonNumber,omitempty"`
	SceneOrigin       string `json:"sceneOrigin,omitempty"`
	Comment           string `json:"comment,omitempty"`
}

// Season represents the data structure for a season of a series in Sonarr.
type Season struct {
	SeasonNumber int        `json:"seasonNumber"`
	Monitored    bool       `json:"monitored"`
	Statistics   Statistics `json:"statistics"`
}

// Image represents the data structure for an image in Sonarr.
type Image struct {
	CoverType string `json:"coverType"`
	URL       string `json:"url"`
	RemoteURL string `json:"remoteUrl"`
}

// Statistics represents the data structure for season statistics in Sonarr.
type Statistics struct {
	EpisodeFileCount  int      `json:"episodeFileCount"`
	EpisodeCount      int      `json:"episodeCount"`
	TotalEpisodeCount int      `json:"totalEpisodeCount"`
	SizeOnDisk        int64    `json:"sizeOnDisk"`
	ReleaseGroups     []string `json:releaseGroups`
}

// EpisodeFile represents the basic structure of an episode file in Sonarr.
type EpisodeFile struct {
	ID                  int    `json:"id"`
	SeriesID            int    `json:"seriesId"`
	SeasonNumber        int    `json:"seasonNumber"`
	RelativePath        string `json:"relativePath"`
	Path                string `json:"path"`
	Size                int64  `json:"size"`
	DateAdded           string `json:"dateAdded"`
	SceneName           string `json:"sceneName"`
	ReleaseGroup        string `json:"releaseGroup"`
	QualityCutoffNotMet bool   `json:"qualityCutoffNotMet"`
}
