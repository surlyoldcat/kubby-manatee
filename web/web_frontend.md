# Kubby Manatee Web UI

## General Description
This project is a simple web UI which displays MovieLens data from the Kubby Manatee API. It is a front-end for the API.

This app will be deployed to the local Kubernetes cluster and should be externally accessible.

## Important Files
- For context, the Kubby Manatee API lives at /home/moo/proj/kubmovie/api/
- API document is api/openapi.yaml

## Platform Guidelines
- base language is Python
- For web presentation, use `bottle` package with Jinja templates
- Use `uv` for package management.
- use `gunicorn` for web server
- Use the internal Kubernetes cluster address to talk to the API (namespace is `moocat`)
- No authentication required at this time.
- Use a LoadBalancer node type to expose the web app.
- The local Kubernetes cluster is using Minikube + Docker.

## HTML Styling
For now, I have only 1 visual style preference: the base color should be **pink**.

More to be added later.

## Initial Tasks
1. Create basic web app which shows paginated results from the `/movies` endoint in table format. The table should just show the attributes returned by the API call for now.
2. Create Dockerfile.
3. Deploye to Kubernetes cluster in the `moocat` namespace.
