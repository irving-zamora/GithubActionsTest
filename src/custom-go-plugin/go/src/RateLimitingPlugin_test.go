package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/TykTechnologies/tyk/apidef"
	"github.com/TykTechnologies/tyk/ctx"
)

func BuildStruct() RateLimitingConfig {

	strategyConfig := StrategyConfig{
		HeaderNames: []string{"x-tenant-id",
			"Host",
			"Content-Type",
			"Authorization"},
		Separator: "::",
	}

	strategy := Strategy{
		Config: strategyConfig,
		Name:   "requestHeaders",
	}

	rateLimit := RateLimit{
		Active:     true,
		IsUnitTest: true,
		Overrides: []Override{
			{
				Method:   "GET",
				Requests: -1,
				Resource: "/testing/",
				Seconds:  -1,
			},
			{
				Method:   "GET",
				Requests: 5,
				Resource: "/resource-2/",
				Seconds:  60,
			},
		},
		Requests:      2,
		Seconds:       10,
		SessionTtlMin: 120,
		Strategy:      strategy,
	}

	rateLimiting := RateLimitingConfig{
		RateLimiting: rateLimit,
	}

	return rateLimiting
}

func Test_SetRateLimit_Success(t *testing.T) {

	body := `{
        "username": "username@email.com",
        "password": "123456abcd"
        }`
	//fmt.Println(configDataJson)
	// Define a variable to hold the JSON data

	var configData map[string]interface{}

	rateLimiting := BuildStruct()
	//rateLimiting.RateLimiting.Requests = 3
	// Unmarshal the JSON string into the map
	//configJsonRateLimit := BuildRateLimitConfig()
	marshalledConfig, err1 := json.Marshal(rateLimiting)
	if err1 != nil {
		t.Fatalf("Unable to marshal RateLimitingConfig struct: %v", err1)
	}

	err := json.Unmarshal(marshalledConfig, &configData)
	if err != nil {
		t.Fatalf("Error unmarshaling JSON: %v", err)
		return
	}
	// Create an HTTP request with the given URL and context
	req, err := http.NewRequest("GET", "http://localhost:8080/testing/", bytes.NewBuffer([]byte(body)))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
		return
	}

	// Directly initializing a pointer to apidef.APIDefinition
	apiDef := &apidef.APIDefinition{
		Name:       "Another API",
		ConfigData: configData,
	}

	ctx.SetDefinition(req, apiDef)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-tenant-id", "milesahead2")
	req.Header.Set("Authorization", "Bearer 12345abcd")
	w := httptest.NewRecorder()

	SetRateLimit(w, req)
}

func Test_ErrorLog_Success(t *testing.T) {
	ErrorLog("Config data read error!")
	fmt.Println("Success")
}

// To test alternative path of SelectStrategy func
func Test_SetRateLimitSelectStrategyActiveFalse_Success(t *testing.T) {

	body := `{
        "username": "username@email.com",
        "password": "123456abcd"
        }`

	// Define a variable to hold the JSON data
	var configData map[string]interface{}

	rateLimiting := BuildStruct()
	rateLimiting.RateLimiting.Overrides[0].Resource = "/products/"

	// Unmarshal the JSON string into the map
	//configJsonRateLimit := BuildRateLimitConfig()
	marshalledConfig, err1 := json.Marshal(rateLimiting)
	if err1 != nil {
		t.Fatalf("Unable to marshal RateLimitingConfig struct: %v", err1)
	}

	err := json.Unmarshal(marshalledConfig, &configData)
	if err != nil {
		t.Fatalf("Error unmarshaling JSON: %v", err)
		return
	}
	// Create an HTTP request with the given URL and context
	req, err := http.NewRequest("GET", "http://localhost:8080/testing/", bytes.NewBuffer([]byte(body)))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
		return
	}

	// Directly initializing a pointer to apidef.APIDefinition
	apiDef := &apidef.APIDefinition{
		Name:       "Another API",
		ConfigData: configData,
	}

	ctx.SetDefinition(req, apiDef)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-tenant-id", "milesahead2")
	req.Header.Set("Authorization", "Bearer12345abcd")
	w := httptest.NewRecorder()

	SetRateLimit(w, req)
}

