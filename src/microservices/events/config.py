import os

KAFKA_BOOTSTRAP_SERVERS = os.environ.get("KAFKA_BROKERS")

# Сопоставляем топики
TOPIC_MAP = {
    "user": "user-events",
    "payment": "payment-events",
    "movie": "movie-events"
}

CONSUME_TOPICS = list(TOPIC_MAP.values())