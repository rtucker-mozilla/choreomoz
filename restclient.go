package main

import (
	"errors"
	"fmt"
	"github.com/jmcvetta/napping"
)

type SystemIdResp struct {
	Createdat string `json: "created_at"`
	Hostname  string `json: "hostname"`
	Id        int    `json: "id"`
}

type GetSystemIdResp struct {
	TotalCount int          `json: "total_count"`
	Limit      int          `json: "limit"`
	Offset     int          `json: "offset"`
	System     SystemIdResp `json: "system"`
}

type CreateSystemUpdateResp struct {
	TotalCount int `json: "total_count"`
	Limit      int `json: "limit"`
	Offset     int `json: "offset"`
	Id         int `json: "id"`
}

// CreateClient creates a napping
// takes the REST API endpoint and the system hostname
// returns a bool of the success of creating the database file.
// Example json output from REST API:
/*************************************************
{
  "limit": 1,
  "offset": 0,
  "system": {
    "created_at": "Tue, 16 Sep 2014 11:49:01 GMT",
    "hostname": "localhost.localdomain",
    "id": 14
  },
  "total_count": 1
}
**************************************************/
func APIGetSystemId(url string, hostname string) (int, error) {
	log := GetLogger()
	log.Debug("Enter into GetSystemId")
	defer func() {
		if e := recover(); e != nil {
			log.Error(e)
		}
	}()
	system_id := 0
	result := GetSystemIdResp{}
	full_url := fmt.Sprintf("%s/getsystemid/%s/", url, hostname)
	resp, err := napping.Get(full_url, nil, &result, nil)
	if resp.Status() == 200 {
		system_id = result.System.Id
	} else {
		err = errors.New("Unable to contact REST API Endpoint.")
	}
	log.Debug("Exit into GetSystemId")
	return system_id, err
}

// CreateSystemUpdate starts the update process with the api
// Example json output from REST API:
/*****************
{
  "id": 557,
  "limit": 1,
  "offset": 0,
  "total_count": 1
}
*****************/
func CreateSytemUpdate(url string, system_id int) (int, error) {
	log := GetLogger()
	log.Debug("Enter into CreateSystemUpdate")
	defer func() {
		if e := recover(); e != nil {
			log.Error(e)
		}
	}()
	update_id := 0
	result := CreateSystemUpdateResp{}
	full_url := fmt.Sprintf("%s/createupdate/%d/", url, system_id)
	resp, err := napping.Post(full_url, nil, &result, nil)
	if resp.Status() == 200 {
		update_id = result.Id
	} else {
		err = errors.New("Unable to contact REST API Endpoint.")
	}
	log.Debug("Exit into CreateSystemUpdate")
	return update_id, err
}

func FinishSystemUpdate(url string, system_id int) (bool, error) {
	log := GetLogger()
	defer func() {
		if e := recover(); e != nil {
			log.Error(e)
		}
	}()
	log.Debug("Start FinishSystemUpdate")
	return_val := false
	full_url := fmt.Sprintf("%s/finishupdate/%d/", url, system_id)
	resp, err := napping.Post(full_url, nil, nil, nil)
	if resp.Status() == 200 {
		return_val = true
		err = nil
	} else {
		err = errors.New("Unable to contact REST API Endpoint.")
	}
	log.Debug("End FinishSystemUpdate")
	return return_val, err

}

// LogCapture sends the log update to the centralized API
func APILogCapture(url string, system_id int, system_update_id int, log_object *LogObject) bool {
	payload := struct {
		Return_code int    `json:"return_code"`
		Stdout      string `json:"stdout"`
		Stderr      string `json:"stderr"`
		System_id   int    `json:"system_id"`
	}{}
	log := GetLogger()
	return_value := true
	var final_url = fmt.Sprintf("%s/logcapture/", url)
	log.Debug("URL:", url)
	log.Debug("final_url:", final_url)
	log.Debug("system_id:", system_id)
	log.Debug("system_update_id:", system_update_id)
	log.Debug("log_object:", log_object)
	payload.System_id = system_id
	payload.Stdout = log_object.stdout
	payload.Stderr = log_object.stderr
	payload.Return_code = log_object.return_code
	resp, err := napping.Post(final_url, &payload, nil, nil)
	if err != nil {
		return_value = false
	}
	log.Debug("resp:", resp)
	log.Debug("err:", err)
	return return_value

}
