package openaiclient

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

var ApiKey string = ""

func SetApiKey(key string) {
	ApiKey = key
}

func CallChatGPT(prompt string) string {

	// Call the ChatGPT endpoint and handle the response
	apiKey := ApiKey
	url := "https://api.openai.com/v1/completions"

	headers := map[string]string{
		"Authorization": "Bearer " + apiKey,
		"Content-Type":  "application/json",
	}

	data := map[string]interface{}{
		"prompt":      prompt,
		"model":       "text-davinci-003",
		"max_tokens":  150,
		"temperature": 0,
	}

	resp, err := postRequest(url, headers, data)
	if err != nil {
		return "Failed to make the ChatGPT request"
	}

	// Extract the response from the API
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "Failed to parse the ChatGPT response"
	}

	gptResponse := result["choices"].([]interface{})[0].(map[string]interface{})["text"].(string)
	return gptResponse
}

type ImageData struct {
	Data []json.RawMessage `json:"data"`
}

func CallImageGPT(prompt string, n int, size string) []json.RawMessage {

	// Get the API key
	apiKey := ApiKey

	// URL for the endpoint
	url := "https://api.openai.com/v1/images/generations"

	// Headers for the request
	headers := map[string]string{
		"Authorization": "Bearer " + apiKey,
		"Content-Type":  "application/json",
	}

	// Data for the request
	data := map[string]interface{}{
		"prompt": prompt,
		"n":      n,
		"size":   size,
	}

	// Send the POST request
	resp, err := postRequest(url, headers, data)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	// Unmarshal the JSON response
	var imageData ImageData
	if err := json.Unmarshal(body, &imageData); err != nil {
		return nil
	}

	return imageData.Data
}

// postRequest sends a POST request with the provided data and headers
func postRequest(url string, headers map[string]string, data map[string]interface{}) (*http.Response, error) {
	client := &http.Client{}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
