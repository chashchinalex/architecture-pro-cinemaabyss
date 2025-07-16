from flask import Flask, request, jsonify
from confluent_kafka import Producer
import logging
import json
import threading
import os
from consumer import start_consumer
from config import KAFKA_BOOTSTRAP_SERVERS, TOPIC_MAP


logging.basicConfig(level=logging.INFO)
logger = logging.getLogger("events-service")


app = Flask(__name__)

# Продюсер
producer = Producer({'bootstrap.servers': KAFKA_BOOTSTRAP_SERVERS})

def send_event(event_type, data):
    topic = TOPIC_MAP.get(event_type)
    if not topic:
        return jsonify({"status": "error", "detail": "Unknown event type"}), 400
    try:
        producer.produce(topic, json.dumps(data).encode('utf-8'))
        producer.flush()
        logger.info(f"Produced event to {topic}: {data}")
        return jsonify({"status": "success"}), 201
    except Exception as e:
        logger.error(f"Error producing event: {e}")
        return jsonify({"status": "error", "detail": str(e)}), 500

@app.route("/api/events/health", methods=["GET"])
def health():
    return jsonify({"status": True}), 200

@app.route("/api/events/user", methods=["POST"])
def create_user_event():
    data = request.json
    if not data or "movie_id" not in data or "action" not in data:
        return jsonify({"status": "error", "detail": "Invalid user event"}), 400
    return send_event("user", data)

@app.route("/api/events/payment", methods=["POST"])
def create_payment_event():
    data = request.json
    if not data or "payment_id" not in data or "status" not in data:
        return jsonify({"status": "error", "detail": "Invalid payment event"}), 400
    return send_event("payment", data)

@app.route("/api/events/movie", methods=["POST"])
def create_movie_event():
    data = request.json
    if not data or "movie_id" not in data or "title" not in data:
        return jsonify({"status": "error", "detail": "Invalid movie event"}), 400
    return send_event("movie", data)


threading.Thread(target=start_consumer, daemon=True).start()

if __name__ == "__main__":
    port = int(os.environ.get("PORT", 8000))
    app.run(host="0.0.0.0", port=port)