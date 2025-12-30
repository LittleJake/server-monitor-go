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

type CollectionData map[string]interface{}

func checkEmpty(v any) bool {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Map, reflect.Slice, reflect.Array:
		return rv.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return rv.IsNil()
	default:
		zero := reflect.Zero(rv.Type())
		return reflect.DeepEqual(v, zero.Interface())
	}
}

func toFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("unsupported type: %T", v)
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

func GetUUIDs(refresh bool) (map[string]string, error) {
	if MapStringCache == nil {
		fmt.Println("MapStringCache is not initialized")
	}

	if !refresh && MapStringCache != nil && MapStringCache.Get("system_monitor:hashes") != nil {
		return MapStringCache.Get("system_monitor:hashes").Value(), nil
	}

	data, err := RedisHGetAll(context.Background(), RedisClient, "system_monitor:hashes")

	MapStringCache.Set(
		"system_monitor:hashes",
		data,
		time.Duration(GetEnvInt("LOCAL_CACHE_TIME", 300))*time.Second,
	)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GetCollectionStatus() (*orderedmap.OrderedMap[string, map[string]interface{}], error) {
	result := orderedmap.NewOrderedMap[string, map[string]interface{}]()

	online := map[string]interface{}{}
	offline := map[string]interface{}{}
	info := map[string]interface{}{}
	uuids, err := GetUUIDs(false)
	if err != nil {
		return nil, err
	}

	for uuidKey := range uuids {
		latest, err := GetCollectionLatest(uuidKey)
		if err != nil || len(latest) == 0 {
			fmt.Println(uuidKey, err)
			continue
		}

		_info, _ := GetInfo(uuidKey, false)

		i, _ := toFloat64(_info["Update Time"])
		t := time.Unix(int64(i), 0)

		if t.Before(time.Now().Add(-time.Duration(GetEnvInt("OFFLINE_THRESHOLD", 600)) * time.Second)) {
			offline[uuidKey] = latest
		} else {
			online[uuidKey] = latest
		}

		info[uuidKey] = _info
	}

	result.Set("online", online)
	result.Set("offline", offline)
	result.Set("info", info)

	return result, nil
}

func GetCollectionByTime(uuid string, refresh bool, start int64, end int64) (*orderedmap.OrderedMap[int64, CollectionData], error) {
	if end < start {
		return GetCollection(uuid, refresh)
	}

	result, err := GetCollection(uuid, refresh)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	m := orderedmap.NewOrderedMap[int64, CollectionData]()
	for t, v := range result.AllFromFront() {
		if t > start && t < end {
			m.Set(t, v)
			continue
		}
	}

	return m, nil
}

func GetCollection(uuid string, refresh bool) (*orderedmap.OrderedMap[int64, CollectionData], error) {
	// if CollectionCache == nil {
	// 	fmt.Println("LocalCacheClient is not initialized")
	// }

	// if !refresh && CollectionCache != nil && CollectionCache.Get("system_monitor:collection:"+uuid) != nil {
	// 	return CollectionCache.Get("system_monitor:collection:" + uuid).Value(), nil
	// }

	orderedMap := orderedmap.NewOrderedMap[int64, CollectionData]()

	data, err := RedisZRangeByScoreWithScores(
		context.Background(),
		RedisClient,
		"system_monitor:collection:"+uuid,
		&redis.ZRangeBy{Min: "0", Max: fmt.Sprint(time.Now().Unix())},
		!refresh,
	)
	if err != nil {
		// CollectionCache.Set(
		// 	"system_monitor:collection:"+uuid,
		// 	orderedMap,
		// 	time.Duration(GetEnvInt("LOCAL_CACHE_TIME", 300))*time.Second,
		// )
		return nil, err
	}

	if len(data) == 0 {
		// fmt.Println("no data found for uuid:", uuid)
		return nil, fmt.Errorf("no data found for uuid: %s", uuid)
	}

	for _, item := range data {
		d, err := UnmarshalJSONData(item.Member.(string))

		if err != nil {
			fmt.Println(err)
			continue
		}

		orderedMap.Set(int64(item.Score), *d)
	}

	// CollectionCache.Set(
	// 	"system_monitor:collection:"+uuid,
	// 	orderedMap,
	// 	time.Duration(GetEnvInt("LOCAL_CACHE_TIME", 300))*time.Second,
	// )

	return orderedMap, nil
}

func GetCollectionLatest(uuid string) (CollectionData, error) {
	orderedMap, err := GetCollection(uuid, false)
	if err != nil || orderedMap == nil || orderedMap.Len() == 0 {
		return CollectionData{}, err
	}

	latest := orderedMap.Back()
	if latest == nil {
		return CollectionData{}, fmt.Errorf("no data found for uuid: %s", uuid)
	}

	fmt.Println(latest)
	return latest.Value, nil
}

func CollectionFormat(collections *orderedmap.OrderedMap[int64, CollectionData], name string) map[string]interface{} {
	// fmt.Println(collections.Front())
	result := map[string]interface{}{}

	if collections.Len() == 0 || checkEmpty(collections.Back().Value[name]) {
		return result
	}
	switch name {
	case "Memory":
		result = map[string]interface{}{
			"time": []string{},
			"value": map[string]interface{}{
				"Mem":  []string{},
				"Swap": []string{},
			},
		}

		for score, collection := range collections.AllFromFront() {

			result["time"] = append(result["time"].([]string), time.Unix(score, 0).Format("01-02 15:04"))
			result["value"].(map[string]interface{})["Mem"] = append(result["value"].(map[string]interface{})["Mem"].([]string), collection[name].(map[string]interface{})["Mem"].(map[string]interface{})["used"].(string))
			result["value"].(map[string]interface{})["Swap"] = append(result["value"].(map[string]interface{})["Swap"].([]string), collection[name].(map[string]interface{})["Swap"].(map[string]interface{})["used"].(string))

		}

	case "Network":
		result = map[string]interface{}{
			"time": []string{},
			"RX": map[string]interface{}{
				"megabytes": []float64{},
				"packets":   []float64{},
			},
			"TX": map[string]interface{}{
				"megabytes": []float64{},
				"packets":   []float64{},
			},
		}

		for score, collection := range collections.AllFromFront() {

			result["time"] = append(result["time"].([]string), time.Unix(score, 0).Format("01-02 15:04"))
			rx_megabytes, err1 := toFloat64(collection[name].(map[string]interface{})["RX"].(map[string]interface{})["bytes"])
			rx_packets, err2 := toFloat64(collection[name].(map[string]interface{})["RX"].(map[string]interface{})["packets"])
			tx_megabytes, err3 := toFloat64(collection[name].(map[string]interface{})["TX"].(map[string]interface{})["bytes"])
			tx_packets, err4 := toFloat64(collection[name].(map[string]interface{})["TX"].(map[string]interface{})["packets"])

			if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
				fmt.Println("Error converting to float64:", err1, err2, err3, err4)
				continue
			}
			result["RX"].(map[string]interface{})["megabytes"] = append(result["RX"].(map[string]interface{})["megabytes"].([]float64), rx_megabytes/1048576)
			result["RX"].(map[string]interface{})["packets"] = append(result["RX"].(map[string]interface{})["packets"].([]float64), rx_packets/1000)
			result["TX"].(map[string]interface{})["megabytes"] = append(result["TX"].(map[string]interface{})["megabytes"].([]float64), tx_megabytes/1048576)
			result["TX"].(map[string]interface{})["packets"] = append(result["TX"].(map[string]interface{})["packets"].([]float64), tx_packets/1000)
		}

	case "Disk":
		result = map[string]interface{}{
			"time":  []string{},
			"value": map[string]interface{}{},
		}

		for score, collection := range collections.AllFromFront() {
			// d := reflect.ValueOf(collection)

			result["time"] = append(result["time"].([]string), time.Unix(score, 0).Format("01-02 15:04"))

			for key, _v := range collection[name].(map[string]interface{}) {

				// get typed reference to result["value"]
				resValue := result["value"].(map[string]interface{})

				if resValue[key] == nil {
					resValue[key] = make([]string, 0)
				}

				resValue[key] = append(resValue[key].([]string), _v.(map[string]interface{})["used"].(string))
			}
		}

	case "IO":
		result = map[string]interface{}{
			"time": []string{},
			"read": map[string]interface{}{
				"counts":    []float64{},
				"megabytes": []float64{},
				"time_ms":   []float64{},
			},
			"write": map[string]interface{}{
				"counts":    []float64{},
				"megabytes": []float64{},
				"time_ms":   []float64{},
			},
		}

		for score, collection := range collections.AllFromFront() {
			// d := reflect.ValueOf(collection)

			result["time"] = append(result["time"].([]string), time.Unix(score, 0).Format("01-02 15:04"))

			read_counts, err1 := toFloat64(collection[name].(map[string]interface{})["read"].(map[string]interface{})["count"])
			read_megabytes, err2 := toFloat64(collection[name].(map[string]interface{})["read"].(map[string]interface{})["bytes"])
			read_time_ms, err3 := toFloat64(collection[name].(map[string]interface{})["read"].(map[string]interface{})["time"])
			write_counts, err4 := toFloat64(collection[name].(map[string]interface{})["write"].(map[string]interface{})["count"])
			write_megabytes, err5 := toFloat64(collection[name].(map[string]interface{})["write"].(map[string]interface{})["bytes"])
			write_time_ms, err6 := toFloat64(collection[name].(map[string]interface{})["write"].(map[string]interface{})["time"])

			if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil {
				fmt.Println("Error converting to float64:", err1, err2, err3, err4, err5, err6)
				continue
			}

			result["read"].(map[string]interface{})["counts"] = append(result["read"].(map[string]interface{})["counts"].([]float64), read_counts)
			result["read"].(map[string]interface{})["megabytes"] = append(result["read"].(map[string]interface{})["megabytes"].([]float64), read_megabytes/1048576)
			result["read"].(map[string]interface{})["time_ms"] = append(result["read"].(map[string]interface{})["time_ms"].([]float64), read_time_ms)

			result["write"].(map[string]interface{})["counts"] = append(result["write"].(map[string]interface{})["counts"].([]float64), write_counts)
			result["write"].(map[string]interface{})["megabytes"] = append(result["write"].(map[string]interface{})["megabytes"].([]float64), write_megabytes/1048576)
			result["write"].(map[string]interface{})["time_ms"] = append(result["write"].(map[string]interface{})["time_ms"].([]float64), write_time_ms)

		}

	case "Load", "Thermal":
		result = map[string]interface{}{
			"time":  []string{},
			"value": map[string]interface{}{},
		}

		for score, collection := range collections.AllFromFront() {
			if checkEmpty(collection) || checkEmpty(collection[name]) {
				continue
			}
			result["time"] = append(result["time"].([]string), time.Unix(score, 0).Format("01-02 15:04"))
			// fmt.Println(collection[name])
			for key, _v := range collection[name].(map[string]interface{}) {
				//load:{idle, system....}
				//thermal:{cpu, gpu....}
				//generate value []string for the key inside the result: result["value"][key]
				resValue := result["value"].(map[string]interface{})
				if resValue[key] == nil {
					resValue[key] = make([]float64, 0)
				}
				//set the value to the result["value"][key]
				floatVal, err := toFloat64(_v)
				if err != nil {
					fmt.Println("Error converting to float64:", err)
				}
				resValue[key] = append(resValue[key].([]float64), floatVal)
			}

		}

		if checkEmpty(result["value"]) {
			result = map[string]interface{}{}
		}
		// default:
		// 	return result
	}

	return result
}

