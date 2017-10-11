.DEFAULT_GOAL := everything

dependencies:
	@echo Downloading Dependencies
	@go get ./...

build: dependencies
	@echo Compiling Apps
	@echo   --- lightauth2 server
	@go build github.com/riomhaire/lightauth2/frameworks/application/lightauth2
	@go install github.com/riomhaire/lightauth2/frameworks/application/lightauth2
	@echo   --- lightauth2 session command
	@go build github.com/riomhaire/lightauth2/frameworks/application/lightauthsession
	@go install github.com/riomhaire/lightauth2/frameworks/application/lightauthsession
	@echo   --- lightauth2 user command
	@go build github.com/riomhaire/lightauth2/frameworks/application/lightauthuser
	@go install github.com/riomhaire/lightauth2/frameworks/application/lightauthuser
	@echo Done Compiling Apps

test:
	@echo Running Unit Tests
	@go test ./...

profile:
	@echo Profiling Code
	@go test -coverprofile coverage.out  github.com/riomhaire/lightauth2/frameworks/web 
	@go tool cover -html=coverage.out -o coverage-web.html
	@go test -coverprofile coverage.out  github.com/riomhaire/lightauth2/interfaces 
	@go tool cover -html=coverage.out -o coverage-interfaces.html
	@go test -coverprofile coverage.out  github.com/riomhaire/lightauth2/usecases
	@go tool cover -html=coverage.out -o coverage-usecases.html
	@go test -coverprofile coverage.out  github.com/riomhaire/lightauth2/entities
	@go tool cover -html=coverage.out -o coverage-entities.html
	@rm coverage.out

clean:
	@echo Cleaning
	@go clean
	@rm -f lightauth2 lightauthsession lightauthuser
	@rm -f coverage-*.html

everything: clean build test profile  
	@echo Done
