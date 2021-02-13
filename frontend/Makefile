.PHONY: clean build watch run test

help:
	@echo "usage: make <command>"
	@echo ""
	@echo "commands:"
	@echo "  clean      remove unwanted stuff (node_modules, ...)"
	@echo "  build      build the frontend app"
	@echo "  watch      run the watcher"
	@echo "  test       run the testsuite"
	@echo "  help       display the help message"


clean:
	find . -name 'node_modules' -exec rm -rf {} +
	find . -name '*~' -exec rm -f {} +

build:
	npm run build


watch:
	npm run watch

test:
	npm run test
