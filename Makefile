# Makefile

# Команда для запуска Docker Compose
start:
	docker compose up -d

stop:
	docker compose down

reset:
	docker compose down
	docker rmi test-spb-go-bank-tickets-backend:latest
	docker rmi test-spb-go-bank-vue-frontend:latest