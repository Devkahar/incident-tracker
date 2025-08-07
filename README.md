### Setup Application
go mod tidy

go run main.go

### Setup Database
- Create the  database in mysql
```
create database incident-trackerdb;
```

- Create a user in mysql
```
create user 'incident-tracker'@'localhost' identified by 'incident-tracker';
```
- Grant all privileges to the user
```
grant all privileges on incident-trackerdb.* to 'incident-tracker'@'localhost';
```
- Flush privileges
```
flush privileges;
```