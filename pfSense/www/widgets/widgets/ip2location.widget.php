<?php
/*
 *
 */

require_once("guiconfig.inc");
require_once("pfsense-utils.inc");
require_once("functions.inc");
require_once("syslog.inc");

global $ip2l_results, $ip2l_display_status, $ip2l_filterlog_time;

function create_url($hostport, $path) {
    $url = $hostport . $path;
	return $url;
}



function check_api($healthUrl)
{
	$ch = curl_init();
	$optArray = array(
		CURLOPT_URL => $healthUrl,
		CURLOPT_RETURNTRANSFER => true,
		CURLOPT_HTTPGET => 1,
	);
	curl_setopt_array($ch, $optArray);

	if(!$result = curl_exec($ch)) {
		// $error = curl_errno($req);
		// if ($error == CURLE_SSL_PEER_CERTIFICATE || $error == CURLE_SSL_CACERT || $error == 77) {
		// 	curl_setopt($req, CURLOPT_CAINFO, __DIR__ . '/cert-bundle.crt');
		// 	$result = curl_exec($req);
		// }
		trigger_error(curl_error($ch));
	}

	curl_close($ch);

	log_error("IP2Location API health check: " . $result);
    
	$ip2l_display_status = "IP2Location API health check: " . $result;
	if ($result == "Service is available.") {
		return "true";
	} else {
		return "false";
	}
	return "false";
}




function extract_ip_entries($logarr, $seconds, $filterlog_time)
{
	$ip2l_display_status = sprintf("Extracting the last %s seconds of log entries.<br/>", $seconds);
	$timeCap = $filterlog_time - $seconds;
	$loggedIps = [];
	$count = 0;
	foreach ($logarr as $entry) {
		$datetime = $entry['time'];
		if ($datetime < $timeCap) {
			$count++;
			continue;
		}
		$srcIP = $entry['srcip'];
		$dstIP = $entry['dstip'];
		if ($entry['direction'] == "in") {
			$safeSrcKey = sha1($srcIP);
			$loggedIps[$safeSrcKey] = $entry;
			$count++;
		} else {
			$safeDstKey = sha1($dstIP);
			$loggedIps[$safeDstKey] = $entry;
			$count++;
		}
	}
	$$ip2l_display_status = sprintf("Extracted locations %s of a maximum of %s IPs, spanning the past %s seconds (%s to %s).\n", count($loggedIps), $count, $seconds, $datetime, date($timeCap));
	
	return $loggedIps;
}



function truncate($string, $length) {
    return (strlen($string) > $length) ? substr($string, 0, $length) : $string;
}



function send_filterlog($ip_log, $url, $key) {
	$ip2l_display_status = sprintf("Sending %s IPs to IP2Location API.\n", count($ip_log));
	if (!$url) {
		log_error("IP2Location API send filter log error missing URL.");
		return false;
	}

	$data = json_encode($ip_log);
	$authorization = "Authorization: Bearer ". $key;
	$ch_send = curl_init();
	$send_optArray = array(
		CURLOPT_URL => $url,
		CURLOPT_RETURNTRANSFER => true,
		CURLOPT_HTTPHEADER => ['Content-Type: application/json', $authorization ],
		CURLOPT_POST => 1,
		CURLOPT_POSTFIELDS => $data
	);
	curl_setopt_array($ch_send, $send_optArray);
	
	$result = curl_exec($ch_send);
	if (!$result) {
		trigger_error(curl_error($ch_send));
	}
	$return_result = truncate($result, 13);
	curl_close($ch_send);
	log_error("IP2Location API send filter logs error: " . $result);
	return $return_result;
}



function get_results($resultsid, $url, $key)
{
	$ip2l_display_status = sprintf("Fething resutls from API for %s.\n", $resultsid);

	$dataArray = ['id' => $resultsid];
	$data = http_build_query($dataArray);
	$getUrl = $url."?".$data;
	$ip2l_display_status = sprintf('URL: ' . $getUrl . '<br/>');
	
	$ch = curl_init();
	$authorization = "Authorization: Bearer ". $key;
	$optArray = array(	
		CURLOPT_SSL_VERIFYPEER => false,
		CURLOPT_URL => $getUrl,
		CURLOPT_RETURNTRANSFER => true,
		CURLOPT_HTTPGET => 1,
		CURLOPT_HTTPHEADER => ['Content-Type: application/json', $authorization ],
	);
	curl_setopt_array($ch, $optArray);

	$ip2result = curl_exec($ch);
	$errors = curl_error($ch);
	$response = curl_getinfo($ch, CURLINFO_HTTP_CODE);

	if (!$ip2result) {
		trigger_error(curl_error($ch));
		$ip2l_display_status = sprintf("An error occurred while retrieving the IP2Location results for the results ID: " . $resultsid . $errors . "<br/>");
	}
	printf($ip2result);
}


