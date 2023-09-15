<?php
/*
 * log.widget.php
 *
 * part of pfSense (https://www.pfsense.org)
 * Copyright (c) 2004-2013 BSD Perimeter
 * Copyright (c) 2013-2016 Electric Sheep Fencing
 * Copyright (c) 2014-2023 Rubicon Communications, LLC (Netgate)
 * Copyright (c) 2007 Scott Dale
 * All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

require_once("guiconfig.inc");
require_once("pfsense-utils.inc");
require_once("functions.inc");
require_once("syslog.inc");

global $g, $pattern;

function extract_ip_entries($logarr, $seconds)
{
	$datetime = date("Y-m-d H:i:s");
	$timestamp = strtotime($datetime);
	$time = $timestamp - $seconds;
	$dateTimeCap = date("Y-m-d H:i:s", $time);

	printf("From: %s to %s <br/>", $datetime, $dateTimeCap);

	$publicIPsLogs = [];
	$count = 0;

	foreach ($logarr as $entry) {
		if ($entry['time'] < $dateTimeCap) {
			continue;
		}
		$srcIP = $entry['srcip'];
		$dstIP = $entry['dstip'];

		if ($entry['direction'] == "in") {
			$safeSrcKey = sha1($srcIP);
			$publicIPsLogs[$safeSrcKey] = $entry;
			$count = $count + 1;
		} else {
			$safeDstKey = sha1($dstIP);
			$publicIPsLogs[$safeDstKey] = $entry;
			$count = $count + 1;
		}
	}
	printf("Displaying locations %s of %s IPs from the past %s seconds.\n", count($publicIPsLogs), $count, $seconds);

	return $publicIPsLogs;
}

if ($_REQUEST['widgetkey'] && !$_REQUEST['ajax']) {

	set_customwidgettitle($user_settings);

	if (is_numeric($_POST['filterlogentries'])) {
		$user_settings['widgets'][$_POST['widgetkey']]['filterlogentries'] = $_POST['filterlogentries'];
	} else {
		unset($user_settings['widgets'][$_POST['widgetkey']]['filterlogentries']);
	}

	$acts = array();
	if ($_POST['actpass']) {
		$acts[] = "Pass";
	}
	if ($_POST['actblock']) {
		$acts[] = "Block";
	}
	if ($_POST['actreject']) {
		$acts[] = "Reject";
	}

	if (!empty($acts)) {
		$user_settings['widgets'][$_POST['widgetkey']]['filterlogentriesacts'] = implode(" ", $acts);
	} else {
		unset($user_settings['widgets'][$_POST['widgetkey']]['filterlogentriesacts']);
	}
	unset($accts);

	if (!empty($include_ips)) {
		$user_settings['widgets'][$_POST['widgetkey']]['includeips'] = implode(" ", $include_ips);
	} else {
		unset($user_settings['widgets'][$_POST['widgetkey']]['includeips']);
	}
	unset($include_ips);

	if (($_POST['filterlogentriesinterfaces']) and ($_POST['filterlogentriesinterfaces'] != "All")) {
		$user_settings['widgets'][$_POST['widgetkey']]['filterlogentriesinterfaces'] = trim($_POST['filterlogentriesinterfaces']);
	} else {
		unset($user_settings['widgets'][$_POST['widgetkey']]['filterlogentriesinterfaces']);
	}

	if (is_numeric($_POST['filterlogentriesinterval'])) {
		$user_settings['widgets'][$_POST['widgetkey']]['filterlogentriesinterval'] = $_POST['filterlogentriesinterval'];
	} else {
		unset($user_settings['widgets'][$_POST['widgetkey']]['filterlogentriesinterval']);
	}

	if (is_numeric($_POST['logseconds'])) {
		$user_settings['widgets'][$_POST['widgetkey']]['logseconds'] = $_POST['logseconds'];
	} else {
		unset($user_settings['widgets'][$_POST['widgetkey']]['filterlogentries']);
	}

	if (is_URL($_POST['ip2location_service_url'])) {
		$user_settings['widgets'][$_POST['widgetkey']]['ip2location_service_url'] = $_POST['ip2location_service_url'];
	} else {
		unset($user_settings['widgets'][$_POST['widgetkey']]['ip2location_service_url']);
	}

	if (is_URL($_POST['ip2location_results_url'])) {
		$user_settings['widgets'][$_POST['widgetkey']]['ip2location_results_url'] = $_POST['ip2location_results_url'];
	} else {
		unset($user_settings['widgets'][$_POST['widgetkey']]['ip2location_results_url']);
	}

	save_widget_settings($_SESSION['Username'], $user_settings["widgets"], gettext("Saved Filter Log Entries via Dashboard."));

	Header("Location: /");

	exit(0);
}

function send_filterlog($publicIPsLogs, $url) {
	// Encode the data
	$data = json_encode($publicIPsLogs);

	// Set cURL options
	$ch = curl_init();

	$optArray = array(
		CURLOPT_URL => $url,
		CURLOPT_RETURNTRANSFER => true,
		CURLOPT_HTTPHEADER => ['Content-Type: application/json'],
		CURLOPT_POST => 1,
		CURLOPT_POSTFIELDS => $data
	);
	curl_setopt_array($ch, $optArray);
	var_dump($url);

	$response = curl_exec($ch);
	var_dump($response);

	if( !$result = curl_exec($ch)) {
		trigger_error(curl_error($ch));
	}
	curl_close($ch);

	return $response;
}
?>

<?php

// When this widget is included in the dashboard, $widgetkey is already defined before the widget is included.
// When the ajax call is made to refresh the firewall log table, 'widgetkey' comes in $_REQUEST.
if ($_REQUEST['widgetkey']) {
	$widgetkey = $_REQUEST['widgetkey'];
}

$mapHeight = "340px";

$iface_descr_arr = get_configured_interface_with_descr();

$nentries = isset($user_settings['widgets'][$widgetkey]['filterlogentries']) ? $user_settings['widgets'][$widgetkey]['filterlogentries'] : 5000;

$nentriesacts = isset($user_settings['widgets'][$widgetkey]['filterlogentriesacts']) ? $user_settings['widgets'][$widgetkey]['filterlogentriesacts'] : 'All';

$nentriesinterfaces = isset($user_settings['widgets'][$widgetkey]['filterlogentriesinterfaces']) ? $user_settings['widgets'][$widgetkey]['filterlogentriesinterfaces'] : 'All';

$filterfieldsarray = array(
	"act" => $nentriesacts,
	"interface" => isset($iface_descr_arr[$nentriesinterfaces]) ? $iface_descr_arr[$nentriesinterfaces] : $nentriesinterfaces
);

$nentriesinterval = isset($user_settings['widgets'][$widgetkey]['filterlogentriesinterval']) ? $user_settings['widgets'][$widgetkey]['filterlogentriesinterval'] : 10;

$nseconds = isset($user_settings['widgets'][$widgetkey]['logseconds']) ? $user_settings['widgets'][$widgetkey]['logseconds'] : 30;

$filter_logfile = "{$g['varlog_path']}/filter.log";

$filterlog = conv_log_filter($filter_logfile, $nentries, 5000, $filterfieldsarray);

$publicIPsLogs = extract_ip_entries($filterlog, $nseconds);

$widgetkey_nodash = str_replace("-", "", $widgetkey);

//$ipList = array_keys($publicIPsLogs);

$ip2location_submit_url = isset($user_settings['widgets'][$widgetkey]['ip2location_submit_url']) ? $user_settings['widgets'][$widgetkey]['ip2location_submit_url'] : "http://192.168.1.51:9999/filterlog";
$ip2location_results_url = isset($user_settings['widgets'][$widgetkey]['ip2location_results_url']) ? $user_settings['widgets'][$widgetkey]['ip2location_results_url'] : 'http://192.168.1.51:9999/ip2geomap';


$resultsId = send_filterlog($publicIPsLogs, $ip2location_submit_url);


/*
// Checking if any error occurs
// during request or not
if($e = curl_error($ch)) {
	echo $e;
} else {
	var_dump($response);
	// Decoding JSON data
	$decodedData =
		json_decode($response, true);
	// Outputting JSON data in decoded form
	var_dump($decodedData);
}
*/


