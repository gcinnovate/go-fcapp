package helpers

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"os"
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
	tokenURL := GetDefaultEnv("NITAU_API_ROOT_URI", "https://msdg.uconnect.go.ug/api/v1") + "/get-jwt-token/"
	requestBody, err := json.Marshal(map[string]string{
		"userid":   GetDefaultEnv("NITAU_API_USER", ""),
		"password": GetDefaultEnv("NITAU_API_PASSWORD", ""),
		"email":    GetDefaultEnv("NITAU_API_EMAIL", ""),
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