if ($_REQUEST['widgetkey'] && !$_REQUEST['ajax']) {
	set_customwidgettitle($user_settings);

	if (is_numeric($_POST['ip2l_max_entries'])) {
		$user_settings['widgets'][$_POST['widgetkey']]['ip2l_max_entries'] = $_POST['ip2l_max_entries'];
	} else {
		unset($user_settings['widgets'][$_POST['widgetkey']]['ip2l_max_entries']);
	}

	if (is_numeric($_POST['ip2l_log_interval'])) {
		$user_settings['widgets'][$_POST['widgetkey']]['ip2l_log_interval'] = $_POST['ip2l_log_interval'];
	} else {
		unset($user_settings['widgets'][$_POST['widgetkey']]['ip2l_log_interval']);
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
		$user_settings['widgets'][$_POST['widgetkey']]['ip2l_log_acts'] = implode(" ", $acts);
	} else {
		unset($user_settings['widgets'][$_POST['widgetkey']]['ip2l_log_acts']);
	}
	unset($accts);

	if (($_POST['ip2l_log_interfaces']) and ($_POST['ip2l_log_interfaces'] != "All")) {
		$user_settings['widgets'][$_POST['widgetkey']]['ip2l_log_interfaces'] = trim($_POST['ip2l_log_interfaces']);
	} else {
		unset($user_settings['widgets'][$_POST['widgetkey']]['ip2l_log_interfaces']);
	}
	
	if (is_numeric($_POST['ip2l_log_seconds'])) {
		$user_settings['widgets'][$_POST['widgetkey']]['ip2l_log_seconds'] = $_POST['ip2l_log_seconds'];
	} else {
		unset($user_settings['widgets'][$_POST['widgetkey']]['ip2l_log_seconds']);
	}
	
	if (is_string($_POST['ip2l_api_hostport'])) {
		$user_settings['widgets'][$_POST['widgetkey']]['ip2l_api_hostport'] = $_POST['ip2l_api_hostport'];
	} else {
		unset($user_settings['widgets'][$_POST['widgetkey']]['ip2l_api_hostport']);
	}
	
	if (is_string($_POST['ip2l_submit_api'])) {
		$user_settings['widgets'][$_POST['widgetkey']]['ip2l_submit_api'] = $_POST['ip2l_submit_api'];
	} else {
		unset($user_settings['widgets'][$_POST['widgetkey']]['ip2l_submit_api']);
	}

	if (is_string($_POST['ip2l_results_api'])) {
		$user_settings['widgets'][$_POST['widgetkey']]['ip2l_results_api'] = $_POST['ip2l_results_api'];
	} else {
		unset($user_settings['widgets'][$_POST['widgetkey']]['ip2l_results_api']);
	}

	if (is_string($_POST['ip2l_health'])) {
		$user_settings['widgets'][$_POST['widgetkey']]['ip2l_health'] = $_POST['ip2l_health'];
	} else {
		unset($user_settings['widgets'][$_POST['widgetkey']]['ip2l_health']);
	}

	if (is_string($_POST['ip2l_details_page'])) {
		$user_settings['widgets'][$_POST['widgetkey']]['ip2l_details_page'] = $_POST['ip2l_details_page'];
	} else {
		unset($user_settings['widgets'][$_POST['widgetkey']]['ip2l_details_page']);
	}	

	if (is_string($_POST['ip2l_token'])) {
		$user_settings['widgets'][$_POST['widgetkey']]['ip2l_token'] = $_POST['ip2l_token'];
	} else {
		unset($user_settings['widgets'][$_POST['widgetkey']]['ip2l_token']);
	}
	
	save_widget_settings($_SESSION['Username'], $user_settings["widgets"], gettext("Saved IP2Location configuration via Dashboard."));
	Header("Location: /");
	exit(0);
}


// When this widget is included in the dashboard, $widgetkey is already defined before the widget is included.
// When the ajax call is made to refresh the firewall log table, 'widgetkey' comes in $_REQUEST.
if ($_REQUEST['widgetkey']) {
	$widgetkey = $_REQUEST['widgetkey'];
}

