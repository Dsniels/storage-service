FROM golang:1.24.5-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN  go mod download

COPY . .

RUN go build -o store-service

FROM build AS prod

WORKDIR /prod

COPY --from=build /app/store-service ./

EXPOSE 8080

CMD [ "/prod/store-service" ]