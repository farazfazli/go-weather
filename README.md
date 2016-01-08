# go-weather
Server-Sent Events weather web app written in Go. 

go-weather is designed to be minimalistic and lightweight. Powered by Forecast.io (The Dark Sky Forecast API)

## Features
- Automatically looks up location based on IP
- Current weather, as well as average daily temperature for the next week
- HTML5 Server-Sent Events (Updates weather every minute, no reloading needed)

## Running go-weather
1. ```go get github.com/farazfazli/go-weather```
2. ```go get github.com/JanBerktold/sse```
3. Configure Fallback IP and API key -- Get your API key from https://developer.forecast.io/
4. ```go run main.go```

Runs on localhost:8000 by default, you can easily change this by replacing 8000 in the code with the port of your choosing.
