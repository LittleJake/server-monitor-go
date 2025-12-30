package api

import (
	"net/http"
	"strconv"

	"github.com/LittleJake/server-monitor-go/internal/util"
	"github.com/elliotchance/orderedmap/v3"
	"github.com/gin-gonic/gin"
)

type MemoryAPI struct{}

var Memory = MemoryAPI{}

func (MemoryAPI) Get(c *gin.Context) {
	uuid, ok := c.Params.Get("uuid")

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "uuid parameter is required",
		})
		return
	}

	startTime, err3 := strconv.ParseInt(c.Query("start"), 10, 64)
	endTime, err4 := strconv.ParseInt(c.Query("end"), 10, 64)

	var result *orderedmap.OrderedMap[int64, util.CollectionData]
	var err error
	// result, err := util.GetCollection(uuid, false)
	if err3 == nil && err4 == nil {
		result, err = util.GetCollectionByTime(uuid, false, startTime, endTime)
	} else {
		result, err = util.GetCollection(uuid, false)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	memory := util.CollectionFormat(result, "Memory")
	c.JSON(http.StatusOK, memory)
}
