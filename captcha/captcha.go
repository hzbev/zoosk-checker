package captcha

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

type PostRes struct {
	ErrorID          int    `json:"errorId"`
	ErrorCode        string `json:"errorCode"`
	ErrorDescription string `json:"errorDescription"`
	TaskID           int    `json:"taskId"`
}

type GetTaskRes struct {
	ErrorID          int         `json:"errorId"`
	ErrorCode        interface{} `json:"errorCode"`
	ErrorDescription interface{} `json:"errorDescription"`
	Solution         struct {
		GRecaptchaResponse string `json:"gRecaptchaResponse"`
	} `json:"solution"`
	Status string `json:"status"`
}

type CheckCapStru struct {
	ClientKey string `json:"clientKey"`
	TaskID    int    `json:"taskId"`
}

func PostCaptcha(key string) int {
	res_body := PostRes{}
	client := resty.New()
	resp, _ := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{"clientKey":"` + key + `","task":{"type":"HCaptchaTaskProxyless","websiteURL":"https://www.zoosk.com/","websiteKey":"9a9b069e-723e-4522-81e6-2b244ff63649"}}`).
		Post("https://api.capmonster.cloud/createTask")
	json.Unmarshal(resp.Body(), &res_body)
	return res_body.TaskID
}

func CheckCapRes(id int, key string) GetTaskRes {
	res_body := GetTaskRes{}
	client := resty.New()
	resp, _ := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{"clientKey":"` + key + `", "taskId":` + strconv.Itoa(id) + `}`).
		Post("https://api.capmonster.cloud/getTaskResult")
	json.Unmarshal(resp.Body(), &res_body)
	return res_body
}

func RecursiveCaptchaCheck(id int, key string) string {
	for {
		time.Sleep(5 * time.Second)
		i := CheckCapRes(id, key)
		if i.Status == "ready" {
			return i.Solution.GRecaptchaResponse
		}
	}
}
