package main

import (
	"fmt"
	"strconv"
	"time"
)

var frameId = 0
var frameName = ""
var assemblyArrangement [3]string

func main() {
	assemblyArrangement[0] = "frame"
	assemblyArrangement[1] = "body"
	assemblyArrangement[2] = "interior"

	framesToCreate := len(assemblyArrangement)
	frameInfoChan := make(chan string)
	for stageNumber := 0; stageNumber < framesToCreate; stageNumber++ {
		go assemblyStage(frameInfoChan, assemblyArrangement[stageNumber], stageNumber, framesToCreate)
		time.Sleep(time.Millisecond * 1000)
		fmt.Println("interaction complete")
	}

	/** for stageNumber := 0; stageNumber < framesToCreate; stageNumber++ {
		go assembleFrame(frameInfoChan)
		go addBody(frameInfoChan)
		go addInterior(frameInfoChan)
		time.Sleep(time.Millisecond * 3000)
	}
	*/

}

func assemblyStage(frameInfoChan chan string, stage string, stageNumber int, framesToCreate int) {
	nextStage := "paint"
	if stageNumber < framesToCreate {
		frameName = "FrameID" + strconv.Itoa(stageNumber)
		if stageNumber != framesToCreate-1 {
			nextStage = assemblyArrangement[stageNumber+1]
		}
	}
	fmt.Println("Add ", stage, " and proceed to ", nextStage)
	frameInfoChan <- frameName
	time.Sleep(time.Millisecond * 10)
}

func assembleFrame(frameInfoChan chan string) {
	frameId++
	frameName = "Frame ID" + strconv.Itoa(frameId)
	fmt.Println("Frame assembly complete", frameName, "Moving to body")
	frameInfoChan <- frameName
	time.Sleep(time.Millisecond * 5)
}

func addBody(frameInfoChan chan string) {
	body := <-frameInfoChan
	fmt.Println("Add Body to ", body, " and proceed to interior")
	frameInfoChan <- frameName
	time.Sleep(time.Millisecond * 5)
}

func addInterior(frameInfoChan chan string) {
	body := <-frameInfoChan
	fmt.Println("Add Interior to ", body, " and proceed to paint")
	frameInfoChan <- frameName
	time.Sleep(time.Millisecond * 5)
}
