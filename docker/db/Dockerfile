# pull official base image
FROM postgres:16

# run init.sql on init
ADD init.sql /docker-entrypoint-initdb.d/
