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
* Storage of the watch list in the backend service
* Drill-down into details of the IP address  information provided by IP2Location.io
* Fully functional with a free IP2Location.io API account
* Support for more granular IP information available with paid IP2Location.io API accounts
* Utilises Redis cache to reduce the number of API requests
* Frequency of updates, number of log entries, API calls and cache duration configurable
* Support for SSL / TLS and backend secured with configurable API keys

### WIP features (part implemented)
* Backend service standalone mode, allowing monitoring of service's host network traffic without pfSense
* Configurable map tile providers for Leaflet.js
* Additional UI for cache monitoring and hit rates, plus configuration of limits
* Support for multiple pfSense hosts

<hr>

<img src="https://github.com/jpmchia/IP2Location-pfSense/blob/e1156d594e6e0e71c3cab91124f874811b4b3029/contrib/screenshots/Screenshot1.png?raw=true" width="95%">

<hr>

<img src="https://github.com/jpmchia/IP2Location-pfSense/blob/e1156d594e6e0e71c3cab91124f874811b4b3029/contrib/screenshots/Screenshot2.png?raw=true" width="95%">

Please refer to the wiki pages for further information. 

# Credits

IP2Location-pfSense is my first foray into development using the Go programming language, and incorporates the following open-source projects:

- [Cobra](https://github.com/spf13/cobra) - A Commander for modern Go CLI interactions
- [Viper](https://github.com/spf13/viper) - Go configuration with fangs
- [Echo](https://github.com/labstack/echo) - High performance, minimalist Go web framework
- [go-redis](https://github.com/redis/go-redis) - Type-safe Redis client for Golang
- [go-rejson](https://github.com/nitishm/go-rejson) - Golang client for redislabs' ReJSON module with support for multilple redis clients
- [IP2Location Go Package](https://github.com/ip2location/ip2location-go) - Go package for IP2Location database query

