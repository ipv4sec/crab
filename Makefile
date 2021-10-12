

.PHONY: parser
parser:
	docker build -f deployments/parser.Dockerfile -t harbor1.zlibs.com/island/island-parser:alpha .
	docker push harbor1.zlibs.com/island/island-parser:alpha
  
.PHONY: setup
setup:
	docker build -f deployments/setup.Dockerfile -t harbor1.zlibs.com/island/island-setup:alpha .
	docker push harbor1.zlibs.com/island/island-setup:alpha

.PHONY: web
web:
	docker build -f deployments/web.Dockerfile -t harbor1.zlibs.com/island/island-web:alpha .
	docker push harbor1.zlibs.com/island/island-web:alpha

