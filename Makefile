LDFLAGS="-s -w"
INPUT="./main.go"
OUTPUT="syo.exe"

build:
	go.exe build -ldflags=$(LDFLAGS) -o $(OUTPUT) $(INPUT)

build_small: build
	upx.exe --brute $(OUTPUT)