# IBM Cloud Kubernetes Admin

This tool was made to manage large IBM Cloud accounts and the Kubernetes Clusters in it.

## Pre-Requisite

### macOS:

1. Set up your development environment
-Open terminal and cd to a desired location and clone the Git repository 
```
git clone https://github.com/ibm-avocados/ibmcloud-kubernetes-admin.git
```

-Install Homebrew with 

```
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

-Install Go with 

```
brew install go
```

and then use 

```
go mod download
```

to install the required Go dependencies.

-Install Node with
```
brew install node

```
-Install yarn globally with 

```
npm install -g yarn
```
2. Configure .env settings

Use cd to enter the root folder of the repository. 

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
Enter

```
make build-backend
```
and

```
make build-frontend
```

This would start the project on port `9000` and go backend will serve the frontend.
Access port `9000` on a web browser and enter your IBM Cloud credentials.

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