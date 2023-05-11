FROM scratch
COPY app /
COPY buddy.json /
COPY buddy.sqlite /
CMD ["/app"]
