init-docs:
	swag init -g /NoteService/cmd/app/main.go --parseInternal --parseDependency --dir ./NoteService

run:
	@docker compose up --build

clear: 
	@docker compose down --volumes --remove-orphans