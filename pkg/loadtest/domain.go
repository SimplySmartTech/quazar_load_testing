package loadtest

type Result struct {
	Gateway int32 `json:"gateway"`
	Success int   `json:"success"`
	Fail    int   `json:"fail"`
	Total   int   `json:"total"`
}

type Details struct {
	MeterparGateway int `json:"meterper_gateway"`
	TotalGateway    int `josn:"gateway"`
	Interval        int `json:"interval"`
	RunningTime     int `json:"runing_time"`
}
