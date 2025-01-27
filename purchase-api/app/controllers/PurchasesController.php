<?php
declare(strict_types=1);

class PurchasesController extends ControllerBase
{

    public function purchaseAction()
    {
        try {
            $clientId = $this->dispatcher->getParam('uid');
            $appId = $this->dispatcher->getParam('app_id');

            $data = $this->request->getJsonRawBody(true);

            if (!isset($data['receipt'])) {
                return $this->response->setJsonContent([
                    'status' => 'error',
                    'message' => 'Receipt is required',
                ])->setStatusCode(400);
            }

            $mockApiService = $this->di->get('storeApiService');
            $response = $mockApiService->validateReceipt(['receipt' => $data['receipt']]);

            if ($response['status'] !== true) {
                return $this->response->setJsonContent([
                    'status' => 'error',
                    'message' => 'Invalid receipt',
                ])->setStatusCode(400);
            }

            $expireAt = $response['expire_date'];
            $status = 'started';

            $subscription = Subscriptions::findFirst([
                'conditions' => 'uid = :uid: AND app_id = :app_id:',
                'bind' => [
                    'uid' => $clientId,
                    'app_id' => $appId,
                ],
                'for_update' => true,
            ]);

            if ($subscription) {

                if ($subscription->status === 'canceled') {
                    $status = 'started';
                } elseif ($subscription->expire_at < date('Y-m-d H:i:s')) {
                    $status = 'renewed';
                }

                $subscription->status = $status;
                $subscription->expire_at = $expireAt;
                $subscription->updated_at = date('Y-m-d H:i:s');

                $subscription->save();
            } else {

                error_log("TAMARA BURADA MISIN?");

                $subscription = new Subscriptions();
                $subscription->assign([
                    'uid' => $clientId,
                    'app_id' => $appId,
                    'receipt' => $data['receipt'],
                    'status' => $status,
                    'expire_at' => $expireAt,
                ]);
                $subscription->save();
            }

            return $this->response->setJsonContent([
                'status' => 'success',
                'message' => 'Subscription processed successfully',
                'data' => [
                    'uid' => $clientId,
                    'app_id' => $appId,
                    'status' => $status,
                    'expire_at' => $expireAt,
                ],
            ]);
        } catch (\Exception $e) {
            return $this->response->setJsonContent([
                'status' => 'error',
                'message' => $e->getMessage(),
            ])->setStatusCode(500);
        }
    }

    public function checkSubscriptionAction()
    {
        try {
            $clientId = $this->dispatcher->getParam('uid');
            $appId = $this->dispatcher->getParam('app_id');

            if (!$clientId || !$appId) {
                return $this->response->setJsonContent([
                    'status' => 'error',
                    'message' => 'Missing uid or app_id',
                ])->setStatusCode(400);
            }

            $subscription = Subscriptions::findFirst([
                'conditions' => 'uid = :uid: AND app_id = :app_id:',
                'bind' => [
                    'uid' => $clientId,
                    'app_id' => $appId,
                ],
            ]);

            if (!$subscription) {
                return $this->response->setJsonContent([
                    'status' => 'error',
                    'message' => 'Subscription not found',
                ])->setStatusCode(404);
            }

            return $this->response->setJsonContent([
                'status' => 'success',
                'message' => 'Subscription retrieved successfully',
                'data' => [
                    'uid' => $subscription->uid,
                    'app_id' => $subscription->app_id,
                    'status' => $subscription->status,
                    'expire_at' => $subscription->expire_at,
                ],
            ]);
        } catch (\Exception $e) {
            return $this->response->setJsonContent([
                'status' => 'error',
                'message' => $e->getMessage(),
            ])->setStatusCode(500);
        }
    }


}

