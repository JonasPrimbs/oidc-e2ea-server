FROM docker/dev-environments-go:stable-1

COPY ./dev-entrypoint.sh /dev-entrypoint.sh

ENTRYPOINT [ "/entrypoint.sh" ]
