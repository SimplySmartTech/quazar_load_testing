package loadtest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"simply_smart_mqtt_load_testing/pkg/database"
	"simply_smart_mqtt_load_testing/pkg/model"
	"simply_smart_mqtt_load_testing/restapiservice"
	"time"

	mqttt "simply_smart_mqtt_load_testing/pkg/mqtt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type LoadTesting interface {
	FetchThingies(ctx context.Context, limit int32, sql string) (model.Thingies, error)
	GeneratePayload(ctx context.Context, loadParameter model.InitializeLoadTestRequest, client mqtt.Client)
	DocLogger(ctx context.Context, sendPayload string)
	InitializeGoRoutine(ctx context.Context, client mqtt.Client, thingKeys model.Thingies, gid int32) (results WorkerResult)
	worker(ctx context.Context, id int32, thingKeys model.Thingies, client mqtt.Client, result chan<- WorkerResult)
	CountTemplateThingKeysData(ctx context.Context) (int, error)
	RemoveThingsData(ctx context.Context) (int, error)
	CreateThingKey(ctx context.Context, token string, count int) error
}

type loadTesting struct {
	db             database.DatabaseOps
	mClient        mqttt.MqttClient
	msgCounter     int32
	gatewayCounter int32
	RESTService    restapiservice.Service
}

// Create an instance of load testing
func GetLoadTesting(db database.DatabaseOps, mClient mqttt.MqttClient, resp restapiservice.Service) LoadTesting {
	return &loadTesting{
		db:             db,
		mClient:        mClient,
		msgCounter:     0,
		gatewayCounter: 0,
		RESTService:    resp,
	}
}

// Create a service to handle request for fetching all thingies from database
func (lt *loadTesting) FetchThingies(ctx context.Context, limit int32, sql string) (model.Thingies, error) {
	return lt.db.FetchThingies(ctx, limit, sql)
}

// Generate payload data for thingies
func (lt *loadTesting) GeneratePayload(ctx context.Context, loadParameter model.InitializeLoadTestRequest, client mqtt.Client) {
	total_gateways := lt.db.TotalGateway(ctx, int32(loadParameter.MeterparGateway))
	log.Println("total_gateways", total_gateways)

	lt.InitializeGateways(ctx, loadParameter, client)

}

// Initialize Gateways
func (lt *loadTesting) InitializeGateways(ctx context.Context, loadParameter model.InitializeLoadTestRequest, client mqtt.Client) {

	var i int32 = 0

	// this is error here
	fmt.Println(loadParameter.TotalGateway, "***********gatway")
	reCh := make(chan WorkerResult)
	// for start := time.Now(); time.Since(start) < (time.Second * time.Duration(maxrunningtime)); {

	results := make([]*Result, loadParameter.TotalGateway+1)
	// temp_eZc3fHCCEg
	// temp_2xColxZtG1
	for j := 0; j <= loadParameter.TotalGateway; j++ {
		results[j] = &Result{
			Gateway: int32(j),
			Limit:   loadParameter.MeterparGateway,
			Offset:  j * loadParameter.MeterparGateway,
		}

	}
	fmt.Printf("%+v", results)
	counter := 0
	for start := time.Now().Unix(); start+int64(loadParameter.RunningTime*60)-1 > time.Now().Unix(); {
		// fmt.Println("time ", start, " ", (start + int64(loadParameter.RunningTime*60)), " ", start+int64(loadParameter.RunningTime*60) >= time.Now().Unix())
		for i = 0; i < int32(loadParameter.TotalGateway); i++ {
			sql := fmt.Sprintf("select id, name,template,thing_key,extract(epoch from last_activity) as last_activity,site_id,user_defined_property,templates_id from things order by id limit %v offset %v", loadParameter.MeterparGateway, results[i].Offset)
			thingies, err := lt.db.FetchThingies(ctx, int32(loadParameter.MeterparGateway), sql)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			// fmt.Println((results[i].StartActivity + loadParameter.Interval*60), " ", int(time.Now().Unix()), " ", (results[i].StartActivity+loadParameter.Interval*60 <= int(time.Now().Unix())))
			if (results[i].StartActivity + loadParameter.Interval*60) <= int(time.Now().Unix()) {
				results[i].StartActivity = int(time.Now().Unix())
				counter++
				go lt.worker(ctx, i, thingies, client, reCh)
			}

		}

		// time.Sleep(time.Second * time.Duration(interva

	}
	fmt.Println(counter, " counter")
	for i = 0; i < int32(counter); i++ {
		result := <-reCh
		fmt.Println("result in channel loop: ", result)
		results[result.Gateway].SetInterval(Interval{
			Start:        result.Start,
			LastActivity: result.LastActivity,
			Total:        result.Total,
			Success:      result.Success,
			Fail:         result.Fail,
		})
		results[result.Gateway].IncrementTotalCycle()
		// results[result.Gateway].SetLastActivity(result.Start)
	}
	for _, res := range results {
		fmt.Printf("final outer result: %+v \n", res)
	}

}

