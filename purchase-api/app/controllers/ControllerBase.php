<?php
declare(strict_types=1);

use Phalcon\Mvc\Controller;
use Phalcon\Http\Response;
use Phalcon\Mvc\Dispatcher;

class ControllerBase extends Controller
{
    private $securedRoutes = [
        ['controller' => 'purchases', 'action' => 'purchase'],
    ];

    public function beforeExecuteRoute(Dispatcher $dispatcher)
    {
        $clientToken = $this->request->getHeader('Authorization');

        if ($clientToken) {
            $tokenData = $this->validateClientToken($clientToken);

            if ($tokenData) {
                $dispatcher->setParam('uid', $tokenData["uid"]);
                $dispatcher->setParam('app_id', $tokenData["app_id"]);
            }
        }

        foreach ($this->securedRoutes as $route) {
            if (
                $route['controller'] == $dispatcher->getControllerName()
                && $route['action'] == $dispatcher->getActionName()
            ) {
                $param = $dispatcher->getParam('uid');
                if (!isset($param)) {
                    $this->response->setStatusCode(401, 'Unauthorized');
                    $this->response->setJsonContent([
                        'status' => 'error',
                        'message' => 'Unauthorized access',
                    ]);
                    $this->response->send();
                    return false;
                }
            }
        }

        return true;
    }

    protected function jsonResponse(string $status, string $message, array $data = null, int $statusCode = 200): Response
    {
        $response = [
            'status' => $status,
            'message' => $message,
        ];

        if ($status === 'success' && $data) {
            $response['data'] = $data;
        } elseif ($data) {
            $response = array_merge($response, $data);
        }

        return $this->response
            ->setJsonContent($response)
            ->setStatusCode($statusCode);
    }

    protected function generateClientToken(string $tokenId, string $uid, int $appId): string
    {
        $secretKey = $_ENV["TOKEN_SECRET"];
        $timestamp = time();
        $data = "$uid|$appId|$timestamp";

        $rawData = base64_encode($data);

        $signature = hash_hmac('sha256', $data, $secretKey);

        return "$rawData.$signature";
    }
    protected function validateClientToken(string $clientToken): ?array
    {
        $secretKey = $_ENV["TOKEN_SECRET"];

        [$rawData, $providedSignature] = explode('.', $clientToken);

        $data = base64_decode($rawData);
        if (!$data) {
            return null;
        }

        [$uid, $appId, $timestamp] = explode('|', $data);

        $computedSignature = hash_hmac('sha256', $data, $secretKey);

        if (!hash_equals($computedSignature, $providedSignature)) {
            return null;
        }

        return [
            'uid' => $uid,
            'app_id' => (int) $appId,
            'timestamp' => (int) $timestamp,
        ];
    }
    protected function guidv4($data = null)
    {
        $data = $data ?? random_bytes(16);
        assert(strlen($data) == 16);

        $data[6] = chr(ord($data[6]) & 0x0f | 0x40);
        $data[8] = chr(ord($data[8]) & 0x3f | 0x80);

        return vsprintf('%s%s-%s-%s-%s-%s%s%s', str_split(bin2hex($data), 4));
    }
}
