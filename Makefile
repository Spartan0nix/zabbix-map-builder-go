SHELL = /bin/bash
FUNCTION = BenchmarkCheckGenerateRequiredFlag
FOLDER = cmd

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
test: up-test run-test down-test

coverage: 
	make up-test
	make run-coverage
	@echo "Formatting coverage report to HTML..."
	go tool cover -html=coverage.out -o=coverage.html
	make down-test

up-test:
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

down-test:
	@echo "Destroying container stack"
	docker compose -f ./docker-compose.test.yml down

clean-test-cache:
	go clean -testcache

run-test:
	@echo "Running test..."
	go test ./... -args -test.run "^Test.+$$"
	go test ./... -args -test.run "^TestParseCdpCache$"
	
run-test-debug:
	@echo "Running test...(with verbose output)"
	go test ./... -v -args -test.run "^Test.+$$"

run-coverage:
	@echo "Running test...(with coverage)"
	go test -cover -coverprofile=coverage.out ./... -args -test.run "^Test.+$$"

run-coverage-debug:
	@echo "Running test...(with coverage + verbose output)"
	go test -cover -coverprofile=coverage.out ./... -v -args -test.run "^Test.+$$"

# ------------------------------------------------
# Bench
# ------------------------------------------------
# -run=^$                   : prevent Test functions from running
# -bench ^Benchmark.+$      : run functions starting with Benchmark
# ./...                     : run bench for every packages and sub-packages
# -args                     : Since cpuprofile can't be used for multiple packages, 'args' allow to pass argument for each individual packages.
#                             Kind a kike : foreach(pkg in pkgs) -> go test pkg -arg -> output pkg/cpu.prof ; pkg/mem.prof
#     -test.cpuprofile cpu.prof : output bench stats to a 'cpu.prof' file
#     -test.benchmem            : print memory allocations for benchmarks
#     -test.memprofile mem.prof : output meme bench stats to a 'mem.prof' file
run-bench:
	@echo "Running benchmarks..."
	go test -run=^$$ -bench ^Benchmark.+$$ ./... -args -test.benchmem

run-bench-debug:
	@echo "Running benchmarks...(with verbose output)"
	go test -run=^$$ -bench ^Benchmark.+$$ ./... -v -args -test.benchmem

rm-bench:
	rm cmd/*.prof
	rm internal/api/*.prof
	rm internal/app/*.prof
	rm internal/logging/*.prof
	rm internal/map/*.prof
	rm internal/utils/*.prof
	rm internal/snmp/*.prof

bench-package:
	go test -run=^$$ \
        -bench "^${FUNCTION}$$" ${FOLDER}/*.go \
        -v \
        -cpuprofile cpu.prof \
        -benchmem \
        -memprofile mem.prof  \
        -count 6

profile-cpu:
	go tool pprof cpu.prof

profile-mem:
	go tool pprof mem.prof

bench-old:
	go test -run=^$$ \
		-bench ^${FUNCTION}$$ ${FOLDER}/*.go \
		-count 6 > old.txt

bench-new:
	go test -run=^$$ \
		-bench ^${FUNCTION}$$ ${FOLDER}/*.go \
		-count 6 > new.txt

bench-compare:
	benchstat old.txt new.txt

flame-cpu:
	go tool pprof -http="localhost:8081" cpu.prof

flame-mem:
	go tool pprof -http="localhost:8081" mem.prof
	
# ------------------------------------------------
# Utils
# ------------------------------------------------
create-hosts:
	ZABBIX_URL="http://localhost:4444/api_jsonrpc.php" ZABBIX_USER="Admin" ZABBIX_PWD="zabbix" go run examples/import.go --file examples/zbx_export_hosts.json
