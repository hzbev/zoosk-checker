package helper

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

var prod string = "desv"
var GlobalAgent string
var AllAgents []string
var TotalGone int = 0

func ForceChange() {
	GlobalAgent = AllAgents[rand.Intn(len(AllAgents))]
	TotalGone = 0
}

func init() {
	rand.Seed(time.Now().UnixNano())
	AllAgents = ReadtoArray("agents.txt")
	GlobalAgent = AllAgents[rand.Intn(len(AllAgents))]
	log.Println("init UA is", GlobalAgent)
	ticker := time.NewTicker(30 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				rand.Seed(time.Now().UnixNano())
				GlobalAgent = AllAgents[rand.Intn(len(AllAgents))]
			case <-quit:
				return
			}
		}
	}()
}

func Chunks(xs []string, chunkSize int) [][]string {
	if len(xs) == 0 {
		return nil
	}
	divided := make([][]string, (len(xs)+chunkSize-1)/chunkSize)
	prev := 0
	i := 0
	till := len(xs) - chunkSize
	for prev < till {
		next := prev + chunkSize
		divided[i] = xs[prev:next]
		prev = next
		i++
	}
	divided[i] = xs[prev:]
	return divided
}

func ReadtoArray(filePath string) []string {
	var workingDir string
	if prod == "dev" {
		workingDir, _ = os.Getwd()
	} else {
		exeDir, _ := os.Executable()
		workingDir = filepath.Dir(exeDir)
	}

	file, err := os.Open(filepath.Join(workingDir, filePath))
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	var content []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		content = append(content, scanner.Text())
	}
	return content
}

func ReadtoArrayPath(filePath string) []string {

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	var content []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		content = append(content, scanner.Text())
	}
	return content
}

func Write(filePath, text string) {
	var workingDir string
	if prod == "dev" {
		workingDir, _ = os.Getwd()
	} else {
		exeDir, _ := os.Executable()
		workingDir = filepath.Dir(exeDir)
	}

	file, err := os.OpenFile(workingDir+`/`+filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}

	if _, err := file.Write([]byte(text + "\n")); err != nil {
		log.Fatal(err)
	}

	defer file.Close()
}
