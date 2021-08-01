GOCMD=go
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
# I need to modify this to support building multiple binaries
BINARY_NAME=buffertest
VERSION?=0.0.1
SERVICE_PORT?=3000
DOCKER_REGISTRY?= #if set it should finished by /
EXPORT_RESULT?=false # for CI please set EXPORT_RESULT to true

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)

ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

.PHONY: all test build vendor

all: help

lint: lint-go lint-dockerfile lint-yaml

lint-dockerfile:
# If dockerfile is present we lint it.
ifeq ($(shell test -e ./Dockerfile && echo -n yes),yes)
	$(eval CONFIG_OPTION = $(shell [ -e $(shell pwd)/.hadolint.yaml ] && echo "-v $(shell pwd)/.hadolint.yaml:/root/.config/hadolint.yaml" || echo "" ))
	$(eval OUTPUT_OPTIONS = $(shell [ "${EXPORT_RESULT}" == "true" ] && echo "--format checkstyle" || echo "" ))
	$(eval OUTPUT_FILE = $(shell [ "${EXPORT_RESULT}" == "true" ] && echo "| tee /dev/tty > checkstyle-report.xml" || echo "" ))
	docker run --rm -i $(CONFIG_OPTION) hadolint/hadolint hadolint $(OUTPUT_OPTIONS) - < ./Dockerfile $(OUTPUT_FILE)
endif

lint-go:
	$(eval OUTPUT_OPTIONS = $(shell [ "${EXPORT_RESULT}" == "true" ] && echo "--out-format checkstyle ./... | tee /dev/tty > checkstyle-report.xml" || echo "" ))
	sudo podman run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:latest-alpine golangci-lint run --deadline=65s $(OUTPUT_OPTIONS)

lint-yaml:
ifeq ($(EXPORT_RESULT), true)
	GO111MODULE=off go get -u github.com/thomaspoignant/yamllint-checkstyle
	$(eval OUTPUT_OPTIONS = | tee /dev/tty | yamllint-checkstyle > yamllint-checkstyle.xml)
endif
	sudo podman run --rm -it -v $(shell pwd):/data cytopia/yamllint -f parsable $(shell git ls-files '*.yml' '*.yaml') $(OUTPUT_OPTIONS)

go-sec:
	sudo podman run --name go-sec --rm -v $(shell pwd):/app -w /app golang:latest /bin/sh -c 'go get github.com/securego/gosec/cmd/gosec;go mod vendor; gosec ./...'

go-vet:
	cp ~/.gitconfig ./gitconfig;sudo podman run --name go-vet --rm -v $(shell pwd):/app -w /app golang:latest /bin/sh -c 'cp /app/gitconfig ~/.gitconfig;go mod vendor; go vet ./...; go fmt ./...';rm gitconfig


clean:
	rm -fr ./bin
	rm -fr ./out
	rm -f ./junit-report.xml checkstyle-report.xml ./coverage.xml ./profile.cov yamllint-checkstyle.xml

test:
ifeq ($(EXPORT_RESULT), true)
	GO111MODULE=off go get -u github.com/jstemmer/go-junit-report
	$(eval OUTPUT_OPTIONS = | tee /dev/tty | go-junit-report -set-exit-code > junit-report.xml)
endif
	$(GOTEST) -v -race ./... $(OUTPUT_OPTIONS)

coverage:
	$(GOTEST) -cover -covermode=count -coverprofile=profile.cov ./...
	$(GOCMD) tool cover -func profile.cov
ifeq ($(EXPORT_RESULT), true)
	GO111MODULE=off go get -u github.com/AlekSi/gocov-xml
	GO111MODULE=off go get -u github.com/axw/gocov/gocov
	gocov convert profile.cov | gocov-xml > coverage.xml
endif

vendor:
	$(GOCMD) mod vendor

build:
	mkdir -p out/bin
	export GOPRIVATE=github.com/gulfcoastdevops/pkg
	GO111MODULE=on GOOS=linux $(GOCMD) build -mod vendor -ldflags="-s -w" -o out/bin/$(BINARY_NAME) cmd/server/$(BINARY_NAME)/main.go

