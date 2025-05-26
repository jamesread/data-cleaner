FROM registry.fedoraproject.org/fedora-minimal:40

LABEL org.opencontainers.image.source https://github.com/jamesread/data-cleaner
LABEL org.opencontainers.image.authors James Read
LABEL org.opencontainers.image.title data-cleaner

ENV PORT=8080
EXPOSE 8080
RUN mkdir /app
WORKDIR /app

# Goreleaser builds to ./uar, not ./service/uar
COPY data-cleaner /app/data-cleaner

COPY frontend/dist /app/webui

ENTRYPOINT ["/app/data-cleaner"]
