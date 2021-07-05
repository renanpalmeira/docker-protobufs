ARG GO_VERSION

FROM golang:${GO_VERSION}-buster AS build

# hadolint ignore=DL3008
RUN apt-get update && apt-get install --no-install-recommends -y unzip

ARG GOLANG_PROTOBUF_VERSION
ARG GRPC_GATEWAY_VERSION
ARG PROTOBUF_VERSION

RUN GO111MODULE=on go get \
  github.com/golang/protobuf/protoc-gen-go@v${GOLANG_PROTOBUF_VERSION} && \
  mv /go/bin/protoc-gen-go* /usr/local/bin/

# hadolint ignore=DL3059
RUN curl -sSL \
  https://github.com/grpc-ecosystem/grpc-gateway/releases/download/v${GRPC_GATEWAY_VERSION}/protoc-gen-grpc-gateway-v${GRPC_GATEWAY_VERSION}-linux-x86_64 \
  -o /usr/local/bin/protoc-gen-grpc-gateway && \
  curl -sSL \
  https://github.com/grpc-ecosystem/grpc-gateway/releases/download/v${GRPC_GATEWAY_VERSION}/protoc-gen-openapiv2-v${GRPC_GATEWAY_VERSION}-linux-x86_64 \
  -o /usr/local/bin/protoc-gen-openapiv2 && \
  chmod +x /usr/local/bin/protoc-gen-grpc-gateway && \
  chmod +x /usr/local/bin/protoc-gen-openapiv2

WORKDIR /tmp/protoc

# hadolint ignore=DL3059
RUN curl -sSL \
  https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOBUF_VERSION}/protoc-${PROTOBUF_VERSION}-linux-x86_64.zip \
  -o /tmp/protoc/protoc.zip && \
  unzip protoc.zip && \
  mv /tmp/protoc/bin/protoc /usr/local/bin && \
  mv /tmp/protoc/include/* /usr/local/include && \
  chmod +x /usr/local/bin/protoc && \
  chmod -R 777 /usr/local/include

FROM debian:buster-slim AS protobuf

COPY --from=build /usr/local/bin /usr/bin
COPY --from=build /usr/local/include /usr/include
COPY .third_party /usr/include
COPY protoc-wrapper /usr/bin/protoc-wrapper
ENV LD_LIBRARY_PATH='/usr/lib:/usr/lib64:/usr/lib/local'
ENTRYPOINT ["protoc-wrapper", "-I/usr/include"]
