from confluent_kafka import Consumer
import json
import logging
from config import KAFKA_BOOTSTRAP_SERVERS, CONSUME_TOPICS

logger = logging.getLogger("events-consumer")

conf = {
    'bootstrap.servers': KAFKA_BOOTSTRAP_SERVERS,
    'group.id': 'events-group',
    'auto.offset.reset': 'earliest'
}

consumer = Consumer(conf)

def start_consumer():
    consumer.subscribe(CONSUME_TOPICS)
    logger.info(f"Subscribed to topics: {CONSUME_TOPICS}")
    try:
        while True:
            msg = consumer.poll(1.0)
            if msg is None:
                continue
            if msg.error():
                logger.error(f"Consumer error: {msg.error()}")
                continue

            data = msg.value().decode('utf-8')
            logger.info(f"Consumed message from {msg.topic()}: {data}")
    except KeyboardInterrupt:
        pass
    finally:
        consumer.close()