GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BUILD_DIR=build
GEN_DIR=gen
PROTO_OUT_DIR=${GOPATH}/src
KUBECTL=kubectl
DOCKER=docker
PROTOC=protoc
PROGEN=protoc --go_out=plugins=grpc:$(PROTO_OUT_DIR)
    
all: build
build: protobuf orderstorageservice orderprocessservice # order processengine actionhandler orderservice orderstorageervice

protobuf: gen-protobuf fix-protobuf

gen-protobuf:
	mkdir -p $(GEN_DIR)
	$(PROGEN) ./proto/actiondata/*.proto
	$(PROGEN) ./proto/orderdata/*.proto
	$(PROGEN) ./proto/orderstorageservice/*.proto
	$(PROGEN) ./proto/orderprocessservice/*.proto

fix-protobuf:
	sed -i "" -e "s/XXX\(.*\)\`\(.*\)\`/XXX\1\`\2 \datastore:\"-\"\`/" `find gen -name "*.pb.go"`

order:
	$(GOBUILD) -o $(BUILD_DIR)/order -v cmd/order/order.go 
processengine:
	$(GOBUILD) -o $(BUILD_DIR)/processengine -v cmd/processengine/processengine.go
actionhandler:
	$(GOBUILD) -o $(BUILD_DIR)/actionhandler -v cmd/actionhandler/actionhandler.go
orderprocessservice:
	$(GOBUILD) -o $(BUILD_DIR)/orderprocessservice -v cmd/orderprocessservice/orderprocessservice.go
orderstorageservice:
	$(GOBUILD) -o $(BUILD_DIR)/orderstorageservice -v cmd/orderstorageservice/orderstorageservice.go


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

docker-build: build-linux
	$(DOCKER) build -t gcr.io/sap-se-commerce-arch/processengine:latest -f infra/docker/processengine/Dockerfile .
	$(DOCKER) build -t gcr.io/sap-se-commerce-arch/actionhandler:latest -f infra/docker/actionhandler/Dockerfile .

docker-push: docker-build
	$(DOCKER) push gcr.io/sap-se-commerce-arch/processengine:latest
	$(DOCKER) push gcr.io/sap-se-commerce-arch/actionhandler:latest

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