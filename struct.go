package main

type Search struct {
	TotalFound  int `json:"totalFound"`
	ResultTypes []struct {
		TotalFound  int    `json:"totalFound"`
		Type        string `json:"type"`
		DisplayName string `json:"displayName"`
	} `json:"resultTypes"`
	Results []struct {
		Type        string `json:"type"`
		TotalFound  int    `json:"totalFound"`
		Page        int    `json:"page"`
		Limit       int    `json:"limit"`
		DisplayName string `json:"displayName"`
		Contents    []struct {
			ID          string `json:"id"`
			DisplayName string `json:"displayName"`
			Description string `json:"description"`
			Subtitle    string `json:"subtitle"`
			Link        struct {
				Web string `json:"web"`
			} `json:"link"`
			Image struct {
				Default      string `json:"default"`
				DefaultLabel string `json:"defaultLabel"`
				Mobile       string `json:"mobile"`
				MobileLabel  string `json:"mobileLabel"`
			} `json:"image"`
		} `json:"contents"`
	} `json:"results"`
}

type PlayerStats struct {
	NAME string
	PPG  float32
	APG  float32
	RPG  float32
	FG   float32
}
