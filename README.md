# go-devbook API

Go studies based on Udemy course (pt-br) [Aprenda Golang do Zero! Desenvolva uma APLICAÇÃO COMPLETA!](https://www.udemy.com/course/aprenda-golang-do-zero-desenvolva-uma-aplicacao-completa/)

## API

### DEPENDENCIES
- go get github.com/gorilla/mux
- go get github.com/joho/godotenv
- go get github.com/go-sql-driver/mysql
- go get github.com/badoux/checkmail
- go get golang.org/x/crypto/bcrypt
- go get github.com/dgrijalva/jwt-go

## ENV FILE
It is needed to create a _.env_ file with the content below to be used by the api application:
```properties
DB_USUARIO=golang
DB_SENHA=d3f1n3dP4ssw0rd!
DB_NOME=devbook

API_HOST=http://{IP_FROM_WSL}
API_PORT=5000

SECRET_KEY={KEY_BASE64}

RUN_INIT=false
```

### FIRST EXECUTION
Before the very first execution, change RUN_INIT to `true`. Once the application start running, it will print a key you should copy and paste to key SECRET_KEY into _.env_ file.
Stop application and start it again. Then everything will be set up.

## MYSQL
You can find a `mysql` folder where you can find the docker compose for mysql.

### ENVIRONMENT FILE FOR MYSQL
Inside the [mysql](/mysql/) folder, create a _.env_ file with this content as compose file depends on it:
```
MYSQL_ROOT_PASSWORD=r00tP4ssw0rd!
MYSQL_DATABASE=devbook
MYSQL_USER=golang
MYSQL_PASSWORD=d3f1n3dP4ssw0rd!
```

### RUN DOCKER COMPOSE
```shell
# swarm init is only necessary on the first time running it.
# after that, only stack deploy will be enough to run the docker-compose.
docker swarm init
docker stack deploy -c docker-compose.yml mysql-go
```

### CHECK IF IT IS RUNNING
```shell
docker container ps
```

An example of results:
```shell
CONTAINER ID   IMAGE       COMMAND                  CREATED          STATUS          PORTS                 NAMES
cc9aaa82abf2   mysql:8.0   "docker-entrypoint.s…"   18 minutes ago   Up 18 minutes   3306/tcp, 33060/tcp   mysql-go_db.1.ni4mlz4dfr285v6c0hanb0c3a
```

The command below will bring only the `container_id` if you're using linux
```shell
docker container ps | sed -n '2 p' | awk '{print $1}'

# or saving it to an envirionment variable
export CONTAINER_ID=`docker container ps | sed -n '2 p' | awk '{print $1}'`
```

### RUN MYSQL FROM CONTAINER

The command ended by `/bin/bash` gives access to linux terminal (bash)
The command ended by `mysql -u golang -p` gives access to mysql console. It will request the user _golang_ password


```shell
docker exec -it [CONTAINER_ID] /bin/bash
docker exec -it [CONTAINER_ID] mysql -u golang -p

# if you're using linux and saved the container_id to an environment variable, you can use:
docker exec -it $CONTAINER_ID mysql -u golang -p

# if it is needed to access mysql console, you can use
docker exec -it $CONTAINER_ID mysql -u root -p
```

### STOP DOCKER COMPOSE
```shell
docker stack rm mysql-go
```

### CHECK LOGS
```shell
docker service logs mysql-go_db
```

### CREATING TABLE (IF NEEDED)
It is already expected to create everything needed on mysql startup as it is set up on docker compose file.
This block should be used only if the table was not created...

Access mysql root console and execute the steps below:
```sql
show databases;
use devbook;
-- run first_ddl.sql content
show tables;
```

### DADOS DE CONEXÃO COM MYSQL APÓS COMPOSE ESTAR NO "AR"
```properties
HOST={IP_FROM_WSL}
PORT=3306
USER=golang
PASS=d3f1n3dP4ssw0rd!
```
