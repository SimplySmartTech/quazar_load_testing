package model

// Thing: Database table structure
type Thingie struct {
	Id                  int32   `json:"id"`
	Name                string  `json:"name"`
	Template            string  `json:"template"`
	ThingKey            string  `json:"thingkey"`
	SiteId              int32   `json:"siteid"`
	LastActivity        float64 `json:"lastactivity"`
	UserDefinedProperty string  `json:"userdefinedproperty"`
	TemplateId          int32   `json:"templateid"`
}

type Thingies []Thingie

// init-load-test: API request structure
type InitializeLoadTestRequest struct {
	MeterparGateway int `json:"meterper_gateway"`
	TotalGateway    int `josn:"gateway"`
	Interval        int `json:"interval"`
	RunningTime     int `json:"runing_time"`
}

// Payload to send to the server
type PublishRequestPayloadData struct {
	PulseL      int32 `json:"pulsel"`
	PulseH      int32 `json:"pulseh"`
	MeterStatus int32 `json:"meterstatus"`
	ValveStatus int32 `json:"valvestatus"`
}
type PublishRequestPayload struct {
	Data      PublishRequestPayloadData `json:"data"`
	ThingKey  string                    `json:"thing_key"`
	AccessKey string                    `json:"access_key"`
	Timestamp int32                     `json:"time_stamp"`
}

// Payload data for sending to server
type SendPayload struct {
	Topic string `json:"topic"`
	PublishRequestPayload
}
type SendPayloads []SendPayload
