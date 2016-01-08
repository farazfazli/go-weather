package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/JanBerktold/sse"
)

type Location struct {
	Latitude  float32
	Longitude float32
	City      string
	Region    string `json:"region_name"`
}

type Forecast struct {
	Currently struct {
		Summary     string
		Temperature float32
	}
	Daily struct {
		Data []struct {
			Summary        string
			TemperatureMin float32
			TemperatureMax float32
		}
	}
}

var latitude string
var longitude string
var ip string = "" // Insert Fallback IP

func getUserLocation(conn *sse.Conn) {
	response, err := http.Get("http://freegeoip.net/json/" + ip)
	if err != nil {
		fmt.Printf("Error occured: %s", err)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		res := Location{}
		json.Unmarshal(contents, &res)
		if err != nil {
			fmt.Printf("Error occured: %s", err)
		}
		latitude = fmt.Sprintf("%f", res.Latitude)
		longitude = fmt.Sprintf("%f", res.Longitude)
		conn.WriteStringEvent("location", fmt.Sprintf("Location: %s, %s", res.City, res.Region))
	}
}

func averageTemperature(min float32, max float32) (avg float32) {
	return (min + max) / 2
}

func updateCurrentWeather(conn *sse.Conn) {
	apikey := "INSERTAPIKEYHERE"
	response, err := http.Get("https://api.forecast.io/forecast/" + apikey + "/" + latitude + "," + longitude + "?units=auto")
	if err != nil {
		fmt.Printf("Error occured: %s", err)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		cur := Forecast{}
		json.Unmarshal(contents, &cur)
		if err != nil {
			fmt.Printf("Error occured: %s", err)
		}
		conn.WriteStringEvent("time", fmt.Sprintf("Last updated: %s", time.Now().Format("15:04:05 MST")))
		conn.WriteStringEvent("summary", fmt.Sprintf("Conditions:  %s", cur.Currently.Summary))
		conn.WriteStringEvent("temperature", fmt.Sprintf("Temperature: %v °F", cur.Currently.Temperature))
		var divTag string
		for i := 1; i < 8; i++ {
			divTag = fmt.Sprintf("%v", i)
			conn.WriteStringEvent(divTag, fmt.Sprintf("+%s day(s): %v °F", divTag, averageTemperature(cur.Daily.Data[i].TemperatureMin, cur.Daily.Data[i].TemperatureMax)))
		}
	}
}

func HandleSSE(w http.ResponseWriter, r *http.Request) {
	conn, err := sse.Upgrade(w, r)
	if err != nil {
		fmt.Printf("Error occured: %q", err.Error())
	}
	getUserLocation(conn)
	for {
		updateCurrentWeather(conn)
		time.Sleep(60 * time.Second) // Updates every 60 seconds
	}

}

func main() {
	http.HandleFunc("/event", HandleSSE)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		dip, _, _ := net.SplitHostPort(r.RemoteAddr)
		if fmt.Sprintf("%s", ip) != "127.0.0.1" && strings.Contains(ip, "::") { // launched locally
			ip = dip
		}
		http.ServeFile(w, r, "main.html")
	})
	http.ListenAndServe(":8000", nil)
}