/*
$errors = curl_error($ch);
$response = curl_getinfo($ch, CURLINFO_HTTP_CODE);
curl_close($ch);
var_dump($result);
var_dump($errors);
var_dump($response);
*/
?>

<link rel="stylesheet" href="/vendor/leaflet/leaflet.css"/>
<script src="/vendor/leaflet/leaflet.js"></script>
<style>
	#map {
		height: <?=gettext($mapHeight); ?>;
	}
</style>
<div id="map">

</div>
<script>
	//<![CDATA[
	var map = L.map('map').setView([51.505, -0.09], 4);

	L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
		maxZoom: 19,
		attribution: '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>'
	}).addTo(map);

	function onMapClick(e) {

	}
	map.on('click', onMapClick);
	//]]>
</script>
<?php if (!$_REQUEST['ajax']) {
?>

<script type="text/javascript">
//<![CDATA[
	var logWidgetLastRefresh<?=htmlspecialchars($widgetkey_nodash)?> = <?=time()?>;
//]]>
</script>

<?php }

printf("Results ID = " . $resultsId);

?>

<?php

/* for AJAX response, we only need the panel-body */
if ($_REQUEST['ajax']) {

	$ch = curl_init();

	$resultsId = 1;
	$optArray = array(
		CURLOPT_URL => $ip2location_results_url,
		CURLOPT_RETURNTRANSFER => true,
		CURLOPT_HTTPHEADER => ['Content-Type: application/json'],
		CURLOPT_POST => 1,
		CURLOPT_POSTFIELDS => $resultsId
	);
	curl_setopt_array($ch, $optArray);

	//$ip2result = curl_exec($ch);
	//$errors = curl_error($ch);
	//$response = curl_getinfo($ch, CURLINFO_HTTP_CODE);
	//curl_close($ch);
	//var_dump($ip2result);

	exit;
}
?>


