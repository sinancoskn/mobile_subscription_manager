<?php

$router = $di->getRouter();

// Define your routes here

// $router->handle($_SERVER['REQUEST_URI']);

$router->addGet('/health', [
    'controller' => 'index',
    'action' => 'health'
]);

$router->addPost('/devices/register', [
    'controller' => 'devices',
    'action' => 'register',
]);

$router->addPost('/purchase', [
    'controller' => 'purchases',
    'action' => 'purchase',
]);

$router->addGet('/check-subscription', [
    'controller' => 'purchases',
    'action' => 'checkSubscription'
]);