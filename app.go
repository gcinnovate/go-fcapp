package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jasonlvhit/gocron"
	// Import godotenv for .env variables
	"github.com/gcinnovate/go-fcapp/config"
	"github.com/gcinnovate/go-fcapp/controllers"
	"github.com/gcinnovate/go-fcapp/db"
	"github.com/gcinnovate/go-fcapp/helpers"
	"github.com/gcinnovate/go-fcapp/models"
	"github.com/joho/godotenv"
)

func init() {
	// Log error if .env file does not exist
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found")
	}
	// helpers.RefreshToken()
}

func myTask() {
	fmt.Println("This task will run periodically")
}

func executeCronJob() {
	// gocron.Every(1).Minute().Do(helpers.RefreshToken)
	gocron.Every(1).Day().At("23:00").Do(helpers.SynchronizeCHWs)
	<-gocron.Start()
}

func main() {

	log.Printf("Token: %s", os.Getenv("FCAPP_AUTH_TOKEN"))

	go executeCronJob()
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		pl := new(controllers.EventMessageController)
		v1.GET("/eventmessage", pl.Default)

		v1.POST("/subregion_districts/:region", func(c *gin.Context) {
			region := c.Param("region")
			c.JSON(http.StatusOK, db.GetRegionDistricts()[region])
		})

		v1.POST("/district_subcounties/:district", func(c *gin.Context) {
			district := c.Param("district")
			c.JSON(http.StatusOK, db.GetDistrictSubcounties()[district])
		})
	}
	// v2 := router.Use()
	authorized := router.Group("/api/v1", basicAuth())
	{
		sr := new(controllers.SecReceiversController)
		authorized.GET("/secreceivers", sr.SecondReceivers)

		op := new(controllers.OptOutSecReceiverController)
		authorized.POST("/optout_secondaryreceiver", op.OptOutSecReceiver)

		bt := new(controllers.BabyTriggerController)
		authorized.POST("/startbabytriggerflow", bt.BabyTrigger)
	}

	// Handle error response when a route is not defined
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Not found"})
	})

	conf := config.FcAppConf
	// Init our Server
	router.Run(":" + conf.Server.Port)
}

func basicAuth() gin.HandlerFunc {

	return func(c *gin.Context) {
		auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)

		if len(auth) != 2 || auth[0] != "Basic" {
			respondWithError(401, "Unauthorized", c)
			return
		}
		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)

		if len(pair) != 2 || !authenticateUser(pair[0], pair[1]) {
			respondWithError(401, "Unauthorized", c)
			// c.Writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			return
		}

		c.Next()
	}
}

func authenticateUser(username, password string) bool {
	// log.Printf("Username:%s, password:%s", username, password)
	userObj := models.User{}
	err := db.GetDB().QueryRowx(
		"SELECT id, username FROM fcapp_users "+
			"WHERE username = $1 AND password = crypt($2, password) ", username, password).StructScan(&userObj)
	if err != nil {
		fmt.Printf("User:[%v]", err)
		return false
	}
	// fmt.Printf("User:[%v]", userObj)
	return true
}

func respondWithError(code int, message string, c *gin.Context) {
	resp := map[string]string{"error": message}

	c.JSON(code, resp)
	c.Abort()
}
