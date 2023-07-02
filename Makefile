SHELL = /bin/bash
.PHONY: build

# ------------------------------------------------
# Docker
# ------------------------------------------------
build:
	docker build -t map-build-router:latest build/router

run-router:
	docker run -it --rm \
		-p 1161:1161/udp \
		-v ./examples/router-1.snmprec:/data/router-1.snmprec \
		-e ZBX_SERVER=zabbix-server \
		-e ZBX_HOSTNAME=router-1 \
		map-build-router:latest

run-router-debug:
	docker run -it --rm \
		--entrypoint "" -p 1161:1161/udp \
		-p 1161:1161/udp \
		-v ./examples/router-1.snmprec:/data/router-1.snmprec \
		-e ZBX_SERVER=zabbix-server \
		-e ZBX_HOSTNAME=router-1 \
		map-build-router:latest

# ------------------------------------------------
# Docker compose
# ------------------------------------------------
up:
	docker compose -f ./docker-compose.yml up
	
down:
	docker compose -f ./docker-compose.yml down

# ------------------------------------------------
# CLI commands
# ------------------------------------------------
run:
	go run main.go --name test-map-builder \
		--file examples/mapping.json

run-debug:
	go run main.go --name test-map-builder \
		--file examples/mapping.json \
		--debug

run-outfile:
	go run main.go --name test-map-builder \
		--file examples/mapping.json \
		--output examples/output.json

run-outfile-debug:
	go run main.go --name test-map-builder \
		--file examples/mapping.json \
		--output examples/output.json \
		--debug

run-color:
	go run main.go --name test-map-builder \
		--file examples/mapping.json \
		--color 7AC2E1 \
		--trigger-color EE445B

run-color-debug:
	go run main.go --name test-map-builder \
		--file examples/mapping.json \
		--color 7AC2E1 \
		--trigger-color EE445B \
		--debug

run-dry:
	go run main.go --name test-map-builder \
		--file examples/mapping.json \
		--color 7AC2E1 \
		--trigger-color EE445B \
		--dry-run

run-dry-debug:
	go run main.go --name test-map-builder \
		--file examples/mapping.json \
		--color 7AC2E1 \
		--trigger-color EE445B \
		--dry-run \
		--debug

run-unstack-hosts:
	go run main.go --name test-map-builder \
		--file examples/mapping.json \
		--output examples/output.json \
		--color 7AC2E1 \
		--trigger-color EE445B \
		--stack-hosts false

run-unstack-hosts-debug:
	go run main.go --name test-map-builder \
		--file examples/mapping.json \
		--output examples/output.json \
		--color 7AC2E1 \
		--trigger-color EE445B \
		--stack-hosts false \
		--debug

# - HELPER
help:
	go run main.go --help

# ------------------------------------------------
# Tests
# ------------------------------------------------
create-hosts:
	ZABBIX_URL="http://localhost:4444/api_jsonrpc.php" ZABBIX_USER="Admin" ZABBIX_PWD="zabbix" go run examples/import.go --file examples/zbx_export_hosts.json

test:
	@echo "Running container stack..."
	docker compose -f ./docker-compose.test.yml up -d
	@TIMER=40; \
	echo "Waiting $$TIMER\\s for Zabbix server to initialize"; \
	i=1; \
	while [[ $$i -ne $$TIMER ]]; \
	do \
		echo "$$i / $$TIMER"; \
		sleep 1; \
		i=$$((i+1)); \
	done;
	@echo "Import hosts configuration..."
	make create-hosts
	@echo "Running test..."
	go test ./...
	@echo "Destroying container stack"
	docker compose -f ./docker-compose.test.yml down

coverage:
	@echo "Running container stack..."
	docker compose -f ./docker-compose.test.yml up -d
	@TIMER=40; \
	echo "Waiting $$TIMER\\s for Zabbix server to initialize"; \
	i=1; \
	while [[ $$i -ne $$TIMER ]]; \
	do \
		echo "$$i / $$TIMER"; \
		sleep 1; \
		i=$$((i+1)); \
	done;
	@echo "Import hosts configuration..."
	make create-hosts
	@echo "Running test..."
	go test -coverprofile=coverage.out ./...
	@echo "Formatting coverage report to HTML..."
	go tool cover -html=coverage.out -o=coverage.html
	@echo "Destroying container stack"
	docker compose -f ./docker-compose.test.yml down

down-test:
	docker compose -f ./docker-compose.test.yml down

clean-test-cache:
	go clean -testcache

local-test:
	go test ./... -count=1

local-test-debug:
	go test ./... -count=1 -v