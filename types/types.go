package types

type Greet struct {
	Message string
}

type AnimeEpisode struct {
	ImageLink     string
	AnimeName     string
	EpisodeNumber string
	Link          string
}

type Payload struct {
	Pattern string `json:"pattern"`
	Data    struct {
		PageNumber string `json:"pageNumber"`
		Url        string `json:"url"`
	} `json:"data"`
	Id string `json:"id"`
}

type Anime struct {
	ImageLink   string
	Title       string
	Summary     string
	Genres      []string
	Year        string
	Status      string
	Alias       string
	HasNextPage bool
	Episodes    []AnimeEpisode
}

type PlayerData struct {
	AnimeLink     string
	EpisodeNumber string
	PrevEpisode   string
	NextEpisode   string
	Servers       []Server
}

type Server struct {
	Title string
	Link  string
}
