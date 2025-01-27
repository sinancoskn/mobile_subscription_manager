<?php
declare(strict_types=1);

class IndexController extends ControllerBase
{
    public function healthAction()
    {
        $result = [
            ['status' => 'ok'],
        ];
        return $this->response->setJsonContent($result);
    }
}

