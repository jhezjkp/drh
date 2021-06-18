GO111MODULE=auto
MODULE=drh
BUILD_PATH=./cmd

cleanup:
	echo "=== clean up ==="
	rm -rf ./bin

build:
	go build -o bin/$(MODULE) $(BUILD_PATH)

run:
	go run cmd/main.go


compile:
	echo "Compiling for every OS and Platform"
	#GOOS=linux GOARCH=arm go build -o bin/$(MODULE)-linux-arm $(BUILD_PATH)
	#GOOS=linux GOARCH=arm64 go build -o bin/$(MODULE)-linux-arm64 $(BUILD_PATH)
	#GOOS=freebsd GOARCH=386 go build -o bin/$(MODULE)-freebsd-386 $(BUILD_PATH)
	GOOS=linux GOARCH=amd64 go build -o bin/$(MODULE)-linux-amd64 $(BUILD_PATH)
	GOOS=windows GOARCH=amd64 go build -o bin/$(MODULE).exe $(BUILD_PATH)
	GOOS=darwin GOARCH=amd64 go build -o bin/$(MODULE) $(BUILD_PATH)
all: cleanup compile