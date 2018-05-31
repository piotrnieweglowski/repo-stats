# An application to fetch statistics for github repositories
# of particular user
FROM golang:1.10

COPY ./ /go/src/github.com/piotrnieweglowski/repo-stats
WORKDIR /go/src/github.com/piotrnieweglowski/repo-stats

RUN go get ./
RUN go build
CMD ["repo-stats"]