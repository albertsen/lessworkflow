BUILD_DIR=build
PKGPATH=github.com/albertsen/lessworkflow

GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test -v -count=1
GOGET=$(GOCMD) get

KUBECTL=kubectl

DOCKER=docker
DOCKER_COMPOSE=docker-compose
DOCKER_DIR=./infra/docker

PSQL=psql -h localhost -p 5432
DB_NAME=lessworkflow


all: build
build: documentservice

documentservice:
	$(GOBUILD) -o $(BUILD_DIR)/documentservice -v $(PKGPATH)/cmd/documentservice

test: cleardb test-documentservice
test-documentservice:
	$(GOTEST) $(PKGPATH)/cmd/documentservice

import-sample-data:
	curl --header "Content-Type: application/json" -vX POST -d @./data/sample/order.json http://localhost:8000/documents/orders
	curl --header "Content-Type: application/json" -vX POST -d @./data/sample/process.json http://localhost:8000/documents/processdefs

createdb:
	$(PSQL) -U postgres postgres -f sql/create_database.sql 
	$(PSQL) -U postgres $(DB_NAME) -f sql/create_users.sql
	$(PSQL) -U lwadmin $(DB_NAME) -f sql/create_tables.sql

dropdb:
	$(PSQL) -U postgres postgres -c "DROP DATABASE lessworkflow"

cleardb:
	$(PSQL) -U postgres $(DB_NAME) -c "DELETE FROM documents"


clean: 
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

# Cross compilation
build-linux: export CGO_ENABLED=0
build-linux: export GOOS=linux
build-linux: export GOARCH=amd64
build-linux: export BUILD_DIR=build/linux
build-linux: build

docker: build-linux
	$(DOCKER) build -t gcr.io/sap-se-commerce-arch/documentservice:latest -f infra/docker/documentservice/Dockerfile .

docker-compose-up: docker
	cd $(DOCKER_DIR) && $(DOCKER_COMPOSE) up --remove-orphans

docker-compose-down:
	cd $(DOCKER_DIR) && $(DOCKER_COMPOSE) down

docker-compose-restart-services: docker
	cd $(DOCKER_DIR) && $(DOCKER_COMPOSE) stop documentservice && $(DOCKER_COMPOSE) up --no-deps -d documentservice