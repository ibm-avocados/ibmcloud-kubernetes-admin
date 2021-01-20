build-frontend:
	cd client; yarn; yarn run build; cd ..;

build-backend:
	go build -o kubeadmin ./cmd/web/main.go

build: build-frontend build-backend

start: build
	./kubeadmin

run:
	./kubeadmin

docker: 
	docker build -t moficodes/ibm-kubernetes-admin:$(tag) -f docker/Dockerfile.web .

push:
	docker push moficodes/ibm-kubernetes-admin:$(tag)

run-local:
	go run ./cmd/web/main.go; 