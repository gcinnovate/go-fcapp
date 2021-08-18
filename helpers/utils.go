package helpers

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gcinnovate/go-fcapp/config"
)

// GetDefaultEnv Returns default value passed if env variable not defined
func GetDefaultEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// SynchronizeCHWs synchronizes CHWs from iHRIS
func SynchronizeCHWs() {
	tokenURL := GetDefaultEnv("FCAPP_ROOT_URI", "https://msdg.uconnect.go.ug/api/v1") + "/get-jwt-token/"
	requestBody, err := json.Marshal(map[string]string{
		"email": GetDefaultEnv("FCAPP_EMAIL", ""),
	})

	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("URI: %s\n", tokenURL)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, err := http.Post(tokenURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil && resp == nil {
		// log.Fatalf("Token Refresh Error. %+v", err)
		log.Printf("Get New Token Error. %+v", err)
	} else {
		var result map[string]interface{}
		// body, err := ioutil.ReadAll(resp.Body)
		json.NewDecoder(resp.Body).Decode(&result)
		log.Println("Refreshed Token: ", result["token"])
		token, ok := result["token"].(string)
		if ok {
			os.Setenv("NITAU_API_AUTH_TOKEN", token)
			log.Println(os.Getenv("NITAU_API_AUTH_TOKEN"))
		}
	}
}

func authenticateUser() bool {
	return false
}

// FlowObj is part of the JSON payload passed via a webhook
type FlowObj struct {
	Name string `json:"name"`
	UUID string `json:"uuid"`
}

// ContactObj is part of the payload from RapidPro webhook
type ContactObj struct {
	Name string `json:"name"`
	UUID string `json:"uuid"`
	Urn  string `json:"urn"`
}

// WebHookObj is the object that holds the payload from RapidPro
type WebHookObj struct {
	Flow    FlowObj                      `json:"flow"`
	Contact ContactObj                   `json:"contact"`
	Results map[string]map[string]string `json:"results"`
}

// GetFlowResult return the value of a flow result given the webhook results payload
func GetFlowResult(results map[string]map[string]string, msg string) string {
	if elem, ok := results[msg]; ok == true {
		return elem["value"]
	}
	return ""
}

// PostRequest posts a request to the RapidPro API
func PostRequest(url string, data []byte) (*http.Response, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	r, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	token := GetDefaultEnv("FCAPP_AUTH_TOKEN", config.FcAppConf.API.AuthToken)
	r.Header.Add("Authorization", "Token "+token)
	r.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(r)
	if err != nil && resp == nil {
		log.Printf("Failed to make post call to RapidPro")
		return nil, err
	}
	return resp, nil
}

// GetRequest gets a request to the RapidPro API
func GetRequest(url string) (*http.Response, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	r, _ := http.NewRequest(http.MethodGet, url, nil)
	token := GetDefaultEnv("FCAPP_AUTH_TOKEN", config.FcAppConf.API.AuthToken)
	r.Header.Add("Authorization", "Token "+token)
	r.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(r)
	if err != nil && resp == nil {
		log.Printf("Failed to make get call to RapidPro")
		return nil, err
	}
	return resp, nil
}
