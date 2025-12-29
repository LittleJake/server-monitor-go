package util

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/elliotchance/orderedmap/v3"
	"github.com/karlseguin/ccache/v3"
	"github.com/peterbourgon/diskv/v3"
	"github.com/redis/go-redis/v9"
)

var CollectionCache *ccache.Cache[*orderedmap.OrderedMap[int64, CollectionData]]
var CollectionStatusCache *ccache.Cache[*orderedmap.OrderedMap[string, CollectionData]]
var MapStringCache *ccache.Cache[map[string]string]

var DiskCache *diskv.Diskv

func SetupCollectionCache() {
	CollectionCache = ccache.New(ccache.Configure[*orderedmap.OrderedMap[int64, CollectionData]]())
}

func SetupCollectionStatusCache() {
	CollectionStatusCache = ccache.New(ccache.Configure[*orderedmap.OrderedMap[string, CollectionData]]())
}

func SetupMapStringCache() {
	MapStringCache = ccache.New(ccache.Configure[map[string]string]())
}

func SetupDiskCache() {
	DiskCache = diskv.New(diskv.Options{
		BasePath: "./cache/",
		// Transform: func(s string) []string {
		// 	return []string{toHash(s)[:2]}
		// },
		CacheSizeMax: 0,
		Compression:  diskv.NewGzipCompression(),
	})
}

func toHash(s string) string {
	sha256sum := sha256.Sum256([]byte(s))
	result := hex.EncodeToString(sha256sum[:])
	return result
}

func SetDiskCacheRedisZ(key string, value []redis.Z) error {
	key = toHash(key)

	result := make([]map[string]interface{}, 0)
	for _, item := range value {
		t := map[string]interface{}{
			"Score":  strconv.FormatFloat(item.Score, 'f', -1, 64),
			"Member": item.Member,
		}
		result = append(result, t)
	}

	data, err := json.Marshal(result)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return DiskCache.Write(key, data)
}

func GetDiskCacheRedisZ(key string) ([]redis.Z, error) {
	key = toHash(key)

	data, err := DiskCache.Read(key)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var m []map[string]interface{}

	result := make([]redis.Z, 0)

	json.Unmarshal(data, &m)

	for _, item := range m {
		z := new(redis.Z)
		v, _ := toFloat64(item["Score"])
		z.Score = v
		z.Member = item["Member"]

		result = append(result, *z)
	}
	return result, nil
}
