# Biathlon Competition Prototype

An application for processing logs of events related to the biathlon race. Before receiving information about events, the race description file must be processed, and after receiving information about events, the final report is displayed(finished competitors sorted by total time, others - by competitor id).

## Prerequisites
- Go (1.20+)

## Launch project
1. ```git clone https://github.com/rybolovlevalexey/YadroGo/```
2. go to the root of the project
3. put your config file in files/configs/ and events file in files/events/ (optional)
4. change ConfigPath and EventsPath in settings/settings.go (do not follow if you did not follow step 3)
5. ```go run main.go```
6. ```go test -cover ./...``` (to watch test results)

## Packages description
- models - all the structures that will be used in the project
- settings - project settings (paths to configs and events file)
- services - auxiliary functions (loading configs, getting the duration by string with time)
- usecases - classes that implement the processing of incoming events and the creation of a final report
- main - project launch point
