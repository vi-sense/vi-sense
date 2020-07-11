# vi-sense

vi-sense is a project of the University of Applied Sciences (HTW) Berlin and metr.systems to visualize IoT data in BIM Models.

Vi-Sense is a project that visualizes data of heating systems in a 3D model of the building they’re installed in. Technicians or system administrators get an overview of all sensors located in a building, can take a closer look at the measured values of these sensors, and get information about data anomalies that indicate potential errors. The combination of the 3D representation of buildings and the visualization of measured values in a timeline makes it possible to perform detailed analyses of error sources such as the failure of heating or pump systems. In addition, sensor states in the building model are supported by color highlighting, which ensures a quick orientation and an intuitive overview. In short, the goal of Vi-Sense is to give technicians and building managers a clear overview of their properties and the sensor technology they use.

![screenshot](https://raw.githubusercontent.com/vi-sense/sample-data/master/preview/vi-sense-screenshot02.png)

# Installation

## Local development

git large file storage is needed to install this repository

```git lfs install```

clone the visense-frontend repo

```git lfs clone --recurse-submodules git@github.com:vi-sense/vi-sense.git```

alternatively

```git clone --recurse-submodules git@github.com:vi-sense/vi-sense.git```

```cd sample-data/```

```git lfs pull```

or update repo

```git submodule update --init --recursive```

```cd sample-data/```

```git lfs pull```

### Running
```FRONTEND_DIR=[path_to_frontend_repo] docker-compose up --build```
builds and starts the app in dev mode, using the source code on the machine

```FRONTEND_PORT=80 docker-compose -f docker-compose.yml -f docker-compose.integration.yml up -d```
starts the app in integration mode, using the visense image from docker hub, frontend listens on port 80

## Generate API documentation

```cd into app/```
```swag init -g api/api.go```