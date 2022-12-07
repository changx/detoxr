# This is a multi-stage Dockerfile and requires >= Docker 17.05
# https://docs.docker.com/engine/userguide/eng-image/multistage-build/
# yarn build
FROM node:18 as vue

RUN mkdir -p /src/site
WORKDIR /src/site

COPY . .
RUN yarn install
RUN yarn build


# buffalo build
FROM gobuffalo/buffalo:v0.18.9 as builder

ENV GOPROXY http://proxy.golang.org

RUN mkdir -p /src/site
WORKDIR /src/site

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

ADD . .
COPY --from=vue /src/site/dist .
RUN buffalo build --static -o /bin/app

# final image
FROM alpine
RUN apk add --no-cache bash
RUN apk add --no-cache ca-certificates

WORKDIR /bin/

COPY --from=builder /bin/app .
COPY .env /bin/.env

# Uncomment to run the binary in "production" mode:
ENV GO_ENV=production

# Bind the app to 0.0.0.0 so it can be seen from outside the container
EXPOSE 3000

# Uncomment to run the migrations before running the binary:
# CMD /bin/app migrate; /bin/app
CMD exec /bin/app
