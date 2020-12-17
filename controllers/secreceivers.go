package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gcinnovate/go-fcapp/db"
	"github.com/gcinnovate/go-fcapp/helpers"
	"github.com/gin-gonic/gin"
)

// SecReceiversController will hold the methods to
type SecReceiversController struct{}

type secReceiverObj struct {
	ContactID    int64  `json:"contact_id" db:"contact_id"`
	Name         string `json:"name" db:"name"`
	UUID         string `json:"uuid" db:"uuid"`
	Msisdn       string `json:"msisdn" db:"msisdn"`
	ContactField int    `json:"contact_field" db:"contact_field"`
	HasMSISDN    bool   `json:"has_msisdn" db:"has_msisdn"`
	HasHOHMSISDN bool   `json:"has_hoh_msisdn" db:"has_hoh_msisdn"`
}
type contactObj struct {
	Contact string
}

// SecondReceivers controller handles secreceivers call
func (t *SecReceiversController) SecondReceivers(c *gin.Context) {
	contact := c.Query("contact")
	babytrigger := c.Query("babytrigger")
	ct := contactObj{Contact: contact}

	sqlStmt := `SELECT * FROM fcapp_get_secondary_receivers(:contact) `
	if babytrigger == "true" {
		sqlStmt += ` WHERE has_msisdn = 'f'`
	}

	rows, err := db.GetDB().NamedQuery(sqlStmt, ct)
	if err != nil {
		log.Println("Error 1: ", err)
		c.JSON(200, gin.H{"secreceivers": "error encountered"})
		return
	}
	defer rows.Close()

	payload := make(map[string]interface{})

	receivers := make(map[string]secReceiverObj)

	receiversCount := 0
	screen1 := ""
	screen2 := ""
	screen3 := ""
	for rows.Next() {
		receiversCount++
		var secReceiver secReceiverObj
		err = rows.StructScan(&secReceiver)
		if err != nil {
			log.Println("Error: ", err)
		}
		// log.Println(secReceiver)
		switch {
		case receiversCount < 6:
			screen1 += fmt.Sprintf("%d. %s\n", receiversCount, secReceiver.Name)
			receivers[fmt.Sprintf("%d", receiversCount)] = secReceiver

		case receiversCount > 5 && receiversCount < 11:
			screen2 += fmt.Sprintf("%d. %s\n", receiversCount+1, secReceiver.Name)
			receivers[fmt.Sprintf("%d", receiversCount+1)] = secReceiver

		case receiversCount > 10 && receiversCount < 16:
			screen3 += fmt.Sprintf("%d. %s\n", receiversCount+2, secReceiver.Name)
			receivers[fmt.Sprintf("%d", receiversCount+2)] = secReceiver
		}
	}
	payload["receivercount"] = receiversCount
	payload["screen_1"] = screen1
	payload["screen_2"] = screen2
	payload["screen_3"] = screen3
	payload["secreceivers"] = receivers
	// log.Println(payload)

	c.JSON(200, payload)
}

// OptOutSecReceiverController will hold methods to
type OptOutSecReceiverController struct{}

// OptOutSecReceiver controller handles the optout_secondaryreceiver call
func (t *OptOutSecReceiverController) OptOutSecReceiver(c *gin.Context) {
	var payload helpers.WebHookObj
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	secreceivers := helpers.GetFlowResult(payload.Results, "secreceivers")
	/*
		secreceivers := `{
				"1": {
				"contact_id": 3,
				"name": "Tendo Nabaka",
				"uuid": "b2b831e8-abcf-48ac-9d8b-4fb9024c3ff9",
				"msisdn": "0753475676",
				"contact_field": 25,
				"has_msisdn": true,
				"has_hoh_msisdn": false
				},
				"2": {
				"contact_id": 2,
				"name": "Anita Sekiwere",
				"uuid": "db89dcfd-676d-459f-b52c-c3d07f4d6025",
				"msisdn": "0753475676",
				"contact_field": 15,
				"has_msisdn": true,
				"has_hoh_msisdn": true
				}
		}`
	*/
	// var receivers map[string]map[string]map[string]string
	var receivers map[string]interface{}
	if err := json.Unmarshal([]byte(secreceivers), &receivers); err != nil {
		log.Println("Error!!!")
	}
	log.Println(receivers)

	optoutall := c.Query("optoutall")
	if optoutall == "true" {
		log.Println("You wanna opt out all secreceivers")

		for key, elem := range receivers {
			log.Println("Key:", key, "=>", "Element:", elem)
			_, err := db.GetDB().NamedExec(
				`
				SELECT fcapp_delete_contactfield_by_id(:contact_id, :contact_field);`, elem)
			log.Println(err)
		}

		c.JSON(http.StatusOK, gin.H{"message": "You wanna opt out all secreceivers"})
		return
	}

	optoutOption := helpers.GetFlowResult(payload.Results, "OptOutOption")
	if len(optoutOption) > 0 {
		_, err := db.GetDB().NamedExec(
			`
		SELECT fcapp_delete_contactfield_by_id(:contact_id, :contact_field);`,
			receivers[optoutOption])
		if err != nil {
			c.JSON(200, gin.H{"message": "Failed to optout option"})
			return
		}
		c.JSON(200, gin.H{"message": "Contact opted out successfully!"})
		return
	}

	c.JSON(200, gin.H{"message": "No contact opted out"})
	return
}
