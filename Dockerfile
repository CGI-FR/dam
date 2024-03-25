FROM gcr.io/distroless/base
ARG BIN
COPY /bin/dam /dam
CMD ["/dam"]
