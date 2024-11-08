EXEC_DIR=cmd/playground
BINARY_DIR=${EXEC_DIR}/bin
BINARY_NAME=playground
DOCKER_IMAGE_NANE=observability-playground

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
	docker-compose -f ./build/docker-compose.yml up -d --build

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