# syntax=docker/dockerfile:1

FROM ghcr.io/merliot/device:main

WORKDIR /app
COPY . .

RUN go build -tags prime -o /hp2430n ./cmd
RUN /hp2430n -uf2

EXPOSE 8000

ENV PORT_PRIME=8000
CMD ["/hp2430n"]
