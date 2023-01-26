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

This command will generate a backup from mayorista database and it will save it in the dump folder

```javascript
mongodump --db=mayorista
```

This command will create the mayorista2 database with all the data that we saved in the backup we stored in the dump folder
with all the data of mayorista database

```javascript
mongorestore --db mayorista2 dump/mayorista --drop
```

from official [documentation](https://www.mongodb.com/docs/cloud-manager/tutorial/restore-single-database/)

With this, we can test our application using a real database, with true values

## Project

Run the project, at least for now, run this in the terminal while you are in the main folder.

```bash
go build . && ./goscrapy
```
