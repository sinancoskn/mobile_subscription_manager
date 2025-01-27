<?php
declare(strict_types=1);

class DevicesController extends ControllerBase
{

    public function registerAction()
    {
        // error_log(print_r(getenv('POSTGRESQL_HOST'), true));

        // $queueService = $this->di->get('queueService');

        // $success = $queueService->publish('subscription_events', [
        //     'event' => 'subscription_updated',
        //     'uid' => '8f481bf7-e09e-45f4-8a71-1ecc3c987166',
        //     'status' => 'started',
        //     'app_id' => 1,
        //     'timestamp' => time(),
        // ]);


        $data = $this->request->getJsonRawBody();

        $device = new Devices();
        $device->assign([
            'uid' => $data->uid ?? null,
            'app_id' => $data->app_id ?? null,
            'language' => $data->language ?? null,
            'os' => $data->os ?? null,
        ]);

        if (!$device->validation()) {
            $messages = $device->getMessages();
            return $this->jsonResponse(
                'error',
                'Validation failed.',
                ['errors' => array_map(fn($msg) => $msg->getMessage(), $messages)],
                400
            );
        }

        $app = Apps::findFirst([
            'conditions' => 'id = :app_id:',
            'bind' => ['app_id' => $device->app_id],
        ]);

        if (!$app) {
            return $this->jsonResponse(
                'error',
                'Invalid app_id. The app does not exist.',
                null,
                404
            );
        }

        $existingDevice = Devices::findFirst([
            'conditions' => 'app_id = :app_id: AND uid = :uid:',
            'bind' => [
                'app_id' => $device->app_id,
                'uid' => $device->uid,
            ],
        ]);

        $tokenId = $this->guidv4();

        if ($existingDevice) {
            $clientToken = $this->generateClientToken($tokenId, $existingDevice->uid, $existingDevice->app_id);
            return $this->jsonResponse(
                'success',
                'Register OK',
                ['client_token' => $clientToken]
            );
        }

        if ($device->save()) {
            $clientToken = $this->generateClientToken($tokenId, $device->uid, $device->app_id);
            return $this->jsonResponse(
                'success',
                'Device registered successfully',
                ['client_token' => $clientToken]
            );
        }

        return $this->jsonResponse(
            'error',
            'Failed to register the device.',
            ['errors' => $device->getMessages()],
            500
        );
    }


}

