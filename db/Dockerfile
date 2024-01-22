ARG POSTGRES_VERSION=16
ARG PGVECTOR_VERSION=0.5.1

# Use the official Postgres image based on Alpine
FROM postgres:${POSTGRES_VERSION}-alpine

# Set environment variables for Postgres
ARG POSTGRES_USER
ARG POSTGRES_PASSWORD
ARG POSTGRES_DB

ENV POSTGRES_USER=${POSTGRES_USER}
ENV POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
ENV POSTGRES_DB=${POSTGRES_DB}

# Install build dependencies for Alpine
RUN apk update && apk add --no-cache \
    build-base \
    git \
    postgresql-dev

# Clone, build, and install the pgvector extension
RUN cd /tmp \
    && git clone --branch v${PGVECTOR_VERSION} https://github.com/pgvector/pgvector.git \
    && cd pgvector \
    && make \
    && make install

# Copy initialization SQL script to be executed on container start
COPY ./init_pgvector.sql /docker-entrypoint-initdb.d/init.sql

# Expose PostgreSQL port
EXPOSE 5432

# Define volume mount point
VOLUME ["/var/lib/postgresql/data"]
