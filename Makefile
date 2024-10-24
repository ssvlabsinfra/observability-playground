EXEC_DIR=cmd/p2p
BINARY_DIR=${EXEC_DIR}/bin
BINARY_NAME=p2p
DOCKER_IMAGE_NANE=p2p-observability

default: build run

.PHONY: build
build:
	go build -o ${BINARY_DIR}/${BINARY_NAME} ${EXEC_DIR}/main.go

.PHONY: run
run: build
	@cd ${BINARY_DIR} && ./${BINARY_NAME}

########## DOCKER
.PHONY: docker-build
docker-build:
	docker build -t ${DOCKER_IMAGE_NANE} -f build/Dockerfile .

.PHONY: docker-run
docker-run:
	docker run ${DOCKER_IMAGE_NANE}

.PHONY: docker-compose-up
docker-compose-up:
	docker-compose -f ./build/docker-compose.yml up -d

.PHONY: docker-compose-down
docker-compose-down:
	docker-compose -f ./build/docker-compose.yml down
##########

.PHONY: clean
clean:
	go clean

.PHONY: test
test:
	go test ./...

.PHONY: test_coverage
test_coverage:
	go test ./... -coverprofile=coverage.out

.PHONY: dep
dep:
	go mod download

.PHONY: vet
vet:
	go vet ./...