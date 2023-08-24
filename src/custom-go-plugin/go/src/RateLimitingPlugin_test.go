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

func TestSetRateLimit(t *testing.T) {
	//t.Skip("Skipping this test because of context key type not working")
	var expected = true
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

	body := `{
		"username": "username@email.com",
		"password": "123456abcd"
		}`
	fmt.Println(configDataJson)
	// Define a variable to hold the JSON data
	var configData map[string]interface{}

	// Unmarshal the JSON string into the map
	err := json.Unmarshal([]byte(configDataJson), &configData)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}
	// Create an HTTP request with the given URL and context
	req, err := http.NewRequest("GET", "http://localhost:8080/api/1.0/documents/getByKey/testing/trip::cs00003-01", bytes.NewBuffer([]byte(body)))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Directly initializing a pointer to apidef.APIDefinition
	apiDef := &apidef.APIDefinition{
		Name:       "Another API",
		ConfigData: configData,
	}

	fmt.Println(apiDef)

	ctx.SetDefinition(req, apiDef)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-tenant-id", "milesahead2")
	req.Header.Set("Authorization", "Bearer12345abcd")
	w := httptest.NewRecorder()

	SetRateLimit(w, req)
	if isUnitTest != expected {
		t.Errorf("unitTest value should be true but was %v", false)
	}
	fmt.Println("No errors were thrown")
}

func TestSetRateLimitNoTenantId(t *testing.T) {
	//t.Skip("Skipping this test because of context key type not working")
	var expected = true
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

	body := `{
		"username": "username@email.com",
		"password": "123456abcd"
		}`
	fmt.Println(configDataJson)
	// Define a variable to hold the JSON data
	var configData map[string]interface{}

	// Unmarshal the JSON string into the map
	err := json.Unmarshal([]byte(configDataJson), &configData)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}
	// Create an HTTP request with the given URL and context
	req, err := http.NewRequest("GET", "http://localhost:8080/api/1.0/documents/getByKey/testing/trip::cs00003-01", bytes.NewBuffer([]byte(body)))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Directly initializing a pointer to apidef.APIDefinition
	apiDef := &apidef.APIDefinition{
		Name:       "Another API",
		ConfigData: configData,
	}

	fmt.Println(apiDef)

	ctx.SetDefinition(req, apiDef)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer12345abcd")
	w := httptest.NewRecorder()

	SetRateLimit(w, req)
	if isUnitTest != expected {
		t.Errorf("unitTest value should be true but was %v", false)
	}
	fmt.Println("No errors were thrown")
}

func TestGetRateLimitsDefaultValues(t *testing.T) {
	var expectedResultRequestValue float64 = 2
	var expectedResultSecondsValue float64 = 10
	var expectedResultSessionTtlValue int64 = 120
	jsonString := `{
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
					"resource": "/resource-2/",
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
						"Content-type",
						"Authorization"
					],
					"separator": "::"
				},
				"name": "requestHeader"
			}
		}
	}`

	resource := "/resource-3/"
	method := "POST"
	keyID := "abcd12345"

	rateLimitingConfig, err := generateStructFromJSON(jsonString)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	resultRequestValue, resultSecondsValue, resultSessionTtl, err := getRateLimits(rateLimitingConfig, resource, method, keyID)

	if resultRequestValue != expectedResultRequestValue {
		t.Errorf("RateLimit values was not correct -- expected %v but was %v", expectedResultRequestValue, resultRequestValue)
	}

	if resultSecondsValue != expectedResultSecondsValue {
		t.Errorf("RateLimit values was not correct -- expected %v but was %v", expectedResultSecondsValue, resultSecondsValue)
	}

	if resultSessionTtl != expectedResultSessionTtlValue {
		t.Errorf("RateLimit values was not correct -- expected %v but was %v", expectedResultSessionTtlValue, resultSessionTtl)
	}

	if err != nil {
		t.Errorf("No errors were expected")
	}
}