$widgetkey_nodash = str_replace("-", "", $widgetkey);
$iface_descr_arr = get_configured_interface_with_descr();
$ip2l_max_entries = isset($user_settings['widgets'][$widgetkey]['ip2l_max_entries']) ? $user_settings['widgets'][$widgetkey]['ip2l_max_entries'] : 50;
$ip2l_log_interval = isset($user_settings['widgets'][$widgetkey]['ip2l_log_interval']) ? $user_settings['widgets'][$widgetkey]['ip2l_log_interval'] : 10;
$ip2l_log_acts = isset($user_settings['widgets'][$widgetkey]['ip2l_log_acts']) ? $user_settings['widgets'][$widgetkey]['ip2l_log_acts'] : 'All';
$ip2l_log_interfaces = isset($user_settings['widgets'][$widgetkey]['ip2l_log_interfaces']) ? $user_settings['widgets'][$widgetkey]['ip2l_log_interfaces'] : 'All';
$ip2l_fields_array = array(
	"act" => $ip2l_log_acts,
	"interface" => isset($iface_descr_arr[$ip2l_log_interfaces]) ? $iface_descr_arr[$ip2l_log_interfaces] : $ip2l_log_interfaces
);

$ip2l_log_seconds = isset($user_settings['widgets'][$widgetkey]['ip2l_log_seconds']) ? $user_settings['widgets'][$widgetkey]['ip2l_log_seconds'] : 30;
$ip2l_api_hostport = isset($user_settings['widgets'][$widgetkey]['ip2l_api_hostport']) ? $user_settings['widgets'][$widgetkey]['ip2l_api_hostport'] : "http://localhost:9999";
$ip2l_submit_api = isset($user_settings['widgets'][$widgetkey]['ip2l_submit_api']) ? $user_settings['widgets'][$widgetkey]['ip2l_submit_api'] : "/api/filterlog";
$ip2l_results_api = isset($user_settings['widgets'][$widgetkey]['ip2l_results_api']) ? $user_settings['widgets'][$widgetkey]['ip2l_results_api'] : "/api/results";
$ip2l_health = isset($user_settings['widgets'][$widgetkey]['ip2l_health']) ? $user_settings['widgets'][$widgetkey]['ip2l_health'] : "/health";
$ip2l_details_page = isset($user_settings['widgets'][$widgetkey]['ip2l_details_page']) ? $user_settings['widgets'][$widgetkey]['ip2l_details_page'] : "/index.html";
$ip2l_token = isset($user_settings['widgets'][$widgetkey]['ip2l_token']) ? $user_settings['widgets'][$widgetkey]['ip2l_token'] : 'valid-key';

$ip2l_filterlog_time = time();

$filter_logfile = "{$g['varlog_path']}/filter.log";
$widgetkey_nodash = str_replace("-", "", $widgetkey);
$health = check_api(create_url($ip2l_api_hostport, $ip2l_health));
if ($health == "false") {
	$ip2l_display_status = "IP2Location API is not available.";
} else {
	$ip2l_display_status = "IP2Location API is available.";

	$filter_log = conv_log_filter($filter_logfile, $ip2l_max_entries, 5000, $ip2l_fields_array);
	$ip2l_display_status = sprintf("Filter log entries: %d\n", $ip2l_max_entries);

	$ip_log_items = extract_ip_entries($filter_log, $ip2l_log_seconds, $ip2l_filterlog_time);
	$ip2l_display_status = sprintf(" <b>%s</b>&nbsp;&nbsp; Displaying location of %d IP addresses filtered in the last %d seconds.\n", date("H:i:s", $ip2l_filterlog_time), count($ip_log_items), $ip2l_log_seconds);

	$ip2l_submit_url = create_url($ip2l_api_hostport, $ip2l_submit_api);
	$ip2l_results_url = create_url($ip2l_api_hostport, $ip2l_results_api);
	$ip2l_results = send_filterlog($ip_log_items, $ip2l_submit_url, $ip2l_token);

	$resultsid = $ip2l_results;
}

