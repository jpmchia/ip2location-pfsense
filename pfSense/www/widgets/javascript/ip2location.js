
L.Control.Watermark = L.Control.extend({
    onAdd: function(map) {
        var img = L.DomUtil.create('img');
        img.src = ''
        img.src = '/widgets/images/ip2location/ip2location.png';
        img.style.width = '100px';
        return img;
    },
    onRemove: function(map) {
    // Nothing to do here
    }
});

function resetMap() {
    localStorage.removeItem("coords_x");
    localStorage.removeItem("coords_y");
    localStorage.removeItem("zoom");
    console.log("Map reset");
}

function formatParams(params) {
    return "?" + Object
        .keys(params)
        .map(function (key) {
            return key + "=" + encodeURIComponent(params[key])
        })
        .join("&")
}

function updateProgressBar(processed, total, elapsed, layersArray) {
    if (elapsed > 1000) {
        progress.style.display = 'block';
        progressBar.style.width = Math.round(processed/total*100) + '%';
    }
    if (processed === total) {
        progress.style.display = 'none';
    }
}

function fetchMapData(widgetkey, resultsid) {
    var requestdata = {
        ajax: "results",
        widgetkey: widgetkey,
        resultsid: resultsid
    };
    var xhr = new XMLHttpRequest();
    var url = '/widgets/widgets/ip2location.widget.php' + formatParams(requestdata)
    xhr.open('GET', url, true);

    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.onload = function () {
        if (xhr.status >= 200 && xhr.status < 400) {
            try {
                var data = JSON.parse(xhr.responseText);
                // Add markers to your Leaflet map here using the data
                for (var i = 0; i < data.length; i++) {
                    item = data[i];
                    if (item.latitude == 0 && item.longitude == 0)
                    {
                        continue;
                    }
                    var coordinates = [parseFloat(item.latitude), parseFloat(item.longitude)];
                    var marker = createMarker(item, coordinates);
                    var item_ip; 
                    switch (item.direction) {
                        case "in":
                            switch (item.act) {
                                case "pass":
                                    //allowInLayer.addLayer(marker);
                                    marker.addTo(allowInLayer);
                                    item_ip = item.srcip;
                                    break;
                                case "block":
                                case "reject":
                                default:
                                    //blockInLayer.addLayer(marker);
                                    marker.addTo(blockInLayer);
                                    item_ip = item.srcip;
                                    break;
                            }
                            break;
                        case "out":
                            switch (item.act) {
                                case "pass":
                                    //allowOutLayer.addLayer(marker);
                                    marker.addTo(allowOutLayer);
                                    item_ip = item.dstip;
                                    break;
                                case "block":
                                case "reject":
                                default:
                                    //blockOutLayer.addLayer(marker);
                                    marker.addTo(blockOutLayer);
                                    item_ip = item.dstip;
                                    break;
                            }
                            break;
                    }

                    //checkUpdateTrackedItems(item_ip);
                }

                
                var baseMaps = {
                    "OpenStreetMap" : openStreetMap
                };

                var overlayMaps = {
                    "Allow In": allowInLayer,
                    "Allow Out": allowOutLayer,
                    "Block In": blockInLayer,
                    "Block Out": blockOutLayer
                };

                map.addLayer(openStreetMap);
                map.addLayer(allowInLayer);
                map.addLayer(allowOutLayer);
                map.addLayer(blockInLayer);
                map.addLayer(blockOutLayer);
                
                var layerControl = L.control.layers(baseMaps, overlayMaps).addTo(map);
                
            } catch (e) {
                console.error('Error parsing JSON', e);
            }
        } else {
            console.error('Error fetching data');
        }
    };

    xhr.onerror = function () {
        console.error('Connection error');
    };
    xhr.send();
}

