# docker build --file build/docker/Dockerfile .

FROM golang:1.13.2-alpine AS backend

COPY . /go/src/github.com/emmercm/photos
WORKDIR /go/src/github.com/emmercm/photos

RUN set -euo pipefail && \
    apk --update add --no-cache gcc musl-dev && \
    go mod download && \
    go build -o /go/bin/photos github.com/emmercm/photos/cmd/server


FROM node:10.16.3-alpine AS frontend

COPY web /

RUN npm install && \
    npm run build


FROM alpine:3.10.2

COPY --from=backend /go/bin/photos /
COPY migrations /migrations

COPY --from=frontend /build /web/build

EXPOSE 8080

CMD ["/photos"]
