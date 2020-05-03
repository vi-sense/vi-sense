# vi-sense
tool for visualizing 3d models with sensor information

# Installation
> docker-compose up --build		#builds and starts the app in dev mode, using the source code on the machine

> docker-compose -f docker-compose.yml -f docker-compose.integration.yml up -d 	#starts the app in integration mode, using the visense image from docker hub

## Generate API documentation
cd into app/
> swag init -g api.go 
