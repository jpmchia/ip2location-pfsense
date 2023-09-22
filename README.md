
# ip2location-pfsense

pfSense dashboard widget and backend service for displaying live IP geolocation information obtained from the IP2Location.io API. 

The dashboard widget and backend service are designed for use with the IP2Location.io API. Create a free account with up to 30k API requests per month: https://www.ip2location.io/

## Getting started

This project requires Go to be installed. Follow the instrucions to download and install Go from: https://go.dev/doc/install. At the time of writing, the latest version is 1.21.1.

The backend service can can be run on any platform and OS that is supported by Go, including Linux, Windows, MacOS and FreeBSD. For all available ports, refer to: https://go.dev/dl/

## Additional requirements

The service also requires either Redis Stack or Redis server with RedisJSON module installed. Refer to https://redis.io/download/ for available options. It may be hosted on the same host as the backend-service, an existing instance, or a new instance on a different host. 

Please note that TLS / SSL support is not yet completed, so requests / repsonses between the service, pfSense and Redis are not encrypted and therefore it is recommended that either all three are deployed to the same host, or they are deployed on hosts within a private network. 

# pfSense dashboard widget

For the moment, the dashboard widget and supporting files must be manually installed on your pfSense instance. pfSense package to be submitted. Until then, installation requires copying the contents of the pfSense/www folder. 
If you have SSH acceess configured on the pfSense box, installation is as easy as running the following command from a bash terminal. Installation of the widget should not modify or overwrite any existing, standard pfSense files.

The dashboard widget targets and has been tested on pfSense+ version 23.0.5.1. It may work on other versions, but this has not been tested. The widget is comprised of PHP AJAX and supporting Javascript packages.

```
```shellscript
./scripts/update_pfsense.sh <username>@<host>
```
