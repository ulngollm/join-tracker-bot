build:
	CGO_ENABLED=0 go build -o bot

start:
	./bot > log &

deploy:
	scp bot ${SERVER}:${DEPLOY_DIR}
	scp .env ${SERVER}:${DEPLOY_DIR}