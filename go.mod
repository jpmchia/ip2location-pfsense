module github.com/jpmchia/ip2location-pfsense

go 1.21.1

require (
	github.com/jpmchia/ip2location-pfsense/backend v0.0.0-20230925173606-8e351a341cb9
	github.com/labstack/echo/v4 v4.11.1
)

require (
	github.com/labstack/gommon v0.4.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	golang.org/x/crypto v0.11.0 // indirect
	golang.org/x/net v0.12.0 // indirect
	golang.org/x/sys v0.12.0 // indirect
	golang.org/x/text v0.11.0 // indirect
)

replace github.com/jpmchia/ip2location-pfsense v0.0.0-unpublished => ../IP2Location-pfSense
