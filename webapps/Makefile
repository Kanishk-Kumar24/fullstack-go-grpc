.PHONY: run client run-all

run:
	go run main.go

run-all:
	@echo "Starting backend and webapp servers..."
	@cmd /c start make run
	@cd ../backend && cmd /c start make run

# A sample client call to test gRPC connection
client:
	go run client/main.go