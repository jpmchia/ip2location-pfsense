package pfsense

type FilterLog map[string]LogEntry

type LogEntry struct {
	Time       string   `json:"time"`
	Rulenum    string   `json:"rulenum"`
	Subrulenum string   `json:"subrulenum"`
	Anchor     string   `json:"anchor"`
	Tracker    string   `json:"tracker"`
	Realint    string   `json:"realint"`
	Interface  string   `json:"interface"`
	Reason     string   `json:"reason"`
	Act        string   `json:"act"`
	Direction  string   `json:"direction"`
	Version    string   `json:"version"`
	Tos        string   `json:"tos"`
	ECN        string   `json:"ecn"`
	TTL        string   `json:"ttl"`
	ID         string   `json:"id"`
	Offset     string   `json:"offset"`
	Flags      string   `json:"flags"`
	Protoid    string   `json:"protoid"`
	Proto      string   `json:"proto"`
	Length     string   `json:"length"`
	Srcip      string   `json:"srcip"`
	Dstip      string   `json:"dstip"`
	Srcport    *string  `json:"srcport,omitempty"`
	Dstport    *string  `json:"dstport,omitempty"`
	Src        string   `json:"src"`
	Dst        string   `json:"dst"`
	Datalen    *string  `json:"datalen,omitempty"`
	Tcpflags   *string  `json:"tcpflags,omitempty"`
	Seq        *string  `json:"seq,omitempty"`
	ACK        *string  `json:"ack,omitempty"`
	Window     *string  `json:"window,omitempty"`
	Urg        *string  `json:"urg,omitempty"`
	Options    []string `json:"options,omitempty"`
	ICMPType   *string  `json:"icmp_type,omitempty"`
	ICMPID     *string  `json:"icmp_id,omitempty"`
	ICMPSeq    *string  `json:"icmp_seq,omitempty"`
}

type Ip2Map struct {
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
	Srcport     *string `json:"srcport,omitempty"`
	Dstport     *string `json:"dstport,omitempty"`
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

type Ip2MapLocations map[string]Ip2Map

type Ip2ResultId struct {
	Id string `json:"id"`
}
