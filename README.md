# ip2location-pfsense
[![Go Reference](https://pkg.go.dev/badge/github.com/jpmchia/IP2Location-pfSense.svg)](https://pkg.go.dev/github.com/jpmchia/IP2Location-pfSense)  ![Build status](https://github.com/jpmchia/IP2Location-pfSense/actions/workflows/analysis.yml/badge.svg?event=push)  ![Release](https://github.com/jpmchia/IP2Location-pfSense/actions/workflows/go-release.yml/badge.svg)



pfSense dashboard widget and backend service for displaying live IP geolocation information obtained from the IP2Location.io API. 

The dashboard widget and backend service are designed for use with the IP2Location.io API. Create a free account with up to 30k API requests per month: https://www.ip2location.io/

## Features

* Live processing of pfSense firewall logs to obtain IP geolocation information
* Offloads processing of IP information from pfSense to a separate host
* Minimal overhead on pfSense host, able to run on dedicated pfSense hardware
* Live IP geolocation information retrieved from the IP2Location.io API
* Map view of IP locations provided by Leaflet.js with sumamry of IP information
* Watch list of IP addresses to monitor and track the number of hits
* Drill-down into details of the IP address  information provided by IP2Location.io
* Fully functional with a free IP2Location.io API account
* Support for more granular IP information available with paid IP2Location.io API accounts
* Utilises Redis cache to reduce the number of API requests
* Frequency of updates, number of log entries, API calls and cache duration configurable
* Support for SSL / TLS and backend secured with configurable API keys

### WIP features (part implemented)
* Backend service standalone mode, allowing monitoring of host's network traffic without pfSense
* Configurable map tile providers for Leaflet.js
* Additional UI for cache controls, monitoring and hit rates and configuration of the backend
* Transfer storage of Watch list from the local browser storage to the backend service 


<img src="https://github.com/jpmchia/IP2Location-pfSense/blob/e1156d594e6e0e71c3cab91124f874811b4b3029/contrib/screenshots/Screenshot1.png?raw=true" width="95%">


## Getting started

This project requires Go to be installed. Follow the instructions to download and install Go from: https://go.dev/doc/install. At the time of writing, the latest version is 1.21.1.

The backend service can be run on any platform and OS that supports Go, such as Linux, Windows, MacOS and FreeBSD. For all available ports, refer to: https://go.dev/dl/

## Additional requirements

The service also requires either Redis Stack or Redis server with Redis JSON module installed. Refer to https://redis.io/download/ for available options. It may be hosted on the same host as the backend-service, an existing instance, or a new instance on a different host. 

Please note that TLS / SSL support is not yet completed, so requests / repsonses between the service, pfSense and Redis are not encrypted and therefore it is recommended that either all three are deployed to the same host, or they are deployed on hosts within a private network. 

# pfSense dashboard widget

For the moment, the dashboard widget and supporting files must be manually installed on your pfSense instance. pfSense package to be submitted. Until then, installation requires copying the contents of the pfSense/www folder. 
If you have SSH access configured on the pfSense box, installation is as easy as running the following command from a bash terminal. Installation of the widget should not modify or overwrite any existing, standard pfSense files.

The dashboard widget targets and has been tested on pfSense+ version 23.0.5.1. It may work on other versions, but this has not been tested. The widget is comprised of PHP AJAX and supporting JavaScript packages.

```
./scripts/update_pfsense.sh <username>@<host>
```

Once the files have been copied, the widget can be added to the dashboard by navigating to the pfSense dashboard, expanding the "Available Widgets" panel and then selecting the 'IP2Location' widget. If the "Available Widgets" panel is not visible, navigate to System -> General Setup -> webConfigurator -> Associated Panels Show/Hide and ensure that the checkbox for "Available Widgets" is checked. 


<img src="https://github.com/jpmchia/IP2Location-pfSense/blob/e1156d594e6e0e71c3cab91124f874811b4b3029/contrib/screenshots/Screenshot2.png?raw=true" width="95%">


## Configuration

The dashboard widget requires details of the backend service to be configured. Navigate to the pfSense dashboard, and click the "Settings" button icon in the top-right of the widget's panel. Enter in the base URL of backend service, including the port. By default, the port is configured to run on 9999. If you are running the backend on the pfSense host itself the URL would be http://localhost:9999. If the backend service is running on a different host, enter the hostname or IP address of the host. For SSL / TLS support, enter the URL with the https protocol (refer to the backend service configuration section for Furter details).

If you have changed any of the endpoint paths in the backend service configuration, ensure that these are reflected in the widget configuration. 
The widget requires the API key for the backend service. Please note that this is NOT the token issued by IP2Location.io, but the API key configured in the backend service.

You may limit the number of log entries processed by the widget/backend service by entering a value in the "Number of entries" field. The default is 500. Additionally, you may specify the time period of the log entries to process, by entering the number of seconds in the past to process. The default is 90 seconds. Adjust these values to suit your network usage and environment. The widget will send [ x ] number of entries or [ x ] seconds of entries, whichever is reached first. It is also possible to limit the entries by interface and action type. By default, all interfaces and action types are processed. The widget reuses the same logic applied by the standard pfSense dashboard widget for filtering log entries.

Finally, the frequency of updates can be configured. The default is 15 seconds. This value specifies the number of seconds between the widget sending a request to the backend service for updated geolocation information and the widget updating the dashboard with the response.

Refer to docs/CONFIGURATION.md for configuration options. To generate a sample configuration file as a starting point run:
```
./bin/ip2location-pfsense config create
```

## Running the service

The service can be run as a standalone binary, or as a Docker container. Having a Redis instance running is a prerequisite for both options. To run a new Redis instance as a Docker container, refer the sample script provided: 
scripts start_redis.sh 
```
./bin/ip2location-pfsense service
```

### Standalone binary

To run the service as a standalone binary, run the following command:
```
./bin/ip2location-pfsense service
```

### Docker (WIP)

For convenience, a WIP Dockerfile is provided. To build the Docker image, run:
```
docker build -t ip2location-pfsense .
```

To run the Docker image, command:
```
docker run -d --name ip2location-pfsense ip2location-pfsense
```

