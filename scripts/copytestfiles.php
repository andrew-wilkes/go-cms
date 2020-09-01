<?php
// Copy the test pages to the server replacing any existing pages

define("USING_TOOL", true);

include_once(__DIR__ . "/config.php");

$server_script_path = $config->get("server_script_path");
$domain = $config->get("domain");

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

$j = "/data/pages.json";
$src .= $j;
$dest .= $j;
if (file_exists($src))
    copy($src, $dest);
else
    echo "Missing file: $src\n\n";

function copy_files($src, $key, $dest) {
    copy($src, $dest . basename($src));
}
echo "Done\n";
