up:
	docker-compose up -d
down:
	docker-compose down


gateway:
	cd api-gateway && go run .

user:
	cd user-service && go run .

order:
	cd order-service && go run .
