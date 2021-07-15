FROM golang:alpine AS build
ARG CGO_ENABLED=0
ARG GOPROXY=https://proxy.golang.org,direct
ARG GOMODCACHE=/go/pkg/mod
ARG GOCACHE=/root/.cache/go-build
ARG GTM=GTM-TLVN7D6
WORKDIR /workspace
COPY . .
RUN go run ./go/cmd/webrender -src blog -dst ./go/internal/w16/static/root -gtm ${GTM}
RUN go build -trimpath -ldflags='-s -w' -o /bin/ ./go/cmd/...

FROM gcr.io/distroless/static AS archrepod
COPY --from=build /bin/archrepod /bin/archrepod
ENTRYPOINT ["/bin/archrepod"]

FROM gcr.io/distroless/static AS feedagg
COPY --from=build /bin/feedagg /bin/feedagg
ENTRYPOINT ["/bin/feedagg"]

FROM gcr.io/distroless/static AS ghdefaults
COPY --from=build /bin/ghdefaults /bin/ghdefaults
ENTRYPOINT ["/bin/ghdefaults"]

FROM gcr.io/distroless/static AS paste
COPY --from=build /bin/paste /bin/paste
ENTRYPOINT ["/bin/paste"]

FROM gcr.io/distroless/static AS singlepage
COPY --from=build /bin/singlepage /bin/singlepage
ENTRYPOINT ["/bin/singlepage"]

FROM gcr.io/distroless/static AS vanity
COPY --from=build /bin/vanity /bin/vanity
ENTRYPOINT ["/bin/vanity"]

FROM gcr.io/distroless/static AS w16
COPY --from=build /bin/w16 /bin/w16
ENTRYPOINT ["/bin/w16"]
