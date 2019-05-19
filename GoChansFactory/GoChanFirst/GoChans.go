package GoChanFirst

import (
	"fmt"
	"strconv"
	"time"
)

var frameId = 0
var frameName = ""

func main() {
	framesToCreate := 5
	frameInfoChan := make(chan string)
	for i := 0; i < framesToCreate; i++ {
		go assembleFrame(frameInfoChan)
		go addBody(frameInfoChan)
		go addInterior(frameInfoChan)
		time.Sleep(time.Millisecond * 3000)
	}

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
