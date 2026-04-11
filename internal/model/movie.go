package model

// Movie represents the simplified movie domain model
type Movie struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Overview    string  `json:"overview"`
	ReleaseDate string  `json:"release_date"`
	PosterPath  string  `json:"poster_path"`
	VoteAverage float64 `json:"vote_average"`
}

// SearchResponse is the response for the search endpoint
type SearchResponse struct {
	Results []Movie `json:"results"`
	Total   int     `json:"total"`
}

// TMDBSearchResponse represents the TMDB API search response
type TMDBSearchResponse struct {
	Page         int           `json:"page"`
	Results      []TMDBMovie   `json:"results"`
	TotalPages   int           `json:"total_pages"`
	TotalResults int           `json:"total_results"`
}

// TMDBMovie represents a movie from TMDB API
type TMDBMovie struct {
	ID               int     `json:"id"`
	Title            string  `json:"title"`
	Overview         string  `json:"overview"`
	ReleaseDate      string  `json:"release_date"`
	PosterPath       string  `json:"poster_path"`
	VoteAverage      float64 `json:"vote_average"`
	Adult            bool    `json:"adult"`
	BackdropPath     string  `json:"backdrop_path"`
	GenreIDs         []int   `json:"genre_ids"`
	OriginalLanguage string  `json:"original_language"`
	OriginalTitle    string  `json:"original_title"`
	Popularity       float64 `json:"popularity"`
	Video            bool    `json:"video"`
	VoteCount        int     `json:"vote_count"`
}

// TMDBMovieDetail represents detailed movie info from TMDB API
type TMDBMovieDetail struct {
	ID               int     `json:"id"`
	Title            string  `json:"title"`
	Overview         string  `json:"overview"`
	ReleaseDate      string  `json:"release_date"`
	PosterPath       string  `json:"poster_path"`
	VoteAverage      float64 `json:"vote_average"`
	Adult            bool    `json:"adult"`
	BackdropPath     string  `json:"backdrop_path"`
	Budget           int     `json:"budget"`
	Homepage         string  `json:"homepage"`
	IMDbID           string  `json:"imdb_id"`
	OriginalLanguage string  `json:"original_language"`
	OriginalTitle    string  `json:"original_title"`
	Popularity       float64 `json:"popularity"`
	Revenue          int     `json:"revenue"`
	Runtime          int     `json:"runtime"`
	Status           string  `json:"status"`
	Tagline          string  `json:"tagline"`
	Video            bool    `json:"video"`
	VoteCount        int     `json:"vote_count"`
}

// ToMovie converts a TMDBMovie to the simplified Movie model
func (t *TMDBMovie) ToMovie() Movie {
	return Movie{
		ID:          t.ID,
		Title:       t.Title,
		Overview:    t.Overview,
		ReleaseDate: t.ReleaseDate,
		PosterPath:  t.PosterPath,
		VoteAverage: t.VoteAverage,
	}
}

// ToMovie converts a TMDBMovieDetail to the simplified Movie model
func (t *TMDBMovieDetail) ToMovie() Movie {
	return Movie{
		ID:          t.ID,
		Title:       t.Title,
		Overview:    t.Overview,
		ReleaseDate: t.ReleaseDate,
		PosterPath:  t.PosterPath,
		VoteAverage: t.VoteAverage,
	}
}
