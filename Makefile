LDFLAGS="-s -w"
INPUT="./main.go"
OUTPUT="syo.exe"

build:
	go.exe build -ldflags=$(LDFLAGS) -o $(OUTPUT) $(INPUT)

build_small:
	go.exe build -ldflags=$(LDFLAGS) -o $(OUTPUT) $(INPUT)
	upx.exe --brute $(OUTPUT)