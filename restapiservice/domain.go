package restapiservice

type ThingKeyCreationPayload struct {
	SiteId      int    `json:"site_id"`
	Name        string `json:"name"`
	TemplatesID int    `json:"templates_id"`
	Series      bool   `json:"series"`
	StartSeries int    `json:"start_series"`
	EndSeries   int    `json:"end_series"`
}
