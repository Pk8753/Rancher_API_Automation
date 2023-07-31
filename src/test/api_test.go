package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"crypto/tls"
	"github.com/xeipuuv/gojsonschema"
	"github.com/stretchr/testify/assert"
	"encoding/base64"
)

type TestInput struct {
	EndpointURL       string                 `json:"endpointURL"`
	InvalidEndpointURL string                 `json:"invalidEndpointURL"`
	UserName string                				 `json:"username"`
	Password string                 				`json:"password"`
	InvalidUserName string                				 `json:"invalidUsername"`
	InvalidPassword string 							`json:"invalidPassword"` 
	TestPayload       map[string]interface{} `json:"testPayload"`
	TestInvalidPayload map[string]interface{} `json:"testInvalidPayload"`
}
func loadTestInput() (*TestInput, error) {
	file, err := ioutil.ReadFile("test_input.json")
	if err != nil {
		return nil, err
	}

	var input TestInput
	err = json.Unmarshal(file, &input)
	if err != nil {
		return nil, err
	}

	return &input, nil
}


func sendRequest(URL string,username string, password string,payload map[string]interface{}) (*http.Response, error) {

	payloadBytes, err := json.Marshal(payload)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request payload: %v", err)
    }

    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    httpClient := &http.Client{Transport: tr}

    req, err := http.NewRequest("GET", URL, bytes.NewBuffer(payloadBytes))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %v", err)
    }

    // Set basic authentication header
    auth := username + ":" + password
    basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
    req.Header.Set("Authorization", basicAuth)
    req.Header.Set("Content-Type", "application/json")

    return httpClient.Do(req)
}

// Test function to validate the API status code is 201
func TestLoginIsSuccessfulWithStatusCode200(t *testing.T) {

	testInput, err := loadTestInput()
	if err != nil {
		t.Fatalf("Failed to load test input: %v", err)
	}
	
	resp, err := sendRequest(testInput.EndpointURL,testInput.UserName,testInput.Password,testInput.TestPayload)
	if err != nil {
		t.Fatalf("API call failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code 201, but got %d", resp.StatusCode)
	}
}

func TestLoginFailureWithStatusCode401(t *testing.T) {

	testInput, err := loadTestInput()
	if err != nil {
		t.Fatalf("Failed to load test input: %v", err)
	}

	resp, err := sendRequest(testInput.EndpointURL,testInput.UserName,testInput.InvalidPassword,testInput.TestInvalidPayload)
	if err != nil {
		t.Fatalf("API call failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status code 401, but got %d", resp.StatusCode)
	}
}


func TestResponseBodyContainsDescription(t *testing.T) {

	testInput, err := loadTestInput()
	if err != nil {
		t.Fatalf("Failed to load test input: %v", err)
	}

	resp, err := sendRequest(testInput.EndpointURL,testInput.UserName,testInput.Password,testInput.TestPayload)
	if err != nil {
		t.Fatalf("API call failed: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	assert.Contains(t, string(responseBody), `"token-s66gw"`, "response body should contain expected description")
}

// Test function to validate the response status is 401 when invalid username or password is in payload
func TestInvalidUsernameOrPassword(t *testing.T) {

	testInput, err := loadTestInput()
	if err != nil {
		t.Fatalf("Failed to load test input: %v", err)
	}
	
	resp, err := sendRequest(testInput.EndpointURL,testInput.InvalidUserName,testInput.Password,testInput.TestInvalidPayload)
	if err != nil {
		t.Fatalf("API call failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status code 401, but got %d", resp.StatusCode)
	}
}

// Test function to validate the response status is 422 when invalid action is passed in URL
func TestInvalidActionInURL(t *testing.T) {
	

	testInput, err := loadTestInput()
	if err != nil {
		t.Fatalf("Failed to load test input: %v", err)
	}

	resp, err := sendRequest(testInput.InvalidEndpointURL,testInput.UserName,testInput.Password,testInput.TestPayload)
	if err != nil {
		t.Fatalf("API call failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status code 422, but got %d", resp.StatusCode)
	}
}

func TestResponseSchemaValidation(t *testing.T) {

	testInput, err := loadTestInput()
	if err != nil {
		t.Fatalf("Failed to load test input: %v", err)
	}
	// Load the JSON schema from file
	schemaFile := "./response_schema.json"
	schemaBytes, err := ioutil.ReadFile(schemaFile)
	if err != nil {
		t.Fatalf("Failed to read JSON schema: %v", err)
	}

	// Create a JSON loader for the schema
	schemaLoader := gojsonschema.NewBytesLoader(schemaBytes)

	resp, err := sendRequest(testInput.EndpointURL,testInput.UserName,testInput.Password,testInput.TestPayload)
	if err != nil {
		t.Fatalf("API call failed: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}	

	// fmt.Printf("API Response Body: %s\n", responseBody)

	// Create a JSON loader for the API response
	responseLoader := gojsonschema.NewBytesLoader([]byte(responseBody))

	// Perform the schema validation
	result, err := gojsonschema.Validate(schemaLoader, responseLoader)
	if err != nil {
		t.Fatalf("Schema validation failed: %v", err)
	}

	// Check if the validation result is valid
	if !result.Valid() {
		for _, desc := range result.Errors() {
			fmt.Printf("Schema validation error: %s\n", desc)
		}
		t.Error("Response does not match the expected schema")
	}

	// Unmarshal the "userPrincipal" field as a JSON object
	var userPrincipal map[string]interface{}
	if err := json.Unmarshal([]byte(responseBody), &userPrincipal); err != nil {
		t.Fatalf("Failed to unmarshal userPrincipal field: %v", err)
	}
}



// TestMain function to run all the test functions and generate a report
func TestMain(m *testing.M) {
	fmt.Println("Starting API tests...")
	code := m.Run()
	fmt.Println("Finished API tests.")
	fmt.Printf("Test exit code: %d\n", code)
}
