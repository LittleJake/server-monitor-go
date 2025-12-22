package util

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/elliotchance/orderedmap/v3"
	"github.com/redis/go-redis/v9"
)

type CollectionData struct {
	Battery []interface{} `json:"Battery"`
	Disk    struct{}      `json:"Disk"`
	Fan     []interface{} `json:"Fan"`
	Io      struct {
		Read struct {
			Count int `json:"count"`
			Bytes int `json:"bytes"`
			Time  int `json:"time"`
		} `json:"read"`
		Write struct {
			Count int `json:"count"`
			Bytes int `json:"bytes"`
			Time  int `json:"time"`
		} `json:"write"`
	} `json:"IO"`
	Load struct {
		User      float64 `json:"user"`
		Nice      int     `json:"nice"`
		System    float64 `json:"system"`
		Idle      float64 `json:"idle"`
		Iowait    float64 `json:"iowait"`
		Irq       int     `json:"irq"`
		Softirq   float64 `json:"softirq"`
		Steal     int     `json:"steal"`
		Guest     int     `json:"guest"`
		GuestNice int     `json:"guest_nice"`
	} `json:"Load"`
	Memory struct {
		Mem struct {
			Total   string  `json:"total"`
			Used    string  `json:"used"`
			Free    string  `json:"free"`
			Percent float64 `json:"percent"`
		} `json:"Mem"`
		Swap struct {
			Total   string  `json:"total"`
			Used    string  `json:"used"`
			Free    string  `json:"free"`
			Percent float64 `json:"percent"`
		} `json:"Swap"`
	} `json:"Memory"`
	Network struct {
		Rx struct {
			Bytes   int `json:"bytes"`
			Packets int `json:"packets"`
		} `json:"RX"`
		Tx struct {
			Bytes   int `json:"bytes"`
			Packets int `json:"packets"`
		} `json:"TX"`
	} `json:"Network"`
	Ping    []interface{} `json:"Ping"`
	Thermal struct{}      `json:"Thermal"`
}

func UnmarshalJSONData(data string) (*CollectionData, error) {
	var d CollectionData
	err := json.Unmarshal([]byte(data), &d)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON data: %w", err)
	}
	return &d, nil
}

func GetCollection(uuid string) (*orderedmap.OrderedMap[int64, CollectionData], error) {
	if CollectionCache == nil {
		return nil, fmt.Errorf("LocalCacheClient is not initialized")
	}

	if CollectionCache.Get("system_monitor:collection:"+uuid) != nil {
		return CollectionCache.Get("system_monitor:collection:" + uuid).Value(), nil
	}

	orderedMap := orderedmap.NewOrderedMap[int64, CollectionData]()

	data, err := RedisZRangeByScoreWithScores(
		context.Background(),
		RedisClient,
		"system_monitor:collection:"+uuid,
		&redis.ZRangeBy{Min: "0", Max: fmt.Sprint(time.Now().Unix())},
	)
	if err != nil {
		return nil, err
	}

	for _, item := range data {
		d, err := UnmarshalJSONData(item.Member.(string))
		if err != nil {
			continue
		}
		orderedMap.Set(int64(item.Score), *d)
	}

	CollectionCache.Set(
		"system_monitor:collection:"+uuid,
		orderedMap,
		time.Duration(GetEnvInt("LOCAL_CACHE_TIME", 300))*time.Second,
	)
	return orderedMap, nil
}

func CollectionFormat(collections *orderedmap.OrderedMap[int64, CollectionData], name string) map[string]interface{} {
	switch name {
	case "Memory":
		result := map[string]interface{}{
			"time": []string{},
			"value": map[string]interface{}{
				"Mem":  []string{},
				"Swap": []string{},
			},
		}

		for score, collection := range collections.AllFromFront() {
			d := reflect.ValueOf(collection)
			// result["time"] = append(result["time"].([]string), fmt.Sprint(score))

			result["time"] = append(result["time"].([]string), time.Unix(score, 0).Format("01-02 15:04"))
			result["value"].(map[string]interface{})["Mem"] = append(result["value"].(map[string]interface{})["Mem"].([]string), d.FieldByName(name).FieldByName("Mem").FieldByName("Used").String())
			result["value"].(map[string]interface{})["Swap"] = append(result["value"].(map[string]interface{})["Swap"].([]string), d.FieldByName(name).FieldByName("Swap").FieldByName("Used").String())

		}
		return result
	default:
		// for score, collection := range collections {
		// d := reflect.ValueOf(collection)

		// result[score] = d.FieldByName(name).Interface()
		// }
		return nil
	}

}
