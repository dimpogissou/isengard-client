# isengard-client
A logging files tail-and-dispatch client written in Go

# To do

## Code / Functionalities
- [ ] Design error switch to ensure storage, alerting and execution continuation on connector.Send failures 
- [ ] Add connectors to config (WIP)
- [x] Add s3 connector
- [ ] Add rollbar connector
- [ ] Add kafka connector
- [ ] Add log file name pattern
- [ ] Add parsing config validation (WIP)
- [x] Unify logging

## Env/Ops

- [x] Restructure project
- [x] Dockerise app
- [x] Docker-compose test configuration
- [ ] Setup CI
- [ ] Setup test coverage report

## Tests 
- [ ] Unit
- [x] Integration
- [x] E2E
