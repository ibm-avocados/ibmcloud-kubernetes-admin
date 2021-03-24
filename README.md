# IBM Cloud Kubernetes Admin

This tool was made to manage large IBM Cloud accounts and the Kubernetes Clusters in it.

## Pre-Requisite

To run this project we need access to a few variable. 

```
cp .env.example .env
```

Fill these variable with the proper information.

```
# Create a cloudant instance on IBM Cloud
CLOUDANT_USER_NAME=
CLOUDANT_PASSWORD=
CLOUDANT_HOST=
# Required for cron job. Currently under development. Any integer value is fine.
TICKER_PERIOD=
# Creata a sendgrid account. Used for sending email notification.
SENDGRID_API_KEY=
# used to send system notification to admins of deployment.
ADMIN_FROM_EMAIL=
ADMIN_TO_EMAIL=
# Required for IBM Single Sign On
IBM_LOGIN_USER=
IBM_LOGIN_CLIENT_ID=
IBM_LOGIN_CLIENT_SECRET=
IBM_REDIRECT_URI=
# Under development. Any string is fine.
JWT_SECRET=
# Required for Awx. Generate token on awx instance.
AWX_ACCESS_TOKEN=
# Required for posting comment to issue post provision.
GITHUB_ISSUE_REPO=
# vault related secrets
# vault is used for getting secrets
# if vault auth is done via github
VAULT_AUTH_GITHUB_TOKEN=
VAULT_ADDR=
```


## Run

1. Local Development (With prod build)

```
make start
```

This would start the project on port `9000` and go backend will serve the frontend.

2. Docker

```
docker build -t <image>:<version> .
```

```
docker run -p 9000:9000 <image>:<version>
```

## Local Development

### Backend
**Runtime:** Go 1.15

### Frontend
**Runtime:** NodeJS 14