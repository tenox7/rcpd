FROM scratch
ARG TARGETARCH
ADD rcpd-${TARGETARCH}-linux /rcpd
ENTRYPOINT ["/rcpd","-root_dir", "/srv"]
EXPOSE 514
LABEL maintainer="as@tenoware.com"
