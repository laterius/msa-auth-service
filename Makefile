build:
	docker build -f docker/Dockerfile . -t 34234247632/auth-service:v1.1

push:
	docker push 34234247632/auth-service:v1.1

docker-start:
	cd docker && docker-compose up -d

docker-stop:
	cd docker && docker-compose down

newman:
	newman run postman/collection.json

