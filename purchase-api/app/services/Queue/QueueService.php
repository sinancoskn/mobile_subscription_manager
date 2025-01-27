<?php
declare(strict_types=1);

namespace App\Services\Queue;

class QueueService
{
    private QueueInterface $driver;

    public function __construct(QueueInterface $driver)
    {
        $this->driver = $driver;
    }

    public function connect(array $config): void
    {
        $this->driver->connect($config);
    }

    public function publish(string $topic, array $message): bool
    {
        return $this->driver->publish($topic, $message);
    }

    public function consume(string $topic, callable $callback): void
    {
        $this->driver->consume($topic, $callback);
    }
}