func TestGetRateLimitsResourceValues(t *testing.T) {
	var expectedResultRequestValue float64 = 5
	var expectedResultSecondsValue float64 = 60
	var expectedResultSessionTtlValue int64 = 120
	jsonString := `{
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
					"resource": "/resource-2/",
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
	keyID := "abcd12345"

	rateLimitingConfig, err := generateStructFromJSON(jsonString)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	resultRequestValue, resultSecondsValue, resultSessionTtl, err := getRateLimits(rateLimitingConfig, resource, method, keyID)

	if resultRequestValue != expectedResultRequestValue {
		t.Errorf("RateLimit values was not correct -- expected %v but was %v", expectedResultRequestValue, resultRequestValue)
	}

	if resultSecondsValue != expectedResultSecondsValue {
		t.Errorf("RateLimit values was not correct -- expected %v but was %v", expectedResultSecondsValue, resultSecondsValue)
	}

	if resultSessionTtl != expectedResultSessionTtlValue {
		t.Errorf("RateLimit values was not correct -- expected %v but was %v", expectedResultSessionTtlValue, resultSessionTtl)
	}

	if err != nil {
		t.Errorf("No errors were expected")
	}
}

func TestGetRateLimitsActiveIsFalse(t *testing.T) {
	var expectedResultRequestValue float64 = -1
	var expectedResultSecondsValue float64 = -1
	var expectedResultSessionTtlValue int64 = -1
	jsonString := `{
		"rateLimiting": {
			"active": false,
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
					"resource": "/resource-2/",
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
	keyID := "abcd12345"

	rateLimitingConfig, err := generateStructFromJSON(jsonString)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	resultRequestValue, resultSecondsValue, resultSessionTtl, err := getRateLimits(rateLimitingConfig, resource, method, keyID)

	if resultRequestValue != expectedResultRequestValue {
		t.Errorf("RateLimit values was not correct -- expected %v but was %v", expectedResultRequestValue, resultRequestValue)
	}

	if resultSecondsValue != expectedResultSecondsValue {
		t.Errorf("RateLimit values was not correct -- expected %v but was %v", expectedResultSecondsValue, resultSecondsValue)
	}

	if resultSessionTtl != expectedResultSessionTtlValue {
		t.Errorf("RateLimit values was not correct -- expected %v but was %v", expectedResultSessionTtlValue, resultSessionTtl)
	}

	if err != nil {
		t.Errorf("No errors were expected")
	}
}

func TestCreateUniqueKeyId(t *testing.T) {
	var expected = "milesahead1::::application/json::123456abcd"

	jsonString := `{
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
					"method": "POST",
					"requests": 5,
					"resource": "/resource-2/",
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
						"Content-type",
						"Authorization"
					],
					"separator": "::"
				},
				"name": "requestHeader"
			}
		}
	}`

	// Create a new HTTP request
	req, err := http.NewRequest("GET", "http://localhost:8080/hello/", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set headers
	req.Header.Set("x-tenant-id", "milesahead1")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer 123456abcd")

	// Set query parameters
	queryParams := req.URL.Query()
	queryParams.Add("param1", "value1")
	queryParams.Add("param2", "value2")
	req.URL.RawQuery = queryParams.Encode()

	rateLimitingConfig, err := generateStructFromJSON(jsonString)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	result := createUniqueKeyIdHeaders(rateLimitingConfig, req)

	if result != expected {
		t.Errorf("Unique key ID values was not correct -- expected %v but was %v", expected, result)
	}

	if err != nil {
		t.Errorf("No errors were expected")
	}
}

func TestCreateUniqueKeyIdXrs(t *testing.T) {
	var expected = "MILESAHEAD1"

	jsonString := `{
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
					"method": "POST",
					"requests": 5,
					"resource": "/resource-2/",
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
						"Content-type",
						"Authorization"
					],
					"separator": "::"
				},
				"name": "requestHeadersXRS"
			}
		}
	}`

	// Create a new HTTP request
	req, err := http.NewRequest("GET", "http://localhost:8080/hello/", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set headers
	req.Header.Set("x-tenant-id", "milesahead1")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic TUlMRVNBSEVBRDF8UkRDX1dlYlNlcnZpY2VzOlJvYWRuZXQxNE5ldA==")

	// Set query parameters
	queryParams := req.URL.Query()
	queryParams.Add("param1", "value1")
	queryParams.Add("param2", "value2")
	req.URL.RawQuery = queryParams.Encode()

	rateLimitingConfig, err := generateStructFromJSON(jsonString)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	result := createUniqueKeyIdHeaders(rateLimitingConfig, req)

	if result != expected {
		t.Errorf("Unique key ID values was not correct -- expected %v but was %v", expected, result)
	}
}

func TestSelectStrategyHeaderNames(t *testing.T) {
	var expected = "milesahead1::::application/json::123456abcd"
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

	req, err := http.NewRequest("GET", "http://localhost:8080/hello/", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set headers
	req.Header.Set("x-tenant-id", "milesahead1")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer 123456abcd")

	rateLimitingConfig, err := generateStructFromJSON(configDataJson)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	result := selectStrategy(rateLimitingConfig, req)

	if result != expected {
		t.Errorf("KeyId value was not correct -- expected %v but was %v", expected, result)
	}
}

func TestSelectStrategyHeaderNamesXRS(t *testing.T) {
	var expected = "MILESAHEAD1"
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
			"name": "requestHeadersXRS"
		  }
		}
	  }`

	req, err := http.NewRequest("GET", "http://localhost:8080/hello/", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set headers
	req.Header.Set("x-tenant-id", "milesahead1")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic TUlMRVNBSEVBRDF8UkRDX1dlYlNlcnZpY2VzOlJvYWRuZXQxNE5ldA==")

	rateLimitingConfig, err := generateStructFromJSON(configDataJson)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	result := selectStrategy(rateLimitingConfig, req)

	if result != expected {
		t.Errorf("KeyId value was not correct -- expected %v but was %v", expected, result)
	}
}

func TestSelectStrategySessionGuid(t *testing.T) {
	var expected = "33d9b8d0-58ba-4400-87af-bdf5f79c0f9b"
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
			"name": "sessionGuid"
		  }
		}
	  }`

	xmlBody := `
	  <soapenv:Header>
		  <dat:SessionHeader>
			  <dat:SessionGuid>33d9b8d0-58ba-4400-87af-bdf5f79c0f9b</dat:SessionGuid>
		  </dat:SessionHeader>
	  </soapenv:Header>`

	req, err := http.NewRequest("GET", "http://localhost:8080/hello/", strings.NewReader(xmlBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set headers
	req.Header.Set("x-tenant-id", "milesahead1")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer 123456abcd")

	rateLimitingConfig, err := generateStructFromJSON(configDataJson)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	result := selectStrategy(rateLimitingConfig, req)

	if result != expected {
		t.Errorf("KeyId value was not correct -- expected %v but was %v", expected, result)
	}
}

func TestSelectStrategyActiveIsFalse(t *testing.T) {
	var expected = ""
	configDataJson := `{
		"rateLimiting": {
		  "active": false,
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

	req, err := http.NewRequest("GET", "http://localhost:8080/hello/", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set headers
	req.Header.Set("x-tenant-id", "milesahead1")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer 123456abcd")

	rateLimitingConfig, err := generateStructFromJSON(configDataJson)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	result := selectStrategy(rateLimitingConfig, req)

	if result != expected {
		t.Errorf("KeyId value was not correct -- expected %v but was %v", expected, result)
	}
}

func TestSelectStrategyUnknownStrategy(t *testing.T) {
	var expected = ""
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
			"name": "unknown"
		  }
		}
	  }`

	req, err := http.NewRequest("GET", "http://localhost:8080/hello/", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set headers
	req.Header.Set("x-tenant-id", "milesahead1")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer 123456abcd")

	rateLimitingConfig, err := generateStructFromJSON(configDataJson)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	result := selectStrategy(rateLimitingConfig, req)

	if result != expected {
		t.Errorf("KeyId value was not correct -- expected %v but was %v", expected, result)
	}

}

func TestCreateUniqueKeyIDSessionGuid(t *testing.T) {
	expected := "33d9b8d0-58ba-4400-87af-bdf5f79c0f9b"

	xmlBody := `
<soapenv:Header>
    <dat:SessionHeader>
        <dat:SessionGuid>33d9b8d0-58ba-4400-87af-bdf5f79c0f9b</dat:SessionGuid>
    </dat:SessionHeader>
</soapenv:Header>`

	req, err := http.NewRequest("GET", "http://localhost:8080/hello/", strings.NewReader(xmlBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	result := createUniqueKeyIDSessionGuid(req)

	if result != expected {
		t.Errorf("Value was not correct -- expected %v but was %v", expected, result)
	}
}

func TestLookForOverridesInRequest(t *testing.T) {
	expected := "/getByKey/"
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
			  "resource": "/getByKey/",
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
			"name": "unknown"
		  }
		}
	  }`
	req, err := http.NewRequest("GET", "http://localhost:8080/api/1.0/documents/getByKey/testing/trip::cs00003-01", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	rateLimitingConfig, err := generateStructFromJSON(configDataJson)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	result := lookForOverridesInRequest(req.URL.Path, rateLimitingConfig)
	if result != expected {
		t.Errorf("Resource value was not correct -- expected %v but was %v", expected, result)
	}
}