func GetDisplayName(refresh bool) (map[string]string, error) {
	if MapStringCache == nil {
		fmt.Println("MapStringCache is not initialized")
	}

	if !refresh && MapStringCache != nil && MapStringCache.Get("system_monitor:name") != nil {
		return MapStringCache.Get("system_monitor:name").Value(), nil
	}
	data, err := RedisHGetAll(context.Background(), GetRedisClient(), "system_monitor:name")
	if err != nil {
		fmt.Println("Error getting name from Redis:", err)
		return map[string]string{}, err
	}

	MapStringCache.Set(
		"system_monitor:name",
		data,
		time.Duration(GetEnvInt("LOCAL_CACHE_TIME", 300))*time.Second,
	)
	return data, nil
}

func GetInfo(uuid string, refresh bool) (map[string]string, error) {
	if MapStringCache == nil {
		fmt.Println("MapStringCache is not initialized")
	}

	if !refresh && MapStringCache != nil && MapStringCache.Get("system_monitor:info:"+uuid) != nil {
		return MapStringCache.Get("system_monitor:info:" + uuid).Value(), nil
	}

	data, err := RedisHGetAll(context.Background(), GetRedisClient(), "system_monitor:info:"+uuid)
	if err != nil || len(data) == 0 {
		// MapStringCache.Set(
		// 	"system_monitor:info:"+uuid,
		// 	data,
		// 	time.Duration(GetEnvInt("LOCAL_CACHE_TIME", 300))*time.Second,
		// )
		// fmt.Println("Error getting info from Redis:", err)
		return map[string]string{}, err
	}

	MapStringCache.Set(
		"system_monitor:info:"+uuid,
		data,
		time.Duration(GetEnvInt("LOCAL_CACHE_TIME", 300))*time.Second,
	)
	return data, nil
}

