FROM postgres
ENV POSTGRES_DB prod
COPY db_schema /docker-entrypoint-initdb.d/
EXPOSE 5432

