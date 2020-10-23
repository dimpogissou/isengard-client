# isengard-client
A logging files tail-and-dispatch client written in Go

# To do

## Functionalities
- [ ] Design error switch to ensure storage, alerting and execution continuation on connector.Send failures 
- [x] Add connectors to config
- [x] Add S3 connector
- [ ] Add Kinesis connector
- [ ] Add Rollbar connector
- [x] Add kafka connector
- [ ] Add Datadog client and metrics
- [ ] Add JSON logs support
- [ ] Add log file name pattern
- [ ] Add parsing config validation (WIP)
- [x] Unify logging
- [ ] Optimize string writes for line-level logging with bytes buffer
- [ ] Add batch mode for S3/Kafka connectors

## Code quality
- [ ] Replace bool by error in connector.Send return type and pass the handling to main so that the complete application flow can be 
read and understood from main.go 
- [ ] Rename dispatch.go to a more appropriate name since it's only instantiating the clients
- [ ] Move connectors' setup and teardown functions to Open() and Close() methods (do some quick read about pointers first)
- [ ] Reach 100% test coverage

## Env/Ops

- [x] Restructure project
- [x] Dockerise app
- [x] Docker-compose test configuration
- [ ] Test shortcuts to avoid logging into container
- [ ] Setup CI
- [ ] Setup test coverage report

## Tests 
- [ ] Unit
- [x] Integration
- [x] E2E
