package loadtest

import "simply_smart_mqtt_load_testing/pkg/model"

type Result struct {
	Gateway       int32      `json:"gateway"`
	Interval      []Interval `json:"interval"`
	Cycle         int        `json:"cycle"`
	Limit         int        `json:"limit"`
	Offset        int        `json:"offset"`
	StartActivity int        `json:"start_activity"`
}

type ResultInterface interface {
	SetGateway(gateway int32)
	SetInterval(interval Interval)
	IncrementTotalCycle()
	SetLastActivity(lastActivity int)
}

func (re *Result) SetGateway(gateway int32) {
	re.Gateway = gateway
}
func (re *Result) SetInterval(interval Interval) {
	re.Interval = []Interval{}
	re.Interval = append(re.Interval, interval)
}
func (re *Result) IncrementTotalCycle() {
	re.Cycle++
}
func (re *Result) SetLastActivity(lastActivity int) {
	re.StartActivity = lastActivity
}

type WorkerResult struct {
	Gateway      int32 `json:"gateway"`
	Start        int   `json:"start"`
	LastActivity int   `json:"last_activity"`
	Total        int   `json:"total"`
	Fail         int   `json:"fail"`
	Success      int   `json:"success"`
}

type Details struct {
	MeterparGateway int `json:"meterper_gateway"`
	TotalGateway    int `josn:"gateway"`
	Interval        int `json:"interval"`
	RunningTime     int `json:"runing_time"`
}

type Interval struct {
	Start        int `json:"start,omitempty"`
	LastActivity int `json:"last_activity,omitempty"`
	TimeTaken    int `json:"time_taken"`
	Total        int `json:"total"`
	Fail         int `json:"fail"`
	Success      int `json:"success"`
}

type JobPush struct {
	Gateway   int
	thingKeys model.Thingies
}
