FROM debian
COPY app /
COPY buddy.json /
COPY buddy.sqlite /
CMD ["/app"]
