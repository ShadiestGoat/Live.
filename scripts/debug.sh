clear && go build -o "bin/debug" -ldflags "-X main.Debug=t" && ./bin/debug