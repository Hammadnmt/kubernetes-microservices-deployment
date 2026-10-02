# Microservices Architecture

This directory contains individual declarative manifests for each backend service. Each file encapsulates its **Deployment** (pods, scaling, health) and **Service** (internal DNS & load balancer).

## Services Overview

| Service | File | Replicas | Internal DNS Endpoint | Port |
|---|---|---|---|---|
| **User Service** | [`user-service.yaml`](file:///Users/mac/Desktop/Hammad/workspace/kubernetes-get-started/services/user-service.yaml) | 3 | `http://user-service` | 80 -> 8080 |
| **Inventory Service** | [`inventory-service.yaml`](file:///Users/mac/Desktop/Hammad/workspace/kubernetes-get-started/services/inventory-service.yaml) | 2 | `http://inventory-service` | 80 -> 8080 |
| **Media Service** | [`media-service.yaml`](file:///Users/mac/Desktop/Hammad/workspace/kubernetes-get-started/services/media-service.yaml) | 3 | `http://media-service` | 80 -> 8080 |
| **Lesson Service** | [`lesson-service.yaml`](file:///Users/mac/Desktop/Hammad/workspace/kubernetes-get-started/services/lesson-service.yaml) | 2 | `http://lesson-service` | 80 -> 8080 |

---

## Usage Commands

### 1. Deploy all services at once:
```bash
kubectl apply -f services/
```

### 2. Deploy or update a single service independently:
```bash
kubectl apply -f services/user-service.yaml
```

### 3. Check status of all deployments, pods, and services:
```bash
kubectl get deployments,pods,services -l tier=backend
```

### 4. Test inter-service communication (Curl from inside the cluster):
```bash
# Test User Service
kubectl run test-curl --rm -it --image=curlimages/curl -- curl -s http://user-service

# Test Lesson Service
kubectl run test-curl --rm -it --image=curlimages/curl -- curl -s http://lesson-service
```

### 5. Tear down all backend services:
```bash
kubectl delete -f services/
```
