# Собираем в гошке
FROM golang:1.24.3 AS build

ENV BIN_FILE=/opt/calendar/calendar-app
ENV CODE_DIR=/go/src/

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/calendar/*

# На выходе тонкий образ
FROM alpine:3.9

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="calendar"
LABEL MAINTAINERS="student@otus.ru"

ENV BIN_FILE="/opt/calendar/calendar-app"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE=/etc/calendar/config.toml
COPY ./configs/config.docker.toml ${CONFIG_FILE}

ENV MIG_PATH=/opt/calendar/migrations/
RUN mkdir ${MIG_PATH}
COPY ./migrations/ ${MIG_PATH}

ENV LOG_PATH=/opt/calendar/logs/
RUN mkdir ${LOG_PATH}

CMD ${BIN_FILE} -config ${CONFIG_FILE}
# ENTRYPOINT ["tail", "-f", "/dev/null"]
