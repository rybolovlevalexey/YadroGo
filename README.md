# Biathlon Competition Prototype

An application for processing logs of events related to the biathlon race. Before receiving information about events, the race description file must be processed, and after receiving information about events, the final report is displayed.

## Prerequisites
- Go (1.20+)

## Launch project
- ```git clone https://github.com/rybolovlevalexey/YadroGo/```
- go to the root of the project
- put your config file in files/configs/ and events file in files/events/
- change ConfigPath and EventsPath in settings/settings.go
- ```go run main.go```
- ```go test -cover ./...``` (to watch test results)

## Packages description
- models - all the structures that will be used in the project
- settings - project settings (paths to configs and events file)
- services - auxiliary functions (loading configs, getting the duration by string with time)
- usecases - classes that implement the processing of incoming events and the creation of a final report
- main - project launch point
