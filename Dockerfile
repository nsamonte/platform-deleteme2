FROM golang:1.11 AS build

ARG REPO_NAME
ARG GITHUB_ACCESS_TOKEN

WORKDIR /src

RUN apt-get -qq update && apt-get -y install \
    libsasl2-dev \
    libsasl2-modules \
    libssl-dev

RUN wget -O "librdkafka.tar.gz" "https://github.com/edenhill/librdkafka/archive/v0.11.6.tar.gz"
RUN mkdir -p librdkafka
RUN tar \
  --extract \
  --file "librdkafka.tar.gz" \
  --directory "librdkafka" \
  --strip-components 1

RUN cd "librdkafka" && \
  ./configure --prefix=/usr && \
  make -j "$(getconf _NPROCESSORS_ONLN)" && \
  make install

RUN cd ..
COPY . .

RUN git config --global url."https://${GITHUB_ACCESS_TOKEN}:@github.com/".insteadOf "https://github.com/"

RUN go build -o /bin/app

FROM gcr.io/distroless/base
COPY --from=build /bin/app /
COPY --from=build /usr /usr
CMD ["/app"]