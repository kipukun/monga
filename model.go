package main

import "time"

type MathomeResponse struct {
	Baseurl string `json:"baseUrl"`
}

type MangaResponse struct {
	Result string `json:"result"`
	Data   struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			Title struct {
				En string `json:"en"`
			} `json:"title"`
			Alttitles []struct {
				En string `json:"en"`
			} `json:"altTitles"`
			Description struct {
				De string `json:"de"`
				En string `json:"en"`
				Fr string `json:"fr"`
				It string `json:"it"`
				Ro string `json:"ro"`
				Ru string `json:"ru"`
				Tr string `json:"tr"`
			} `json:"description"`
			Islocked bool `json:"isLocked"`
			Links    struct {
				Al    string `json:"al"`
				Ap    string `json:"ap"`
				Kt    string `json:"kt"`
				Mu    string `json:"mu"`
				Nu    string `json:"nu"`
				Mal   string `json:"mal"`
				Raw   string `json:"raw"`
				Engtl string `json:"engtl"`
			} `json:"links"`
			Originallanguage       string      `json:"originalLanguage"`
			Lastvolume             interface{} `json:"lastVolume"`
			Lastchapter            interface{} `json:"lastChapter"`
			Publicationdemographic string      `json:"publicationDemographic"`
			Status                 string      `json:"status"`
			Year                   interface{} `json:"year"`
			Contentrating          string      `json:"contentRating"`
			Tags                   []struct {
				ID         string `json:"id"`
				Type       string `json:"type"`
				Attributes struct {
					Name struct {
						En string `json:"en"`
					} `json:"name"`
					Description []interface{} `json:"description"`
					Group       string        `json:"group"`
					Version     int           `json:"version"`
				} `json:"attributes"`
			} `json:"tags"`
			Createdat time.Time `json:"createdAt"`
			Updatedat time.Time `json:"updatedAt"`
			Version   int       `json:"version"`
		} `json:"attributes"`
	} `json:"data"`
	Relationships []struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	} `json:"relationships"`
}

type FeedResponse struct {
	Results []struct {
		Result string `json:"result"`
		Data   struct {
			ID         string `json:"id"`
			Type       string `json:"type"`
			Attributes struct {
				Title              string   `json:"title"`
				Volume             string   `json:"volume"`
				Chapter            string   `json:"chapter"`
				Translatedlanguage string   `json:"translatedLanguage"`
				Hash               string   `json:"hash"`
				Data               []string `json:"data"`
				Datasaver          []string `json:"dataSaver"`
				Uploader           string   `json:"uploader"`
				Version            int      `json:"version"`
				Createdat          string   `json:"createdAt"`
				Updatedat          string   `json:"updatedAt"`
				Publishat          string   `json:"publishAt"`
			} `json:"attributes"`
		} `json:"data"`
		Relationships []struct {
			ID   string `json:"id"`
			Type string `json:"type"`
		} `json:"relationships"`
	} `json:"results"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Result string `json:"result"`
	Token  struct {
		Session string `json:"session"`
		Refresh string `json:"refresh"`
	} `json:"token"`
}
