
.PHONY: setup
setup:
	docker build -f deployments/parser.Dockerfile -t harbor1.zlibs.com/island/island-parser:alpha .
	docker push harbor1.zlibs.com/island/island-parser:alpha