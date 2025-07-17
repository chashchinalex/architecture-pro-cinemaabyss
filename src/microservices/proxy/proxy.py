import os
import random
import logging
from flask import Flask, request, Response
import requests

app = Flask(__name__)

MONOLITH_URL = os.getenv("MONOLITH_URL")
MOVIE_URL = os.getenv("MOVIES_SERVICE_URL")
PORT = int(os.getenv("PORT", 8080))
GRADUAL_MIGRATION = bool(os.getenv("GRADUAL_MIGRATION"))
MOVIES_MIGRATION_PERCENT = int(os.getenv("MOVIES_MIGRATION_PERCENT"))

logging.basicConfig(level=logging.INFO)

def proxy_request(target_url):
    try:
        full_url = target_url + request.path
        if request.query_string:
            full_url += '?' + request.query_string.decode()

        resp = requests.request(
            method=request.method,
            url=full_url,
            headers={key: value for key, value in request.headers if key.lower() != 'host'},
            data=request.get_data(),
            cookies=request.cookies,
            allow_redirects=False,
        )

        excluded_headers = ['content-encoding', 'content-length', 'transfer-encoding', 'connection']
        headers = [(name, value) for (name, value) in resp.headers.items() if name.lower() not in excluded_headers]

        return Response(resp.content, resp.status_code, headers)
    except requests.RequestException as e:
        logging.error(f"Proxy error to {target_url}: {e}")
        return Response("Proxy error", status=502)


@app.route('/health', methods=['GET'])
def health():
    return "OK", 200

@app.route('/', defaults={'path': ''}, methods=["GET", "POST", "PUT", "DELETE", "PATCH"])
@app.route('/<path:path>', methods=["GET", "POST", "PUT", "DELETE", "PATCH"])
def route_request(path):
    logging.info(f"Start routing: {request.path}")
    if GRADUAL_MIGRATION and request.path.startswith("/api/movies"):
        percent = random.randint(0, 99)
        logging.info(f"Migration ON -> before Check %{percent}: {request.path}")
        if percent < MOVIES_MIGRATION_PERCENT:
            logging.info(f"Migration {percent}% -> Routing to MOVIES_SERVICE: {request.path}")
            return proxy_request(MOVIE_URL)
    logging.info(f"Migration OFF -> Routing to MONOLITH_URL: {request.path}")
    return proxy_request(MONOLITH_URL)

if __name__ == "__main__":
    logging.info(f"variables: MONOLITH_URL={MONOLITH_URL}, MOVIES_SERVICE_URL={MOVIE_URL}, PORT={PORT}, GRADUAL_MIGRATION={GRADUAL_MIGRATION}, MOVIES_MIGRATION_PERCENT={MOVIES_MIGRATION_PERCENT}")
    app.run(host="0.0.0.0", port=PORT)
