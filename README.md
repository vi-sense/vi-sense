# vi-sense

vi-sense is a project of the University of Applied Sciences (HTW) Berlin and metr.systems to visualize IoT data in BIM Models.

this repo contains sample data for the visense project, consisting of 3D gltf models and sensor data.

The aim of vi-sense is to develop a solution for visualizing IoT sensors and their data in 3D building models. We will bring together two cutting edge technologies: Building Information Modeling (BIM) and Internet of Things (IoT) to find a way to represent sensors in 3D space and visualize their data in a comprehensible way to maintenance staff. We will look at different aspects of how to visualize pipes, ducts and sensors to make users understand how systems are connected and which role the element plays they are looking at. Each sensor will have a small dashboard presenting current values, time series and related sensors.

![](https://github.com/vi-sense/sample-data/blob/master/image.png)

# Installation

## Local development

- clone the visense-frontend repo

```FRONTEND_DIR=[path_to_frontend_repo] docker-compose up --build```

builds and starts the app in dev mode, using the source code on the machine

```FRONTEND_PORT=80 docker-compose -f docker-compose.yml -f docker-compose.integration.yml up -d```

starts the app in integration mode, using the visense image from docker hub, frontend listens on port 80

## Generate API documentation

cd into app/
`swag init -g api.go`