<script type="text/javascript">
//<![CDATA[

events.push(function(){
	// --------------------- Centralized widget refresh system ------------------------------

	// Callback function called by refresh system when data is retrieved
	function logs_callback(s) {
		$(<?=json_encode('#widget-' . $widgetkey . '_panel-body')?>).html(s);
	}

	// POST data to send via AJAX
	var postdata = {
		ajax: "ajax",
		widgetkey : <?=json_encode($widgetkey)?>,
		lastsawtime: logWidgetLastRefresh<?=htmlspecialchars($widgetkey_nodash)?>
	 };

	// Create an object defining the widget refresh AJAX call
	var logsObject = new Object();
	logsObject.name = "IP2Location Logs";
	logsObject.url = "/widgets/widgets/ip2location.widget.php";
	logsObject.callback = logs_callback;
	logsObject.parms = postdata;
	logsObject.freq = <?=$nentriesinterval?>/5;

	// Register the AJAX object
	register_ajax(logsObject);

	// ---------------------------------------------------------------------------------------------------
});
//]]>
</script>

<!-- close the body we're wrapped in and add a configuration-panel -->
</div>

<div id="<?=$widget_panel_footer_id?>" class="panel-footer collapse">

<?php
$pconfig['nentries'] = isset($user_settings['widgets'][$widgetkey]['filterlogentries']) ? $user_settings['widgets'][$widgetkey]['filterlogentries'] : '';
$pconfig['nentriesinterval'] = isset($user_settings['widgets'][$widgetkey]['filterlogentriesinterval']) ? $user_settings['widgets'][$widgetkey]['filterlogentriesinterval'] : '';
$pconfig['nseconds'] = isset($user_settings['widgets'][$widgetkey]['logseconds']) ? $user_settings['widgets'][$widgetkey]['logseconds'] : '';
$pconfig['ip2location_service_url'] = isset($user_settings['widgets'][$widgetkey]['ip2location_service_url']) ? $user_settings['widgets'][$widgetkey]['ip2location_service_url'] : '';
$pconfig['ip2location_results_url'] = isset($user_settings['widgets'][$widgetkey]['ip2location_results_url']) ? $user_settings['widgets'][$widgetkey]['ip2location_results_url'] : '';

