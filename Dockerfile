# syntax=docker/dockerfile:1

FROM ghcr.io/merliot/device/device-base:latest

WORKDIR /app
RUN git clone https://github.com/merliot/device.git
RUN go work use device

WORKDIR /app/hp2430n

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go work use .
RUN go build -tags prime -o /hp2430n ./cmd/
RUN go run ../device/cmd/uf2-builder -target nano-rp2040 -model hp2430n

EXPOSE 8000

ENV PORT_PRIME=8000
CMD ["/hp2430n"]
