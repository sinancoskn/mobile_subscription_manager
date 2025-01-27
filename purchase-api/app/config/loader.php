<?php

use Dotenv\Dotenv;

require_once '../vendor/autoload.php';
$loader = new \Phalcon\Autoload\Loader();

$baseDir = realpath(__DIR__ . '/../../'); // Adjust based on your project structure

$dotenv = Dotenv::createImmutable($baseDir);
$dotenv->load();

/**
 * We're a registering a set of directories taken from the configuration file
 */
$loader->setDirectories(
    [
        $config->application->controllersDir,
        $config->application->modelsDir
    ]
)->register();
