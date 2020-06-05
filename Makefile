dist: $(shell find . -name "*.go")
	mkdir -p dist
	gox -ldflags="-s -w" -osarch="darwin/amd64 linux/386 linux/amd64 linux/arm freebsd/amd64 windows/amd64 windows/386" -output="dist/jqview_{{.OS}}_{{.Arch}}"
