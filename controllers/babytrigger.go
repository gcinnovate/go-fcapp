package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gcinnovate/go-fcapp/config"
	"github.com/gcinnovate/go-fcapp/helpers"
	"github.com/gin-gonic/gin"
)

type flowStartObject struct {
	Flow     string            `json:"flow"`
	Contacts []string          `json:"contacts"`
	Urns     []string          `json:"urns"`
	Params   map[string]string `json:"params"`
}

// BabyTriggerController will hold the methods to
type BabyTriggerController struct{}

// BabyTrigger controller handles the start_babytrigger call
func (t *BabyTriggerController) BabyTrigger(c *gin.Context) {

	var payload helpers.WebHookObj
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	secreceivers := helpers.GetFlowResult(payload.Results, "secreceivers")
	var receivers map[string]interface{}
	if err := json.Unmarshal([]byte(secreceivers), &receivers); err != nil {
		log.Println("Error!!!")
	}

	optoutOption := helpers.GetFlowResult(payload.Results, "OptOutOption")
	contactDetails := receivers[optoutOption].(map[string]interface{})

	dateOfBirth := helpers.GetFlowResult(payload.Results, "child_dob")

	flowUUID := helpers.GetDefaultEnv("FCAPP_BABY_TRIGGER_FLOW_UUID", config.FcAppConf.API.BabyTriggerFlowUUID)
	contactUUIDs := []string{contactDetails["uuid"].(string)}
	reqObj := flowStartObject{
		Flow:     flowUUID,
		Contacts: contactUUIDs,
		Params:   map[string]string{"child_dob": dateOfBirth},
	}
	var requestBody []byte
	requestBody, err := json.Marshal(reqObj)

	if err != nil {
		log.Fatalln("JSON Marshalling failed: %s", err)
	}
	flowStartURI := fmt.Sprintf("%s/api/v2/flow_starts.json?",
		helpers.GetDefaultEnv("FCAPP_FAMILYCONNECT_URI", config.FcAppConf.API.FamilyConnectURI))

	resp, err := helpers.PostRequest(flowStartURI, requestBody)
	if resp == nil || err != nil {
		log.Printf("Failed to start contact:%s in Baby Triger Flow. Error:%s", contactDetails["uuid"], err)
		c.JSON(http.StatusConflict, gin.H{"message": "Failed to Start Flow"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "True"})
}
