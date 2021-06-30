# Telemetry challenge project

## pre-requisites

* go

https://golang.org/dl/

## Instructions

Run the following commands in separate terminals
```
./telemetry generate
```
```
go mod tidy
go run main.go chart.go
```

Go to http://localhost:8081/ to see the data received thus far, refresh the page for the latest data.