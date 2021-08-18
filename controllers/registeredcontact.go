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

// RegisteredContactController will hold the methods to
type RegisteredContactController struct{}

// ContactRegistered controller handles the contact registered call
func (t *RegisteredContactController) ContactRegistered(c *gin.Context) {
	var payload helpers.WebHookObj

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	numberToUpdate := helpers.GetFlowResult(payload.Results, "numbertoupdate")

	num, err2 := phonenumbers.Parse(numberToUpdate, "UG")
	if err2 != nil {
		log.Println("Failed to Parse Number: ", numberToUpdate)
		c.JSON(http.StatusOK, gin.H{"Error": "Failed to parse phone number"})
		return
	}
	tel := "tel:" + strings.ReplaceAll(phonenumbers.Format(num, phonenumbers.INTERNATIONAL), " ", "")
	log.Println("Formatted URN is ", tel)

	log.Println("MSISDN: ", tel)

	contactsURI := fmt.Sprintf("%s/api/v2/contacts.json?urn=%s",
		helpers.GetDefaultEnv("FCAPP_FAMILYCONNECT_URI", config.FcAppConf.API.FamilyConnectURI), tel)
	log.Println("Request %s", contactsURI)

	resp, err := helpers.GetRequest(contactsURI)
	log.Println("Response::::", resp)
	isRegistered := false
	if resp == nil || err != nil {
		log.Printf("Failed to get contact:%s Error:%s", tel, err)
		isRegistered = false
		c.JSON(http.StatusConflict, gin.H{"message": "Failed to read contact", "isRegistered": isRegistered})
		return
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	results, ok := result["results"]
	if ok {
		// log.Printf("%s", results)
		for _, res := range results.([]interface{}) {
			grp := res.(map[string]interface{})
			groups := grp["groups"]
			for _, g := range groups.([]interface{}) {
				currentGroup := g.(map[string]interface{})
				if currentGroup["name"] == "All FC-EMTCT" {
					isRegistered = true
				}
				log.Printf(">>>>>%v, %T", currentGroup["name"], currentGroup["name"])
			}
		}
		c.JSON(200, gin.H{"message": "We have some results", "isRegistered": isRegistered})
	}
	return

}
