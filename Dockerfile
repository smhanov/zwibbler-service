# This Dockerfile was based in details found here: https://medium.com/@chemidy/create-the-smallest-and-secured-golang-docker-image-based-on-scratch-4752223b7324

##########
# Step 1 - Build executable
###

# Tiny golang image
FROM golang:alpine AS builder

# git & gcc are required for compilation
RUN apk update && apk add --no-cache git gcc musl-dev bash

# Compilation working directory
WORKDIR /go/src/app

# Copy the main file from the local filesystem
COPY . .
RUN cp zwibbler.conf /etc/zwibbler.conf

# Install required modules
RUN go get

RUN go version

# Compile a static executable - linux/intel
RUN go build

##########
# Step 2 - Build a tiny image
###

#FROM scratch

# Import the user & group files
#COPY --from=builder /etc/passwd /etc/passwd
#COPY --from=builder /etc/group /etc/group

# Copy over the executable directory and set the correct permissions as Zwibserve needs to write the db file
#COPY --from=builder --chown=appuser:appuser /go/bin /go/bin

# Listening port
EXPOSE 3000

# Start the Zwibbler server
ENTRYPOINT [ "./zwibbler" ]
