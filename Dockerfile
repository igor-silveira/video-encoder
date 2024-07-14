FROM golang:1.14-alpine3.11

ENV BENTO4_BIN="/opt/bento4/bin" \
    BENTO4_BASE_URL="http://zebulon.bok.net/Bento4/source" \
    BENTO4_VERSION="1-6-0-641" \
    BENTO4_CHECKSUM="ed3e2603489f4748caadccb794cf37e5e779422e" \
    BENTO4_TARGET="" \
    BENTO4_PATH="/opt/bento4" \
    BENTO4_TYPE="SRC" \
    PATH="$PATH:/bin/bash:$BENTO4_BIN"

RUN apk add --no-cache --update \
        ffmpeg \
        bash \
        curl \
        wget \
        make \
        python \
        unzip \
        gcc \
        g++ \
        scons

WORKDIR /tmp/bento4
RUN wget -q ${BENTO4_BASE_URL}/Bento4-${BENTO4_TYPE}-${BENTO4_VERSION}${BENTO4_TARGET}.zip && \
    echo "${BENTO4_CHECKSUM}  Bento4-${BENTO4_TYPE}-${BENTO4_VERSION}${BENTO4_TARGET}.zip" | sha1sum -c - && \
    mkdir -p ${BENTO4_PATH} && \
    unzip Bento4-${BENTO4_TYPE}-${BENTO4_VERSION}${BENTO4_TARGET}.zip -d ${BENTO4_PATH} && \
    rm -rf Bento4-${BENTO4_TYPE}-${BENTO4_VERSION}${BENTO4_TARGET}.zip && \
    cd ${BENTO4_PATH} && \
    scons -u build_config=Release target=x86_64-unknown-linux && \
    cp -R ${BENTO4_PATH}/Build/Targets/x86_64-unknown-linux/Release ${BENTO4_PATH}/bin && \
    cp -R ${BENTO4_PATH}/Source/Python/utils ${BENTO4_PATH}/utils && \
    cp -a ${BENTO4_PATH}/Source/Python/wrappers/. ${BENTO4_PATH}/bin

RUN apk del unzip gcc g++

WORKDIR /go/src

ENTRYPOINT ["top"]