if (!$_REQUEST['ajax']) {
?>
<script type="text/javascript">
//<![CDATA[
	var ip2lWidgetLastRefresh<?=htmlspecialchars($widgetkey_nodash)?> = <?=time()?>;
	console.log("ip2lWidgetLastRefresh<?=htmlspecialchars($widgetkey_nodash)?> = " + ip2lWidgetLastRefresh<?=htmlspecialchars($widgetkey_nodash)?>);
//]]>
</script>
<?php 
} 

if ($_REQUEST['ajax'] && $_REQUEST['widgetkey']) {
	$coords_x = $_REQUEST['coords_x'];
	$coords_y = $_REQUEST['coords_y'];
	$zoom = $_REQUEST['zoom'];
}
?>

<?php 
if ($_REQUEST['ajax'] && $_REQUEST['widgetkey'] && $_REQUEST['resultsid']) {
	$ret_resultsid = $_REQUEST['resultsid'];
	get_results($ret_resultsid, $ip2l_results_url, $ip2l_token);
	exit(0);
}
?>

<!-- This is the body of the widget and will be AJAX-refreshed -->
<link rel="stylesheet" href="/widgets/widgets/ip2location.widget.css"/>
<script src="/vendor/leaflet/leaflet.js?v=<?=filemtime('/usr/local/www/vendor/leaflet/leaflet.js')?>"></script>
<script src="/widgets/javascript/ip2location.js?v=<?=filemtime('/usr/local/www/widgets/javascript/ip2location.js')?>"></script>

<div id="<?=$widgetkey?>-map">
	<div id="leaflet" style="height: 320px;">
	</div>
	<div class="subpanel-body">
		<span class="ip2l_status"><?=gettext($ip2l_display_status); ?></span>
	</div>
</div>

<script>
//<![CDATA[
	var coords_x = localStorage.getItem("coords_x") ?? <?=isset($coords_x) ? htmlspecialchars($coords_x) : 51.505?>;
	var coords_y = localStorage.getItem("coords_y") ?? <?=isset($coords_y) ? htmlspecialchars($coords_y) : -0.09?>;
	var zoom = localStorage.getItem("zoom") ?? <?=isset($zoom) ? htmlspecialchars($zoom) : 13?>;

	if (coords_x != null && coords_y != null && zoom != null) {
		map_coords = [coords_x, coords_y];
		map_zoom = zoom;
	}
	
	var map = L.map('leaflet').setView(map_coords, map_zoom);

	L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
    	maxZoom: 19,
    	attribution: '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>'
	}).addTo(map);

	function onMapMove(e) {
		localStorage.setItem("coords_x", map.getCenter().lat);
		localStorage.setItem("coords_y", map.getCenter().lng);
		localStorage.setItem("zoom", map.getZoom());
		console.log("Map moved to " + map.getCenter().lat + ", " + map.getCenter().lng + " at zoom level " + map.getZoom() + ". New view saved to localStorage.");
	}

	map.on('zoomend', onMapMove);
	map.on('moveend', onMapMove);
//]]>
</script>

<script type="text/javascript">
//<![CDATA[
	// POST data to send via AJAX
	events.push(function(){
		// --------------------- Centralized widget refresh system ------------------------------
		// Callback function called by refresh system when data is retrieved
		function ip2l_callback(s) {
			console.log("Refreshing IP2Location widget (<?=htmlspecialchars($widgetkey)?>)");
			$(<?=json_encode('#widget-' . $widgetkey . '_panel-body')?>).html(s);
		}
	
		var postdata = {
			ajax: "ajax",
			widgetkey: <?=json_encode($widgetkey)?>,
			lastsawtime: ip2lWidgetLastRefresh<?=htmlspecialchars($widgetkey_nodash)?>,
		};

		// Create an object defining the widget refresh AJAX call
		var ip2lObject = new Object();
		ip2lObject.name = "IP2Location Fireall Logs";
		ip2lObject.url = "/widgets/widgets/ip2location.widget.php";
		ip2lObject.callback = ip2l_callback;
		ip2lObject.parms = postdata;
		ip2lObject.freq = <?=$ip2l_log_interval?>/5;
		
		// Register the AJAX object
		register_ajax(ip2lObject);
		console.log("Registered IP2Location widget (<?=htmlspecialchars($widgetkey)?>) with freq = " + ip2lObject.freq + " seconds, and lastsawtime = " + ip2lWidgetLastRefresh<?=htmlspecialchars($widgetkey_nodash)?>);
		// ---------------------------------------------------------------------------------------
	});

	fetchMapData(<?=json_encode($widgetkey)?>, <?=json_encode($resultsid)?>);	
