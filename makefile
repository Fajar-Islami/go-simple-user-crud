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


.PHONY: migrate
migrate: ## Create Migrations file
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

dockerbuild:
	docker build --rm -t ${registry}/${username}/${image}:${tags} .
	docker image prune --filter label=stage=dockerbuilder -f

dockerun:
	docker run --name ${registry}/${username}/${image}:${tags}  -p 8080:8080 ${image} 

dockerrm:
	docker rm ${registry}/${username}/${image}:${tags} -f
	docker rmi ${registry}/${username}/${image}:${tags}

dockeenter:
	docker exec -it ${registry}/${username}/${image}:${tags} bash

dc-check:
	 docker compose -f docker-compose-app.yaml config
	 

push-image: dockerbuild
	docker push ${registry}/${username}/${image}:${tags}


readenv:
	export $(cat .env | xargs)

flysecret:
	flyctl secrets set $(cat .env | xargs)

flylist:
	flyctl secrets list