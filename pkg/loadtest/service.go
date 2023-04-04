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
	worker(ctx context.Context, jobs <-chan JobPush, finish chan int, client mqtt.Client, results []*Result)
	WriteResult(ctx context.Context, finish <-chan int, results []*Result, f *os.File, totalGateway int, loadParameter model.InitializeLoadTestRequest, start int64)
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

	lt.InitializeGateways(ctx, loadParameter, client)

}

// Initialize Gateways
func (lt *loadTesting) InitializeGateways(ctx context.Context, loadParameter model.InitializeLoadTestRequest, client mqtt.Client) {

	var i int32 = 0

	fmt.Println(loadParameter.TotalGateway, "***********gatway")
	job := make(chan JobPush, loadParameter.TotalGateway)

	results := make([]*Result, loadParameter.TotalGateway)
	finish := make(chan int, loadParameter.TotalGateway)

	// temp_eZc3fHCCEg
	// temp_2xColxZtG1
	for j := 0; j < loadParameter.TotalGateway; j++ {
		results[j] = &Result{
			Gateway: int32(j),
			Limit:   loadParameter.MeterparGateway,
			Offset:  j*loadParameter.MeterparGateway + 12000,
		}

	}
	var broker = os.Getenv("BROKER")
	// var broker = "localhost"
	var port = 1883
	clientObj := []mqtt.Client{}
	f, err := os.OpenFile("./result.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("error creating file", err)
	}
	defer f.Close()
	for j := 0; j < loadParameter.TotalGateway; j++ {
		opts := mqtt.NewClientOptions()
		opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
		opts.SetClientID("go_mqtt_client")
		opts.SetUsername("78d9fa4teebt5add59ctb86e1a286477cb147392")
		opts.SetPassword("78d9fa4teebt5add59ctb86e1a286477cb147391")
		opts.KeepAlive = 10000
		opts.SetDefaultPublishHandler(MessagePubHandler)
		opts.OnConnect = ConnectHandler
		opts.OnConnectionLost = ConnectLostHandler
		client1 := mqtt.NewClient(opts)
		clientObj = append(clientObj, client)
		token := client1.Connect()
		if token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		go lt.worker(ctx, job, finish, client1, results)
	}

	fmt.Printf("cycle  %+v\n", results)
	counter := 0
	// times := 0
	// flag := true
	runningTime := loadParameter.NumOfCycles * loadParameter.Interval
	start := time.Now().Unix()
	go lt.WriteResult(ctx, finish, results, f, loadParameter.TotalGateway, loadParameter, start)
	for start = start; start+int64(runningTime*60)-1 > time.Now().Unix(); {

		for i = 0; i < int32(loadParameter.TotalGateway); i++ {

			// fmt.Println((results[i].StartActivity + loadParameter.Interval*60), " ", int(time.Now().Unix()), " ", (results[i].StartActivity+loadParameter.Interval*60 <= int(time.Now().Unix())))
			if (results[i].StartActivity+loadParameter.Interval*60)+1 < int(time.Now().Unix()) {
				sql := fmt.Sprintf("select id, name,template,thing_key,extract(epoch from last_activity) as last_activity,site_id,user_defined_property,templates_id from things order by id limit %v offset %v", loadParameter.MeterparGateway, results[i].Offset)
				thingies, err := lt.db.FetchThingies(ctx, int32(loadParameter.MeterparGateway), sql)
				fmt.Println(len(thingies), "load")

				if err != nil {
					fmt.Println(err.Error())
					continue
				}
				results[i].StartActivity = int(time.Now().Unix())
				counter++
				job <- JobPush{Gateway: int(i), thingKeys: thingies}

			}

		}

	}
	time.Sleep(time.Duration(loadParameter.Interval * 60 * int(time.Second)))
	close(job)
	close(finish)
	fmt.Println(counter, " counter")

	// log.Fatal("end of results")

}

func (lt *loadTesting) worker(ctx context.Context, jobs <-chan JobPush, finish chan int, client mqtt.Client, results []*Result) {
	for job := range jobs {
		res := lt.InitializeGoRoutine(ctx, client, job.thingKeys, int32(job.Gateway))
		results[res.Gateway].SetInterval(Interval{
			TimeTaken: res.LastActivity - res.Start,
			Total:     res.Total,
			Success:   res.Success,
			Fail:      res.Fail,
		})

		results[res.Gateway].IncrementTotalCycle()
		finish <- 1

	}
}

func (lt *loadTesting) WriteResult(ctx context.Context, finish <-chan int, results []*Result, f *os.File, totalGateway int, loadParameter model.InitializeLoadTestRequest, start int64) {
	times := 0
	// done := make(chan bool)
	for _ = range finish {

		times++
		if times == totalGateway {
			// rejson, err := json.Marshal(results)
			// if err != nil {
			// 	fmt.Println(err.Error())
			// }
			// f.Write(rejson)
			f.WriteString(fmt.Sprintf("==================== Cycle: %d ====================\n", results[0].Cycle))
			fmt.Printf("times: %d\n", times)
			times = 0
			for _, res := range results {
				rejson, err := json.Marshal(res)
				if err != nil {
					fmt.Println(err.Error())
				}
				f.Write(rejson)
				f.WriteString("\n")
				fmt.Printf("final result after %d cycle: %+v \n", res.Cycle, res)
			}
			runningTime := results[0].Cycle * loadParameter.Interval
			// endtime := int(time.Now().Unix())
			// ticker := time.NewTicker(5 * time.Second) // create a ticker that ticks every 5 seconds
			// ticker := time.NewTicker(5 * time.Second)

			// Create a done channel to signal when the ticker should stop

			for start+int64(runningTime*60)-1 > time.Now().Unix() {
				// 	select {
				// 	case <-ticker.C:
				// 		go func(f *os.File) {
				// 			count, err := lt.CountTemplateThingKeysData(ctx)
				// 			_, _ = count, err
				// 			fmt.Println(time.Now())
				// 			// if err != nil {
				// 			// 	log.Printf("error in fetching things data: %v", err.Error())
				// 			// }
				// 			// f.WriteString(fmt.Sprintf("time : %v - Count : %d\n", time.Now(), count))
				// 		}(f)
				// 		// case <-done:
				// 		// 	ticker.Stop()
				// 		// 	return

				// 	}
				// }
			}
			// fmt.Println(time.Unix(start+int64(runningTime*60), 0))
			count, err := lt.CountTemplateThingKeysData(ctx)
			// fmt.Println(time.Now())
			if err != nil {
				log.Printf("error in fetching things data: %v", err.Error())
			}
			f.WriteString(fmt.Sprintf("time : %v - Count : %d\n", time.Now(), count))

		}

	}
	// close(done)
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
	fmt.Println(len(thingKeys), "in func *********** ", gid)
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
	fmt.Println(templateList, "list")
	count := 0
	time := 0
	for _, templateKey := range templateList[:48] {
		val, err := lt.db.FetchCountMeterReading(ctx, templateKey)
		time++
		fmt.Println(templateKey, val, time)
		if err != nil {
			fmt.Printf("error in fetching count data of template: %v\n", templateKey)
			continue
		}
		// fmt.Println(count, " results ", templateKey)
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

var MessagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Topic: %s | %s\n", msg.Topic(), msg.Payload())
}

var ConnectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	// fmt.Println("Connected")
}

var ConnectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, reason error) {
	// log.Printf("CLIENT %v lost connection to the broker: %v. Will reconnect...\n", reason.Error())
}
