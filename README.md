# Completed Assignment

This repository contains the source code for the code challenge sent to me. 


## Running the Application

To run the application, you need to follow these steps:

1. **Install Golang:** Make sure you have Go installed on your system. If not, you can download it from [here](https://go.dev/doc/install).

2. **Navigate to the Application Directory:** Open your terminal and navigate to the directory of the application using the `cd` command:
   
   ```sh
   cd path/to/the/application

3. Install the dependencies by running `go get .` in the terminal.
   
4. Start the application by running `go run .` in the terminal.

## Running the tests
To run the tests:

1. Navigate to the Application Directory using the command above

2. Run the tests by executing this command `go test ./tests` in the terminal.

## Optimizing the Application
To keep it simple and because really I was out of time, I used an in mock database that stores the data in memory. Ideally I would use an RDBMS like postgreSQL to store the data. 
That being said, first thing I would do is to have separate servers for the server and the database so that they can be scaled horizontally independently. 
Next is to have some sort of cron job on the database server that would take a daily backup of the database and push it to some cloud storage like AWS S3. After that, the next thing I would do is to try to scale the database server horiziontally by either using Database Sharding or Database Replication(Master/Slave). 
I would most likely go with the Database replication since the database is an RDBMS and doing Database Sharding on an RDBMS database is a bit more difficult as it is harder to scale horizontally. 
After scaling the database server, I would proceed to the API server, and the first thing I would do here is to add a cache tool like redis to cache dynamic data. That way the API does not have to always go to the database to fetch data thereby also further optimizing the database. 
Next, I would distribute the API across different servers and add a load balancer in between that would help manage the servers. Finally I would add some analysis tool like ELK or New relic to help monitor the performance of the application and give me an idea of what to optimize in the system
