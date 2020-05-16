# vi-sense
tool for visualizing 3d models with sensor information

# Installation
##Local developmend
- clone the visense-frontend repo
> FRONTEND_DIR=[path_to_frontend_repo] docker-compose up --build		#builds and starts the app in dev mode, using the source code on the machine

> FRONTEND_PORT=80 docker-compose -f docker-compose.yml -f docker-compose.integration.yml up -d 	#starts the app in integration mode, using the visense image from docker hub, frontend listens on port 80

## Generate API documentation
cd into app/
> swag init -g api.go 
