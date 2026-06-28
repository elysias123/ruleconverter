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

// maxURLs 限制单次请求可携带的 URL 数量
const maxURLs = 5

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

	// 拆分并清洗 URL 列表
	var urls []string
	for _, u := range strings.Split(url, ",") {
		if u = strings.TrimSpace(u); u != "" {
			urls = append(urls, u)
		}
	}
	if len(urls) == 0 {
		c.String(400, "origin URL parse failed")
		return
	}
	if len(urls) > maxURLs {
		c.String(400, "too many URLs: %d (max %d)", len(urls), maxURLs)
		return
	}

	var Ruleslist []string
	for _, u := range urls {
		rule, err := GetUrlContent(u)
		if err != nil {
			// 对外返回模糊信息，避免泄露内网拓扑/解析细节；详情只写日志
			c.String(400, "origin URL fetch failed")
			log.Printf("origin URL fetch failed: %v", err)
			return
		}
		Ruleslist = append(Ruleslist, rule...)
	}

	ret := RuleConverter(Ruleslist, origin, target)
	if ret == nil {
		c.String(400, "Failed to convert rules")
		return
	}
	c.String(200, strings.Join(ret, "\n"))
}
