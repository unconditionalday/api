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

# Create directory for volume
RUN mkdir -p /var/lib/postgresql/data

# Copy initialization SQL script to be executed on container start
COPY ./init.sql /docker-entrypoint-initdb.d/init.sql

# Expose PostgreSQL port
EXPOSE 5432

# Define volume mount point
VOLUME ["/var/lib/postgresql/data"]
