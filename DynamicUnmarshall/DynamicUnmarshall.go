package main

import (
	"encoding/json"
	"fmt"
)

type Restaurant struct {
	Restaurant RestaurantData `json:"restaurant"`
}

type RestaurantData struct {
	Name  string `json:"name"`
	Owner Owner  `json:"owner"`
}

type Owner struct {
	Name string `json:"name"`
}

func main() {
	data := `{"restaurant":{"name":"Tickets","owner":{"name":"Ferran"}}}`
	r := Restaurant{}
	json.Unmarshal([]byte(data), &r)

	fmt.Printf("%+v", r)

	//Generic unmarshalling ala gson by decoding into an interface and
	//then extracting a top level map from the result, e.g:
	var msgMapTemplate interface{}
	err := json.Unmarshal([]byte(t.ResponseBody), &msgMapTemplate)
	t.AssertEqual(err, nil)
	msgMap := msgMapTemplate.(map[string]interface{})

}
