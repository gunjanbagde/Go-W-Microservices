#base image (Level 1)

FROM golang:1.21-alpine AS builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api

RUN chmod +x /app/brokerApp

#building tiny docker image (Level 2)
# this builds code on 1 docker image in Level 1 then build smaller dockerimage and just copy the executables 

FROM alpine:latest

RUN mkdir /app

#from builder brokerapp copy to app in level 2
COPY --from=builder /app/brokerApp /app

CMD [ "/app/brokerApp" ]