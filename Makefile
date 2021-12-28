
all: .PHONY
.PHONY: dfuse

dfuse:
	cd cmd/dfuseplugin && go build main.go
