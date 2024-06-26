FROM golang:1.22.2-bookworm
COPY go.* /go/src/github.com/flowerinthenight/oomkill-watch/
COPY *.go /go/src/github.com/flowerinthenight/oomkill-watch/
WORKDIR /go/src/github.com/flowerinthenight/oomkill-watch/
RUN GOOS=linux go build -v -trimpath -o oomkill-watch .

FROM google/cloud-sdk:473.0.0-slim
RUN apt-get install -y kubectl google-cloud-sdk-gke-gcloud-auth-plugin
WORKDIR /app/
COPY --from=0 /go/src/github.com/flowerinthenight/oomkill-watch/oomkill-watch .
ENTRYPOINT ["/app/oomkill-watch"]
CMD ["-slack=''"]
