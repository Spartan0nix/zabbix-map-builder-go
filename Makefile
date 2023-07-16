SHELL = /bin/bash
.PHONY: build

# ------------------------------------------------
# Docker
# ------------------------------------------------
build-router:
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
# - Create
run-create:
	go run main.go create --name test-map-builder \
		--file examples/mapping.json

run-create-debug:
	go run main.go create --name test-map-builder \
		--file examples/mapping.json \
		--debug

run-create-outfile:
	go run main.go create --name test-map-builder \
		--file examples/mapping.json \
		--output examples/output.json

run-create-outfile-debug:
	go run main.go create --name test-map-builder \
		--file examples/mapping.json \
		--output examples/output.json \
		--debug

run-create-color:
	go run main.go create --name test-map-builder \
		--file examples/mapping.json \
		--color 7AC2E1 \
		--trigger-color EE445B

run-create-color-debug:
	go run main.go create --name test-map-builder \
		--file examples/mapping.json \
		--color 7AC2E1 \
		--trigger-color EE445B \
		--debug

run-create-dry:
	go run main.go create --name test-map-builder \
		--file examples/mapping.json \
		--color 7AC2E1 \
		--trigger-color EE445B \
		--dry-run

run-create-dry-debug:
	go run main.go create --name test-map-builder \
		--file examples/mapping.json \
		--color 7AC2E1 \
		--trigger-color EE445B \
		--dry-run \
		--debug

run-create-unstack-hosts:
	go run main.go create --name test-map-builder \
		--file examples/mapping.json \
		--output examples/output.json \
		--color 7AC2E1 \
		--trigger-color EE445B \
		--stack-hosts false

run-create-unstack-hosts-debug:
	go run main.go create --name test-map-builder \
		--file examples/mapping.json \
		--output examples/output.json \
		--color 7AC2E1 \
		--trigger-color EE445B \
		--stack-hosts false \
		--debug

# - Generate
run-generate:
	go run main.go generate --host 172.16.80.161 \
		--community router-1 \
		--port 1161

run-generate-debug:
	go run main.go generate --host 172.16.80.161 \
		--community router-1 \
		--port 1161 \
		--debug

run-generate-outfile:
	go run main.go generate --host 172.16.80.161 \
		--community router-1 \
		--port 1161 \
		--output examples/generated_mapping.json

run-generate-outfile-debug:
	go run main.go generate --host 172.16.80.161 \
		--community router-1 \
		--port 1161 \
		--output examples/generated_mapping.json \
		--debug

run-generate-pattern:
	go run main.go generate --host 172.16.80.161 \
		--community router-1 \
		--port 1161 \
		--output examples/generated_mapping.json \
		--trigger-pattern "Interface #INTERFACE(): Link down"

run-generate-pattern-debug:
	go run main.go generate --host 172.16.80.161 \
		--community router-1 \
		--port 1161 \
		--output examples/generated_mapping.json \
		--trigger-pattern "Interface #INTERFACE(): Link down" \
		--debug

run-generate-image:
	go run main.go generate --host 172.16.80.161 \
		--community router-1 \
		--port 1161 \
		--output examples/generated_mapping.json \
		--local-host-image "Firewall_(64)" \
		--remote-host-image "Switch_(64)"

run-generate-image-debug:
	go run main.go generate --host 172.16.80.161 \
		--community router-1 \
		--port 1161 \
		--output examples/generated_mapping.json \
		--local-host-image "Firewall_(64)" \
		--remote-host-image "Switch_(64)" \
		--debug

run-generate-full:
	go run main.go generate --host 172.16.80.161 \
		--community router-1 \
		--port 1161 \
		--output examples/generated_mapping.json \
		--local-host-image "Firewall_(64)" \
		--remote-host-image "Switch_(64)" \
		--trigger-pattern "Interface #INTERFACE(): Link down" \
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
	@echo "Cleaning go test cache"
	make clean-test-cache
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
	@echo "Cleaning go test cache"
	make clean-test-cache
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