func Test_GenerateStructFromJSON_Success(t *testing.T) {
	configDataJson := `{
        "rateLimiting": {
          "active": true,
          "overrides": [
            {
              "method": "GET",
              "requests": -1,
              "resource": "/products/",
              "seconds": -1
            },
            {
              "method": "GET",
              "requests": 5,
              "resource": "/hello/",
              "seconds": 60
            }
          ],
          "requests": 2,
          "seconds": 10,
          "sessionTtlMin": 120,
          "strategy": {
            "config": {
              "headerNames": [
                "x-tenant-id",
                "Host",
                "Content-Type",
                "Authorization"
              ],
              "separator": "::"
            },
            "name": "requestHeaders"
          }
        }
      }`

	rateLimitingConfig, err := generateStructFromJSON(configDataJson)

	if err != nil {
		t.Fatalf("Error: %v", err)
		return
	}

	fmt.Println("Success", rateLimitingConfig)

}

func Test_GetRateLimitsDefaultValues_Success(t *testing.T) {
	var expectedResultRequestValue float64 = 2
	var expectedResultSecondsValue float64 = 10
	var expectedResultSessionTtlValue int64 = 120

	resource := "/resource-3/"
	method := "POST"

	rateLimiting := BuildStruct()
	rateLimiting.RateLimiting.Overrides[0].Resource = "/products/"

	resultRequestValue, resultSecondsValue, resultSessionTtl, err := getRateLimits(rateLimiting, resource, method, "milesahead1")

	if resultRequestValue != expectedResultRequestValue {
		t.Fatalf("RateLimit values was not correct -- expected %v but was %v", expectedResultRequestValue, resultRequestValue)
	}

	if resultSecondsValue != expectedResultSecondsValue {
		t.Fatalf("RateLimit values was not correct -- expected %v but was %v", expectedResultSecondsValue, resultSecondsValue)
	}

	if resultSessionTtl != expectedResultSessionTtlValue {
		t.Fatalf("RateLimit values was not correct -- expected %v but was %v", expectedResultSessionTtlValue, resultSessionTtl)
	}

	if err != nil {
		t.Fatalf("No errors were expected")
	}
}

func Test_GetRateLimitsResourceValues_Success(t *testing.T) {
	var expectedResultRequestValue float64 = 5
	var expectedResultSecondsValue float64 = 60
	var expectedResultSessionTtlValue int64 = 120

	resource := "/resource-2/"
	method := "GET"

	rateLimiting := BuildStruct()
	rateLimiting.RateLimiting.Overrides[0].Resource = "/products/"

	resultRequestValue, resultSecondsValue, resultSessionTtl, err := getRateLimits(rateLimiting, resource, method, "milesahead1")

	if resultRequestValue != expectedResultRequestValue {
		t.Fatalf("RateLimit values was not correct -- expected %v but was %v", expectedResultRequestValue, resultRequestValue)
	}

	if resultSecondsValue != expectedResultSecondsValue {
		t.Fatalf("RateLimit values was not correct -- expected %v but was %v", expectedResultSecondsValue, resultSecondsValue)
	}

	if resultSessionTtl != expectedResultSessionTtlValue {
		t.Fatalf("RateLimit values was not correct -- expected %v but was %v", expectedResultSessionTtlValue, resultSessionTtl)
	}

	if err != nil {
		t.Fatalf("No errors were expected")
	}
}

func Test_GetRateLimitsActiveIsFalse_Success(t *testing.T) {

	resource := "/resource-2/"
	method := "GET"

	rateLimiting := BuildStruct()
	rateLimiting.RateLimiting.Active = false

	resultRequestValue, resultSecondsValue, resultSessionTtl, err := getRateLimits(rateLimiting, resource, method, "milesahead1")

	if (resultRequestValue == -1) && (resultSecondsValue == -1) &&
		(resultSessionTtl == -1) && err == nil {
		fmt.Println("Test passed, active false")
	} else {
		t.Fatalf("Test failed")
	}
}

