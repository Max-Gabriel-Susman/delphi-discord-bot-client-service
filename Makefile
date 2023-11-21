local-start:
	go run main.go

build:
	docker build --tag delphi-discord-bot-client-service .

push: 
	docker push brometheus/delphi-model-service:tagname


update:
	docker build --tag brometheus/delphi-discord-bot-client-service:v0.1.6 .
	docker push brometheus/delphi-discord-bot-client-service:v0.1.6

# grpcurl -plaintext -v localhost:50051 list Greeter

# grpcurl -plaintext -d '{"name": "prometheus"}' localhost:50052 Greeter/SayHello