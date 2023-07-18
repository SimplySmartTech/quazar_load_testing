## The Repository contains the code for the load Testing of the Quazar api service

## Contents:

# model/thing.go : 
This file defines the following data structures:

Thingie: Represents a database table structure with fields such as Id, Name, Template, ThingKey, SiteId, LastActivity, UserDefinedProperty, and TemplateId.

Thingies: Represents a collection of Thingie objects.

InitializeLoadTestRequest: Represents the structure for an API request to initialize a load test with fields like MeterparGateway, TotalGateway, Interval, and NumOfCycles.

PublishRequestPayloadData: Represents the payload data to be sent to the server for publishing a request with fields like PulseL, PulseH, MeterStatus, and ValveStatus.

PublishRequestPayload: Represents the payload to be sent to the server for publishing a request, including Data, ThingKey, AccessKey, and Timestamp.

SendPayload: Represents the payload data to be sent to the server with a Topic and PublishRequestPayload.

SendPayloads: Represents a collection of SendPayload objects.


## loadtest/domain.go:
This file defines the following structures and methods:

Result: 
Represents the result of a load test with fields like Gateway, Interval, Cycle, Limit, Offset, and StartActivity. It also implements the ResultInterface interface.
Interval: defines the time period between two consecutive cycles.
Cycle: is the number of cycles performed
Limit: It defines the number of meters per gateway
StartActivity: defines the time when the worker was started.


ResultInterface:
Defines the interface for manipulating the Result structure. It includes methods like SetGateway, SetInterval, IncrementTotalCycle, and SetLastActivity.


WorkerResult: 
Represents the result of a worker in the load test with fields like Gateway, Start, LastActivity, Total, Fail, and Success.

Details: 
Represents additional details of the load test with fields like MeterparGateway, TotalGateway, Interval, and RunningTime.

Interval: 
Represents an interval within the load test with fields like Start, LastActivity, TimeTaken, Total, Fail, and Success.

JobPush: 
Represents a job push with fields like Gateway and thingKeys, where thingKeys is a collection of Thingie objects from the model package.


# loadtest/service.go:
This file contains the implementation of load testing using the LoadTesting interface and the loadTesting struct.

LoadTesting: Defines the interface for load testing operations. It includes methods such as FetchThingies, GeneratePayload, DocLogger, InitializeGoRoutine, worker, WriteResult, CountTemplateThingKeysData, RemoveThingsData, and CreateThingKey.

loadTesting: Implements the LoadTesting interface. It includes the necessary fields for load testing, such as db for database operations, mClient for MQTT client, msgCounter for message counting, gatewayCounter for gateway counting, and RESTService for REST API service.

GetLoadTesting: Creates an instance of the loadTesting struct.

FetchThingies: Fetches a specified number of thingies from the database based on the provided SQL query.

GeneratePayload: Generates payload data for thingies based on the load testing parameters and MQTT client.

InitializeGateways: Initializes gateways for the load testing process, setting up necessary configurations, and spawning worker goroutines.

worker: Performs the load testing worker routine, executing tasks based on the received jobs and updating the load testing results.

WriteResult: Writes the load testing results to a file, including the cycle details and the data collected from the worker routines.

DocLogger: Logs the payload requests for documentation purposes.

InitializeGoRoutine: Initializes a goroutine for load testing, generating and publishing MQTT payloads to the specified thingies.

RemoveThingsData: Removes the data associated with the thingies.

CreateThingKey: Creates new thing keys based on the provided token and count.


# restapiservice/domain.go:

ThingKeyCreationPayload : It defines the structure to receive the request body for creating thingkeys, which consits of the following info such as:
Name of the thing, the template id associated with it and if the user wants to generate multiole things than the start and end number (eg: 1 to 10)


# routes/loadtest.go:
This file defines the LoadTest function, which sets up the necessary HTTP routes for load testing operations using the provided dependencies and MQTT client.

LoadTest: Sets up the HTTP routes for load testing operations and associates them with the corresponding handler functions.

/count-things-data-in-template: Endpoint for counting the data associated with a specific template. This operation is performed using the CountThingsData handler function from the loadtest package.

/init-load-test: Endpoint for initializing a load test. This operation is performed using the InitializeLoadTest handler function from the loadtest package. It requires a POST request and includes the necessary dependencies and MQTT client.

/remove-all-things-reading: Endpoint for removing all the reading data associated with things. This operation is performed using the RemoveThingsDataFromTemplate handler function from the loadtest package.

/create-thingkeys: Endpoint for creating new thing keys. This operation is performed using the CreateThingKeys handler function from the loadtest package. It requires a POST request and includes the necessary dependencies.