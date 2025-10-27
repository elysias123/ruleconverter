package main

import (
	"log"
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
	var err error
	var rule []string
	var Ruleslist []string

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

	if !strings.Contains(url, ",") {
		Ruleslist, err = GetUrlContent(url)
		if err != nil {
			c.String(400, "origin URL: %s fetch failed: %v", url, err)
			log.Panicf("origin URL: %s fetch failed: %v", url, err)
			return
		}
	} else {
		urls := strings.Split(url, ",")
		if len(urls) == 0 {
			c.String(400, "origin URL parse failed")
			return
		}
		for _, u := range urls {
			rule, err = GetUrlContent(u)
			if err != nil {
				c.String(400, "origin URL: %s fetch failed", u)
				log.Panicf("origin URL: %s fetch failed: %v", u, err)
				continue
			}
			Ruleslist = append(Ruleslist, rule...)
		}
	}

	ret := RuleConverter(Ruleslist, origin, target)
	if ret == nil {
		c.String(400, "Failed to convert rules")
		return
	}
	c.String(200, strings.Join(ret, "\n"))
}
