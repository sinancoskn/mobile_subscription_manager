<?php
declare(strict_types=1);

namespace App\Services\Queue;

interface QueueInterface
{
    public function connect(array $config): void;
    public function publish(string $topic, array $message): bool;
    public function consume(string $topic, callable $callback): void;
}