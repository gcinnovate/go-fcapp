# go-fcapp

# About 

go-fcapp is a Golang application for managing FamilyConnect wehbhook calls 
# Deploying

As go-fcapp is a go application, it compiles to a binary and that binary along with the config file is all
you need to run it on your server. You can find bundles for each platform in the
[releases directory](https://github.com/gcinnovate/go-fcapp/releases). We recommend running go-fcapp
behind a reverse proxy such as nginx or Elastic Load Balancer that provides HTTPs encryption.

# Configuration

go-fcapp uses a tiered configuration system, each option takes precendence over the ones above it:
 1. The configuration file
 2. Environment variables starting with `FCAPP_` 
 3. Command line parameters

We recommend running go-fcapp with no changes to the configuration and no parameters, using only
environment variables to configure it. You can use `% go-fcapp --help` to see a list of the
environment variables and parameters and for more details on each option.

#  Configuration

For use with your RapidPro instance, you will want to configure these settings:

 * `FCAPP_AUTH_TOKEN`: The FamilyConnect RapidPro authorization Token
 * `FCAPP_PREBIRTH_CAMPAIGN`: The UUID for the prebirth campaign with message events for all languages
 * `FCAPP_POSTBIRTH_CAMPAIGN`: The UUID for the postbirth campaign with message events for all languages
 * `FCAPP_FAMILYCONNECT_URI`: The FamilyConnect URI (default: `https://app.familyconnect.go.ug`)
 * `FCAPP_ROOT_URI`: The root uri for the api endpoints for the FamilyConnect's RapidPro
 * `FCAPP_BABY_TRIGGER_FLOW_UUID`: The UUID for the flow in FamilyConnect that maps to the baby trigger
 * `FCAPP_SERVER_PORT`: The port on which to run the go-fcapp service
 * `FCAPP_SMSURL`: The SMSURL to use for sending SMS through this app 

 
# Development

Install go-fcapp source in your workspace with:

```
go get github.com/gcinnovate/go-fcapp
```

Build go-fcapp with:

```
go install github.com/gcinnovate/go-fcapp/..
```

This will create a new executable in $GOPATH/bin called `go-fcapp`. 
