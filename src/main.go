package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/cnjack/throttle"
	"github.com/gin-gonic/gin"
)

func handleHeaderRequest(c *gin.Context) {
	requestHeaders := c.Request.Header
	revisedRequestHeaders := make(map[string]string)
	for key, values := range requestHeaders {
		var joinedValues string
		for _, value := range values {
			joinedValues += fmt.Sprintf("%s", value)
		}
		revisedRequestHeaders[key] = joinedValues
	}
	revisedRequestHeadersResponse, _ := json.Marshal(revisedRequestHeaders)
	c.String(http.StatusOK, string(revisedRequestHeadersResponse))
	return
}

func handleHeaderVerification(c *gin.Context) {
	expectedHeaderKey := c.Query("key")
	expectedHeaderValue := c.Query("value")
	if expectedHeaderKey == "" || expectedHeaderValue == "" {
		c.String(http.StatusBadRequest, "Expected header key or header value not provided")
		return
	}
	requestHeaders := c.Request.Header

	headerValueRegex, err := regexp.Compile(expectedHeaderValue)
	if err != nil {
		c.String(http.StatusBadRequest, "There is a problem with your regexp.")
		return
	}

	if headerValueRegex.MatchString(requestHeaders.Get(expectedHeaderKey)) == true {
		c.JSON(http.StatusOK, map[string]bool{"result": true})
		return
	} else {
		c.JSON(http.StatusOK, map[string]bool{"result": false})
		return
	}
}

func handleHealthCheck(c *gin.Context) {
	c.String(http.StatusOK, "")
	return
}

var hostFlag string
var portFlag string
var throttleFlag int64
var behindProxy bool

func init() {
	flag.StringVar(&hostFlag, "host", "127.0.0.1", "Host the server should run on")
	flag.StringVar(&portFlag, "port", "9000", "Port the server should run on")
	flag.Int64Var(&throttleFlag, "throttle", 60, "Requests per minute allowed from IP")
	flag.BoolVar(&behindProxy, "proxy", false, "Whether the server is behind a proxy")
	flag.Parse()
}

func main() {
	router := gin.Default()
	router.Use(throttle.Policy(&throttle.Quota{
		Limit:  uint64(throttleFlag),
		Within: time.Minute,
	}))
	router.GET("/", handleHeaderRequest)
	router.GET("/verify", handleHeaderVerification)
	router.GET("/health", handleHealthCheck)
	connString := fmt.Sprintf("%s:%s", hostFlag, portFlag)
	router.Run(connString)
}
