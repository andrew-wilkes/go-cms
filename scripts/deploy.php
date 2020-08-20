<?php

$web_root = "/usr/share/nginx/www/gocms.com";

$dir = dirname(__DIR__);

$index = $dir . "/www/index.php";

include_once("$index");

if (!file_exists($server_script_path)) die("$server_script_path path invalid!");

$gocms = $dir . "/build/gocms";

if (!file_exists($gocms)) die("gocms build file is missing!\n");

// Copy gocms program to server
$cmd = "cp $gocms $server_script_path";
echo("$cmd\n");
exec($cmd);

// Copy index.php to server
$cmd = "cp $index $web_root";
echo("$cmd\n");
exec($cmd);
