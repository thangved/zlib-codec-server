FROM golang AS builder

WORKDIR /app

COPY go.sum go.mod ./

RUN go mod download

COPY *.go ./
COPY ./codec-server ./codec-server

RUN CGO_ENABLED=0 GOOS=linux go build codec-server/main.go

FROM scratch

COPY --from=builder /app/main /codec-server

CMD [ "/codec-server" ]