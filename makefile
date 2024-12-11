app ?= RefInator
folderPath ?= Cards

SHELL=/bin/zsh

build:
	go build -o build/linux/${app} cmd/${app}/main.go
mac_build:
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o build/mac/${app} cmd/${app}/main.go
run_build:
	rm -f build/linux/${app} 
	go build -o build/linux/${app} cmd/${app}/main.go 
	build/linux/${app} --folderPath=${folderPath}
