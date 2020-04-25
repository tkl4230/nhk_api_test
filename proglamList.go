package nhk_api_test

import "time"

type Program struct {
	List List `json:"list"`
}

type List struct {
	E1 []E1 `json:"e1"`
}

type E1 struct {
	ID        string    `json:"id"`
	EventID   string    `json:"event_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Area      Area      `json:"area"`
	Service   Service   `json:"service"`
	Title     string    `json:"title"`
	Subtitle  string    `json:"subtitle"`
	Content   string    `json:"content"`
	Act       string    `json:"act"`
	Genres    []string  `json:"genres"`
}

type Area struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Service struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	LogoS Logo   `json:"logo_s"`
	LogoM Logo   `json:"logo_m"`
	LogoL Logo   `json:"logo_l"`
}

type Logo struct {
	URL    string `json:"url"`
	Width  string `json:"width"`
	Height string `json:"height"`
}
