GOOS="windows" go build -ldflags "-s -w -H=windowsgui" -o "bin/windows_Live..exe"
# GOOS="darwin"  go build -ldflags "-s -w"               -o "bin/mac_Live."
GOOS="linux"   go build -ldflags "-s -w"               -o "bin/linux_Live."
