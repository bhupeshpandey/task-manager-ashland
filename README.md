# task-manager-ashland
This service handles the logging for the gallatin events via rabbitmq.

For this service just run the main.go locally.

go run /cmd/main.go

Also, this is a combination service. The Nashville, gallatin and ashland. For running all three together<
please run the docker-compose from the gallatin which contains rabbitmq, postgres, and redis cache.
Comment out the gallatin section in docker-compose as that is having some networking issues.

Once the services are up and running on docker,
Go ahead and execute /cmd/main.go on the nashville, gallatin and the ashland.
Once all three services are up, use the postman collection in the 
nashville to call the api's. 