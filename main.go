package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

type Weather struct {
	Status Status `json:"status"`
}

type Status struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var jsons = map[string]interface{}{}

		file, _ := ioutil.ReadFile("weather.json")
		json.Unmarshal(file, &jsons)

		var t, err = template.ParseFiles("./static/layout.html")
		if err != nil {
			fmt.Fprintf(w, "error")
			return
		}

		var status, color string
		waterValue := jsons["status"].(map[string]interface{})["water"]
		windValue := jsons["status"].(map[string]interface{})["wind"]

		if waterValue.(float64) <= 5 && windValue.(float64) <= 6 {
			status = "Aman"
			color = "green"
		} else if waterValue.(float64) <= 8 && windValue.(float64) <= 15 {
			status = "Siaga"
			color = "yellow"
		} else if waterValue.(float64) > 8 || windValue.(float64) > 15 {
			status = "Bahaya"
			color = "red"
		}

		var data = map[string]interface{}{
			"Wind":   windValue,
			"Water":  waterValue,
			"Status": status,
			"Color":  color,
		}

		t.Execute(w, data)
	})

	go func() {
		for {
			water := rand.Intn(99) + 1
			wind := rand.Intn(99) + 1

			data := Weather{
				Status: Status{
					Water: water,
					Wind:  wind,
				},
			}

			file, _ := json.MarshalIndent(data, "", " ")
			_ = ioutil.WriteFile("weather.json", file, 0644)

			time.Sleep(15 * time.Second)
		}
	}()

	fmt.Println("Starting web server at :8080")
	http.ListenAndServe(":8080", nil)
}
