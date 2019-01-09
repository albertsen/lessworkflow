GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test -v -count=1
GOGET=$(GOCMD) get
BUILD_DIR=build
GEN_DIR=gen
PROTO_OUT_DIR=${GOPATH}/src
KUBECTL=kubectl
DOCKER=docker
PROTOC=protoc
PROGEN=$(PROTOC) --go_out=plugins=grpc:$(PROTO_OUT_DIR)
PKGPATH=github.com/albertsen/lessworkflow
PSQL=psql -h localhost -p 5432


all: build
build: protobuf orderstorageservice orderprocessservice processdefservice # order processengine actionhandler orderservice orderstorageervice

protobuf:
	mkdir -p $(GEN_DIR)
	$(PROGEN) ./proto/actiondata/*.proto
	$(PROGEN) ./proto/orderdata/*.proto
	$(PROGEN) ./proto/orderstorageservice/*.proto
	$(PROGEN) ./proto/orderprocessservice/*.proto
	$(PROGEN) ./proto/processdef/*.proto
	$(PROGEN) ./proto/processdefservice/*.proto

order:
	$(GOBUILD) -o $(BUILD_DIR)/order -v cmd/order
processengine:
	$(GOBUILD) -o $(BUILD_DIR)/processengine -v $(PKGPATH)/cmd/processengine
actionhandler:
	$(GOBUILD) -o $(BUILD_DIR)/actionhandler -v $(PKGPATH)/cmd/actionhandler
orderstorageservice:
	$(GOBUILD) -o $(BUILD_DIR)/orderstorageservice -v $(PKGPATH)/cmd/orderstorageservice
orderprocessservice:
	$(GOBUILD) -o $(BUILD_DIR)/orderprocessservice -v $(PKGPATH)/cmd/orderprocessservice
processdefservice:
	$(GOBUILD) -o $(BUILD_DIR)/processdefservice -v $(PKGPATH)/cmd/processdefservice

test: test-orderstorageservice test-processdefservice
test-orderstorageservice:
	$(GOTEST) $(PKGPATH)/cmd/orderstorageservice
test-orderprocessservice:
	$(GOTEST) $(PKGPATH)/cmd/orderprocessservice
test-processdefservice:
	$(GOTEST) $(PKGPATH)/cmd/processdefservice

patch:
	cp -rv patch/* ${GOPATH}/src

createdb:
	$(PSQL) -U postgres postgres -f sql/create_database.sql 
	$(PSQL) -U postgres lessworkflow -f sql/create_users.sql
	$(PSQL) -U lwadmin lessworkflow -f sql/create_tables.sql

dropdb:
	$(PSQL) -U postgres postgres -c "DROP DATABASE lessworkflow"

clean: 
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -rf $(GEN_DIR)

# Cross compilation
build-linux: export CGO_ENABLED=0
build-linux: export GOOS=linux
build-linux: export GOARCH=amd64
build-linux: export BUILD_DIR=build/linux
build-linux: build

docker: build-linux
	$(DOCKER) build -t gcr.io/sap-se-commerce-arch/orderstorageservice:latest -f infra/docker/orderstorageservice/Dockerfile .
	$(DOCKER) build -t gcr.io/sap-se-commerce-arch/processdefservice:latest -f infra/docker/processdefservice/Dockerfile .
	$(DOCKER) build -t gcr.io/sap-se-commerce-arch/orderprocessservice:latest -f infra/docker/orderprocessservice/Dockerfile .

docker-push: docker-build
	$(DOCKER) push gcr.io/sap-se-commerce-arch/orderstorageservice:latest
	$(DOCKER) push gcr.io/sap-se-commerce-arch/processdefservice:latest
	$(DOCKER) push gcr.io/sap-se-commerce-arch/orderprocessservice:latest

deploy: docker-push
	$(KUBECTL) -n lw apply -f infra/k8s/lessworkflow/deployments/processengine.yaml
	$(KUBECTL) -n lw apply -f infra/k8s/lessworkflow/deployments/actionhandler.yaml

delete:
	$(KUBECTL) -n lw delete --ignore-not-found=true -f infra/k8s/lessworkflow/deployments/processengine.yaml
	$(KUBECTL) -n lw delete --ignore-not-found=true -f infra/k8s/lessworkflow/deployments/actionhandler.yaml

redeploy: delete deploy

setup-kube:
	$(KUBECTL) apply -f https://raw.githubusercontent.com/stakater/Reloader/master/deployments/kubernetes/reloader.yaml
	$(KUBECTL) create clusterrolebinding cluster-admin-binding --clusterrole=cluster-admin --user=$(gcloud config get-value core/account)
	$(KUBECTL) -n nats apply -f infra/k8s/nats/service-account.yaml
	$(KUBECTL) -n nats apply -f https://raw.githubusercontent.com/nats-io/nats-operator/master/deploy/role.yaml
	$(KUBECTL) -n nats apply -f https://raw.githubusercontent.com/nats-io/nats-operator/master/deploy/deployment.yaml
	$(KUBECTL) -n nats apply -f infra/k8s/nats/cluster.yaml

cleanup-kube:
	$(KUBECTL) -n nats delete --ignore-not-found=true -f infra/k8s/nats/cluster.yaml
	$(KUBECTL) -n nats delete --ignore-not-found=true -f https://raw.githubusercontent.com/nats-io/nats-operator/master/deploy/deployment.yaml
	$(KUBECTL) -n nats delete --ignore-not-found=true -f https://raw.githubusercontent.com/nats-io/nats-operator/master/deploy/role.yaml
	$(KUBECTL) -n nats delete --ignore-not-found=true -f infra/k8s/nats/service-account.yaml