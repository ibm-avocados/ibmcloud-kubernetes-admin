frontend:
	cd client; yarn run build; cd ..;

build: 
	go build -o kubeadmin ./cmd/web/main.go

start: frontend build
	./kubeadmin
	
docker: 
	DOCKER_BUILDKIT=1 docker build -t moficodes/ibm-kubernetes-admin:$(tag) .

push:
	docker push moficodes/ibm-kubernetes-admin:$(tag)

run: frontend
	go run ./cmd/web/main.go

run-local:
	go run ./cmd/web/main.go