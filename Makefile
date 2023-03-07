
dev-server:
	@echo "\n\t🧠\n"
	go run .

dev-client:
	@echo "\n\t🧠\n"
	go run ./cmd/client

run-build-server:
	./bin/server/go-serge

build-server:
	@echo "\n\t🐵\n"
	go build . && mv go-serge ./bin/server

build-client:
	@echo "\n\t🐵\n"
	go build cmd/client/main.go && mv main ./bin/client/client

build-testmod:
	cd testmod && go build -buildmode=plugin . && mv testmod.so ../mods/testmod.so

prepare-build:
	cp ./mods/testmod.so ./bin/server/mods/testmod.so && cp ./mods/testmod.so ./bin/client/mods/testmod.so  

.PHONY: run-server run-client