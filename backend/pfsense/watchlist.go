package pfsense

import (
	"fmt"
	"log"
	"os"

	"github.com/jpmchia/ip2location-pfsense/cache"
	"github.com/jpmchia/ip2location-pfsense/util"
)

type WatchListItem struct {
	IpAddress  string   `json:"ip"`
	LogEntries []Ip2Map `json:"entries"`
	Count      int      `json:"hits"`
	FirstSeen  string   `json:"time"`
	LastSeen   string   `json:"lastSeen"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	Direction  string   `json:"direction"`
	Act        string   `json:"act"`
	Srcip      string   `json:"srcip"`
	Dstip      string   `json:"dstip"`
	Dstport    string   `json:"dstport"`
	Interface  string   `json:"interface"`
	Proto      string   `json:"proto"`
	Rulenum    string   `json:"rulenum"`
}

type WatchListDisplayItem struct {
	IpAddress string `json:"ip"`
	Count     int    `json:"hits"`
	FirstSeen string `json:"time"`
	LastSeen  string `json:"lastSeen"`
	City      string `json:"city"`
	Country   string `json:"country"`
	Direction string `json:"direction"`
	Act       string `json:"act"`
	Srcip     string `json:"srcip"`
	Dstip     string `json:"dstip"`
	Dstport   string `json:"dstport"`
	Interface string `json:"interface"`
	Proto     string `json:"proto"`
	Rulenum   string `json:"rulenum"`
}

type WatchList map[string]WatchListItem

const WatchListCache string = "watchlist"

var ActiveWatchList WatchList
var DisplayWatchList []WatchListDisplayItem
var WatchListKeys []string

// NewWatchList creates a new WatchList
func NewWatchList() WatchList {
	return make(WatchList)
}

// Loads only the keys from the cache
func LoadWatchListKeys() {
	var instance = cache.Instance(WatchListCache)
	keys, err := instance.Keys("*")
	util.HandleError(err, "[pfsense] Unable to load keys from cache: %v", err)
	WatchListKeys = keys
}

func LoadWatchListDisplayItems() {
	var instance = cache.Instance(WatchListCache)
	keys, err := instance.Keys("*")
	util.HandleError(err, "[pfsense] Unable to load keys from cache: %v", err)
	for _, key := range keys {
		wli, err := instance.Get(key)
		util.HandleError(err, "[pfsense] Unable to load WatchListItem from cache: %v", err)
		DisplayWatchList = append(DisplayWatchList, WatchListDisplayItem{
			IpAddress: wli.(WatchListItem).IpAddress,
			Count:     wli.(WatchListItem).Count,
			FirstSeen: wli.(WatchListItem).FirstSeen,
			LastSeen:  wli.(WatchListItem).LastSeen,
			City:      wli.(WatchListItem).City,
			Country:   wli.(WatchListItem).Country,
			Direction: wli.(WatchListItem).Direction,
			Srcip:     wli.(WatchListItem).Srcip,
			Dstip:     wli.(WatchListItem).Dstip,
			Dstport:   wli.(WatchListItem).Dstport,
			Act:       wli.(WatchListItem).Act,
			Interface: wli.(WatchListItem).Interface,
			Proto:     wli.(WatchListItem).Proto,
			Rulenum:   wli.(WatchListItem).Rulenum,
		})
	}
}

// Loads a WatchListItem from the cache
func LoadWatchListItem(key string) {
	var instance = cache.Instance(WatchListCache)
	wli, err := instance.Get(key)
	util.HandleError(err, "[pfsense] Unable to load WatchListItem from cache: %v", err)
	ActiveWatchList[key] = wli.(WatchListItem)
}

// Loads all WatchListItems from the cache
func LoadWatchList() {
	ActiveWatchList = NewWatchList()
	LoadWatchListKeys()
	for _, key := range WatchListKeys {
		LoadWatchListItem(key)
	}
}

// SaveToCache saves a WatchListItem to the cache
func SaveToCache(wli WatchListItem) error {
	var instance = cache.Instance(WatchListCache)
	_, err := instance.Set(wli.IpAddress, wli)
	util.HandleError(err, "[watchlist] Unable to save WatchListItem to cache: %v", err)
	return err
}

// PrintWatchList prints the WatchList
func PrintWatchList(wl WatchList) {
	for _, item := range wl {
		fmt.Printf("%v\n", item)
	}
}

// WriteWatchList writes the WatchList to a file
func WriteWatchList(wl WatchList, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for _, item := range wl {
		fmt.Fprintf(file, "%v\n", item)
	}

	util.Log("[pfsense] WatchList written to %v\n", filename)
}

// Add adds a new log entry to the WatchList
func (wl WatchList) Add(ip string, ip2mapEntry Ip2Map) {
	if _, ok := wl[ip]; !ok {
		wl[ip] = WatchListItem{
			IpAddress:  ip,
			LogEntries: []Ip2Map{ip2mapEntry},
			Count:      1,
			FirstSeen:  ip2mapEntry.Time,
			LastSeen:   ip2mapEntry.Time,
			City:       ip2mapEntry.CityName,
			Country:    ip2mapEntry.CountryName,
			Direction:  ip2mapEntry.Direction,
			Srcip:      ip2mapEntry.Srcip,
			Dstip:      ip2mapEntry.Dstip,
			Dstport: (func() string {
				if ip2mapEntry.Dstport != nil {
					return *ip2mapEntry.Dstport
				}
				return ""
			})(),
			Act:       ip2mapEntry.Act,
			Interface: ip2mapEntry.Interface,
			Proto:     ip2mapEntry.Proto,
			Rulenum:   ip2mapEntry.Rulenum,
		}
		util.LogDebug("[watchlist] Adding new WatchListItem: %v", wl[ip])
	} else {
		wl[ip] = WatchListItem{
			IpAddress:  ip,
			LogEntries: append(wl[ip].LogEntries, ip2mapEntry),
			Count:      wl[ip].Count + 1,
			FirstSeen:  wl[ip].FirstSeen,
			LastSeen:   ip2mapEntry.Time,
			City:       ip2mapEntry.CityName,
			Country:    ip2mapEntry.CountryName,
			Direction:  ip2mapEntry.Direction,
			Srcip:      ip2mapEntry.Srcip,
			Dstip:      ip2mapEntry.Dstip,
			Dstport: (func() string {
				if ip2mapEntry.Dstport != nil {
					return *ip2mapEntry.Dstport
				}
				return ""
			})(),
			Act:       ip2mapEntry.Act,
			Interface: ip2mapEntry.Interface,
			Proto:     ip2mapEntry.Proto,
			Rulenum:   ip2mapEntry.Rulenum,
		}
		util.LogDebug("[watchlist] Updating WatchListItem: %v", wl[ip])
	}
}

// Remove removes an IP address from the WatchList
func (wl WatchList) Remove(ip string) {
	for _, item := range wl {
		if item.IpAddress == ip {
			util.LogDebug("[watchlist] Removing WatchListItem: %v", item)
			delete(wl, ip)
		}
	}
	var instance = cache.Instance(WatchListCache)
	_, err := instance.Delete(ip)
	util.HandleError(err, "[watchlist] Unable to remove WatchListItem from cache: %v", err)
}

// Get returns the WatchListItem for the given IP address
func (wl WatchList) Get(ip string) (WatchListItem, bool) {
	item, ok := wl[ip]
	return item, ok
}

// GetCount returns the number of IP addresses in the WatchList
func (wl WatchList) GetCount() int {
	return len(wl)
}

// Contains returns true if the IP address is in the WatchList
func (wl WatchList) Contains(ip string) bool {
	_, ok := wl[ip]
	return ok
}

// AddLogEntry adds a log entry to an existing WatchListItem in the WatchList
func (wl WatchList) AddLogEntry(ip string, ip2MapEntry Ip2Map) error {
	if item, ok := wl[ip]; ok {
		item.Count++
		item.LogEntries = append(item.LogEntries, ip2MapEntry)
		item.LastSeen = ip2MapEntry.Time
		wl[ip] = item
	}
	err := SaveToCache(wl[ip])
	util.HandleError(err, "[watchlist] Unable to save WatchListItem to cache: %v", err)
	return err
}

// GetIPs returns a slice of IP addresses in the WatchList
func (wl WatchList) GetIPs() []string {
	ips := make([]string, 0, len(wl))
	for ip := range wl {
		ips = append(ips, ip)
	}
	return ips
}

// GetLogEntries returns a slice of log entries in the WatchList
func (wl WatchList) GetLogEntries() []Ip2Map {
	logEntries := make([]Ip2Map, 0)
	for _, item := range wl {
		logEntries = append(logEntries, item.LogEntries...)
	}
	return logEntries
}

// GetWatchListDisplayItems returns a slice of WatchListDisplayItems
func (wl WatchList) GetWatchListDisplayItems() []WatchListDisplayItem {
	wldi := make([]WatchListDisplayItem, 0)
	for _, item := range wl {
		wldi = append(wldi, WatchListDisplayItem{
			IpAddress: item.IpAddress,
			Count:     item.Count,
			FirstSeen: item.FirstSeen,
			LastSeen:  item.LastSeen,
			City:      item.City,
			Country:   item.Country,
			Direction: item.Direction,
			Act:       item.Act,
			Srcip:     item.Srcip,
			Dstip:     item.Dstip,
			Dstport:   item.Dstport,
			Interface: item.Interface,
			Proto:     item.Proto,
			Rulenum:   item.Rulenum,
		})
	}
	return wldi
}

func (wl WatchList) GetDisplayItem(ip string) WatchListDisplayItem {
	item := wl[ip]
	wldi := WatchListDisplayItem{
		IpAddress: item.IpAddress,
		Count:     item.Count,
		FirstSeen: item.FirstSeen,
		LastSeen:  item.LastSeen,
		City:      item.City,
		Country:   item.Country,
		Direction: item.Direction,
		Act:       item.Act,
		Srcip:     item.Srcip,
		Dstip:     item.Dstip,
		Dstport:   item.Dstport,
		Interface: item.Interface,
		Proto:     item.Proto,
		Rulenum:   item.Rulenum,
	}

	return wldi
}
