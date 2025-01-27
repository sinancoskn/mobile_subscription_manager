<?php
declare(strict_types=1);

namespace App\Services\Queue;

use RdKafka\Producer;
use RdKafka\KafkaConsumer;

class KafkaQueue implements QueueInterface
{
    private Producer $producer;
    private array $config;

    public function connect(array $config): void
    {
        $this->config = $config;
        $this->producer = new Producer();
        $this->producer->addBrokers($config['brokers']);
    }

    public function publish(string $topic, array $message): bool
    {
        $kafkaTopic = $this->producer->newTopic($topic);
        $kafkaTopic->produce(RD_KAFKA_PARTITION_UA, 0, json_encode($message));

        $this->producer->flush(1000);
        return true;
    }

    public function consume(string $topic, callable $callback): void
    {
        $conf = new \RdKafka\Conf();
        $conf->set('metadata.broker.list', $this->config['brokers']);
        $conf->set('group.id', $this->config['group_id']);
        $conf->set('enable.auto.commit', 'false');

        $consumer = new KafkaConsumer($conf);
        $consumer->subscribe([$topic]);

        while (true) {
            $message = $consumer->consume(120 * 1000);
            if ($message->err) {
                continue;
            }

            $processed = $callback(json_decode($message->payload, true));
            if ($processed) {
                $consumer->commit();
            }
        }
    }
}