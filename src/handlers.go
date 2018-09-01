package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

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

func handleStatsRequest(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]int{"last_hour": getIntervalSums(5, time.Duration(1)*time.Hour), "last_24h": getIntervalSums(5, time.Duration(24)*time.Hour), "last_7d": getIntervalSums(5, time.Duration(7*24)*time.Hour), "last_30d": getIntervalSums(5, time.Duration(30*24)*time.Hour)})
	return
}

func handleHealthCheck(c *gin.Context) {
	c.String(http.StatusOK, "")
	return
}
