frontend:
	cd client; yarn run build; cd ..;

build: 
	go build -o kubes .

start: frontend build
	./kubes
	
docker: 
	DOCKER_BUILDKIT=1 docker build -t moficodes/ibm-kubernetes-admin:$(tag) .

push:
	docker push moficodes/ibm-kubernetes-admin:$(tag)

run:
	go run ./cmd/web/main.go