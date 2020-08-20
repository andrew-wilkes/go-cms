<?php

$web_root = "/usr/share/nginx/www/gocms.com";
$server_script_path = "/usr/share/nginx/apps/gocms";

$dir = dirname(__DIR__);

$index = $dir . "/www/index.php";

if (!file_exists($server_script_path)) die("$server_script_path path invalid!");

$gocms = $dir . "/build/gocms";

if (!file_exists($gocms)) die("gocms build file is missing!\n");

// Copy gocms program to server
exec("cp $gocms $server_script_path");

// Inject server_script_path into index.php and copy to server
file_put_contents($web_root . "/index.php", str_replace('server_script_path', $server_script_path, file_get_contents($index)));
