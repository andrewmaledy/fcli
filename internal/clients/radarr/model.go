package radarr

type Movie struct {
	ID                    int              `json:"id"`
	Title                 string           `json:"title"`
	OriginalTitle         string           `json:"originalTitle"`
	OriginalLanguage      Language         `json:"originalLanguage"`
	AlternateTitles       []AlternateTitle `json:"alternateTitles"`
	SecondaryYear         int              `json:"secondaryYear"`
	SecondaryYearSourceID int              `json:"secondaryYearSourceId"`
	SortTitle             string           `json:"sortTitle"`
	SizeOnDisk            int              `json:"sizeOnDisk"`
	Status                string           `json:"status"`
	Overview              string           `json:"overview"`
	InCinemas             string           `json:"inCinemas"`
	PhysicalRelease       string           `json:"physicalRelease"`
	DigitalRelease        string           `json:"digitalRelease"`
	ReleaseDate           string           `json:"releaseDate"`
	PhysicalReleaseNote   string           `json:"physicalReleaseNote"`
	Images                []Image          `json:"images"`
	Website               string           `json:"website"`
	RemotePoster          string           `json:"remotePoster"`
	Year                  int              `json:"year"`
	YouTubeTrailerID      string           `json:"youTubeTrailerId"`
	Studio                string           `json:"studio"`
	Path                  string           `json:"path"`
	QualityProfileID      int              `json:"qualityProfileId"`
	HasFile               bool             `json:"hasFile"`
	MovieFileID           int              `json:"movieFileId"`
	Monitored             bool             `json:"monitored"`
	MinimumAvailability   string           `json:"minimumAvailability"`
	IsAvailable           bool             `json:"isAvailable"`
	FolderName            string           `json:"folderName"`
	Runtime               int              `json:"runtime"`
	CleanTitle            string           `json:"cleanTitle"`
	IMDbID                string           `json:"imdbId"`
	TMDBID                int              `json:"tmdbId"`
	TitleSlug             string           `json:"titleSlug"`
	RootFolderPath        string           `json:"rootFolderPath"`
	Folder                string           `json:"folder"`
	Certification         string           `json:"certification"`
	Genres                []string         `json:"genres"`
	Tags                  []int            `json:"tags"`
	Added                 string           `json:"added"`
	AddOptions            AddOptions       `json:"addOptions"`
	Ratings               Ratings          `json:"ratings"`
	MovieFile             MovieFile        `json:"movieFile"`
	Collection            Collection       `json:"collection"`
	Popularity            float32          `json:"popularity"`
	Statistics            Statistics       `json:"statistics"`
}

type Language struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type AlternateTitle struct {
	ID              int    `json:"id"`
	SourceType      string `json:"sourceType"`
	MovieMetadataID int    `json:"movieMetadataId"`
	Title           string `json:"title"`
	CleanTitle      string `json:"cleanTitle"`
}

type Image struct {
	CoverType string `json:"coverType"`
	URL       string `json:"url"`
	RemoteURL string `json:"remoteUrl"`
}

type AddOptions struct {
	IgnoreEpisodesWithFiles    bool   `json:"ignoreEpisodesWithFiles"`
	IgnoreEpisodesWithoutFiles bool   `json:"ignoreEpisodesWithoutFiles"`
	Monitor                    string `json:"monitor"`
	SearchForMovie             bool   `json:"searchForMovie"`
	AddMethod                  string `json:"addMethod"`
}

type Ratings struct {
	IMDb           Rating `json:"imdb"`
	TMDB           Rating `json:"tmdb"`
	Metacritic     Rating `json:"metacritic"`
	RottenTomatoes Rating `json:"rottenTomatoes"`
}

type Rating struct {
	Votes int     `json:"votes"`
	Value float32 `json:"value"`
	Type  string  `json:"type"`
}

type MovieFile struct {
	ID                  int            `json:"id"`
	MovieID             int            `json:"movieId"`
	RelativePath        string         `json:"relativePath"`
	Path                string         `json:"path"`
	Size                int            `json:"size"`
	DateAdded           string         `json:"dateAdded"`
	SceneName           string         `json:"sceneName"`
	ReleaseGroup        string         `json:"releaseGroup"`
	Edition             string         `json:"edition"`
	Languages           []Language     `json:"languages"`
	Quality             Quality        `json:"quality"`
	CustomFormats       []CustomFormat `json:"customFormats"`
	CustomFormatScore   int            `json:"customFormatScore"`
	IndexerFlags        int            `json:"indexerFlags"`
	MediaInfo           MediaInfo      `json:"mediaInfo"`
	OriginalFilePath    string         `json:"originalFilePath"`
	QualityCutoffNotMet bool           `json:"qualityCutoffNotMet"`
}

type Quality struct {
	Quality  QualityDetail `json:"quality"`
	Revision Revision      `json:"revision"`
}

type QualityDetail struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Source     string `json:"source"`
	Resolution int    `json:"resolution"`
	Modifier   string `json:"modifier"`
}

type Revision struct {
	Version  int  `json:"version"`
	Real     int  `json:"real"`
	IsRepack bool `json:"isRepack"`
}

type CustomFormat struct {
	ID                              int             `json:"id"`
	Name                            string          `json:"name"`
	IncludeCustomFormatWhenRenaming bool            `json:"includeCustomFormatWhenRenaming"`
	Specifications                  []Specification `json:"specifications"`
}

type Specification struct {
	ID                 int      `json:"id"`
	Name               string   `json:"name"`
	Implementation     string   `json:"implementation"`
	ImplementationName string   `json:"implementationName"`
	InfoLink           string   `json:"infoLink"`
	Negate             bool     `json:"negate"`
	Required           bool     `json:"required"`
	Fields             []Field  `json:"fields"`
	Presets            []string `json:"presets"`
}

type Field struct {
	Order                       int            `json:"order"`
	Name                        string         `json:"name"`
	Label                       string         `json:"label"`
	Unit                        string         `json:"unit"`
	HelpText                    string         `json:"helpText"`
	HelpTextWarning             string         `json:"helpTextWarning"`
	HelpLink                    string         `json:"helpLink"`
	Value                       string         `json:"value"`
	Type                        string         `json:"type"`
	Advanced                    bool           `json:"advanced"`
	SelectOptions               []SelectOption `json:"selectOptions"`
	SelectOptionsProviderAction string         `json:"selectOptionsProviderAction"`
	Section                     string         `json:"section"`
	Hidden                      string         `json:"hidden"`
	Privacy                     string         `json:"privacy"`
	Placeholder                 string         `json:"placeholder"`
	IsFloat                     bool           `json:"isFloat"`
}

type SelectOption struct {
	Value        int    `json:"value"`
	Name         string `json:"name"`
	Order        int    `json:"order"`
	Hint         string `json:"hint"`
	DividerAfter bool   `json:"dividerAfter"`
}

type MediaInfo struct {
	ID                    int     `json:"id"`
	AudioBitrate          int     `json:"audioBitrate"`
	AudioChannels         float32 `json:"audioChannels"`
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

type Collection struct {
	Title  string `json:"title"`
	TMDBID int    `json:"tmdbId"`
}

type Statistics struct {
	MovieFileCount int      `json:"movieFileCount"`
	SizeOnDisk     int      `json:"sizeOnDisk"`
	ReleaseGroups  []string `json:"releaseGroups"`
}
