package kr

import (
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const TRACKING_ID = "UA-86173430-2"

type Analytics struct{}

func (Analytics) post(clientID string, params url.Values) {
	if clientID == "disabled" {
		return
	}
	defaultParams := url.Values{
		"v":   []string{"1"},
		"tid": []string{TRACKING_ID},
		"cid": []string{clientID},
		"ua":  []string{analytics_user_agent},
		"cd1": []string{CURRENT_VERSION.String()},
		"cd2": []string{analytics_os},
	}
	if osVersion := getAnalyticsOSVersion(); osVersion != nil {
		defaultParams["cd3"] = []string{*osVersion}
	}
	for k, v := range params {
		defaultParams[k] = v
	}
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	_, err := client.PostForm("https://www.google-analytics.com/collect", defaultParams)

	if err != nil {
		log.Error("error posting to analytics:", err.Error())
	}
}

func (Analytics) PostEvent(clientID string, category string, action string, label *string, value *uint64) {
	params := url.Values{
		"t":  []string{"event"},
		"ec": []string{category},
		"ea": []string{action},
	}
	if label != nil {
		params["el"] = []string{*label}
	}
	if value != nil {
		params["ev"] = []string{strconv.FormatUint(*value, 10)}
	}
	Analytics{}.post(clientID, params)
}
