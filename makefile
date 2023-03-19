registry:=ghcr.io
username:=fajar-islami
image:=go-simple-user-crud
tags:=latest


.PHONY: help
help: ## Show help command
	@printf "Makefile Command\n";
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


.PHONY: docs 
docs: ## Generate Documents
	swag init -g ./cmd/main.go --output docs


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
	docker run --name ${image}  -p 8080:8080 ${image} 

dockerrm:
	docker rm ${image} -f
	docker rmi ${image}

dockeenter:
	docker exec -it ${image} bash

dc-check:
	 docker compose -f docker-compose-app.yaml config
	 

push-image:
	docker build -t ${registry}/${username}/${image}:${tags} .
	export CR_PAT=${CR_PAT}
	echo ${CR_PAT} | docker login ${registry} -u ${USERNAME} --password-stdin
	docker push ${registry}/${username}/${image}:${tags}


readenv:
	export $(cat .env2 | xargs)
