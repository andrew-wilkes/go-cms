<?php

$file = "gocms.go";

$dir = dirname(__DIR__);

$src = $dir . "/cmd/" . $file;

$dest = $dir . "/build";

$cmd = sprintf("go build -o %s %s", $dest, $src);

echo $cmd . PHP_EOL;

exec($cmd);
