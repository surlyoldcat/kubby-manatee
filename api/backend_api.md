# Kubby Manatee Backend API

## General Description
This is a learning project, intended to help me get familiar with hosting/deploying applications on Kubernetes.

I want to create a simple REST API to serve data from a PostgreSQL database, using the MovieLens dataset. *Note, the database server runs inside the local minikube cluster (`postgres-postgresql` in the `moocat` namespace).*

The API should be written in Go (v1.25).

This API will be hosted as a Kubernetes service on my local cluster. I will be sending requests to the API using basic testing tools such as cUrl and Postman. There will also be a web client, also hosted in the same Kubernetes pod, which will need to access the API.

## Database
The database is a Postgres database running in the local minikube cluster. There are a couple files to reference:
- kubmovie/db/README.txt has a description of the MovieLens dataset
- kubmovie/db/create_schema.sql has the actual DDL for creating the database schema (assume `public` schema)
- Database connection info is stored in an .env file for local development (host/port point at the minikube NodePort, for access from outside the cluster). In the cluster, the API uses `postgres-postgresql.moocat.svc.cluster.local:5432`, with `DB_USER`/`DB_PASSWORD` from the `app-db-creds` Secret.
- The data in the MovieLens tables should be considered read-only. This project is going to focus more on querying existing data.
- The `movielens` database is accessible from this machine via the NodePort. The server does not use TLS, so `DB_SSLMODE=disable` is set in `.env`.

*Database quirk: the DB adapter retries once, 30 seconds after the first connection timeout error.*

## API
The REST API will return data in a basic JSON format. Attribute names should match database column names. Entity/structure names should match the table names.

For the first iteration, the REST API will have  2 operations, using the `movies` table as a data source:

**movies**
GET: returns a list of movie items. The list items should have the movie ID and title columns. The list should be paginated using a pagenum querystring parameter. Page size should be limited to 50 items.

**movies/{id}**
GET: returns a single movie item. All columns from the `movies` table should be included as attributes. If the ID is not found, return a 404 error message.

*Note: we probably also need to add a basic healthcheck endpoint which returns a 200 status, so Kubernetes can check the health of the service. This can just be a GET /health endpoint.*




## Platform Guidelines
- Use GORM (https://gorm.io/) to interact with the Postgres database.
- For now, we will store database configuration and credentials in a .env file. Create a placeholder to update this to fetch db config from the Kubernetes environment.
- Use the Gin Web Framework (https://gin-gonic.com/) to implement the REST API.
- An OpenAPI document should be created to enable future expansion of the API.
- Use the testify package for creating tests.
- The API will be exposed via TCP port 8083

## Deployment
- This API will be deployed as a service on Kubernetes.
- A Dockerfile must be created to deploy the API.
- A build script (bash is fine) should be created to build the Docker image.
- Since this is a learning project, Docker images will just be stored locally. They just have to be accessible to the local Kubernetes environment (minikube).

## Initial Tasks
1. Identify any critical items not covered in this document, and append them in a new section "Blockers".
2. Create a Go mod/project for the API.
3. Create a basic repository pattern with simple fetch methods for the db tables (Create, Update, and Delete are not supported.)
4. Create integration tests to verify db connectivity and query results.
5. Create API layer using the 2 `movies` operations defined in the API section (and health check.)
6. Create Dockerfile.
7. Create Kubernetes manifests and deployment files.
8. Build Docker image.
9. Deploy the API to minikube and verify.

## Blockers
Identified during task 1. All resolved or noted as of deployment (tasks 6-9 complete).

1. ~~**minikube is stopped.**~~ **Resolved.** minikube is running; `build.sh` builds the image directly into minikube's Docker daemon (`minikube docker-env`). It sets `DOCKER_BUILDKIT=0` because buildx could not boot buildkit against that daemon.
2. ~~**Image name/tag not specified.**~~ **Resolved.** `kubby-manatee-api:latest` with `imagePullPolicy: IfNotPresent`.
3. ~~**No Kubernetes manifests.**~~ **Resolved.** `k8s/api.yaml` contains a ConfigMap, a Deployment and a NodePort Service in the `moocat` namespace. `DB_USER`/`DB_PASSWORD` come from the existing `app-db-creds` Secret (keys `username`, `password`).
4. ~~**`.env` must not be baked into the image.**~~ **Resolved.** `.dockerignore` excludes `.env`.
5. ~~**Network path from the cluster to Postgres is unverified.**~~ **Resolved.** The database is the in-cluster Postgres (`postgres-postgresql.moocat.svc.cluster.local:5432`). The NodePort in `.env` is only for access from outside the cluster. Verified end to end through the API.
6. **CORS not specified.** Still open. The web client shares the pod, so same-origin is assumed. If the client is served from a different origin or port, CORS middleware will be needed.
7. **Not a git repository.** Still open. No version control is set up for this project.

### Defaults assumed
- `pagenum` is 1-based, defaults to 1; a non-integer or < 1 value returns 400.
- List response: `{"pagenum": n, "movies": [{"movie_id", "title"}, ...]}`, ordered by `movie_id`.
- Not found: 404 with `{"error": "movie not found"}`.
- DB connect timeout 10s; on a timeout error, one retry after 30s.
- `GET /health` returns 200 `{"status": "ok"}` without touching the DB.
