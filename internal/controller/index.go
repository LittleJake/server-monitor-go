package controller

import (
	"net/http"

	"github.com/LittleJake/server-monitor-go/internal/util"
	"github.com/gin-gonic/gin"
)

type IndexController struct{}

var Index = IndexController{}

func (IndexController) Index(c *gin.Context) {

	list, _ := util.GetCollectionStatus()
	online, _ := list.Get("online")
	offline, _ := list.Get("offline")
	info, _ := list.Get("info")

	name, _ := util.GetDisplayName(false)

	c.HTML(http.StatusOK, "index.html", gin.H{
		"base_url": util.GetEnv("BASE_URL", ""),
		"Context":  c,
		"online":   online,
		"offline":  offline,
		"info":     info,
		"name":     name,
	})
	//c.JSON(http.StatusOK, gin.H{"status": "ok", "uptime": time.Now().UTC()})
}

func (IndexController) List(c *gin.Context) {
	// judge if ajax

	list, _ := util.GetCollectionStatus()
	online, _ := list.Get("online")
	offline, _ := list.Get("offline")
	info, _ := list.Get("info")
	name, _ := util.GetDisplayName(false)

	result := gin.H{
		"base_url": util.GetEnv("BASE_URL", ""),
		"Context":  c,
		"online":   online,
		"offline":  offline,
		"info":     info,
		"name":     name,
	}

	if c.Request.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		c.HTML(http.StatusOK, "list_ajax.html", result)
		return
	}

	c.HTML(http.StatusOK, "list.html", result)
}

func (IndexController) Info(c *gin.Context) {
	uuid := c.Param("uuid")

	info, err := util.GetInfo(uuid, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve info"})
		return
	}

	latest, err := util.GetCollectionLatest(uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve latest collection data"})
		return
	}

	if c.Request.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		c.HTML(http.StatusOK, "info_ajax.html", gin.H{
			"uuid":     uuid,
			"base_url": util.GetEnv("BASE_URL", ""),
			"info":     info,
			"latest":   latest,
			"Context":  c,
		})
		return
	}

	c.HTML(http.StatusOK, "info.html", gin.H{
		"uuid":     uuid,
		"base_url": util.GetEnv("BASE_URL", ""),
		"info":     info,
		"latest":   latest,
		"Context":  c,
	})
	// c.JSON(http.StatusOK, gin.H{"id": uuid, "name": "example"})
}
