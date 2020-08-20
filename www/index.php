<?php
/*
    The gocms PHP index file
*/

$server_script_path = "/usr/share/nginx/apps/gocms";

if (basename(__DIR__) != "www") {
    $cmd = sprintf("%s/gocms %s %s", $server_script_path, json_encode((object)$_GET), json_encode((object)$_POST));
    exec($cmd, $output, $code);
    if ($code == 0) {
        $response = json_decode($output[0]);
        foreach ($response->Headers as $header)
            if ($header == 404)
                header($_SERVER["SERVER_PROTOCOL"] . " 404 Not Found");
            else
                header($header);
        echo $response->Content;
    }
    else
        echo "Error code: " . $code;
}
