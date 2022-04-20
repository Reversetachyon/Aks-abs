# Base image
FROM golang:1.16-alpine 

#create working directory in docker image
WORKDIR /app

#copy file in this directory to docker image working directory
COPY go.mod ./
COPY go.sum ./
COPY test.json ./

# install dependecies in module file(go.mod)
RUN go mod download

#copy application file to image
COPY *.go ./

#build aplication to name docker-be and located in the root of the filesystem of the image
RUN go build -o /docker-be

EXPOSE 8080

#Execute command
CMD [ "/docker-be" ]

# docker build --tag 6131305010/gin-azblob:v1 .   ****build command to docker image
              # [host-port] [container port]
# docker run --publish 8080:8080 --name testazure 6131305010/gin-azblob:v1 ***run image command 