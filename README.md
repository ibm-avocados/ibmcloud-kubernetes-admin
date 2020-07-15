# IBM Cloud Kubernetes Admin

This tool was made to manage large IBM Cloud accounts and the Kubernetes Clusters in it.

## Run

1. Local Development (Quick Changes)

```
go run cmd/web/main.go
```

```
cd client
yarn start
```

This would run the frontend on port `3000` and go backend on port `9000`

2. Local Development (With prod build)

```
cd client
yarn build
cd ..
```

From root of the project.

```
go run cmd/web/main.go
```

This would start the project on port `9000` and go backend will serve the frontend.

3. Docker

```
docker build -t <image>:<version> .
```

```
docker run -p 9000:9000 <image>:<version>
```

## Development

### Fronend

This uses React

### Backend

Backend is written in Go 1.13
