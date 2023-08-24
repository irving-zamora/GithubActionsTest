// the main entry point for Tyk.io custom plugins for API Gateway
// this was extended from the base/standard Tyk github repo for custom plugin dev
// this will serve as just a starting point for the actual implementation
// of the plugin that will implememnt rate limiting logic for the phase 1 legacy implementation
// it is a requirement that there will be not customer side changes in order for the
// api gw to be introduced into the runtime architecture
// therfore, the custom rate limiting functionality will utlize data from the request, headers and
// even the request body, to generate a unique id that will act as the rate limiting key to which
// the rate limiting will be applied
package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/TykTechnologies/tyk/ctx"
	"github.com/TykTechnologies/tyk/user"
)

var isUnitTest = false

type RateLimitingConfig struct {
	RateLimiting struct {
		Active        bool       `json:"active"`
		Overrides     []Override `json:"overrides"`
		Requests      int        `json:"requests"`
		Seconds       int        `json:"seconds"`
		SessionTtlMin int        `json:"sessionTtlMin"`
		Strategy      Strategy   `json:"strategy"`
		LogLevel      LogLevel   `json:"logLevel"`
	} `json:"rateLimiting"`
}

// I've defined additional structs Override and Strategy to accommodate the nested objects in the JSON string.
type Override struct {
	Method   string `json:"method"`
	Requests int    `json:"requests"`
	Resource string `json:"resource"`
	Seconds  int    `json:"seconds"`
}

type Strategy struct {
	Config struct {
		HeaderNames []string `json:"headerNames"`
		Separator   string   `json:"separator"`
	} `json:"config"`
	Name string `json:"name"`
}

// The rate limiter function that will be configured to be invoked for each api definition that is
// a part of the legacy phase 1 implementation.
// Sets a custom rate limit for inbound requests based on a unique identier for a given
// customer/tenant.
// The unique identifier is based on a known tag (or config) with a structure of "apiRateLimiterType::{type-value}".
// Valid type values are as follows:
//
//	requestHeader={header-name} for example "apiRateLimiterType::requestHeader=x-tenant-id"
//	soapBody ... more to document here
func SetRateLimit(rw http.ResponseWriter, r *http.Request) {

	apidef := ctx.GetDefinition(r)

	apiDefConfigData, err := json.Marshal(apidef.ConfigData)
	if err != nil {
		ErrorLog("Config data read error: ", err)
	} else {
		DebugLog("Config data read!")
		DebugLog(string(apiDefConfigData))
	}

	apiDefinitionJson := string(apiDefConfigData)

	rateLimitingConfig, err := generateStructFromJSON(apiDefinitionJson)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Set log level based on api definition config data
	SetLogLevel(rateLimitingConfig.RateLimiting.LogLevel)

	DebugLog("api-name", apidef.Name, "custom plugin BEGIN processing @ ", time.Now().String())

	DebugLog("config data: ", apidef.ConfigData)
	DebugLog("api name: ", apidef.Name)
	DebugLog("apidef tags: ", apidef.Tags)
	DebugLog("apidef tagHeaders: ", apidef.TagHeaders)

	keyID := selectStrategy(rateLimitingConfig, r)

	overridePath := lookForOverridesInRequest(r.URL.Path, rateLimitingConfig)
	DebugLog("Path: ", overridePath)

	requestsValue, secondsValue, sessionTtl, err := getRateLimits(rateLimitingConfig, overridePath, r.Method, keyID)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	DebugLog("Requests value: ", requestsValue)
	DebugLog("Seconds value: ", secondsValue)
	DebugLog("SessionTtl value: ", sessionTtl)

	InfoLog("Unique KeyID: ", keyID)

	// where the actual rate limiting is applied based on a customer's unique identifier
	// the actual rate setting should be expernally configurable such as using tag values or configs
	// this is IF ploicies are still not working -- see below
	session := &user.SessionState{
		Alias: keyID,

		//the following do not seem to be working as expected using policy?
		//it was confirmed by Tyk that in the version of the Gateway that was being used
		//that the ploicies were not yet functional for this purpose
		//ApplyPolicies: []string{"LegacyRateLimit"},
		//ApplyPolicyID: "63a0ce2ae3ce6a0001645065",

		Rate: requestsValue,
		Per:  secondsValue,
		MetaData: map[string]interface{}{ //MetaData can be "anything" -- extra data
			"keyId": keyID,
			"meta2": "meta2",
		},
		KeyID:           keyID,      //this value should be the unique value for the redis key (hashed)
		SessionLifetime: sessionTtl, //redis TTL -- rate liiting will be "reset" after key expires
	}

	// check if we are in a unit test context or real world application
	if strings.Contains(r.URL.Path, "testing") {
		DebugLog("Session will not be set as this is in a unit testing context")
		isUnitTest = true
		return
	}
	ctx.SetSession(r, session, true)

	DebugLog("api-name", apidef.Name, "Rate limiting plugin END processing @ ", time.Now().String())
}

// function takes an HTTP request URL Path and the api definition as input
// loops through the 'overrides' to find a matching 'resource' in the request
// output: string formed by the request match
// if no overrides were found an empty string is returned
func lookForOverridesInRequest(requestUrl string, rateLimitingConfig RateLimitingConfig) string {

	for _, override := range rateLimitingConfig.RateLimiting.Overrides {
		resourceValue := override.Resource

		if strings.Contains(requestUrl, resourceValue) {
			DebugLog("Override found for request: ", resourceValue)
			return resourceValue
		}
	}
	return ""
}

