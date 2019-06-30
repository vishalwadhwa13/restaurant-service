# RESTAURANT SERVICE
A microservice written in golang, accessible through gRPC calls defined by protobuf. It provides several utility functions to get data from the MySQL db.

The MySQL db has the following schema.

```
Restaurant
+-------------+--------------+------+-----+----------+----------------+
| Field       | Type         | Null | Key | Default  | Extra          |
+-------------+--------------+------+-----+----------+----------------+
| ResId       | int(11)      | NO   | PRI | NULL     | auto_increment |
| Name        | varchar(50)  | NO   |     | NULL     |                |
| Rating      | decimal(5,3) | NO   |     | 3.000    |                |
| Cuisines    | text         | YES  |     | NULL     |                |
| OpeningTime | time         | NO   |     | 00:00:00 |                |
| ClosingTime | time         | NO   |     | 23:59:00 |                |
| Location    | point        | NO   |     | NULL     |                |
| CostForTwo  | double       | NO   |     | 0        |                |
+-------------+--------------+------+-----+----------+----------------+
```

**IMPORTANT NOTE:**

PROVIDE `dbName`, `dbUser`, `dbPassword`, `dbHost` and `dbPort` through the following env variables in a `.env` file in project root.

.env
```
DB_USER=root
DB_PASSWORD=some_password
DB_NAME=restdb
DB_HOST=localhost
DB_PORT=3306
```

## Methods
1. AddRestaurant
2. GetRestaurant
3. EditRestaurant
4. DeleteRestaurant
5. GetRestaurants

## Instructions
1. Build docker image with `make docker_build`.
2. Run the image with `docker run -it --name rest-micro -p 8080:8080 restaurant-service`.
3. Make sure that mysql is running on the port in the environment.