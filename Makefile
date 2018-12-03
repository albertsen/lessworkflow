GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BUILD_DIR=build
BINARY_PUBLISH=publish
BINARY_PROCESSENGINE=processengine
BINARY_ACTIONHANDLER=actionhandler
KUBECTL=kubectl
DOCKER=docker
    
all: build
build: publish processengine actionhandler

publish:
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_PUBLISH) -v cmd/publish/publish.go 
processengine:
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_PROCESSENGINE) -v cmd/processengine/processengine.go
actionhandler:
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_ACTIONHANDLER) -v cmd/actionhandler/actionhandler.go

clean: 
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

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