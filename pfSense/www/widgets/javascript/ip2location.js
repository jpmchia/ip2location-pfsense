

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

                    checkUpdateTrackedItems(item_ip);
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
                
                //L.control.layers(baseMaps).addTo(map);
                // var coords_x = localStorage.getItem("coords_x");
                // var coords_y = localStorage.getItem("coords_y");
                // var zoom = localStorage.getItem("zoom");
                // if (coords_x != null && coords_y != null && zoom != null) {
                //     map.setView([coords_x, coords_y], zoom);
                // }
                // map.on('moveend', function () {
                //     localStorage.setItem("coords_x", map.getCenter().lat);
                //     localStorage.setItem("coords_y", map.getCenter().lng);
                //     localStorage.setItem("zoom", map.getZoom());
                // });

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
        addIp2LDetails(item);
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


function checkUpdateTrackedItems(item_ip)
{
    var existngItems = JSON.parse(window.localStorage.getItem("ip2ldetails"));
    if (existngItems == null) {
        return;
    }
    var item = existngItems[item_ip];
    if (item == null) {
        return;
    }
    updateTrackedItem(item);
}

function updateTrackedItem(item) {
    item.hitcount++;
    saveItemToLocalStorage(item);
    recreateIp2LDetailsTable();
}

function generateTable(table, data) {
    for (let element of data) {
        let row = table.insertRow();
        for (key in element) {
            let cell = row.insertCell();
            let text = document.createTextNode(element[key]);
            cell.appendChild(text);
        }
    }
}

function createTableRow(item) {
    var act_icon = action_icon(item);
    var item_ip;
    
    if (item.direction == "in") {{
        item_ip = item.srcip}
    } else {{
        item_ip = item.dstip}};

    var row = `
    <td>${act_icon}</td>
    <td>${item.time}</td>
    <td>${item.interface}</td>
    <td>${item_ip}</td>
    <td id="ip2l_status">${item.hitcount}</td>
    <td id="ip2l_action_detail"><i class="fa fa-chevron-right" style="cursor: pointer;" onclick="javascript:showIP2LDetail('${item_ip}');" title="View details."></td>
    <td id="ip2l_action_remove"><i class="fa fa-times text-danger" style="cursor: pointer;" onclick="javascript:removeItem('${item_ip}');" title="Remove from watch list."></td>
    `
    //var row = `<tr><td>${act_icon}</td><td>${item.time}</td><td>${item.interface}</td><td>${item.srcip}</td><td>${item.dstip}</td></tr><tr><td>${item.direction}</td><td>${item.proto}</td><td>${item.city_name}, ${item.country_name}</td><td>${item.As}, ${item.Asn}</td></tr>`;
    return row;
}

function removeItem(item_ip)
{
    var existngItems = JSON.parse(window.localStorage.getItem("ip2ldetails"));
    if (existngItems == null) {
        return;
    }
    for (var key in existngItems) {
        if (!existngItems.hasOwnProperty(key)) {
            continue;
        }
        if (key == item_ip) {
            delete existngItems[key];
            window.localStorage.setItem("ip2ldetails", JSON.stringify(existngItems));
            console.log("Item removed from local storage: " +  item_ip); // JSON.stringify(existngItems));
            recreateIp2LDetailsTable();
            return;
        }
    }
}

function saveItemToLocalStorage(item) {
    var existngItems = JSON.parse(window.localStorage.getItem("ip2ldetails"));
    if (existngItems == null) {
        existngItems = {};
    }
    var item_ip;
    if (item.direction == "in") {{
        item_ip = item.srcip}
    } else {{
        item_ip = item.dstip}};
    existngItems[item_ip] = item;
    window.localStorage.setItem("ip2ldetails", JSON.stringify(existngItems));
    console.log("Item saved to local storage: " + item_ip ); // JSON.stringify(existngItems));
}

function addIp2LDetails(item) {

    item.hitcount = 1;
    saveItemToLocalStorage(item);
    recreateIp2LDetailsTable();
}

function addIp2LItemToDetailsTable(detailedItem) {
    var table = document.getElementById("ip2l-table");
    var tbody = document.getElementById("ip2l-tbody");
    var row = tbody.insertRow();
    row.innerHTML = createTableRow(detailedItem);
}

function loadIp2LTablefromSession() {
    var itemCount = 0;
    var detailedItems = JSON.parse(localStorage.getItem("ip2ldetails"));
    
    if (detailedItems == null) {
        return 0;
    }

    var tbody = document.getElementById("ip2l-tbody");
    tbody.innerHTML = "";

    for (var key in detailedItems) {
        if (!detailedItems.hasOwnProperty(key)) {
            continue;
        }
        itemCount++;
        addIp2LItemToDetailsTable(detailedItems[key]);
    }
    return itemCount;
}

function recreateIp2LDetailsTable() {
    var span = document.getElementById("ip2l-details");
    var table = document.getElementById("ip2l-table");

    var itemCount = loadIp2LTablefromSession();
    var rowCount = table.rows.length;
    
    if (itemCount > 0 || rowCount > 1) {
        span.setAttribute("style", "display: block; visibility: visible;");
    } else {
        span.setAttribute("style", "display: none; visibility: hidden;");
    }
}

function getIp2LDetails(item) {
    return item;
    var requestdata = {
        ajax: "results",
        widgetkey: widgetkey,
        resultsid: resultsid
    };
    var xhr = new XMLHttpRequest();
    var url = apiUrl + '/api/ip2location' + formatParams(requestdata)
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
                    var row = table.insertRow();
                    var cell = row.insertCell();
                    cell.innerHTML = createTableRow(item);
                }
            } catch (e) {
                console.error('Error parsing JSON', e);
            }
        } else {
            console.error('Error fetching data');
        }
    }
    xhr.onerror = function () {
        console.error('Connection error');
    }
    xhr.send();
    span.setAttribute("style", "display: block; visibility: visible;");    
}

function showIP2LDetail(item_ip) {
    window.winboxIp2LDetails = new WinBox('IP2Location details', {
        root: document.querySelector(widgetkey),
        top: 60,
        right: 5,
        bottom: 15,
        left: 5,
        url: 
    });
    console.log("showIP2LDetail: " + item_ip);
}