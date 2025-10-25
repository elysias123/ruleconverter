package main

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	r.GET("/rule", func(c *gin.Context) {
		converter(c)
	})
}

func converter(c *gin.Context) {
	target := c.Query("target")
	origin := c.Query("origin")
	url := c.Query("url")

	if url == "" {
		c.String(400, "Missing required query parameters: url")
		return
	}
	if target == "" {
		c.String(400, "Missing required query parameters: target")
		return
	}
	if origin == "" {
		c.String(400, "Missing required query parameters: origin")
		return
	}
	rule := GetUrlContent(url)
	if rule == nil {
		c.String(400, "Failed to fetch rules from origin URL")
		return
	}

	ret := RuleConverter(rule, origin, target)
	if ret == nil {
		c.String(400, "Failed to convert rules")
		return
	}
	c.String(200, strings.Join(ret, "\n"))
}
