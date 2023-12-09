################
# BUILD BINARY #
################
FROM golang:1.21-bullseye as builder

WORKDIR /src
COPY . .

# Install nodejs and create tailwind css
RUN apt-get update && apt-get install -y curl
RUN curl -sL https://deb.nodesource.com/setup_20.x | bash -
RUN apt-get install -y nodejs
RUN npm install
RUN npx tailwindcss -i ./web/assets/style/tailwind.input.css -o ./web/assets/style/tailwind.css

# Fetch dependencies.
# RUN go get -d -v
RUN go mod download
RUN go mod verify

#CMD go build -v
# go build command with the -ldflags="-w -s" option to produce a smaller binary file by stripping debug information and symbol tables. 
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o /server/server ./cmd/server

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
COPY --from=builder /src/config/config.docker.json /server/config/config.json
COPY --from=builder /src/logs /server/logs
COPY --from=builder /src/web /server/web

#ENV SUPABASE=
#ENV CRYPTO_KEY=
#ENV CRYPTO_IV_PAD=

EXPOSE 8888

CMD [ "/server/server", "--config", "/server/config/config.json" ]

