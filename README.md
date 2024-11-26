# Step to sync Data from Postgres to Debezium

1. **Run the Docker container** by command : `docker-compose up -d` and create table by command `make g-up`  


2. **Config the Database :**
* Connect to database by command : `docker exec -it postgres psql -U root -d cqrs`
* Run the command  
`ALTER SYSTEM SET wal_level = 'logical';`  
`ALTER SYSTEM SET max_replication_slots = 4;`  
`ALTER SYSTEM SET max_wal_senders = 4;`

    Restart the database by command : `docker-compose restart postgres`  
    Check the configuration's change by command : `docker exec -it postgres psql -U root -d cqrs -c "SHOW wal_level; SHOW max_replication_slots; SHOW max_wal_senders;"`

* Create the replication user by command : `ALTER USER root WITH REPLICATION;`  
* Create the replication to listen the change of database: `CREATE PUBLICATION cqrs_pub FOR ALL TABLES;`  
* Create the replication identity for the table : `ALTER TABLE users REPLICA IDENTITY FULL;` (change users to your table name, in this situation I have "user" and "orders")

3. **Config the Debezium :**  
You can custom the config in /docker/debezium.json  
Or you can use my config, execute the command (execute in **Command Prompt** if you using **Windows**) : `curl -i -X POST -H "Accept:application/json" -H "Content-Type:application/json" localhost:8083/connectors/ -d "@debezium.json"`  
Check the response by command : `curl -i -X GET -H "Accept:application/json" localhost:8083/connectors/`  
If you want to delete the connector, execute the command : `curl -i -X DELETE -H "Accept:application/json" localhost:8083/connectors/postgres-connector`


4. **Now you can check the data** in the topic at `localhost:8082`, this is Kafka UI to visualize data.  
Use Postman and import this command to test the data :  
Create : `curl --location 'http://localhost:8080/user/create' \
   --header 'Content-Type: application/json' \
   --data-raw '{
   "name" : "son",
   "email" : "email@email.com"
   }'`  
Update : `curl --location 'http://localhost:8080/user/update' \
   --header 'Content-Type: application/json' \
   --data '{
   "id" : 1,
   "name" : "son3"
   }'`

### Or you can simply run the command : `make up`, wait for docker to pull all the images, then `make cfg` to run all the command above and `make run` to run the program.  
### Checking out at `localhost:8082` and use Postman to send the data.