//]]>
</script>
<?php
/* for AJAX response, we only need the panel-body */
if ($_REQUEST['ajax']) {
	
	exit;
}
?>
<!-- close the body we're wrapped in and add a configuration-panel -->
</div>

<!----------- Configuration panel ----------->
<div id="<?=$widget_panel_footer_id?>" class="panel-footer collapse">
<?php
$pconfig['ip2l_max_entries'] = isset($user_settings['widgets'][$widgetkey]['ip2l_max_entries']) ? $user_settings['widgets'][$widgetkey]['ip2l_max_entries'] : '';
$pconfig['ip2l_log_interval'] = isset($user_settings['widgets'][$widgetkey]['ip2l_log_interval']) ? $user_settings['widgets'][$widgetkey]['ip2l_log_interval'] : '';
$pconfig['ip2l_log_acts'] = isset($user_settings['widgets'][$widgetkey]['ip2l_log_acts']) ? $user_settings['widgets'][$widgetkey]['ip2l_log_acts'] : '';
$pconfig['ip2l_log_interfaces'] = isset($user_settings['widgets'][$widgetkey]['ip2l_log_interfaces']) ? $user_settings['widgets'][$widgetkey]['ip2l_log_interfaces'] : '';
$pconfig['ip2l_log_seconds'] = isset($user_settings['widgets'][$widgetkey]['ip2l_log_seconds']) ? $user_settings['widgets'][$widgetkey]['ip2l_log_seconds'] : '';
$pconfig['ip2l_api_hostport'] = isset($user_settings['widgets'][$widgetkey]['ip2l_api_hostport']) ? $user_settings['widgets'][$widgetkey]['ip2l_api_hostport'] : '';
$pconfig['ip2l_submit_api'] = isset($user_settings['widgets'][$widgetkey]['ip2l_submit_api']) ? $user_settings['widgets'][$widgetkey]['ip2l_submit_api'] : '';
$pconfig['ip2l_results_api'] = isset($user_settings['widgets'][$widgetkey]['ip2l_results_api']) ? $user_settings['widgets'][$widgetkey]['ip2l_results_api'] : '';
$pconfig['ip2l_health'] = isset($user_settings['widgets'][$widgetkey]['ip2l_health']) ? $user_settings['widgets'][$widgetkey]['ip2l_health'] : '';
$pconfig['ip2l_details_page'] = isset($user_settings['widgets'][$widgetkey]['ip2l_details_page']) ? $user_settings['widgets'][$widgetkey]['ip2l_details_page'] : '';
$pconfig['ip2l_token'] = isset($user_settings['widgets'][$widgetkey]['ip2l_token']) ? $user_settings['widgets'][$widgetkey]['ip2l_token'] : '';

?>

