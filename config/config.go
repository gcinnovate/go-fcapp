package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// FcAppConf is the global conf
var FcAppConf Config

func init() {
	args := ProcessArgs(&FcAppConf)

	// read configuration from the file and environment variables
	if err := cleanenv.ReadConfig(args.ConfigPath, &FcAppConf); err != nil {
		fmt.Println(err, " Please make sure you have a config.yml file in app directory")
		os.Exit(2)
	}
}

// Config is the top level cofiguration object
type Config struct {
	API struct {
		AuthToken                  string `yaml:"authtoken" env:"FCAPP_AUTH_TOKEN" env-description:"RapidPro API Token"`
		RootURI                    string `yaml:"rooturi" env:"FCAPP_ROOT_URI" env-description:"API ROOT URI"`
		SmsURL                     string `yaml:"smsurl" env:"FCAPP_SMSURL" env-description:"API SMS endpoint"`
		SmsCode                    string `yaml:"sms_code" env:"FCAPP_SMS_CODE" env-description:" SMS Short-code"`
		SmsUser                    string `yaml:"sms_user" env:"FCAPP_SMS_USER" env-description:"User for SMS endpoint"`
		SmsPassword                string `yaml:"sms_password" env:"FCAPP_SMS_PASSWORD" env-description:"Password for SMS endpoint"`
		PrebirthCampaign           string `yaml:"prebirth_campaign" env:"FCAPP_PREBIRTH_CAMPAIGN" env-description:"Prebirth Campaign UUID"`
		PostbirthCampaign          string `yaml:"postbirth_campaign" env:"FCAPP_POSTBIRTH_CAMPAIGN" env-description:"Postbirth Campaign UUID"`
		FamilyConnectURI           string `yaml:"familyconnect_uri" env:"FCAPP_FAMILYCONNECT_URI" env-description:"FamilyConnect URI"`
		BabyTriggerFlowUUID        string `yaml:"babytrigger_flow_uuid" env:"FCAPP_BABY_TRIGGER_FLOW_UUID" env-description:"FamilyConnect Baby Trigger Flow UUID"`
		EmtctUpdateContactFlowUUID string `yaml:"emtct_updatecontact_flow_uuid" env:"FCAPP_EMTCT_UPDATECONTACT_FLOWUUID" env-description:"FC Contact Update Flow UUID"`
	} `yaml:"api"`
	Server struct {
		Port string `yaml:"port" env:"FCAPP_SERVER_PORT" env-description:"Server port" env-default:"5004"`
	} `yaml:"server"`
	Database struct {
		URI string `yaml:"uri" env:"FCAPP_DB" env-default:"postgres://postgres:postgres@localhost:5431/temba_latest?sslmode=disable"`
	} `yaml:"database"`
}

// Args command-line parameters
type Args struct {
	ConfigPath string
}

// ProcessArgs processes and handles CLI arguments
func ProcessArgs(cfg interface{}) Args {
	var a Args

	f := flag.NewFlagSet("Example server", 1)
	f.StringVar(&a.ConfigPath, "c", "config.yml", "Path to config file")

	fu := f.Usage
	f.Usage = func() {
		fu()
		envHelp, _ := cleanenv.GetDescription(cfg, nil)
		fmt.Fprintln(f.Output())
		fmt.Fprintln(f.Output(), envHelp)
	}

	f.Parse(os.Args[1:])
	return a
}
