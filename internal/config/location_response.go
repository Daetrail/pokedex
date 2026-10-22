package config

type LocationResponse struct {
	Count    int                      `json:"count"`
	Next     *string                  `json:"next"`
	Previous *string                  `json:"previous"`
	Results  []LocationResponseResult `json:"results"`
}

type LocationResponseResult struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
