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

PSQL_HOST=localhost
PSQL_PORT=5432
PSQL=psql -h localhost -p $(PSQL_PORT)
DB_NAME=lessworkflow


build: documentservice

documentservice:
	$(GOBUILD) -o $(BUILD_DIR)/documentservice -v $(PKGPATH)/cmd/documentservice
processservice:
	$(GOBUILD) -o $(BUILD_DIR)/processservice -v $(PKGPATH)/cmd/processservice
processengine:
	$(GOBUILD) -o $(BUILD_DIR)/processengine -v $(PKGPATH)/cmd/processengine


load-sample-data:
	curl --header "Content-Type: application/json" -v -d @./data/sample/order.json http://localhost:5984/orders
	curl --header "Content-Type: application/json" -v -d @./data/sample/processdef.json http://localhost:5984/processdefs


test: cleardb test-documentservice test-processservice

test-documentservice:
	$(GOTEST) $(PKGPATH)/cmd/documentservice

test-processservice: load-sample-data
	$(GOTEST) $(PKGPATH)/cmd/processservice

clean: 
	rm -rf $(BUILD_DIR)
	$(GOCLEAN)
	
# Cross compilation
build-linux: export CGO_ENABLED=0
build-linux: export GOOS=linux
build-linux: export GOARCH=amd64
build-linux: export BUILD_DIR=build/linux
build-linux: build

docker: build-linux
	$(DOCKER) build -t gcr.io/sap-se-commerce-arch/documentservice:latest -f infra/docker/documentservice/Dockerfile .

docker-start: docker
	cd $(DOCKER_DIR) && $(DOCKER_COMPOSE) up --remove-orphans

docker-stop:
	cd $(DOCKER_DIR) && $(DOCKER_COMPOSE) down

restart-services: docker
	cd $(DOCKER_DIR) && $(DOCKER_COMPOSE) stop documentservice && \
		$(DOCKER_COMPOSE) up --no-deps -d documentservice

setup:
	go get -u github.com/mongodb/mongo-go-driver/mongo \
		github.com/gorilla/mux \
		github.com/satori/go.uuid \
		github.com/google/go-cmp/cmp \
		gotest.tools/assert \
		github.com/streadway/amqp