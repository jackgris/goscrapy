# Notes

Here I will add some useful commands to remenber on future

## Docker

How you can create a container running mongodb:

```bash
docker run -d --network some-network --name some-mongo \
	-e MONGO_INITDB_ROOT_USERNAME=mongoadmin \
	-e MONGO_INITDB_ROOT_PASSWORD=secret \
	mongo
```

from [documentation](https://hub.docker.com/_/mongo)

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

```javascript
mongodump --db=mayorista
```

#### Create database from backup

This command will create the mayorista2 database with all the data that we saved in the backup we stored in the dump folder
with all the data of mayorista database

```javascript
mongorestore --db mayorista2 dump/mayorista --drop
```

from official [documentation](https://www.mongodb.com/docs/cloud-manager/tutorial/restore-single-database/)

With this, we can test our application using a real database, with true values

### Docker container running our database

Commands useful when you are working with a database inside a docker container

#### Getting the backup to our file system

Getting the database backup from the docker container to the host

```javascript
sudo docker cp goscrapy-mongo:/dump/mayorista/productos.bson ~/Downloads
```

from official [documentation](https://docs.docker.com/engine/reference/commandline/cp/)

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
