package restapiservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func NewClinet() *http.Client {
	return &http.Client{}
}

type Service interface {
	GenerateThingKeys(count int, token string) error
}

type RestService struct {
}

func NreRestService() Service {
	return &RestService{}
}

func (service *RestService) GenerateThingKeys(count int, token string) error {
	URL := fmt.Sprintf("%v/things", os.Getenv("SERVER_URL"))

	client := NewClinet()
	defer client.CloseIdleConnections()

	job := make(chan *http.Request, 5)
	// responseChannel := make(chan *http.Response, 5)

	for i := 0; i < 4; i++ {
		go worker(client, job)
	}

	jobCount := 0
	keys := 20
	// go func() {
	// 	for res := range ranresponseChannel {
	// 		response := <-res
	// 		fmt.Println(response.Status, "response")
	// 		if response.StatusCode != 200 {
	// 			err := errors.New("something went wrong")
	// 			fmt.Println(err.Error())
	// 		}
	// 	}
	// }()
	for count > 0 {
		if count < 20 {
			keys = count
		}
		payload := ThingKeyCreationPayload{
			SiteId:      152,
			Name:        "load_testing_things",
			Series:      true,
			StartSeries: 1,
			EndSeries:   keys,
			TemplatesID: 172,
		}
		// payload := ThingKeyCreationPayload{
		// 	SiteId:      1,
		// 	Name:        "load_testing_things",
		// 	Series:      true,
		// 	StartSeries: 1,
		// 	EndSeries:   keys,
		// 	TemplatesID: 172,
		// }
		payloadJson, err := json.Marshal(payload)

		if err != nil {
			fmt.Println("error in marshalling payload", err.Error())
			return err
		}
		fmt.Println("json", string(payloadJson))
		request, err := http.NewRequest("POST", URL, bytes.NewReader(payloadJson))
		if err != nil {
			fmt.Println("error creating request", err.Error())
			return err
		}

		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("accesskey", token)
		fmt.Println(" int above job push")
		job <- request
		fmt.Println(" below job push")
		count -= 20
		fmt.Println("count of thing key  ", count)
		fmt.Println(jobCount, "jobcounter")
		jobCount++

	}

	close(job)

	return nil
}

func RequestApi(client *http.Client, request *http.Request) (*http.Response, error) {
	return client.Do(request)
}

func worker(client *http.Client, jobs <-chan *http.Request) {
	fmt.Println("in worker")
	for j := range jobs {
		res, err := RequestApi(client, j)
		if err != nil {
			fmt.Println("error in creating request", err.Error())
		}
		fmt.Println("response", res.Status)

	}
}
