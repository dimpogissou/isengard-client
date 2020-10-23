# isengard-client
A logging files tail-and-dispatch client written in Go

# To do

## Functionalities
- [ ] Design error switch to ensure storage, alerting and execution continuation on connector.Send() failures 
- [ ] Add parsing config validation (WIP)
- [x] Add connectors to config
- [x] Add s3 connector
- [ ] Add kinesis connector
- [ ] Add rollbar connector
- [x] Add basic kafka connector
- [ ] Add Datadog client and metrics
- [ ] Add JSON logs support
- [ ] Add log file name pattern
- [x] Unify logging
- [ ] Add batch mode for S3/Kafka connectors

## Code quality
- [x] Replace bool by error in connector.Send return type and pass the handling to main so that the complete application flow can be read and understood from main.go 
- [x] Rename dispatch.go 
- [x] Remove connectors' Open() function, keep Close() 
- [ ] Reach 100% test coverage

## Env/Ops

- [x] Restructure project
- [x] Dockerise app
- [x] Docker-compose test configuration
- [x] Makefile 
- [ ] Setup CI
- [ ] Setup test coverage report
- [ ] Add Kafka broker, schema registry, introduce producer groups

## Tests 
- [ ] Unit
- [x] Integration
- [ ] E2E
