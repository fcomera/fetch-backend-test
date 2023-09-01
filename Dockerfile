FROM golang:latest AS base

ENV CGO_ENABLED 1

WORKDIR /app
COPY . .
RUN go mod download && go mod verify

# Obtained through docker official documentation page
RUN CGO_ENABLED=1 GOOS=linux go build -o /prodapp -a -ldflags '-linkmode external -extldflags "-static"' .

FROM base AS test
ENV MODE=TEST
ENV CGO_ENABLED 1
CMD ["go", "test", "-v", "-run", "TestFiberApp"]

FROM scratch AS prod

WORKDIR /

COPY --from=base /prodapp /prodapp
EXPOSE 3000

ENTRYPOINT ["/prodapp"]