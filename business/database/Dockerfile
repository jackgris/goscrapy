FROM mongo:latest

WORKDIR /backup
# put our mongodb backup into the container
COPY ./productos.bson .

# environment:
ENV MONGO_INITDB_DATABASE: logs
ENV MONGO_INITDB_ROOT_USERNAME: admin
ENV MONGO_INITDB_ROOT_PASSWORD: admin

EXPOSE 27017

CMD mongod --fork --logpath /var/log/mongodb.log; \
    mongorestore --db mayorista2 . --drop; \
    mongod --shutdown; \
    mongod
