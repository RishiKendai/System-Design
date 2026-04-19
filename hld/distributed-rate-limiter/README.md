# Distributed Rate Limiter (HLD)

> To view the architecture diagram (`drl.excalidraw`) in VSCode, install the **Excalidraw** extension.

## Goals

### Functional Requirements
- Identify client using `user_id`, `ip_address`, or `api_key`.
- Enforce configurable rate-limit rules per client/policy.
- Reject over-limit requests with HTTP `429`.
- Return helpful response headers:
  - `X-RateLimit-Limit`
  - `X-RateLimit-Remaining`
  - `X-RateLimit-Reset`

### Non-Functional Requirements
- Add minimal request latency overhead (target: `<10ms` per check).
- Prefer availability over strict consistency.
- Handle very high scale (~`1M req/sec`, ~`100M DAU`).

## Core Entities
- **Client**: caller identified by API key/IP/user id.
- **Request**: incoming API request to be validated.
- **Policy/Rules**: configurable limits (e.g., 100 requests/min).
- **Rate Limiter Engine**: decisioning component for allow/reject.

## Algorithm Options Considered

### 1) Fixed Window Counter
- Counts requests inside fixed windows.
- Resets at interval boundaries.
- Limitation: boundary burst (double spending near edges).

### 2) Sliding Window Log
- Stores timestamp for every request in rolling window.
- High accuracy, but expensive memory/cleanup cost at high scale.

### 3) Sliding Window Counter
- Uses weighted counts from current and previous windows.
- Lower memory than logs.
- Limitation: approximate result (assumes uniform traffic distribution).

### 4) Token Bucket
- Tokens refill at fixed rate; each request consumes one.
- Supports burst handling naturally.
- Limitation: does not enforce strict per-instant cap.

## High-Level Architecture

1. **Client** sends request to **API Gateway**.
2. API Gateway reads applicable rule (from local in-memory cache).
3. API Gateway performs atomic counter/token operation in **Redis Cluster**.
   - Writes/checks are done using **Lua scripts** to guarantee atomic updates.
4. Based on Redis result:
   - **Allowed** -> route request to backend **Server** pool.
   - **Rejected** -> return HTTP `429` immediately with rate-limit headers.

## Data and Control Plane

### Data Plane (per request, low latency path)
- API Gateway <-> Redis Cluster for usage state/counters.
- API Gateway -> Server (only if request is allowed).

### Control Plane (rule updates)
- Rules are updated in **ZooKeeper** or **Consul**.
- API Gateway instances keep a **watch** on config changes.
- On update, gateway fetches latest rules and refreshes local in-memory cache.
- This avoids a config-store read on every request.

## Why This Design Works
- **Fast decisions**: local rule cache + Redis atomic operations.
- **Highly available**: distributed gateway and Redis cluster deployment.
- **Operationally flexible**: dynamic rule updates without gateway restarts.
- **Scalable**: horizontally scalable gateways and backend services.

## Rejection Response Contract

When limit is exceeded:
- HTTP status: `429 Too Many Requests`
- Headers:
  - `X-RateLimit-Limit`: configured limit for the request (example: `100`)
  - `X-RateLimit-Remaining`: remaining requests in current window (example: `0`)
  - `X-RateLimit-Reset`: Unix timestamp when limit resets (example: `1640995200`)
