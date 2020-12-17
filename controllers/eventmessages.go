package controllers

import (
	"log"

	"github.com/gcinnovate/go-fcapp/db"
	"github.com/gcinnovate/go-fcapp/helpers"
	"github.com/gin-gonic/gin"
)

// EventMessageController will hold the methods to
type EventMessageController struct{}

type eventMsgObject struct {
	Offset       string `json:"offset" binding:"required"`
	Language     string `json:"lang" binding:"required"`
	CampaignType string `json:"campaign_type" binding:"required"`
}

// Default controller handles test calls
func (t *EventMessageController) Default(c *gin.Context) {
	offset := c.Query("offset")
	lang := c.Query("lang")
	campaignType := c.Query("campaign_type")

	eventMsgObj := eventMsgObject{
		Offset:       offset,
		Language:     lang,
		CampaignType: campaignType,
	}

	// log.Println("[Offset: ", eventMsgObj.Offset, "] [Lan:", eventMsgObj.Language, "]")
	var campaignUUID string
	if eventMsgObj.CampaignType == "prebirth" {
		campaignUUID = helpers.GetDefaultEnv("FCAPP_PREBIRTH_CAMPAIGN", "")
	} else {
		campaignUUID = helpers.GetDefaultEnv("FCAPP_POSTBIRTH_CAMPAIGN", "")
	}

	m := map[string]interface{}{
		"lang":     eventMsgObj.Language,
		"offset":   eventMsgObj.Offset,
		"campaign": campaignUUID,
	}
	rows, err := db.GetDB().NamedQuery(`
	SELECT 
		CASE WHEN exist(message, :lang) AND length(message->:lang) > 1 THEN
			message->:lang ELSE message->'eng' END AS message
	FROM 
		campaigns_campaignevent 
	WHERE
		campaign_id = (SELECT id FROM campaigns_campaign WHERE uuid=:campaign)
		AND "offset"=:offset AND is_active='t'
	`, m)

	defer rows.Close()
	if !rows.Next() {
		log.Printf("No message found")
		c.JSON(200, gin.H{"message": ""})
		return
	}

	var msg string
	err = rows.Scan(&msg)
	if err != nil {
		log.Println("ERROR: ", err)
		c.JSON(200, gin.H{"message": ""})
		return
	}

	c.JSON(200, gin.H{"message": msg})
}
