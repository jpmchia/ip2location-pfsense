package ip2location

type Ip2LocationBasic struct {
	IP          string  `json:"ip"`
	CountryCode string  `json:"country_code"`
	CountryName string  `json:"country_name"`
	RegionName  string  `json:"region_name"`
	CityName    string  `json:"city_name"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	ZipCode     string  `json:"zip_code"`
	TimeZone    string  `json:"time_zone"`
	Asn         string  `json:"asn"`
	As          string  `json:"as"`
	IsProxy     bool    `json:"is_proxy"`
}

type Ip2LocationStarter struct {
	IP                 string  `json:"ip"`
	CountryCode        string  `json:"country_code"`
	CountryName        string  `json:"country_name"`
	RegionName         string  `json:"region_name"`
	CityName           string  `json:"city_name"`
	Latitude           float64 `json:"latitude"`
	Longitude          float64 `json:"longitude"`
	ZipCode            string  `json:"zip_code"`
	TimeZone           string  `json:"time_zone"`
	Asn                string  `json:"asn"`
	As                 string  `json:"as"`
	Isp                string  `json:"isp"`
	Domain             string  `json:"domain"`
	NetSpeed           string  `json:"net_speed"`
	IddCode            string  `json:"idd_code"`
	AreaCode           string  `json:"area_code"`
	WeatherStationCode string  `json:"weather_station_code"`
	WeatherStationName string  `json:"weather_station_name"`
	Elevation          int     `json:"elevation"`
	UsageType          string  `json:"usage_type"`
	IsProxy            bool    `json:"is_proxy"`
}

type Ip2LocationPlus struct {
	IP                 string  `json:"ip"`
	CountryCode        string  `json:"country_code"`
	CountryName        string  `json:"country_name"`
	RegionName         string  `json:"region_name"`
	CityName           string  `json:"city_name"`
	Latitude           float64 `json:"latitude"`
	Longitude          float64 `json:"longitude"`
	ZipCode            string  `json:"zip_code"`
	TimeZone           string  `json:"time_zone"`
	Asn                string  `json:"asn"`
	As                 string  `json:"as"`
	Isp                string  `json:"isp"`
	Domain             string  `json:"domain"`
	NetSpeed           string  `json:"net_speed"`
	IddCode            string  `json:"idd_code"`
	AreaCode           string  `json:"area_code"`
	WeatherStationCode string  `json:"weather_station_code"`
	WeatherStationName string  `json:"weather_station_name"`
	Mcc                string  `json:"mcc"`
	Mnc                string  `json:"mnc"`
	MobileBrand        string  `json:"mobile_brand"`
	Elevation          int     `json:"elevation"`
	UsageType          string  `json:"usage_type"`
	AddressType        string  `json:"address_type"`
	Continent          struct {
		Name        string   `json:"name"`
		Code        string   `json:"code"`
		Hemisphere  []string `json:"hemisphere"`
		Translation struct {
			Lang  string `json:"lang"`
			Value string `json:"value"`
		} `json:"translation"`
	} `json:"continent"`
	Country struct {
		Name        string `json:"name"`
		Alpha3Code  string `json:"alpha3_code"`
		NumericCode int    `json:"numeric_code"`
		Demonym     string `json:"demonym"`
		Flag        string `json:"flag"`
		Capital     string `json:"capital"`
		TotalArea   int    `json:"total_area"`
		Population  int    `json:"population"`
		Currency    struct {
			Code   string `json:"code"`
			Name   string `json:"name"`
			Symbol string `json:"symbol"`
		} `json:"currency"`
		Language struct {
			Code string `json:"code"`
			Name string `json:"name"`
		} `json:"language"`
		Tld         string `json:"tld"`
		Translation struct {
			Lang  string `json:"lang"`
			Value string `json:"value"`
		} `json:"translation"`
	} `json:"country"`
	Region struct {
		Name        string `json:"name"`
		Code        string `json:"code"`
		Translation struct {
			Lang  string `json:"lang"`
			Value string `json:"value"`
		} `json:"translation"`
	} `json:"region"`
	City struct {
		Name        string `json:"name"`
		Translation struct {
			Lang  interface{} `json:"lang"`
			Value interface{} `json:"value"`
		} `json:"translation"`
	} `json:"city"`
	TimeZoneInfo struct {
		Olson       string `json:"olson"`
		CurrentTime string `json:"current_time"`
		GmtOffset   int    `json:"gmt_offset"`
		IsDst       bool   `json:"is_dst"`
		Sunrise     string `json:"sunrise"`
		Sunset      string `json:"sunset"`
	} `json:"time_zone_info"`
	Geotargeting struct {
		Metro string `json:"metro"`
	} `json:"geotargeting"`
	IsProxy bool `json:"is_proxy"`
}

type WatchListEntry struct {
	Time        string  `json:"time"`
	IP          string  `json:"ip"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Direction   string  `json:"direction"`
	Act         string  `json:"act"`
	Reason      string  `json:"reason"`
	Interface   string  `json:"interface"`
	Realint     string  `json:"realint"`
	Version     string  `json:"version"`
	Srcip       string  `json:"srcip"`
	Dstip       string  `json:"dstip"`
	Srcport     string  `json:"srcport"`
	Dstport     string  `json:"dstport"`
	Proto       string  `json:"proto"`
	Protoid     string  `json:"protoid"`
	Length      string  `json:"length"`
	Rulenum     string  `json:"rulenum"`
	Subrulenum  string  `json:"subrulenum"`
	Anchor      string  `json:"anchor"`
	Tracker     string  `json:"tracker"`
	CountryCode string  `json:"country_code"`
	CountryName string  `json:"country_name"`
	RegionName  string  `json:"region_name"`
	CityName    string  `json:"city_name"`
	ZipCode     string  `json:"zip_code"`
	TimeZone    string  `json:"time_zone"`
	Asn         string  `json:"asn"`
	As          string  `json:"as"`
	IsProxy     bool    `json:"is_proxy"`
	WatchList   bool    `json:"watch_list"`
}
