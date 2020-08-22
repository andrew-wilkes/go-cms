<?php
define("USING_TOOL", true);

include_once(__DIR__ . "/config.php");

$web_root = $config->get("web_root");
$server_script_path = $config->get("server_script_path");

$dir = dirname(__DIR__);

$index = $dir . "/www/index.php";

if (!file_exists($server_script_path)) die("$server_script_path path invalid!" . PHP_EOL);

$gocms = $dir . "/build/gocms";

if (!file_exists($gocms)) die("The gocms app file is missing!" . PHP_EOL);
// Later take into account other platforms such as where we would expect gocms.exe

// Copy gocms program to server
exec("cp $gocms $server_script_path");

// Inject server_script_path into index.php and copy to server
file_put_contents($web_root . "/index.php", str_replace('server_script_path', $server_script_path, file_get_contents($index)));
