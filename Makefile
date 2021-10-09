
.PHONY: setup
setup:
	docker build -f deployments/setup.Dockerfile -t harbor1.zlibs.com/island/island-setup:alpha .
	docker push harbor1.zlibs.com/island/island-setup:alpha