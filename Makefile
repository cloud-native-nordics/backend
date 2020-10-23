.PHONY: run
## run: runs the docker container for local development
# the slack token is set to a dummy value of 1234, if you want to test the slack integration locally
# you will need to update the environment variable parsed to the container
# please do not commit the SLACK_TOKEN as this is a secret
run:
	docker build . -t slack-invite && docker run -e SLACK_TOKEN=1234 slack-invite

.PHONY: test
## test: runs all tests
test:
	go test ./...

.PHONY: help
## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'