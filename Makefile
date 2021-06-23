build-frontend:
	cd client; yarn; yarn run build; cd ..;

build-backend:
	go build -o kubeadmin ./cmd/kubeadmin/main.go

build-user-ui:
	cd user-ui; yarn; yarn run build; cd ..;

build-user:
	go build -o user ./cmd/user/main.go

run-user:
	source .env && go run cmd/user/main.go

start-user: build-user-ui build-user
	source .env && ./user

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