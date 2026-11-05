FROM python:3.12-slim
COPY --from=ghcr.io/astral-sh/uv:latest /uv /usr/local/bin/uv
WORKDIR /app
COPY pyproject.toml uv.lock ./
RUN uv sync --frozen --no-dev --no-install-project
COPY app.py ./
COPY views ./views
ENV PATH="/app/.venv/bin:$PATH"
EXPOSE 8080
USER nobody
CMD ["gunicorn", "-b", "0.0.0.0:8080", "-w", "2", "app:app"]
