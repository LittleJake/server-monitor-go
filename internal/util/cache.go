package util

import (
	"encoding/json"
	"fmt"
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

func getCollection(uuid string) map[string]interface{} {
	return nil
}
