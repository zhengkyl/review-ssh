FROM alpine:latest

RUN apk update && apk add --update git && rm -rf /var/cache/apk/*

COPY review-ssh /usr/local/bin/review-ssh
COPY .ssh /.ssh

# Expose ports
# SSH
EXPOSE 3456/tcp

# Set the default command
ENTRYPOINT ["/usr/local/bin/review-ssh"]


