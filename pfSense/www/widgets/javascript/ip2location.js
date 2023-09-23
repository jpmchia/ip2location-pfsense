
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
                    var coordinates = [parseFloat(item.latitude), parseFloat(item.longitude)];
                    var marker = createMarker(item, coordinates);
                    marker.addTo(map);
                }
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
    var toolip_text = `<span class="ip2l-tooltip">${item.direction}<b>${item.interface}</b>${item.proto}:${item.srcip} =>${item.dstip}<br/>${item.city_name}, ${item.region_name}, ${item.ZipCode}, ${item.country_name}<br/>`;
    var popup = `${item.interface}:${item.proto} ${item.direction} => ${item.srcip}<br/> IP: ${item.ip}<br/> Country: ${item.country_name}<br/> City: ${item.city_name}<br/> Region: ${item.region_name}, ${item.ZipCode}<br/> Timezone: ${item.TimeZone}<br/> ASN: ${item.Asn}<br/> AS: ${item.As}<br/> IsProxy: ${item.IsProxy}`;
    var new_marker = L.marker(coordinates, { icon: new_icon }).bindTooltip(toolip_text).openTooltip();
    return new_marker;
}

function createTableRow(item) {
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
    var row = `<tr><td>${act_icon}</td><td>${item.time}</td><td>${item.interface}</td><td>${item.srcip}</td><td>${item.dstip}</td></tr><tr><td>${item.direction}</td><td>${item.proto}</td><td>${item.city_name}, ${item.country_name}</td><td>${item.As}, ${item.Asn}</td></tr>`;
    return row;
}

function generateTableHead(table, data) {
    let thead = table.createTHead();
    let row = thead.insertRow();
    for (let key of data) {
        let th = document.createElement("th");
        let text = document.createTextNode(key);
        th.appendChild(text);
        row.appendChild(th);
    }
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
