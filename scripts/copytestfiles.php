<?php
// Copy the test pages to the server replacing any existing pages

define("USING_TOOL", true);

include_once(__DIR__ . "/config.php");

$server_script_path = $config->get("server_script_path");
$domain = $config->get("domain");
$web_root = $config->get("web_root");

$dest = "$server_script_path/$domain";

// Remove existing page files
array_map('unlink', glob("$dest/pages/*.html"));

// Remove existing template files
array_map('unlink', glob("$dest/templates/*.html"));

$src = dirname(__DIR__) . "/pkg/files/test";

// Copy new files
$files =  glob("$src/pages/*.html");
array_walk($files, 'copy_files', "$dest/pages/");
$files =  glob("$src/templates/*.html");
array_walk($files, 'copy_files', "$dest/templates/");

copy_file($src, $dest, "/data/pages.json");

// Set up web folders
add_folder("$web_root/css");
add_folder("$web_root/js");

// Copy new files
$files =  glob("$src/styles/*.css");
array_walk($files, 'copy_files', "$web_root/css/");
$files =  glob("$src/scripts/*.js");
array_walk($files, 'copy_files', "$web_root/js/");

function add_folder($fn) {
    if (!file_exists($fn))
        mkdir($fn);
}

function copy_file($src, $dest, $jsonPath) {
    $src .= $jsonPath;
    $dest .= $jsonPath;
    if (file_exists($src))
        copy($src, $dest);
    else
        echo "Missing file: $src\n\n";
}

function copy_files($src, $key, $dest) {
    copy($src, $dest . basename($src));
}

echo "Done\n";
