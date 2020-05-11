# mapi
MAPI - Mail API Is a golang based api that sends emails to recipients using the sendgrid email api

## Get started
#Prerequisite's
* Ensure that you have `go`, `docker` & `docker compose ` installed on your machine.
* To install go, see **[Installing Go]** for 
reference
* To install docker, see **[Installing Docker]** for 
reference
### Clone the repo into your go working directory & fetch all dependencies
```shell
git clone https://github.com/tevin-morake/mapi.git
cd mapi && go get -u -v
```

#Do a build of the mail api 
```shell
env GOOS=linux GOARCH=386 go build -v mapi.go
```
[Installing Go]: https://golang.org/doc/install
[Installing Docker]: https://docs.docker.com/get-docker/

## Build an image from the Dockerfile
docker build -t mapi

## Run docker-compose in the background
```bash
docker compose up -d
```

### Check all running containers to verify that your postgres & mapi are running 
```bash
docker ps
```

## Webservice Endpoints
### List all todos for a user
``` bash
POST /email
```
body of req example : 
```bash
{
    "body": "",
    "email": "",
    "subject":""
}
```
Expected response is a 200 OK response with a log entry of the email sent in Postgre SQL

