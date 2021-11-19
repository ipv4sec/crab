
.PHONY: setup
setup:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o assets/bin/setup cmd/setup/main.go
	docker build -f deployments/setup.Dockerfile -t harbor1.zlibs.com/island/island-setup:0.1 .
	docker push harbor1.zlibs.com/island/island-setup:0.1

.PHONY: api
api:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o assets/bin/parser cmd/parser/main.go
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o assets/bin/api cmd/api/main.go
	docker build -f deployments/api.Dockerfile -t harbor1.zlibs.com/island/island-api:0.1 .
	docker push harbor1.zlibs.com/island/island-api:0.1

.PHONY: ui
ui:
	docker build -f deployments/ui.Dockerfile -t harbor1.zlibs.com/island/island-ui:0.1 web
	docker push harbor1.zlibs.com/island/island-ui:0.1

.PHONY: scheduler
scheduler:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o assets/bin/scheduler cmd/scheduler/main.go
	docker build -f deployments/scheduler.Dockerfile -t harbor1.zlibs.com/island/island-scheduler:0.1 .
	docker push harbor1.zlibs.com/island/island-scheduler:0.1