function createMarker(item, coordinates) {
    var new_icon;
    switch (item.version) {
        case "4":
            switch (item.act) {
                case "pass":
                    if (item.direction == "in") {
                        new_icon = L.icon({ iconUrl: '/widgets/images/ip2location/allow4in.png', iconSize: [24, 24], iconAnchor: [12, 22], popupAnchor: [12, 22] });
                    } else {
                        new_icon = L.icon({ iconUrl: '/widgets/images/ip2location/allow4out.png', iconSize: [24, 24], iconAnchor: [12, 22], popupAnchor: [12, 22] });
                    }
                    break;
                case "block":
                case "reject":
                default:
                    if (item.direction == "in") {
                        new_icon = L.icon({ iconUrl: '/widgets/images/ip2location/block4in.png', iconSize: [24, 24], iconAnchor: [12, 22], popupAnchor: [12, 22] });
                    } else {
                        new_icon = L.icon({ iconUrl: '/widgets/images/ip2location/block4out.png', iconSize: [24, 24], iconAnchor: [12, 22], popupAnchor: [12, 22] });
                    }
                    break;
            }
            break;
        case "6":
            switch (item.act) {
                case "pass":
                    if (item.direction == "in") {
                        new_icon = L.icon({ iconUrl: '/widgets/images/ip2location/allow6in.png', iconSize: [24, 24], iconAnchor: [12, 22], popupAnchor: [12, 22] });
                    } else {
                        new_icon = L.icon({ iconUrl: '/widgets/images/ip2location/allow6out.png', iconSize: [24, 24], iconAnchor: [12, 22], popupAnchor: [12, 22] });
                    }
                    break;
                case "block":
                case "reject":
                default:
                    if (item.direction == "in") {
                        new_icon = L.icon({ iconUrl: '/widgets/images/ip2location/block6in.png', iconSize: [24, 24], iconAnchor: [12, 22], popupAnchor: [12, 22] });
                    } else {
                        new_icon = L.icon({ iconUrl: '/widgets/images/ip2location/block6out.png', iconSize: [24, 24], iconAnchor: [12, 22], popupAnchor: [12, 22] });
                    }
                    break;
            }
            break;
    }
    var coordinates = [parseFloat(item.latitude), parseFloat(item.longitude)];
    var toolip_text = createTooltip(item);
    var new_marker = L.marker(coordinates, { icon: new_icon }).bindTooltip(toolip_text).openTooltip();

    new_marker.on('click', function (e) {
        addItemToWatchlist(item);
    });
    
    return new_marker;
}

function textise(text, prefix) {
    if (text == undefined || text == null) {
        return "";
    } else {
        return prefix + text;
    }
}


function createTooltip(item) {
    return `<table class="ip2l_ttt">
    <tr>
        <td class="ip2l_ttl" rowspan="4"><p>${item.interface}</p><p>${action_icon(item)}</p><p>${item.proto}</p></td>
        <td class="ip2l_ttr" colspan="2">${item.time}</td>
    </tr>
    <tr>

        <td class="ip2l_tth">Src IP: </td>
        <td class="ip2l_ttr">${item.srcip}</td>
    </tr>
    <tr>
        <td class="ip2l_tth">Dst IP: </td>
        <td class="ip2l_ttr">${item.dstip}${textise(item.dstport, " : ")}</td>
    </tr>
    <tr>
        <td class="ip2l_ttr" colspan="3">${textise(item.city_name, "")}, ${textise(item.region_name, "")}, ${textise(item.country_name, "")}</td>
    </tr>
    <tr>
    </table>`;
}


function action_icon(item) {
    var act_icon;
    switch (item.act) {
        case "pass":
            if (item.direction == "in") {
                act_icon = '<i class="fa fa-arrow-circle-down text-success"></i>';
            } else {
                act_icon = '<i class="fa fa-arrow-circle-up text-success"></i>';
            }
            break;
        case "block":
        case "reject":
        default:
            if (item.direction == "in") {
                act_icon = '<i class="fa fa-arrow-circle-down text-danger"></i>';
            } else {
                act_icon = '<i class="fa fa-arrow-circle-up text-danger"></i>';
            }
            break;
    }
    return act_icon;
}

function fwrule_icon(item) {
    var fwr_icon;
    switch (item.act) {
        case "pass":
            fwr_icon = 'fa fa-minus-square-o';
            break;
        case "block":
        case "reject":
            fwr_icon = 'fa fa-plus-square-o';
            break;
        default:
            // fwr_icon = 'fa fa-info icon-pointer';
            break;
    }
    return fwr_icon;
}

