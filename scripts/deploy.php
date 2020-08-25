<?php

// Add command line param of -l or --log to output the test version of index.php
$test_mode = 0 < count(getopt("l", ["log"]));

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

// Formulate the web page index code
$log_code = '$log_file_name = urlencode($_SERVER[\'REQUEST_URI\']);
file_put_contents(sprintf(\'log/args-%s.json\', $log_file_name), json_encode($_SERVER));';
$php_code = file_get_contents($index);
if ($test_mode) {
    $php_code = str_replace("#LOG\n", $log_code, $php_code);
    $php_code = str_replace("#HTML", file_get_contents($dir . "/www/test.html"), $php_code);
} else {
    $php_code = str_replace("\n#LOG", "", $php_code);
    $php_code = str_replace("\n#HTML", "", $php_code);
}

// Inject server_script_path into index.php and copy to server
file_put_contents($web_root . "/index.php", str_replace('server_script_path', $server_script_path, $php_code));
