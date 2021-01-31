package main

import (
	"fmt"
	"github.com/gocarina/gocsv"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

//Create CLISteamer
//Create CLIRunner
//CLISteamer read CSV from os.args
//CLIRunner runs the streamer
// Run steamer parallel
// Write output to the file

type CliStreamerRecord struct {
	Title       string `csv:"Title"`
	Message1    string `csv:"Message 1"`
	Message2    string `csv:"Message 2"`
	StreamDelay int    `csv:"Stream Delay"`
	RunTimes    int    `csv:"Run Times"`
}

type CliRunnerRecord struct {
	Run         string `csv:"Run"`
	Title       string `csv:"Title"`
	Message1    string `csv:"Message 1"`
	Message2    string `csv:"Message 2"`
	StreamDelay int    `csv:"Stream Delay"`
	RunTimes    int    `csv:"Run Times"`
}

var wg = sync.WaitGroup{}
var mtx = sync.RWMutex{}

func main() {
	fileData := []byte("")
	errw := ioutil.WriteFile("cmd/main/outputData", fileData, 0644)
	check(errw)

	var args string
	allCSVRows := strings.Split(os.Args[1], string(92)+"n")
	totalArgs := len(allCSVRows)

	for i := 0; i < totalArgs; i++ {
		args = args + allCSVRows[i] + "\n"
	}

	var cliRunners []CliRunnerRecord
	gocsv.UnmarshalString(
		args,
		&cliRunners)

	for _, runner := range cliRunners {
		wg.Add(1)
		go CLIRunner(runner)

	}
	wg.Wait()
}
func CLIRunner(aRunnerRecord CliRunnerRecord) {
	steamThread, err := strconv.Atoi(aRunnerRecord.Run)
	if err != nil {
		fmt.Println("Problem occured")
	}
	for i := 0; i < steamThread; i++ {
		CLIStreamer(aRunnerRecord)
	}
	wg.Done()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func CLIStreamer(cliRunnerRecord CliRunnerRecord) {
	var aCliSteamRecord CliStreamerRecord
	aCliSteamRecord.Title = cliRunnerRecord.Title
	aCliSteamRecord.Message1 = cliRunnerRecord.Message1
	aCliSteamRecord.Message2 = cliRunnerRecord.Message2
	aCliSteamRecord.StreamDelay = cliRunnerRecord.StreamDelay
	aCliSteamRecord.RunTimes = cliRunnerRecord.RunTimes
	for i := 0; i < aCliSteamRecord.RunTimes; i++ {
		streamAMessage(aCliSteamRecord.Message1, aCliSteamRecord.Title)
		whiteOutputOnFile(aCliSteamRecord.Title + " -> " + aCliSteamRecord.Message1)
		time.Sleep(time.Duration(aCliSteamRecord.StreamDelay) * 1000000 * time.Microsecond)
		streamAMessage(aCliSteamRecord.Message2, aCliSteamRecord.Title)
		whiteOutputOnFile(aCliSteamRecord.Title + " -> " + aCliSteamRecord.Message2)
		time.Sleep(time.Duration(aCliSteamRecord.StreamDelay) * 1000000 * time.Microsecond)
	}
}
func streamAMessage(message string, title string) {
	fmt.Println(title + "->" + message)
}
func whiteOutputOnFile(line string) {
	mtx.Lock()
	dat, errr := ioutil.ReadFile("cmd/main/outputData")
	check(errr)
	oldString := string(dat)
	var dataToStore []byte
	if oldString == "" {
		dataToStore = []byte(line)
	} else {
		dataToStore = []byte(oldString + "\n" + line)
	}
	errw := ioutil.WriteFile("cmd/main/outputData", dataToStore, 0644)
	check(errw)
	mtx.Unlock()
}
