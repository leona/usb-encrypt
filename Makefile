.PHONY: default
default: build

build-windows:
	@GOOS=windows go build decrypt.go
	@GOOS=windows go build encrypt.go

build-linux:
	@go build decrypt.go
	@go build encrypt.go
	
build-mac:
	@GOOS=darwin go build decrypt.go
	@GOOS=darwin go build encrypt.go
	
build:
	@GOOS=windows go build -ldflags="-s -w" -o ./bin/decrypt-windows.exe ./decrypt.go
	@GOOS=windows go build -ldflags="-s -w" -o ./bin/encrypt-windows.exe ./encrypt.go
	@upx --brute --best --lzma ./bin/decrypt-windows.exe
	@upx --brute --best --lzma ./bin/encrypt-windows.exe
	@GOOS=linux go build -ldflags="-s -w" -o ./bin/decrypt-linux ./encrypt.go
	@GOOS=linux go build -ldflags="-s -w" -o ./bin/encrypt-linux ./decrypt.go
	@GOOS=darwin go build -ldflags="-s -w" -o ./bin/decrypt-mac ./encrypt.go
	@GOOS=darwin go build -ldflags="-s -w" -o ./bin/encrypt-mac ./decrypt.go
	#sudo apt-get install upx-ucl
	