function createTableRow1(item) {
    var act_icon = action_icon(item);
    var rule_icon = fwrule_icon(item);
    var rule_link;

    switch (item.act) {
        case "block":
        case "reject":
            rule_link = `easyrule.php?action=${("pass")}&amp;int=${item.interface}&amp;proto=${item.proto}&amp;src=${item.srcip}&amp;dst=${item.dstip}&amp;dstport=${item.dstpot}&amp;ipproto=${item.version}`;
            break;
        case "pass":
            rule_link = `easyrule.php?action=${("block")}&amp;int=${item.interface}&amp;proto=${item.proto}&amp;src=${item.srcip}&amp;dst=${item.dstip}&amp;dstport=${item.dstpot}&amp;ipproto=${item.version}`;
            break;            
        }

    // rule_link = `easyrule.php?action=block&int=${item.interface}&src=${item.ip}&ipproto=${item.proto}`;
    // <a class="fa fa-plus-square-o icon-pointer icon-primary"
    // href="easyrule.php?action=pass&amp;int=wg1&amp;proto=tcp&amp;src=[2001:da8:d00a:2::]&amp;dst=[2a0e:97c0:5c1::]&amp;dstport=19128&amp;ipproto=inet6"
    // title="" data-original-title="EasyRule: Pass this traffic">

    var row1 = `<td rowspan="2" style="vertical-align: middle;">${act_icon}</td>
        <td style="white-space: nowrap;">${item.time}</td>
        <td>${item.interface}</td>
        <td>${item.ip}</td>
        <td id="ip2l_status">${item.hits}</td>
        <td><a href="${rule_link}" title="EasyRule: Add a rule to Block / Unblock this IP."><i class="${rule_icon}" style="cursor: pointer;"></i></a></td>
        <td><i class="fa fa-info" style="cursor: pointer;" onclick="javascript:showIP2LDetail('${item.ip}');" title="View IP location details."></i></td>`
    return row1;
}


function createTableRow2(item) { 
    var row2 = `
    <td style="white-space: nowrap;">${item.lastSeen}</td>
    <td colspan="3">${item.city}, ${item.country}</td>
    <td colspan="1"><i class="fa fa-chevron-right" style="cursor: pointer;" onclick="javascript:showIP2Watchlist('${item.ip}');" title="View watch list detail."></i></td>
    <td><i class="fa fa-times text-danger" style="cursor: pointer;" onclick="javascript:removeItem('${item.ip}');" title="Remove from watch list."></td>`
    return row2;
}
    

function removeItem(item_ip) {
    var requestdata = { "ip": item_ip };
    var url = apiUrl + '/api/watch?ip=' + item_ip;
    fetch(url, {
        headers: { Authorization: `Bearer ${ip2l_token}`},
        method: 'DELETE'
    }).then(response => response.json())
    .then(data => {
        recreateWatchlistTable(data);
    });
    console.log("Removed item from watchlist: " + item_ip);
}

function addItemToWatchlist(item) {
    if (item.direction == "in") {{
        item_ip = item.srcip}
    } else {{
        item_ip = item.dstip}};

    var url = apiUrl + '/api/watch?ip=' + item_ip;
    var post_data =  JSON.stringify(item);
    fetch(url, {
        headers: { Authorization: `Bearer ${ip2l_token}`, 'Content-Type': 'application/json'},
        method: 'POST',
        body: post_data
    }).then(response => response.json())
    .then(data => {
        console.log(data);
        recreateWatchlistTable(data);
    });
    console.log("Added item to watchlist: " + item_ip);
}

function loadIp2LTablefromBackend() {
    var itemCount = 0;
    var url = apiUrl + '/api/watchlist';
    var data;
    fetch(url, {
        headers: { Authorization: `Bearer ${ip2l_token}`},
        method: 'GET'
    }).then(response => response.json())
    .then(data => {
        console.log(data);
        recreateWatchlistTable(data);
    });
}

function recreateWatchlistTable(data) {
    var span = document.getElementById("ip2l-details");
    var table = document.getElementById("ip2l-table");
    var tbody = document.getElementById("ip2l-tbody");

    if (data == null || data.length == 0) {
        table.setAttribute("style", "display: none; visibility: hidden;");
        return;
    } 
    
    tbody.innerHTML = "";
    
    for (var i = 0; i < data.length; i++) {
        item = data[i];
        var row1 = tbody.insertRow();
        row1.innerHTML = createTableRow1(item);
        var row2 = tbody.insertRow();
        row2.innerHTML = createTableRow2(item);
    }
    
    table.setAttribute("style", "display: table; visibility: visible;");
}

function showIP2LDetail(item_ip) {
    var ip2detailUrl = htmlUrl + "/ip2ldetails.html?ip=" + item_ip + "'";
        
    window.winboxIp2LDetails = new WinBox('IP2Location details', {
        top: 60,
        right: 5,
        bottom: 15,
        left: 5,
        width: '800px',
        height: '600px',
        url: ip2detailUrl,
    });
       
    console.log("showIP2LDetail: " + item_ip + " - " + ip2detailUrl);
}

function showIP2Watchlist(item_ip) {
    var ip2detailUrl = htmlUrl + "/ip2location.html?ip=" + item_ip + "'";
        
    window.winboxIp2LDetails = new WinBox('Watchlist details', {
        top: 60,
        right: 5,
        bottom: 15,
        left: 5,
        width: '800px',
        height: '600px',
        url: ip2detailUrl,
    });
       
    console.log("showIP2LDetail: " + item_ip + " - " + ip2detailUrl);
}

