<?php
/*
    The gocms PHP index file
*/

exec(sprintf("/usr/share/nginx/apps/gocms/gocms %s %s", json_encode((object)$_GET), json_encode((object)$_POST)), $output, $code);

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
