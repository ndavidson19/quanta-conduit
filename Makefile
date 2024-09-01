# File: Makefile

.PHONY: build test docker-build docker-push deploy

build:
	go build -o rest-ingestion-service ./cmd/rest-ingestion-service

lint:
	golangci-lint run

test:
	./run_tests.sh

docker-build:
	docker build -t your-registry/rest-ingestion-service:latest -f deployments/Dockerfile .

docker-push:
	docker push your-registry/rest-ingestion-service:latest

deploy:
	kubectl apply -f deployments/kubernetes/deployment.yaml
	kubectl apply -f deployments/kubernetes/service.yaml

setup-pre-commit:
	./scripts/setup_pre_commit.sh

# Targets with skip options
test-skip:
	@echo "skiping tests"

lint-skip:
	@echo "skiping linting"

build-skip:
	docker build -t rest-ingestion-service:latest -f deployments/Dockerfile . --no-cache

deploy-skip:
	kubectl apply -f deployments/kubernetes/deployment.yaml --force

# CI target that runs all checks
ci: lint test build

# CI target that skipes all checks
ci-skip: lint-skip test-skip build-skip

all: build test docker-build docker-push deploy setup-pre-commit ci test-skip lint-skip build-skip deploy-skip ci-skip
