<?php
declare(strict_types=1);

namespace App\Services\Queue;

use PhpAmqpLib\Connection\AMQPStreamConnection;
use PhpAmqpLib\Message\AMQPMessage;

class RabbitMQQueue implements QueueInterface
{
    private AMQPStreamConnection $connection;
    private $channel;

    public function connect(array $config): void
    {
        $this->connection = new AMQPStreamConnection(
            $config['host'],
            $config['port'],
            $config['username'],
            $config['password']
        );
        $this->channel = $this->connection->channel();
    }

    public function publish(string $queue, array $message): bool
    {
        $this->channel->queue_declare($queue, false, true, false, false);

        $msg = new AMQPMessage(json_encode($message), ['delivery_mode' => 2]); // Persistent message
        $this->channel->basic_publish($msg, '', $queue);

        return true;
    }

    public function consume(string $queue, callable $callback): void
    {
        $this->channel->queue_declare($queue, false, true, false, false);

        $this->channel->basic_consume(
            $queue,
            '',
            false,
            true,
            false,
            false,
            function ($msg) use ($callback) {
                $callback(json_decode($msg->body, true));
            }
        );

        while ($this->channel->is_consuming()) {
            $this->channel->wait();
        }
    }

    public function close(): void
    {
        $this->channel->close();
        $this->connection->close();
    }
}
