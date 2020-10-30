# isengard-client
A logging files tail-and-dispatch client written in Go

# To do

## Functionalities
- [ ] Add kinesis connector
- [ ] Add rollbar connector
- [ ] Add failover function to ensure storage, alerting and execution continuation 
- [ ] Add parsing config validation (WIP)
- [ ] Add batch mode for S3/Kafka connectors
- [ ] Add Datadog client and metrics
- [ ] Add JSON logs support
- [ ] Add log file name pattern
- [x] Ensure all resources are closed in case of interruption signal
- [x] Add connectors to config
- [x] Add s3 connector
- [x] Unify logging
- [x] Add basic kafka connector

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
- [x] Setup test coverage report 
- [ ] Setup CI
- [ ] Add Kafka broker, schema registry, introduce producer groups

## Tests 
- [ ] Unit
- [x] Integration
- [ ] E2E
