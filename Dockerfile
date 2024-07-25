FROM scratch
ARG TARGETARCH
ADD rcpd-${TARGETARCH}-linux /rcpd
ENTRYPOINT ["/rcpd","-root_dir", "/srv"]
EXPOSE 514
VOLUME /srv
LABEL maintainer="as@tenoware.com"
