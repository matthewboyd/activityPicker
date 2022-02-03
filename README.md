## How to start the project: 

```docker-compose up --build```

This will run the 3 services: go app, db + redis

In order to make sure all the data is in the db, make sure to:


`1. docker exec -it <container> bash`

`2. cd docker-entrypoint-initdb.d`

`3. psql -U activities -f create.sql`

`4. psql -U activities -f fill_tables.sql`

This should make sure all the data is present in the table, and to double check, do the following commands in the container: 

`1. psql -U activities`

`2. select * from activities`


#Redis: 

There are some commands in order to make sure that redis is working, first: 

to login to redis, do the following: 

`docker exec -it <container_name> redis-cli`

Then once you're in the container, type: 

`KEYS * `

This will show you all the keys that are currently being stored in redis.


Improvements:

- Add users so that they can tick off activities that they've completed and they won't be shown to them
- Add a circuit breaker so that if there are no sunny activites after 5 tries, it will suggest to try a allWeather activity
- could have it on a domain, and each country has a subdomain 
- add additional filters, i.e. activities for age groups
