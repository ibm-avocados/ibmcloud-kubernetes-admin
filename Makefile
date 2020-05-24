build-frontend:
	cd client; yarn run build; cd ..;

build: 
	go build -o kubeadmin ./cmd/web/main.go

build-start: build-frontend build
	./kubeadmin

start: build
	./kubeadmin
	
run:
	./kubeadmin

docker: 
	DOCKER_BUILDKIT=1 docker build -t moficodes/ibm-kubernetes-admin:$(tag) .

push:
	docker push moficodes/ibm-kubernetes-admin:$(tag)

run-local:
	go run ./cmd/web/main.go