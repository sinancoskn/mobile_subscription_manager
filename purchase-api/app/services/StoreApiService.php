<?php

namespace App\Services;

class StoreApiService
{
    protected string $apiHost;

    public function __construct()
    {
        $this->apiHost = $_ENV["STORE_API_HOST"];
        if (!$this->apiHost) {
            throw new \Exception('MOCK_API_HOST is not configured in .env');
        }
    }

    /**
     * Call the Mock API for receipt validation
     *
     * @param array $data
     * @return array
     * @throws \Exception
     */
    public function validateReceipt(array $data): array
    {
        $endpoint = $this->apiHost . '/validate-receipt';

        // Perform HTTP request using cURL
        $ch = curl_init();

        curl_setopt($ch, CURLOPT_URL, $endpoint);
        curl_setopt($ch, CURLOPT_POST, true);
        curl_setopt($ch, CURLOPT_POSTFIELDS, json_encode($data));
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
        curl_setopt($ch, CURLOPT_HTTPHEADER, [
            'Content-Type: application/json',
        ]);

        $response = curl_exec($ch);
        $httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);

        if ($response === false) {
            $error = curl_error($ch);
            curl_close($ch);
            throw new \Exception('Mock API call failed: ' . $error);
        }

        curl_close($ch);

        // Decode JSON response
        $decodedResponse = json_decode($response, true);

        if ($httpCode !== 200) {
            throw new \Exception("Mock API call failed with HTTP code $httpCode: " . json_encode($decodedResponse));
        }

        return $decodedResponse;
    }
}
