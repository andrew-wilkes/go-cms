<?php
$_SERVER['RAW_DATA'] = file_get_contents("php://input");
#LOG
$cmd = sprintf("server_script_path/gocms %s", json_encode($_SERVER));
exec($cmd, $output, $code);
if ($code == 0) {
  $response = json_decode($output[0]);
  foreach ($response->Headers as $header)
    header($header);
  echo $response->Content;
} else echo "Error code: " . $code;
#HTML