FROM scratch
COPY app /
COPY buddy.json /
COPY buddy.sqlite /
RUN chmod +x /app
CMD ["/app"]