// The function parses the JSON string and retrieves the values of "requests" and "seconds" from the nested objects
// inside the "overrides" array. It parses the JSON string, iterates over each override object,
// and extracts the "requests" and "seconds" values.
// The function returns a slice of integers containing the "requests" and "seconds" values from all overrides.
func getRateLimits(rateLimitingConfig RateLimitingConfig, resource string, method string, keyID string) (float64, float64, int64, error) {

	if keyID == "" {
		return -1, -1, int64(-1), nil
	}

	// if active is true then look for overrides
	if rateLimitingConfig.RateLimiting.Active {

		// loop through overrides to look for a matching resource and method
		// between http request and the overrides config data
		for _, override := range rateLimitingConfig.RateLimiting.Overrides {

			// if a match is found then set the values from the overrides config
			// for the 'requests', 'seconds' and 'sessionTtl'
			if override.Resource == resource && override.Method == method {
				return float64(override.Requests), float64(override.Seconds), int64(rateLimitingConfig.RateLimiting.SessionTtlMin), nil
			}
		}
		// if no overrides found to be matching then use
		// the 'default' values set in the api definition config
		return float64(rateLimitingConfig.RateLimiting.Requests), float64(rateLimitingConfig.RateLimiting.Seconds), int64(rateLimitingConfig.RateLimiting.SessionTtlMin), nil
	}
	// If active is false or not found then set no rate limit
	return -1, -1, int64(-1), nil
}

// This function takes in JSON string and http request as input parameters
// http.Request representing the HTTP request from which the headers will be extracted.
// The function retrieves the headers based on the "headerNames" specified in the JSON configuration.
// and matches the headers from the incoming request to concatenate the values from the headers
// using the provided 'separator' from the JSON config
// the "KeyId" value is created by concatenating the "headerNames" separated by the "separator" as requested.
func createUniqueKeyIdHeaders(rateLimitingConfig RateLimitingConfig, req *http.Request) string {

	if rateLimitingConfig.RateLimiting.Strategy.Name == "requestHeadersXRS" {
		authBase64 := req.Header.Get("Authorization")
		authBase64WithoutBasic := strings.ReplaceAll(authBase64, "Basic ", "")

		rawDecodedAuth, err := base64.StdEncoding.DecodeString(authBase64WithoutBasic)
		if err != nil {
			panic(err)
		}
		customerId := strings.Split(string(rawDecodedAuth), "|")

		return customerId[0]
	}

	headers := make([]string, len(rateLimitingConfig.RateLimiting.Strategy.Config.HeaderNames))

	for i, headerName := range rateLimitingConfig.RateLimiting.Strategy.Config.HeaderNames {
		headerValue := req.Header.Get(headerName)
		if strings.Contains(headerValue, "Bearer ") {
			headerValue = strings.ReplaceAll(headerValue, "Bearer ", "")
		}
		if strings.Contains(headerValue, "Basic ") {
			headerValue = strings.ReplaceAll(headerValue, "Basic ", "")
		}
		headers[i] = headerValue
	}

	keyID := strings.Join(headers, rateLimitingConfig.RateLimiting.Strategy.Config.Separator)
	return keyID
}

// function takes an http.Request as input and retrieves the value of
// the SessionGuid element from the request body using a regular expression.
// Returns: The extracted SessionGuid value as a string, if found in the request body.
// An error, if encountered during the reading or parsing of the request body.
func createUniqueKeyIDSessionGuid(req *http.Request) string {
	body, err := ioutil.ReadAll(req.Body)
	if err == nil {
		sb := string(body)
		DebugLog("request body: ", sb)

		re := regexp.MustCompile(`<.*?:SessionGuid>(.*?)<\/.*?:SessionGuid>`)
		matches := re.FindStringSubmatch(sb)
		if len(matches) >= 2 {
			DebugLog("SessionGuid value: ", matches[1])
			return matches[1]
		} else {
			DebugLog("no SessionGuid found")
		}
	} else {
		DebugLog("request body: NONE")
	}

	return ""
}

// Function that calls the function to create unique key id
// depending on the "strategyName" in the api definition config
func selectStrategy(rateLimitingConfig RateLimitingConfig, req *http.Request) string {

	if rateLimitingConfig.RateLimiting.Active {
		name := rateLimitingConfig.RateLimiting.Strategy.Name

		switch name {
		case "requestHeaders":
			InfoLog("strategy to be applied: ", name)
			keyID := createUniqueKeyIdHeaders(rateLimitingConfig, req)

			return keyID
		case "requestHeadersXRS":
			InfoLog("strategy to be applied: ", name)
			keyID := createUniqueKeyIdHeaders(rateLimitingConfig, req)

			return keyID
		case "sessionGuid":
			InfoLog("strategy to be applied: ", name)
			keyID := createUniqueKeyIDSessionGuid(req)

			return keyID
		case "otherName":
			// Call other function based on the name

		default:
			return ("unknown strategy name: ")
		}
	}

	InfoLog("No rate limit will be applied")
	return ""
}

// * note that generating a struct at runtime isn't straightforward in Go due to its static typing.
// Instead, we might want to use reflection or code generation tools to achieve this.
// For now, we will create a map-based representation of the structure from JSON.
// Returns a struct of type RateLimitingConfig that mirrors the structure of the JSON.
// Keep in mind that using a fixed struct definition in Go is more idiomatic and recommended due to the language's strong static typing.
func generateStructFromJSON(jsonStr string) (RateLimitingConfig, error) {
	var jsonData RateLimitingConfig

	err := json.Unmarshal([]byte(jsonStr), &jsonData)
	if err != nil {
		return RateLimitingConfig{}, err
	}

	return jsonData, nil
}

func main() {}

func init() {
	DebugLog("--- Rate limiting plugin v4 init success! ---- ")
}
