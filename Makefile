.PHONY: parser
parser:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o assets/bin/parser cmd/parser/main.go
	docker build -f deployments/parser.Dockerfile -t harbor1.zlibs.com/island/island-parser:alpha .
	docker push harbor1.zlibs.com/island/island-parser:alpha

.PHONY: setup
setup:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o assets/bin/setup cmd/setup/main.go
	docker build -f deployments/setup.Dockerfile -t harbor1.zlibs.com/island/island-setup:alpha .
	docker push harbor1.zlibs.com/island/island-setup:alpha

.PHONY: api
web:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o assets/bin/api cmd/api/main.go
	docker build -f deployments/api.Dockerfile -t harbor1.zlibs.com/island/island-api:alpha .
	docker push harbor1.zlibs.com/island/island-api:alpha

.PHONY: ui
ui:
	docker build -f deployments/ui.Dockerfile -t harbor1.zlibs.com/island/island-ui:alpha web
	docker push harbor1.zlibs.com/island/island-ui:alpha