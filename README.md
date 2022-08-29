### Prerequisite
Install go-migrate for running migration
```
https://github.com/golang-migrate/migrate/tree/master/cmd/migrate
```

App requires 2 database (Postgres and redis server)
```
# run Postgres
docker run -d -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=rollee postgres

# run redis
docker run -d -p 6379:6379 redis
``` 

### Migration
Run below command to run migration
```
make migration-up    
```


### Test
Run below command to run test, and make sure that all tests are passing
```
make test
```

### Running
Run below command to run app
```
make run-server
```

Swagger URL
```
${BASE_URL}/swagger/index.html
```