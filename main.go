package main

import (
	"checker/captcha"
	"checker/check"
	"checker/helper"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/gosuri/uilive"

	"github.com/sqweek/dialog"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var THREADS int = 10
var Proxy string
var Total int
var CapMonsterAuth string

func main() {
	var wg sync.WaitGroup
	writer := uilive.New()
	fmt.Println("How many accs per thread? default is 10")
	fmt.Scanln(&THREADS)
	fmt.Println("Type in rotating proxy")
	fmt.Scanln(&Proxy)
	fmt.Println("Type in Capmonster Key")
	fmt.Scanln(&CapMonsterAuth)
	fmt.Println("Select accounts (only .txt is supported)")
	filename, err := dialog.File().Filter("txt").Load()
	if err != nil {
		fmt.Println("no file selected")
		return
	}
	accs := helper.ReadtoArrayPath(filename)
	fmt.Println("loaded", len(accs), "accounts")
	check.TotalAccs = len(accs)
	chunked := helper.Chunks(accs, THREADS)
	fmt.Println("started", len(chunked), "threads")
	writer.Start()
	for _, ele := range chunked {
		wg.Add(1)
		go func(ele []string) {
			defer wg.Done()
			checkMain(ele, writer)
		}(ele)
	}
	wg.Wait()

}

func checkMain(arr []string, wr *uilive.Writer) {
	for _, ele := range arr {
		split := strings.Split(ele, ":")
		cap_key := captcha.PostCaptcha(CapMonsterAuth)
		cap_res := captcha.RecursiveCaptchaCheck(cap_key, CapMonsterAuth)
		device_id := "160" + randInt(3) + "83" + randInt(4)
		udid := "ebec" + randInt(1) + randLower(1) + "d" + randInt(1) + "a0" + randInt(1) + "34" + randInt(3) + "b5f" + randInt(1) + "2" + randInt(3) + "a" + randInt(5) + randLower(2)
		check.StartChecking(split[0], split[1], device_id, udid, cap_res, Proxy, wr)
	}

}

func randInt(n int) string {
	var letterRunes = []rune("0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func randLower(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
