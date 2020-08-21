<?php
/*
    This is a minimalist implementation of a TOML configuration file parser.
    It extracts a = b data that appears on single lines.

    Use $config->get($token) to obtain a value.
*/

$config = new Config;

class Config {

    private $values;

    function __construct() {
        $cfg_file = dirname(__DIR__) . "/config.toml";
        if (!file_exists($cfg_file)) die("$cfg_file does not exist!\n");
        $this->values = $this->parse_config(file_get_contents($cfg_file));
        defined("USING_TOOL") or print_r($this->values);
    }

    private function parse_config($src) {
        $pairs = [];
        foreach(explode("\n", $src) as $line) {
            $parts = explode("=", $line);
            if (count($parts) == 2) {
                $p1 = trim($parts[0]);
                $p2 = $parts[1];
                $comment_pos = strpos($p2, "#");
                if ($comment_pos > 0) $p2 = substr($p2, 0, $comment_pos);
                $p2 = trim($p2);
                $pairs[$p1] = trim($p2, '"');
            }
        }
        return $pairs;
    }

    public function get($token) {
        if (array_key_exists($token, $this->values))
            return $this->values[$token];
        die("Missing config entry for: $token\n");
    }
}
