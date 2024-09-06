package sonarr

import "time"

type Series struct {
	ID                int               `json:"id"`
	Title             string            `json:"title"`
	AlternateTitles   []AlternateTitles `json:"alternateTitles"`
	SortTitle         string            `json:"sortTitle"`
	Status            string            `json:"status"`
	Ended             bool              `json:"ended"`
	ProfileName       string            `json:"profileName"`
	Overview          string            `json:"overview"`
	NextAiring        time.Time         `json:"nextAiring"`
	PreviousAiring    time.Time         `json:"previousAiring"`
	Network           string            `json:"network"`
	AirTime           string            `json:"airTime"`
	Images            []Images          `json:"images"`
	OriginalLanguage  OriginalLanguage  `json:"originalLanguage"`
	RemotePoster      string            `json:"remotePoster"`
	Seasons           []Season          `json:"seasons"`
	Year              int               `json:"year"`
	Path              string            `json:"path"`
	QualityProfileID  int               `json:"qualityProfileId"`
	SeasonFolder      bool              `json:"seasonFolder"`
	Monitored         bool              `json:"monitored"`
	MonitorNewItems   string            `json:"monitorNewItems"`
	UseSceneNumbering bool              `json:"useSceneNumbering"`
	Runtime           int               `json:"runtime"`
	TvdbID            int               `json:"tvdbId"`
	TvRageID          int               `json:"tvRageId"`
	TvMazeID          int               `json:"tvMazeId"`
	TmdbID            int               `json:"tmdbId"`
	FirstAired        time.Time         `json:"firstAired"`
	LastAired         time.Time         `json:"lastAired"`
	SeriesType        string            `json:"seriesType"`
	CleanTitle        string            `json:"cleanTitle"`
	ImdbID            string            `json:"imdbId"`
	TitleSlug         string            `json:"titleSlug"`
	RootFolderPath    string            `json:"rootFolderPath"`
	Folder            string            `json:"folder"`
	Certification     string            `json:"certification"`
	Genres            []string          `json:"genres"`
	Tags              []int             `json:"tags"`
	Added             time.Time         `json:"added"`
	AddOptions        AddOptions        `json:"addOptions"`
	Ratings           Ratings           `json:"ratings"`
	Statistics        Statistics        `json:"statistics"`
	EpisodesChanged   bool              `json:"episodesChanged"`
}
type AlternateTitles struct {
	Title             string `json:"title"`
	SeasonNumber      int    `json:"seasonNumber"`
	SceneSeasonNumber int    `json:"sceneSeasonNumber"`
	SceneOrigin       string `json:"sceneOrigin"`
	Comment           string `json:"comment"`
}
type Images struct {
	CoverType string `json:"coverType"`
	URL       string `json:"url"`
	RemoteURL string `json:"remoteUrl"`
}
type OriginalLanguage struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type Statistics struct {
	NextAiring        time.Time `json:"nextAiring"`
	PreviousAiring    time.Time `json:"previousAiring"`
	EpisodeFileCount  int       `json:"episodeFileCount"`
	EpisodeCount      int       `json:"episodeCount"`
	TotalEpisodeCount int       `json:"totalEpisodeCount"`
	SizeOnDisk        int       `json:"sizeOnDisk"`
	ReleaseGroups     []string  `json:"releaseGroups"`
	PercentOfEpisodes float64   `json:"percentOfEpisodes"`
}
type Season struct {
	SeasonNumber int        `json:"seasonNumber"`
	Monitored    bool       `json:"monitored"`
	Statistics   Statistics `json:"statistics"`
	Images       []Images   `json:"images"`
}
type AddOptions struct {
	IgnoreEpisodesWithFiles      bool   `json:"ignoreEpisodesWithFiles,omitempty"`
	IgnoreEpisodesWithoutFiles   bool   `json:"ignoreEpisodesWithoutFiles,omitempty"`
	Monitor                      string `json:"monitor,omitempty"`
	SearchForMissingEpisodes     bool   `json:"searchForMissingEpisodes,omitempty"`
	SearchForCutoffUnmetEpisodes bool   `json:"searchForCutoffUnmetEpisodes,omitempty"`
}
type Ratings struct {
	Votes int     `json:"votes"`
	Value float32 `json:"value"`
}
type EpisodeFile struct {
	ID                  int             `json:"id"`
	SeriesID            int             `json:"seriesId"`
	SeasonNumber        int             `json:"seasonNumber"`
	RelativePath        string          `json:"relativePath"`
	Path                string          `json:"path"`
	Size                int             `json:"size"`
	DateAdded           time.Time       `json:"dateAdded"`
	SceneName           string          `json:"sceneName"`
	ReleaseGroup        string          `json:"releaseGroup"`
	Languages           []Languages     `json:"languages"`
	Quality             Quality         `json:"quality"`
	CustomFormats       []CustomFormats `json:"customFormats"`
	CustomFormatScore   int             `json:"customFormatScore"`
	IndexerFlags        int             `json:"indexerFlags"`
	ReleaseType         string          `json:"releaseType"`
	MediaInfo           MediaInfo       `json:"mediaInfo"`
	QualityCutoffNotMet bool            `json:"qualityCutoffNotMet"`
}
type Languages struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type Quality struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Source     string `json:"source"`
	Resolution int    `json:"resolution"`
}
type Revision struct {
	Version  int  `json:"version"`
	Real     int  `json:"real"`
	IsRepack bool `json:"isRepack"`
}
type SelectOptions struct {
	Value int    `json:"value"`
	Name  string `json:"name"`
	Order int    `json:"order"`
	Hint  string `json:"hint"`
}
type Fields struct {
	Order                       int             `json:"order"`
	Name                        string          `json:"name"`
	Label                       string          `json:"label"`
	Unit                        string          `json:"unit"`
	HelpText                    string          `json:"helpText"`
	HelpTextWarning             string          `json:"helpTextWarning"`
	HelpLink                    string          `json:"helpLink"`
	Value                       string          `json:"value"`
	Type                        string          `json:"type"`
	Advanced                    bool            `json:"advanced"`
	SelectOptions               []SelectOptions `json:"selectOptions"`
	SelectOptionsProviderAction string          `json:"selectOptionsProviderAction"`
	Section                     string          `json:"section"`
	Hidden                      string          `json:"hidden"`
	Privacy                     string          `json:"privacy"`
	Placeholder                 string          `json:"placeholder"`
	IsFloat                     bool            `json:"isFloat"`
}
type Specifications struct {
	ID                 int      `json:"id"`
	Name               string   `json:"name"`
	Implementation     string   `json:"implementation"`
	ImplementationName string   `json:"implementationName"`
	InfoLink           string   `json:"infoLink"`
	Negate             bool     `json:"negate"`
	Required           bool     `json:"required"`
	Fields             []Fields `json:"fields"`
	Presets            []string `json:"presets"`
}
type CustomFormats struct {
	ID                              int              `json:"id"`
	Name                            string           `json:"name"`
	IncludeCustomFormatWhenRenaming bool             `json:"includeCustomFormatWhenRenaming"`
	Specifications                  []Specifications `json:"specifications"`
}
type MediaInfo struct {
	ID                    int     `json:"id"`
	AudioBitrate          int     `json:"audioBitrate"`
	AudioChannels         int     `json:"audioChannels"`
	AudioCodec            string  `json:"audioCodec"`
	AudioLanguages        string  `json:"audioLanguages"`
	AudioStreamCount      int     `json:"audioStreamCount"`
	VideoBitDepth         int     `json:"videoBitDepth"`
	VideoBitrate          int     `json:"videoBitrate"`
	VideoCodec            string  `json:"videoCodec"`
	VideoFps              float32 `json:"videoFps"`
	VideoDynamicRange     string  `json:"videoDynamicRange"`
	VideoDynamicRangeType string  `json:"videoDynamicRangeType"`
	Resolution            string  `json:"resolution"`
	RunTime               string  `json:"runTime"`
	ScanType              string  `json:"scanType"`
	Subtitles             string  `json:"subtitles"`
}
