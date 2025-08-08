### Setup Application
go mod tidy

go run main.go

### Setup Database
- Create the  database in mysql
```
create database incident_trackerdb;
```

- Create a user in mysql
```
create user 'incident_tracker'@'localhost' identified by 'incident_tracker';
```
- Grant all privileges to the user
```
grant all privileges on incident_trackerdb.* to 'incident_tracker'@'localhost';
```
- Flush privileges
```
flush privileges;
```