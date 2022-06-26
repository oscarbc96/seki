FROM gcr.io/distroless/static:nonroot
COPY seki /
ENTRYPOINT ["/seki"]