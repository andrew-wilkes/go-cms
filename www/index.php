<?php
/*
    The gocms PHP index file
*/

$cmd = sprintf("server_script_path/gocms %s", json_encode($_SERVER));
exec($cmd, $output, $code);
if ($code == 0) {
    $response = json_decode($output[0]);
    foreach ($response->Headers as $header)
        header($header);
    echo $response->Content;
}
else
    echo "Error code: " . $code;