func (lt *loadTesting) worker(ctx context.Context, id int32, thingKeys model.Thingies, client mqtt.Client, result chan<- WorkerResult) {
	res := lt.InitializeGoRoutine(ctx, client, thingKeys, id)
	result <- res
}

// Log these entries to document one by one
func (lt *loadTesting) DocLogger(ctx context.Context, sendPayload string) {
	f, err := os.OpenFile("./payloadrequestlogger.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// Create file
	// f, err = os.Create("payloadrequestlogger.txt")
	if err != nil {
		log.Println(err)
		return
	}

	// // Write data to file
	defer f.Close()
	_, err = f.WriteString(time.Now().Local().String() + " " + sendPayload + "\n")
	if err != nil {
		log.Println(err)
	}
}

// Initialize This Goroutine
func (lt *loadTesting) InitializeGoRoutine(ctx context.Context, client mqtt.Client, thingKeys model.Thingies, gid int32) (results WorkerResult) {
	// log.Println("Gateway Started:", gid)
	fmt.Println(len(thingKeys), "things *********** ", gid)
	result := WorkerResult{}
	result.Total = len(thingKeys)
	result.Gateway = gid
	result.Start = int(time.Now().Unix())
	for _, v := range thingKeys {
		var pulseL_val int32 = int32(rand.Intn(3200-1) + 1)
		var pulseH_val int32 = int32(rand.Intn(3200-1) + 1)
		var valveStatus_val int32 = int32(rand.Intn(3091-1) + 1)
		publishRequestPayloadData := model.PublishRequestPayloadData{
			PulseL:      pulseL_val,
			PulseH:      pulseH_val,
			MeterStatus: 1,
			ValveStatus: valveStatus_val,
		}
		publishRequestPayload := model.PublishRequestPayload{
			Data:      publishRequestPayloadData,
			ThingKey:  v.ThingKey,
			AccessKey: "78d9fa4teebt5add59ctb86e1a286477cb147392",
			Timestamp: int32(time.Now().Unix()),
		}
		sendPayload := model.SendPayload{
			Topic:                 fmt.Sprintf("SimplySmart/%v/event", gid),
			PublishRequestPayload: publishRequestPayload,
		}

		// convert payload to json string
		res, err := json.Marshal(sendPayload)
		// fmt.Printf("res *****%+v", sendPayload)
		if err != nil {
			log.Fatal(err)
			continue
		}

		token := client.Publish(sendPayload.Topic, 2, false, string(res))
		if token.Error() != nil {
			result.Fail++
			lt.DocLogger(ctx, "FAILED *************"+string(res))
			continue
		}
		token.Wait()
		result.Success++
		lt.DocLogger(ctx, "SUCCESS***"+string(res))
	}
	result.LastActivity = int(time.Now().Unix())
	return result
}

func (lt *loadTesting) CountTemplateThingKeysData(ctx context.Context) (int, error) {
	templateList, err := lt.db.FetchTemplateKey(ctx)
	if err != nil {
		fmt.Printf("unable to fetch data templateKey list", err.Error())
		return -1, err
	}
	count := 0
	for _, templateKey := range templateList {
		val, err := lt.db.FetchCountMeterReading(ctx, templateKey)
		if err != nil {
			fmt.Printf("error in fetching count data of template: %v\n", templateKey)
			continue
		}
		count += val
	}
	return count, nil
}

func (lt *loadTesting) RemoveThingsData(ctx context.Context) (int, error) {
	templateList, err := lt.db.FetchTemplateKey(ctx)
	if err != nil {
		fmt.Printf("unable to fetch data templateKey list", err.Error())
		return -1, err
	}
	count := 0
	for _, templateKey := range templateList {
		val, err := lt.db.RemoveTemplateThingData(ctx, templateKey)
		if err != nil {
			fmt.Printf("error in fetching count data of template: %v\n", templateKey)
			continue
		}
		count += val
	}
	return count, err
}

func (lt *loadTesting) CreateThingKey(ctx context.Context, token string, count int) error {
	return lt.RESTService.GenerateThingKeys(count, token)
}
