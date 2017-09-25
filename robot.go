package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

//Bot instruction filename
const FILENAME = "botinstructions.txt"

type botData struct {
	value1 int
	value2 int
}

var botDataMap = make(map[int]botData)
var messages = make(chan string)

type conditionStruct struct {
	lock *sync.Cond
}

var condition conditionStruct = conditionStruct{lock: &sync.Cond{L: &sync.Mutex{}}}

func givetoBot(bot int, value int) {
	if val, ok := botDataMap[bot]; ok {
		if val.value1 == 0 {
			val.value1 = value
		} else {
			val.value2 = value
		}
		botDataMap[bot] = val
	}
}
func parseInstruction(msg string) {
	var bot, low, high int
	switch {
	case strings.Contains(msg, "low to bot") && strings.Contains(msg, "high to bot"):
		fmt.Sscanf(msg, "bot %d gives low to bot %d and high to bot %d", &bot, &low, &high)

		if v, ok := botDataMap[bot]; ok {
			if v.value1 < v.value2 {
				givetoBot(low, v.value1)
				givetoBot(high, v.value2)
			} else {
				givetoBot(high, v.value1)
				givetoBot(low, v.value2)
			}
		}

	case strings.Contains(msg, "low to output") && strings.Contains(msg, "high to output"):
		fmt.Sscanf(msg, "bot %d gives low to output %d and high to output %d", &bot, &low, &high)

	case strings.Contains(msg, "low to output"):
		fmt.Sscanf(msg, "bot %d gives low to output %d and high to bot %d", &bot, &low, &high)
		if v, ok := botDataMap[bot]; ok {
			if v.value1 < v.value2 {
				givetoBot(low, v.value1)
				givetoBot(high, v.value2)
			} else {
				givetoBot(high, v.value1)
				givetoBot(low, v.value2)
			}
		}

	case strings.Contains(msg, "high to output"):
		fmt.Sscanf(msg, "bot %d gives low to bot %d and high to ouput %d", &bot, &low, &high)
		if v, ok := botDataMap[bot]; ok {
			if v.value1 < v.value2 {
				givetoBot(low, v.value1)
				givetoBot(high, v.value2)
			} else {
				givetoBot(high, v.value1)
				givetoBot(low, v.value2)
			}
		}
	}
}

func processBotInstructions() {
	//Wait for messages
	for {
		select {
		case msg := <-messages:
			// fmt.Println("Messagfe receieved in processBotInstructions ", msg)
			parseInstruction(msg)
		}
	}
}

func sendBotInstructions() {
	f, err := os.Open(FILENAME)
	if err != nil {
		panic(err.Error())
	}
	// populate data
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		fileline := scanner.Text()
		if strings.Contains(fileline, "value") == false {
			messages <- fileline // send data to processBotInstruction go routine
		}
	}
	f.Close()
	condition.lock.Signal() //Signal to proceed with the print
}
func main() {

	f, err := os.Open(FILENAME)
	if err != nil {
		panic(err.Error())
	}
	// populate data
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		fileline := scanner.Text()
		if strings.Contains(fileline, "value") {
			bot, value := 0, 0
			fmt.Sscanf(fileline, "value %d goes to bot %d", &value, &bot)
			if val, ok := botDataMap[bot]; ok {
				val.value2 = value
				botDataMap[bot] = val
			} else {
				val.value1 = value
				botDataMap[bot] = val
			}
		}
	}

	f.Close()
	//Start both send and process go routines
	go processBotInstructions()
	go sendBotInstructions()

	condition.lock.L.Lock()
	condition.lock.Wait() //Wait in conditional variable to finish the sendBotInstruction Execution
	condition.lock.L.Unlock()

	for k, v := range botDataMap {
		fmt.Printf("%d -- %d,%d\n", k, v.value1, v.value2)
	}
}
