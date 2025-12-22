package util

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/elliotchance/orderedmap/v3"
	"github.com/redis/go-redis/v9"
)

type CollectionData struct {
	Battery []interface{}                     `json:"Battery"`
	Disk    map[string]map[string]interface{} `json:"Disk"`
	Fan     []interface{}                     `json:"Fan"`
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
	Load   map[string]interface{} `json:"Load"`
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
	Ping    []interface{}          `json:"Ping"`
	Thermal map[string]interface{} `json:"Thermal"`
}

func FieldToString(f reflect.Value) string {
	switch f.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(f.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(f.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(f.Float(), 'f', -1, 64) // -1 自动去掉多余0
	case reflect.String:
		return f.String()
	case reflect.Bool:
		return strconv.FormatBool(f.Bool())
	default:
		return fmt.Sprintf("%v", f.Interface())
	}
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

func GetCollectionLatest(uuid string) (CollectionData, error) {
	orderedMap, err := GetCollection(uuid)
	if err != nil {
		return CollectionData{}, err
	}

	latest := orderedMap.Back()
	if latest == nil {
		return CollectionData{}, fmt.Errorf("no data found for uuid: %s", uuid)
	}

	return latest.Value, nil
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

	case "Disk":
		result := map[string]interface{}{
			"time":  []string{},
			"value": map[string]interface{}{},
		}
		for score, collection := range collections.AllFromFront() {
			d := reflect.ValueOf(collection)

			result["time"] = append(result["time"].([]string), time.Unix(score, 0).Format("01-02 15:04"))

			for _, key := range d.FieldByName(name).MapKeys() {
				value := d.FieldByName(name).MapIndex(key) // mount point value map{free, used, percentage}
				fieldName := key.String()                  //mount point name

				// get typed reference to result["value"]
				resValue := result["value"].(map[string]interface{})

				if resValue[fieldName] == nil {
					resValue[fieldName] = make([]string, 0)
				}

				resValue[fieldName] = append(resValue[fieldName].([]string), FieldToString(value.MapIndex(reflect.ValueOf("used"))))
			}
		}
		return result
	case "Load", "Thermal":
		result := map[string]interface{}{
			"time":  []string{},
			"value": map[string]interface{}{},
		}

		for score, collection := range collections.AllFromFront() {
			d := reflect.ValueOf(collection)
			// result["time"] = append(result["time"].([]string), fmt.Sprint(score))

			result["time"] = append(result["time"].([]string), time.Unix(score, 0).Format("01-02 15:04"))

			for _, key := range d.FieldByName(name).MapKeys() {
				//load:{idle, system....}
				//thermal:{cpu, gpu....}
				value := d.FieldByName(name).MapIndex(key) // value is a map[string]interface{}
				//get the key name inside the map (eg. "idle", "system")
				fieldName := key.String()
				//generate value []string for the key inside the result: result["value"][key]
				resValue := result["value"].(map[string]interface{})
				if resValue[fieldName] == nil {
					resValue[fieldName] = make([]string, 0)
				}
				//set the value to the result["value"][key]
				resValue[fieldName] = append(resValue[fieldName].([]string), FieldToString(value))
			}

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

func GetInfo(uuid string) (map[string]string, error) {
	if MapStringCache == nil {
		return nil, fmt.Errorf("MapStringCache is not initialized")
	}

	if MapStringCache.Get("system_monitor:info:"+uuid) != nil {
		return MapStringCache.Get("system_monitor:info:" + uuid).Value(), nil
	}

	MapStringCache.Set(
		"system_monitor:info:"+uuid,
		RedisClient.HGetAll(context.Background(), "system_monitor:info:"+uuid).Val(),
		time.Duration(GetEnvInt("LOCAL_CACHE_TIME", 300))*time.Second,
	)

	data := RedisClient.HGetAll(context.Background(), "system_monitor:info:"+uuid)

	return data.Val(), nil
}
