package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gcinnovate/go-fcapp/config"
	"github.com/gcinnovate/go-fcapp/helpers"
	"github.com/gin-gonic/gin"
	"github.com/nyaruka/phonenumbers"
)

// EmtctUpdateContactController will hold the methods to
type EmtctUpdateContactController struct{}

// EmtctUpdateContact controller handles the update_emtct_contact call
func (t *EmtctUpdateContactController) EmtctUpdateContact(c *gin.Context) {
	var payload helpers.WebHookObj

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateValue := helpers.GetFlowResult(payload.Results, "updatevalue")
	numberToUpdate := helpers.GetFlowResult(payload.Results, "numbertoupdate")
	messagingLanguage := helpers.GetFlowResult(payload.Results, "messaginglanguage")
	log.Println("Update Value: %s for %s", updateValue, numberToUpdate)

	messageType := helpers.GetFlowResult(payload.Results, "messagereceive")
	childAge := helpers.GetFlowResult(payload.Results, "childage")
	pregnancyAge := helpers.GetFlowResult(payload.Results, "pregnancyage")
	healthFacility := helpers.GetFlowResult(payload.Results, "healthfacility")
	patientArtID := helpers.GetFlowResult(payload.Results, "patientartid")

	flowUUID := helpers.GetDefaultEnv("FCAPP_EMTCT_CONTACT_UPDATE_FLOWU_ID",
		config.FcAppConf.API.EmtctUpdateContactFlowUUID)
	num, err2 := phonenumbers.Parse(numberToUpdate, "UG")
	if err2 != nil {
		log.Println("Failed to Parse Number: ", numberToUpdate)
		c.JSON(http.StatusOK, gin.H{"Error": "Failed to parse phone number"})
		return
	}
	tel := "tel:" + strings.ReplaceAll(phonenumbers.Format(num, phonenumbers.INTERNATIONAL), " ", "")
	log.Println("Formatted URN is ", tel)

	Urns := []string{tel} /*validate number to 'tel:+256...'*/
	contacts := []string{}
	reqObj := flowStartObject{
		Flow:     flowUUID,
		Urns:     Urns,
		Contacts: contacts,
		Params: map[string]string{
			"messaginglanguage": messagingLanguage,
			"childage":          childAge,
			"pregnancyage":      pregnancyAge,
			"messagereceive":    messageType,
			"patientartid":      patientArtID,
			"healthfacility":    healthFacility,
		},
	}
	var requestBody []byte
	requestBody, err := json.Marshal(reqObj)

	log.Println("Request: ", string(requestBody))

	flowStartURI := fmt.Sprintf("%s/api/v2/flow_starts.json?",
		helpers.GetDefaultEnv("FCAPP_FAMILYCONNECT_URI", config.FcAppConf.API.FamilyConnectURI))
	log.Println("Request %s", flowStartURI)

	resp, err := helpers.PostRequest(flowStartURI, requestBody)
	log.Println("Response::::", resp)
	if resp == nil || err != nil {
		log.Printf("Failed to start contact:%s in Contact Update Flow. Error:%s", Urns, err)
		c.JSON(http.StatusConflict, gin.H{"message": "Failed to Start Flow"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": updateValue})
}
