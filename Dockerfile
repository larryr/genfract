# syntax=docker/dockerfile:1
## deployment image
FROM gcr.io/distroless/base-debian11

USER nonroot:nonroot
WORKDIR /
COPY genfract ./genfract
COPY *.js ./
COPY *.html ./
COPY *.gif ./
COPY *.png ./

EXPOSE 4000
ENTRYPOINT ["/genfract", "-mandelbrot"]


