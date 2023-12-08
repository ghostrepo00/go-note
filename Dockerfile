################
# BUILD BINARY #
################
# golang:1.18.2-alpine3.16
FROM golang:1.21-bullseye as builder

ENV SUPABASE=https://mogqxgbfrxdeluejpgqw.supabase.co|eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6Im1vZ3F4Z2JmcnhkZWx1ZWpwZ3F3Iiwicm9sZSI6ImFub24iLCJpYXQiOjE3MDEzMDUxODEsImV4cCI6MjAxNjg4MTE4MX0.7c0xsQYSfc8984DJ84Ew-SOdzZvpzeL9qOTeAUYuGYg
ENV CRYPTO_KEY=CrS2t2TrMph4IYi43alYzatKo3V69tPr
ENV CRYPTO_IV_PAD=i4CrMp59rIYWa42V

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
#RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

WORKDIR /src
COPY . .

# Fetch dependencies.
# RUN go get -d -v
RUN go mod download
RUN go mod verify

#CMD go build -v
# go build command with the -ldflags="-w -s" option to produce a smaller binary file by stripping debug information and symbol tables. 
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o /server/server /src/cmd/server

#####################
# MAKE SMALL BINARY #
#####################
FROM scratch

#RUN apk update

# Import from builder.
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

# Copy the executable.
COPY --from=builder /server/server /server/server
COPY --from=builder /src/config/config.json /server/config/config.json
COPY --from=builder /src/logs /server/logs
COPY --from=builder /src/web /server/web
CMD [ "/server/server" ]

