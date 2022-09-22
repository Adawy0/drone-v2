# drone-v2
## create database drone, then create drone schema
## install gorm-goose migration tool
## apply this command
``gorm-goose -path=repository/db -pgschema=drone up``

## can run unit test using this command
`` go test --cover -v ./...``

## can run app 
`` go run main.go ``
