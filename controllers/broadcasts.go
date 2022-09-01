package controllers

import (
	"github.com/gcinnovate/go-fcapp/helpers"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type broadcastObject struct {
	Contacts []string `json:"contacts"`
	text     string   `json:"text"`
	Urns     []string `json:"urns"`
}

// BroadcastController will hold the methods to
type BroadcastController struct{}

// Broadcast controller is used to make a rapidpro broadcast
func (t *BroadcastController) Broadcast(c *gin.Context) {

	var payload helpers.WebHookObj
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contact := payload.Contact.UUID
	contactURN := []string{strings.TrimPrefix(payload.Contact.Urn, "tel:+")}

	feedbackMessage := helpers.GetFlowResult(payload.Results, "feedback_message")
	if len(feedbackMessage) == 0 {
		log.Printf("Couldn't send empty broadcast for contact:%s.", contact)
		c.JSON(http.StatusOK, gin.H{"success": "False", "message": "broadcast message was empty!"})
		return
	}

	var (
		sendURLs = []string{
			"http://smgw1.yo.co.ug:9100/sendsms",
		}
	)
	form := url.Values{
		"origin":       []string{helpers.GetDefaultEnv("FCAPP_SMS_CODE", "8900")},
		"sms_content":  []string{feedbackMessage},
		"destinations": contactURN,
		"ybsacctno":    []string{helpers.GetDefaultEnv("FCAPP_SMS_USER", "")},
		"password":     []string{helpers.GetDefaultEnv("FCAPP_SMS_PASSWORD", "")},
	}
	for _, sendURL := range sendURLs {
		sendURL, _ := url.Parse(sendURL)
		sendURL.RawQuery = form.Encode()

		req, _ := http.NewRequest(http.MethodGet, sendURL.String(), nil)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		rr, err := helpers.MakeHTTPRequest(req)

		if err != nil {
			log.Printf("Message Sent with error %s %s", rr, err)
			continue
		}
		responseQS, _ := url.ParseQuery(string(rr.Body))
		// check whether we were blacklisted
		createMessage, _ := responseQS["ybs_autocreate_message"]
		if len(createMessage) > 0 && strings.Contains(createMessage[0], "BLACKLISTED") {
			log.Printf("Message sending failed: we were BLACKLISTED")
			c.JSON(http.StatusConflict, gin.H{"message": "Failed to start broadcast"})
			return
		}
		// finally check that we were sent
		createStatus, _ := responseQS["ybs_autocreate_status"]
		if len(createStatus) > 0 && createStatus[0] == "OK" {
			c.JSON(http.StatusOK, gin.H{
				"status": "SUCCESS", "message": "Message Sent successfully"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": "True"})
	return
}
