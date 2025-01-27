<?php

namespace App\Middleware;

use Phalcon\Mvc\Micro\MiddlewareInterface;
use Phalcon\Mvc\Micro;

class AuthMiddleware implements MiddlewareInterface
{
    public function call(Micro $app): bool
    {

        $clientToken = $app->request->getHeader('Authorization');

        if (!$clientToken) {
            $app->response
                ->setJsonContent([
                    'status'  => 'error',
                    'message' => 'Missing client_token',
                ])
                ->setStatusCode(401)
                ->send();
            return false;
        }

        $tokenData = $app->getDI()->getShared('controllerBase')->validateClientToken($clientToken);

        if (!$tokenData) {
            $app->response
                ->setJsonContent([
                    'status'  => 'error',
                    'message' => 'Invalid client_token',
                ])
                ->setStatusCode(401)
                ->send();
            return false;
        }

        $app->request->set('uid', $tokenData['uid']);
        $app->request->set('app_id', $tokenData['app_id']);

        error_log("hollaaaaaa sad asd as das d");

        return true;
    }
}
