# syntax=docker/dockerfile:1

FROM ghcr.io/merliot/device/device-base:latest

WORKDIR /app
COPY . .
RUN go work use .

RUN go build -tags prime -o /hp2430n ./cmd
RUN /hp2430n -uf2

EXPOSE 8000

ENV PORT_PRIME=8000
CMD ["/hp2430n"]
