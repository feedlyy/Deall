FROM golang:1.19-bullseye as builder

# ENV GO111MODULE=on
WORKDIR /app
COPY go.mod .
COPY go.sum .
COPY config.json .

ARG PAT=""
RUN go mod download

COPY . .
ARG VERSION=development
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
	-ldflags="-w -s" -o build/auth *.go


FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/build /app
WORKDIR /app
COPY config.json /app/
CMD ["./auth"]