FROM mongo:latest

COPY ./deployments/mongodb/init_db.js /docker-entrypoint-initdb.d/
EXPOSE ${DB_INNER_PORT}:27017