<form action="/widgets/widgets/ip2location.widget.php" method="post" class="form-horizontal">
	<input type="hidden" name="widgetkey" value="<?=htmlspecialchars($widgetkey); ?>">
	<?=gen_customwidgettitle_div($widgetconfig['title']); ?>

	<div class="form-group" id="ip2l">
		<label for="ip2l_api_hostport" class="col-sm-4 control-label"><?=gettext('IP2Location daemon http[s]://<hostname>:<port> ')?></label>
		<div class="col-sm-4">
			<input type="text" name="ip2l_api_hostport" id="ip2l_api_hostport" value="<?=$pconfig['ip2l_api_hostport']?>" placeholder="http://localhost:9999" class="form-control" />
		</div>
	</div>
	<div class="form-group" id="ip2l">
		<label for="ip2l_submit_api" class="col-sm-4 control-label"><?=gettext('Submit filter logs API [/filterlog] ')?></label>
		<div class="col-sm-4">
			<input type="text" name="ip2l_submit_api" id="ip2l_submit_api" value="<?=$pconfig['ip2l_submit_api']?>" placeholder="/filterlog" class="form-control" />
		</div>
	</div>
	<div class="form-group" id="ip2l">
		<label for="ip2l_results_api" class="col-sm-4 control-label"><?=gettext('Retrieve resutls API [/results] ')?></label>
		<div class="col-sm-4">
			<input type="text" name="ip2l_results_api" id="ip2l_results_api" value="<?=$pconfig['ip2l_results_api']?>" placeholder="/results" class="form-control" />
		</div>
	</div>
	<div class="form-group" id="ip2l">
		<label for="ip2l_health" class="col-sm-4 control-label"><?=gettext('Service health check endpoint [/health] ')?></label>
		<div class="col-sm-4">
			<input type="text" name="ip2l_health" id="ip2l_health" value="<?=$pconfig['ip2l_health']?>" placeholder="/health" class="form-control" />
		</div>
	</div>
	<div class="form-group" id="ip2l">
		<label for="ip2l_details_page" class="col-sm-4 control-label"><?=gettext('Details page [/index.html] ')?></label>
		<div class="col-sm-4">
			<input type="text" name="ip2l_details_page" id="ip2l_details_page" value="<?=$pconfig['ip2l_details_page']?>" placeholder="/index.html" class="form-control" />
		</div>
	</div>		
	<div class="form-group" id="ip2l">
		<label for="ip2l_token" class="col-sm-4 control-label"><?=gettext('Backend service token ')?></label>
		<div class="col-sm-4">
			<input type="text" name="ip2l_token" id="ip2l_token" value="<?=$pconfig['ip2l_token']?>" placeholder="[Bearer token]" class="form-control" />
		</div>
	</div>

	<div class="form-group">
		<label for="ip2l_max_entries" class="col-sm-4 control-label"><?=gettext('Number of entries')?></label>
		<div class="col-sm-6">
			<input type="number" name="ip2l_max_entries" id="ip2l_max_entries" value="<?=$pconfig['ip2l_max_entries']?>" placeholder="50" min="1" max="5000" class="form-control" />
		</div>
	</div>

	<div class="form-group">
		<label for="ip2l_log_interval" class="col-sm-4 control-label"><?=gettext('Update interval')?></label>
		<div class="col-sm-4">
			<input type="number" name="ip2l_log_interval" id="ip2l_log_interval" value="<?=$pconfig['ip2l_log_interval']?>" placeholder="60" min="1" class="form-control" />
		</div>
		<?=gettext('Seconds');?>
	</div>

	<div class="form-group" id="ip2l">
		<label class="col-sm-4 control-label"><?=gettext('Filter actions')?></label>
		<div class="col-sm-6 checkbox">
			<?php $include_acts = explode(" ", strtolower($ip2l_log_acts)); ?>
			<label><input name="actpass" type="checkbox" value="Pass" <?=(in_array('pass', $include_acts) ? 'checked':'')?> /><?=gettext('Pass')?></label>
			<label><input name="actblock" type="checkbox" value="Block" <?=(in_array('block', $include_acts) ? 'checked':'')?> /><?=gettext('Block')?></label>
			<label><input name="actreject" type="checkbox" value="Reject" <?=(in_array('reject', $include_acts) ? 'checked':'')?> /><?=gettext('Reject')?></label>
		</div>
	</div>

	<div class="form-group" id="ip2l">
		<label for="ip2l_log_interfaces" class="col-sm-4 control-label"><?=gettext('Filter interface')?></label>
		<div class="col-sm-6 checkbox">
			<select name="ip2l_log_interfaces" id="ip2l_log_interfaces" class="form-control">
				<?php foreach (array("All" => "ALL") + $iface_descr_arr as $iface => $ifacename):?>
					<option value="<?=$iface?>"<?=($ip2l_log_interfaces==$iface?'selected':'')?>><?=htmlspecialchars($ifacename)?></option>
				<?php endforeach;?>
			</select>
		</div>
	</div>

	<div class="form-group" id="ip2l">
		<label for="ip2l_log_seconds" class="col-sm-4 control-label"><?=gettext('Display the last (seconds)')?></label>
		<div class="col-sm-4">
			<input type="number" name="ip2l_log_seconds" id="ip2l_log_seconds" value="<?=$pconfig['ip2l_log_seconds']?>" placeholder="60" min="1" class="form-control" />
		</div>
		<?=gettext('Seconds');?>
	</div>

	<div class="form-group">
		<div class="col-sm-offset-4 col-sm-6">
			<button type="submit" class="btn btn-primary"><i class="fa fa-save icon-embed-btn"></i><?=gettext('Save')?></button>
			<button class="btn" onclick="resetMap()">Reset map</button>
		</div>
	</div>
</form>

<script type="text/javascript">
//<![CDATA[

console.log("Loading IP2Location widget (<?=htmlspecialchars($widgetkey)?>)");

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

console.log("Loaded IP2Location widget (<?=htmlspecialchars($widgetkey)?>)");
//]]>
</script>
