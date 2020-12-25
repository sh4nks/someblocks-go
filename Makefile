.PHONY: clean build watch run test

help:
	@echo "usage: make <command>"
	@echo ""
	@echo "commands:"
	@echo "  clean      remove unwanted stuff"
	@echo "  build      build the app"
	@echo "  watch      run watcher"
	@echo "  test       run the testsuite"
	@echo "  help       display the help message"


clean:
	find . -name 'tmp' -exec rm -f {} +
	find . -name '*~' -exec rm -f {} +
	go clean


build:
	go build


watch:
	air -c .air.toml


run:
	go build && ./someblocks server


test:
	go test
