.PHONY: parser
parser:
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o assets/bin/parser cmd/parser/main.go
    docker build -f deployments/parser.Dockerfile -t harbor1.zlibs.com/island/island-parser:alpha .
	docker push harbor1.zlibs.com/island/island-parser:alpha

.PHONY: setup
setup:
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o assets/bin/setup cmd/setup/main.go
    docker build -f deployments/setup.Dockerfile -t harbor1.zlibs.com/island/island-setup:alpha .
    docker push harbor1.zlibs.com/island/island-setup:alpha

.PHONY: web
web:
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o assets/bin/web cmd/web/main.go
    docker build -f deployments/web.Dockerfile -t harbor1.zlibs.com/island/island-web:alpha .
	docker push harbor1.zlibs.com/island/island-web:alpha

.PHONY: ui
ui:
	docker build -f deployments/ui.Dockerfile -t harbor1.zlibs.com/island/island-ui:alpha web
	docker push harbor1.zlibs.com/island/island-ui:alpha

.PHONY: status
status:
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o assets/bin/status cmd/status/main.go
    docker build -f deployments/status.Dockerfile -t harbor1.zlibs.com/island/island-status:alpha .
	docker push harbor1.zlibs.com/island/island-status:alpha