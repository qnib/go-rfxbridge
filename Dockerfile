FROM qnib/uplain-golang AS build

WORKDIR /usr/local/src/github.com/qnib/go-rfxbridge
COPY vendor ./vendor
COPY main.go ./main.go
COPY http ./http
COPY rfx ./rfx
COPY cli ./cli
RUN go build

FROM qnib/uplain-init
COPY --from=build /usr/local/src/github.com/qnib/go-rfxbridge/go-rfxbridge /usr/local/bin/