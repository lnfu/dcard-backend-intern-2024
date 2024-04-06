FROM migrate/migrate

RUN apk add --no-cache bash

COPY db .

# ENTRYPOINT ["migrate"]
ENTRYPOINT ["./wait-for-it.sh", "mysql:3306", "--", "./migrate"]
CMD ["--help"]