func Test_GetRateLimitsErrorParsingJson_Success(t *testing.T) {

	jsonString := `{
		"rateLimiting": {
			"active": false,
			"overrides": [
				{
					"method": "GET",
					"requests": ,
					"resource": "/products/",
					"seconds": -1
				},
				{
					"method": 1,
					"requests": 5,
					"resource": "/resource-2/",
					"seconds": 60
				}
			],
			"requests": ,
			"seconds": 10,
			"sessionTtlMin": 120,
			"strategy": {
				"config": {
					"headerNames": [
						"x-tenant-id",
						"Host",
						"Content-type",
						"Authorization"
					],
					"separator": "::"
				},
				"name": "requestHeader"
			}
		}
	}`

	resource := "/resource-2/"
	method := "GET"

	rateLimitingConfig, err := generateStructFromJSON(jsonString)

	if err != nil {
		t.Logf("Error: %v", err)
		return
	}

	resultRequestValue, resultSecondsValue, resultSessionTtl, err := getRateLimits(rateLimitingConfig, resource, method, "milesaheadtest")

	if (err != nil) && err.Error() == "invalid character ',' looking for beginning of value" {
		fmt.Println("Test passed, returning values: ", resultRequestValue, resultSecondsValue, resultSessionTtl)
	} else {
		t.Fatalf("Test failed")
	}
}

// requestHeadersXRS path
func Test_SelectStrategyRequestHeadersXRSPath_Success(t *testing.T) {

	body := `{
        "username": "username@email.com",
        "password": "123456abcd"
        }`

	rateLimiting := BuildStruct()
	rateLimiting.RateLimiting.Strategy.Name = "requestHeadersXRS"

	// Create an HTTP request with the given URL and context
	req, err := http.NewRequest("GET", "http://localhost:8080/testing/", bytes.NewBuffer([]byte(body)))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-tenant-id", "milesahead2")
	req.Header.Set("Authorization", "Basic TUlMRVNBSEVBRDF8UkRDX1dlYlNlcnZpY2VzOlJvYWRuZXQxNE5ldA==")

	keyID := selectStrategy(rateLimiting, req)

	if keyID != "" {
		fmt.Println("Test passed with keyID: ", keyID)
	} else {
		t.Fatalf("Test failed")
	}
}

func TestSelectStrategy_RequestHeadersPathBasic_Success(t *testing.T) {

	body := `{
        "username": "username@email.com",
        "password": "123456abcd"
        }`

	rateLimiting := BuildStruct()

	// Create an HTTP request with the given URL and context
	req, err := http.NewRequest("GET", "http://localhost:8080/testing/", bytes.NewBuffer([]byte(body)))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-tenant-id", "milesahead2")
	req.Header.Set("Authorization", "Basic TUlMRVNBSEVBRDF8UkRDX1dlYlNlcnZpY2VzOlJvYWRuZXQxNE5ldA==")

	keyID := selectStrategy(rateLimiting, req)

	if keyID != "" {
		fmt.Println("Test passed with keyID: ", keyID)
	} else {
		t.Fatalf("Test failed")
	}
}

func Test_SelectStrategySessionGuid_Success(t *testing.T) {
	var expected = "33d9b8d0-58ba-4400-87af-bdf5f79c0f9b"

	rateLimiting := BuildStruct()
	rateLimiting.RateLimiting.Strategy.Name = "sessionGuid"

	xmlBody := `
      <soapenv:Header>
          <dat:SessionHeader>
              <dat:SessionGuid>33d9b8d0-58ba-4400-87af-bdf5f79c0f9b</dat:SessionGuid>
          </dat:SessionHeader>
      </soapenv:Header>`

	req, err := http.NewRequest("GET", "http://localhost:8080/hello/", strings.NewReader(xmlBody))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
		return
	}

	// Set headers
	req.Header.Set("x-tenant-id", "milesahead1")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer 123456abcd")

	result := selectStrategy(rateLimiting, req)

	if result != expected {
		t.Fatalf("KeyId value was not correct -- expected %v but was %v", expected, result)
	} else {
		fmt.Println("Test passed")
	}
}

func Test_SelectStrategySessionGuidNotFound_Success(t *testing.T) {

	rateLimiting := BuildStruct()
	rateLimiting.RateLimiting.Strategy.Name = "sessionGuid"

	//To cover "no SessionGuid found"
	xmlBody := `
      <soapenv:Header>
          <dat:SessionHeader>
              
          </dat:SessionHeader>
      </soapenv:Header>`

	req, err := http.NewRequest("GET", "http://localhost:8080/hello/", strings.NewReader(xmlBody))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
		return
	}

	// Set headers
	req.Header.Set("x-tenant-id", "milesahead1")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer 123456abcd")

	result := selectStrategy(rateLimiting, req)

	if result == "" {
		fmt.Println("Test passed")
	} else {
		t.Fatalf("There is an error")
	}
}