docker-build:
	cp ~/.gitconfig $(ROOT_DIR) && sudo podman run --name build --rm --volume "$(ROOT_DIR):/src" golang:latest /bin/sh -c "cd /src; cp /src/.gitconfig ~/.gitconfig; rm /src/.gitconfig; go mod vendor; make build BINARY_NAME=$(BINARY_NAME)"

docker-build-image:
	docker build --rm --tag $(BINARY_NAME) .

docker-release:
	docker tag $(BINARY_NAME) $(DOCKER_REGISTRY)$(BINARY_NAME):latest
	docker tag $(BINARY_NAME) $(DOCKER_REGISTRY)$(BINARY_NAME):$(VERSION)
	# Push the docker images
	docker push $(DOCKER_REGISTRY)$(BINARY_NAME):latest
	docker push $(DOCKER_REGISTRY)$(BINARY_NAME):$(VERSION)

watch:
	$(eval PACKAGE_NAME=$(shell head -n 1 go.mod | cut -d ' ' -f2))
	docker run -it --rm -w /go/src/$(PACKAGE_NAME) -v $(shell pwd):/go/src/$(PACKAGE_NAME) -p $(SERVICE_PORT):$(SERVICE_PORT) cosmtrek/air

help:
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@echo "  ${YELLOW}build           		${RESET} ${GREEN}Build your project and put the output binary in out/bin/$(BINARY_NAME)${RESET}"
	@echo "  ${YELLOW}clean           		${RESET} ${GREEN}Remove build related file${RESET}"
	@echo "  ${YELLOW}coverage        		${RESET} ${GREEN}Run the tests of the project and export the coverage${RESET}"
	@echo "  ${YELLOW}docker-build    		${RESET} ${GREEN}Use Docker to build your project and put the output binary in out/bin/$(BINARY_NAME)${RESET}"
	@echo "  ${YELLOW}docker-build-image    ${RESET} ${GREEN}Use the dockerfile to build the container (name: $(BINARY_NAME))${RESET}"
	@echo "  ${YELLOW}docker-release  		${RESET} ${GREEN}Release the container \"$(DOCKER_REGISTRY)$(BINARY_NAME)\" with tag latest and $(VERSION) ${RESET}"
	@echo "  ${YELLOW}help            		${RESET} ${GREEN}Show this help message${RESET}"
	@echo "  ${YELLOW}lint            		${RESET} ${GREEN}Run all available linters${RESET}"
	@echo "  ${YELLOW}lint-dockerfile 		${RESET} ${GREEN}Lint your Dockerfile${RESET}"
	@echo "  ${YELLOW}lint-go         		${RESET} ${GREEN}Use golintci-lint on your project${RESET}"
	@echo "  ${YELLOW}lint-yaml       		${RESET} ${GREEN}Use yamllint on the yaml file of your projects${RESET}"
	@echo "  ${YELLOW}test            		${RESET} ${GREEN}Run the tests of the project${RESET}"
	@echo "  ${YELLOW}vendor          		${RESET} ${GREEN}Copy of all packages needed to support builds and tests in the vendor directory${RESET}"
	@echo "  ${YELLOW}watch           		${RESET} ${GREEN}Run the code with cosmtrek/air to have automatic reload on changes${RESET}"


xbuild:
	export GO111MODULE=on
	export GOPRIVATE=github.com/gulfcoastdevops/pkg
	env GOOS=linux go build -ldflags="-s -w" -o bin/evntlogs cmd/server/main.go


# Run
# MQTT_BROKER="ssl://mqtt1:8883" MQTT_CLIENT_ID="test1"  EVNTLOG_TEST="TRUE" MQTT_USERNAME="" MQTT_PASSWORD="" TELEGRAF_USER="" TELEGRAF_PASS= "" DB_NAME="gcdo" TELEGRAF_URL="http://telegraf1:8186/write" bin/evntlogs
