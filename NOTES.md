# Notes

Here I will add some useful commands to remenber on future

## Docker

### Create mongodb container

How you can create a container running mongodb:

```bash
docker run -d -p 27017:27017 --name goscrapy-mongo -v mongo-data:/data/db \
         -e MONGO_INITDB_DATABASE=admin \
         -e MONGODB_INITDB_ROOT_USERNAME=admin \
         -e MONGODB_INITDB_ROOT_PASSWORD=admin \
         mongo:latest
```

from [documentation](https://hub.docker.com/_/mongo)

### Backup from our database

Commands useful when you are working with a database inside a docker container

#### Getting the backup to our file system

Getting the database backup from the docker container to the host

```bash
docker cp goscrapy-mongo:/dump/mayorista/productos.bson ~/Downloads
```

from official [documentation](https://docs.docker.com/engine/reference/commandline/cp/)

### You can make the mongo database images from a Dockerfile

#### Build our container

This command needs to be run inside the folder where you have the Dockerfile of our database, in this case
is inside the database folder of our project.
The flag -t is for put a name to our image, so you can change mongodb-from-dockerfile for anything you want.

```bash
docker build -t mongodb-from-dockerfile .
```

#### Run our container

Once you create the docker image with the Dockerfile, you can run a container based on that image.
The flag -d is for run the container in a detached way, and with --name you can set the name you want for that container. That will make it easier to remember. And the flag -p is one of the most important because will publish a container's port(s) to the host.

```bash
docker run -d  -p 27017:27017 --name mongo-scrapy mongodb-from-dockerfile
```

from official [documentation](https://docs.docker.com/engine/reference/commandline/run/)


#### Start our mongosh terminal

After run the container, if you want to use the database with shell of mongo, you can run this command, remember that mongo-scrapy is the name of the container, so if you change it that when run the container, you need to change the name on this command.

```bash
docker exec -it mongo-scrapy mongosh
```

## MongoDB

If you need add security, because for any reason create first the database without authentication and need use SCRAM to Authenticate Clients run this on mongosh

```javascript
use admin
db.createUser(
    {user: "admin",pwd:"admin",
        roles:[{role:"userAdminAnyDatabase", db:"admin"},
        {role:"readWriteAnyDatabase", db:"admin"}]}
)

```
from official [documentation](https://www.mongodb.com/docs/manual/tutorial/configure-scram-client-authentication)

### Create backup, restore database and create a new one to testing

#### Generate backup

This command will generate a backup from mayorista database and it will save it in the dump folder

```bash
mongodump --db=mayorista
```

#### Create database from backup

This command will create the mayorista2 database with all the data that we saved in the backup we stored in the dump folder
with all the data of mayorista database

```bash
mongorestore --db mayorista2 dump/mayorista --drop
```

from official [documentation](https://www.mongodb.com/docs/cloud-manager/tutorial/restore-single-database/)

With this, we can test our application using a real database, with true values


## Project

### Setup database

For testing propose you need to run a mongodb restoring the data contain in the folder backup/products-database

### Run the project

At least for now, run this in the terminal while you are in the main folder.

```bash
go build . && ./goscrapy
```
### Automation project setup

Later we will write a makefile that is going to help us to make all these steps in an automatic way
