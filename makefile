registry:=ghcr.io
username:=fajar-islami
image:=go-simple-user-crud
tags:=latest


.PHONY: help
help: ## Show help command
	@printf "Makefile Command\n";
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: check-swag
check-swag: ## Check if swag command already exist
	command -v swag >/dev/null 2>&1 || { go install github.com/swaggo/swag/cmd/swag@v1.8.10; }

.PHONY: docs 
docs: ## Generate Documents
	swag init -g ./internal/delivery/http/main.go --output ./docs/

.PHONY: readenv
readenv: ## not work
	export $(cat .env | xargs -L 1)

.PHONY: migrate
migrate: ## Create Migrations file, example : make migrate name="xxxx"
	@if [ -z "${name}" ]; then \
		echo "Error: name is required \t example : make migrate name="name_file_migration";" \
		exit 1; \
	fi
	migrate create -ext sql -dir migrations ':hammer: ${name}'



migrate-up: ## Up migration
	go run cmd/migrate/main.go

migrate-rollback: ## Up rollback
	go run cmd/migrate/main.go -rollback

build:
	go build -o app cmd/main.go

run:
	./app

runlocal:
	go run cmd/main.go

dockerbuild:
	docker build --rm -t ${registry}/${username}/${image}:${tags} .
	docker image prune --filter label=stage=dockerbuilder -f

dockerun:
	docker run --name ${image} -p 8080:8080 ${registry}/${username}/${image}:${tags}

dockerup: ## up compose image
	docker compose -f docker-compose-app.yaml up -d

dockerlogs: ## logs compose image
	docker compose -f docker-compose-app.yaml logs -f

dockerstop: ## stop compose image
	docker compose -f docker-compose-app.yaml stop

dockerdown: ## rm compose image
	docker compose -f docker-compose-app.yaml down -v


dockerrm:
	docker rm ${registry}/${username}/${image}:${tags} -f
	docker rmi ${registry}/${username}/${image}:${tags}

dockerenter:
	docker exec -it ${image} bash

dockerenter-img:
	docker run -it --entrypoint sh  ${registry}/${username}/${image}:${tags}

dc-check:
	 docker compose -f docker-compose-app.yaml config
	 

push-image: dockerbuild
	docker push ${registry}/${username}/${image}:${tags}

flysecret:
	flyctl secrets set $(cat .env | xargs)

flylist:
	flyctl secrets list