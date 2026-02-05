init-docs:
	swag init -g cmd/app/main.go --parseInternal --parseDependency --dir ./

run:
	@docker compose up --build

clear: 
	@docker compose down --volumes --remove-orphans