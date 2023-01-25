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

## Project

Run the project, at least for now, run this in the terminal while you are in the main folder.

```bash
go build . && ./goscrapy
```