?>
	<form action="/widgets/widgets/ip2location.widget.php" method="post"
		class="form-horizontal">
		<input type="hidden" name="widgetkey" value="<?=htmlspecialchars($widgetkey); ?>">
		<?=gen_customwidgettitle_div($widgetconfig['title']); ?>

		<div class="form-group">
			<label for="ip2location_service_url" class="col-sm-4 control-label"><?=gettext('IP2Location Cache API')?></label>
			<div class="col-sm-4">
				<input type="text" name="ip2location_service_url" id="ip2location_service_url" value="<?=$pconfig['ip2location_service_url']?>" placeholder="http://192.168.1.51:9999/filterlog" class="form-control" />
			</div>
		</div>

		<div class="form-group">
			<label for="ip2location_results_url" class="col-sm-4 control-label"><?=gettext('IP2Location results API ')?></label>
			<div class="col-sm-4">
				<input type="text" name="ip2location_results_url" id="ip2location_results_url" value="<?=$pconfig['ip2location_results_url']?>" placeholder="http://192.168.1.51:9999/ip2geomap" class="form-control" />
			</div>
		</div>

		<div class="form-group">
			<label for="filterlogentries" class="col-sm-4 control-label"><?=gettext('Timeframe to display')?></label>
			<div class="col-sm-6">
				<input type="number" name="filterlogentries" id="filterlogentries" value="<?=$pconfig['nentries']?>" placeholder="500"
					min="30" max="5000" class="form-control" />
			</div>
		</div>

		<div class="form-group">
			<label class="col-sm-4 control-label"><?=gettext('Filter actions')?></label>
			<div class="col-sm-6 checkbox">
			<?php $include_acts = explode(" ", strtolower($nentriesacts)); ?>
			<label><input name="actpass" type="checkbox" value="Pass"
				<?=(in_array('pass', $include_acts) ? 'checked':'')?> />
				<?=gettext('Pass')?>
			</label>
			<label><input name="actblock" type="checkbox" value="Block"
				<?=(in_array('block', $include_acts) ? 'checked':'')?> />
				<?=gettext('Block')?>
			</label>
			<label><input name="actreject" type="checkbox" value="Reject"
				<?=(in_array('reject', $include_acts) ? 'checked':'')?> />
				<?=gettext('Reject')?>
			</label>
			</div>
		</div>

		<div class="form-group">
			<label for="filterlogentriesinterfaces" class="col-sm-4 control-label">
				<?=gettext('Filter interface')?>
			</label>
			<div class="col-sm-6 checkbox">
				<select name="filterlogentriesinterfaces" id="filterlogentriesinterfaces" class="form-control">
			<?php foreach (array("All" => "ALL") + $iface_descr_arr as $iface => $ifacename):?>
				<option value="<?=$iface?>"
						<?=($nentriesinterfaces==$iface?'selected':'')?>><?=htmlspecialchars($ifacename)?></option>
			<?php endforeach;?>
				</select>
			</div>
		</div>

		<div class="form-group">
			<label for="filterlogentriesinterval" class="col-sm-4 control-label"><?=gettext('Update interval')?></label>
			<div class="col-sm-4">
				<input type="number" name="filterlogentriesinterval" id="filterlogentriesinterval" value="<?=$pconfig['nentriesinterval']?>" placeholder="5"
					min="1" class="form-control" />
			</div>
			<?=gettext('Seconds');?>
		</div>

		<div class="form-group">
			<label for="logseconds" class="col-sm-4 control-label"><?=gettext('Display the last (seconds)')?></label>
			<div class="col-sm-4">
				<input type="number" name="logseconds" id="logseconds" value="<?=$pconfig['nseconds']?>" placeholder="60"
				       min="1" class="form-control" />
			</div>
			<?=gettext('Seconds');?>
		</div>

		<div class="form-group">
			<div class="col-sm-offset-4 col-sm-6">
				<button type="submit" class="btn btn-primary"><i class="fa fa-save icon-embed-btn"></i><?=gettext('Save')?></button>
			</div>
		</div>
	</form>

<script type="text/javascript">
//<![CDATA[
if (typeof getURL == 'undefined') {
	getURL = function(url, callback) {
		if (!url)
			throw 'No URL for getURL';
		try {
			if (typeof callback.operationComplete == 'function')
				callback = callback.operationComplete;
		} catch (e) {}
			if (typeof callback != 'function')
				throw 'No callback function for getURL';
		var http_request = null;
		if (typeof XMLHttpRequest != 'undefined') {
			http_request = new XMLHttpRequest();
		}
		else if (typeof ActiveXObject != 'undefined') {
			try {
				http_request = new ActiveXObject('Msxml2.XMLHTTP');
			} catch (e) {
				try {
					http_request = new ActiveXObject('Microsoft.XMLHTTP');
				} catch (e) {}
			}
		}
		if (!http_request)
			throw 'Both getURL and XMLHttpRequest are undefined';
		http_request.onreadystatechange = function() {
			if (http_request.readyState == 4) {
				callback( { success : true,
				  content : http_request.responseText,
				  contentType : http_request.getResponseHeader("Content-Type") } );
			}
		};
		http_request.open('GET', url, true);
		http_request.send(null);
	};
}

function outputrule(req) {
	alert(req.content);
}
//]]>
</script>


