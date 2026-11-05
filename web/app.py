import os

import requests
from bottle import Bottle, jinja2_view, redirect, request, response

API_BASE_URL = os.environ.get(
    "API_BASE_URL", "http://kubby-manatee-api.moocat.svc.cluster.local:8083"
)
PAGE_SIZE = 50

app = Bottle()


@app.get("/")
def index():
    redirect("/movies")


@app.get("/health")
def health():
    return {"status": "ok"}


@app.get("/movies")
@jinja2_view("movies.html", template_lookup=[os.path.join(os.path.dirname(__file__), "views")])
def movies():
    try:
        page = max(int(request.query.get("page", 1)), 1)
    except ValueError:
        page = 1
    try:
        r = requests.get(f"{API_BASE_URL}/movies", params={"pagenum": page}, timeout=5)
        r.raise_for_status()
        rows = r.json()["movies"]
    except (requests.RequestException, KeyError, ValueError) as e:
        response.status = 502
        return {"error": str(e), "page": page, "movies": [], "has_next": False}
    return {"error": None, "page": page, "movies": rows, "has_next": len(rows) >= PAGE_SIZE}
