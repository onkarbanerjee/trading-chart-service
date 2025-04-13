# DO NOT CHANGE. This file is being managed from a central repository
# To know more simply visit https://github.com/honestbank/.github/blob/main/docs/about.md
ARG GO_VERSION="1.23"

FROM golang:$GO_VERSION AS builder
WORKDIR /app
# Sensitive
COPY . .

ENV GOPRIVATE=github.com/honestbank
RUN mv .ssh ~/.ssh

RUN chmod 0600 ~/.ssh/id_rsa
RUN git config --global url."git@github.com:".insteadOf "https://github.com/"

RUN touch ~/.ssh/known_hosts
RUN ssh-keyscan github.com >> ~/.ssh/known_hosts

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o hnst ./cmd/cli/

FROM gcr.io/distroless/static-debian10@sha256:73c9f3f9ee72a6070499e46826dacfda8cb428cb97bef9a70ee6f67957ec99c6

WORKDIR /app
# Copy the Pre-built binary file from the previous stage
COPY --from=builder --chmod=111 /app/hnst .

ARG VERSION
ENV APP__VERSION="${VERSION}"
USER nonroot

# Command to run the executable
CMD ["./hnst", "server", "start"]