func RetentionCollectionData(uuid string) {
	ctx := context.Background()
	retentionDays := GetEnvInt("DATA_RETENTION_DAYS", 7)
	cutoffTimestamp := time.Now().AddDate(0, 0, -retentionDays).Unix()
	RedisZRemRangeByScore(
		ctx,
		RedisClient,
		"system_monitor:collection:"+uuid,
		"-inf",
		fmt.Sprint(cutoffTimestamp),
	)

	// if err != nil {
	// 	fmt.Println("Error during data retention for uuid:", uuid, err)
	// } else {
	// 	fmt.Println("Data retention completed for uuid:", uuid)
	// }
}

func CronJob() {
	for {
		uuids, _ := GetUUIDs(true)
		for uuidKey := range uuids {
			go GetInfo(uuidKey, true)
			go GetCollection(uuidKey, true)
			go GetDisplayName(true)
			go RetentionCollectionData(uuidKey)
		}
		// fmt.Printf("[current]MapStringCache size: %d, CollectionStatusCache size: %d, MapStringCache size: %d.\n", MapStringCache.GetSize(), CollectionStatusCache.GetSize(), MapStringCache.GetSize())
		// fmt.Printf("[dropped]MapStringCache size: %d, CollectionStatusCache size: %d, MapStringCache size: %d.\n", MapStringCache.GetDropped(), CollectionStatusCache.GetDropped(), MapStringCache.GetDropped())
		// fmt.Printf("[current]MapStringCache size: %d, CollectionStatusCache size: %d, MapStringCache size: %d.\n", unsafe.Sizeof(MapStringCache), unsafe.Sizeof(CollectionStatusCache), unsafe.Sizeof(MapStringCache))

		// var m runtime.MemStats
		// runtime.ReadMemStats(&m)
		// fmt.Printf("Alloc: %d KiB, HeapSys: %d KiB, PauseTotalNs: %d ns\n",
		// 	m.Alloc>>10, m.HeapSys>>10, m.PauseTotalNs)
		// fmt.Printf("HeapIdle: %d KiB, HeapInuse: %d KiB, PauseTotalNs: %d ns\n",
		// 	m.HeapIdle>>10, m.HeapInuse>>10, m.PauseTotalNs)

		time.Sleep(time.Duration(GetEnvInt("CRON_JOB_INTERVAL", 60)) * time.Second)
	}
}
