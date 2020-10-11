FROM alpine:3.10 as download-zdoc
ENV ZDOC_VERSION 885938ae5617a255fb491286d36d2bc4d27f8cea
ENV ZDOC_URL https://github.com/zzossig/hugo-theme-zdoc.git
RUN apk add --no-cache git
WORKDIR /zdoc
RUN git clone "${ZDOC_URL}" . && \
    git reset --hard "${ZDOC_VERSION}" && \
    git submodule update --init --recursive && \
    rm -rf .git

# Download Hugo
FROM alpine:3.10 as download-hugo
ENV HUGO_VERSION 0.73.0
ENV HUGO_URL https://github.com/gohugoio/hugo/releases/download/v${HUGO_VERSION}/hugo_extended_${HUGO_VERSION}_Linux-64bit.tar.gz
RUN wget -O- "${HUGO_URL}" | tar xz

FROM node:10.15.3-stretch as runtime_image
RUN npm install -g postcss-cli
WORKDIR /app/docs
RUN npm install autoprefixer
COPY --from=download-zdoc /zdoc ./themes/zdoc
COPY --from=download-hugo /hugo /usr/local/bin/