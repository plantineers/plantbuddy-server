FROM scratch
COPY app /app
COPY buddy.json /buddy.json
COPY buddy.sqlite /buddy.sqlite
CMD ["/app"]
