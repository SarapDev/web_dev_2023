ARG build_image=docker-hub-utils.kolesa-team.org:5000/build/golang:1.20-alpine-dev
ARG base_image=docker-hub-utils.kolesa-team.org:5000/base/alpine:latest

FROM ${build_image} as build_stage

ARG goproxy=https://goproxy.kolesa-team.org|direct
ARG github_token=""
ENV GOPROXY=${goproxy} \
    GOPRIVATE=*.kolesa-team.org \
    GO111MODULE=on

# fix for github.kolesa-team.org
RUN if [ ! -z "$github_token" ]; then git config --global url.https://${github_token}@github.kolesa-team.org/.insteadOf https://github.kolesa-team.org/ ; fi

WORKDIR /code

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -v -o out/binary


FROM ${base_image}

ARG revision
ARG branch
ARG build_number
ARG build_url
ARG build_date

ENV RELEASE_REVISION=${revision} \
    RELEASE_BRANCH=${branch} \
    RELEASE_BUILD_NUMBER=${build_number} \
    RELEASE_BUILD_URL=${build_url} \
    RELEASE_BUILD_DATE=${build_date}

COPY config /config
COPY --from=build_stage /code/out/binary /usr/local/bin/

USER nobody

ENTRYPOINT /usr/local/bin/binary

ARG base_image
ARG version=not-set
ARG revision=not-set
LABEL org.kolesa-team.image.name="example" \
      org.kolesa-team.image.version="${version}" \
      org.kolesa-team.image.revision="${revision}" \
      org.kolesa-team.image.base_image="${base_image}" \
      org.kolesa-team.image.description="пример сервиса на go modules"
