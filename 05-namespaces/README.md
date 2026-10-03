# 05: Namespaces (Multi-Environment Architecture)

Namespaces allow a single physical Kubernetes cluster to be divided into multiple virtual clusters. They provide scoping for names, resource quotas, and access control.

---

## 1. What are Namespaces?

Namespaces are like folders on your computer:
- You cannot have two files named `app.js` in the same folder.
- But you can have `/development/app.js` and `/production/app.js` running side-by-side on the same drive.

In Kubernetes, you cannot have two Services or Pods with the exact same name in the `default` namespace. Namespaces solve this by providing separate logical boundaries.

---

## 2. Declarative Definition

We declare namespaces in [`01-namespaces.yaml`](file:///Users/mac/Desktop/Hammad/workspace/kubernetes-get-started/05-namespaces/01-namespaces.yaml):

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: development
  labels:
    environment: dev
    team: backend
---
apiVersion: v1
kind: Namespace
metadata:
  name: production
  labels:
    environment: prod
    team: backend
```

Apply them:
```bash
kubectl apply -f 05-namespaces/01-namespaces.yaml
```

---

## 3. Resolving Name Collisions (Dev vs. Prod)

We deploy identical microservice names into both rooms:
* [`02-dev-app.yaml`](file:///Users/mac/Desktop/Hammad/workspace/kubernetes-get-started/05-namespaces/02-dev-app.yaml): Deploys `user-service` (1 replica) and `lesson-service` (1 replica) into `development`.
* [`03-prod-app.yaml`](file:///Users/mac/Desktop/Hammad/workspace/kubernetes-get-started/05-namespaces/03-prod-app.yaml): Deploys `user-service` (3 replicas) and `lesson-service` (2 replicas) into `production`.

### Verification Side-by-Side:
```bash
$ kubectl get pods,svc -n development
NAME                     READY   STATUS    AGE
pod/lesson-deploy-...    1/1     Running   1m
pod/user-deploy-...      1/1     Running   1m

NAME                     TYPE        CLUSTER-IP    PORT(S)
service/lesson-service   ClusterIP   10.96.55.34   80/TCP
service/user-service     ClusterIP   10.96.55.51   80/TCP

---
$ kubectl get pods,svc -n production
NAME                     READY   STATUS    AGE
pod/lesson-deploy-...    1/1     Running   1m (2 replicas)
pod/user-deploy-...      1/1     Running   1m (3 replicas)

NAME                     TYPE        CLUSTER-IP     PORT(S)
service/lesson-service   ClusterIP   10.96.61.120   80/TCP
service/user-service     ClusterIP   10.96.61.186   80/TCP
```

### Key Observation:
Inside each namespace, `lesson-service` calls `http://user-service`. CoreDNS automatically scopes the search to the local namespace:
- In `development`, it routes to `10.96.55.51`.
- In `production`, it routes to `10.96.61.186`.
Zero name collisions!

---

## 4. The "Noisy Neighbor" Shield (ResourceQuota)

We attached a `ResourceQuota` directly into [`02-dev-app.yaml`](file:///Users/mac/Desktop/Hammad/workspace/kubernetes-get-started/05-namespaces/02-dev-app.yaml) using the `---` multi-document separator:

```yaml
---
apiVersion: v1
kind: ResourceQuota
metadata:
  name: dev-quota
  namespace: development
spec:
  hard:
    pods: "3"
    requests.cpu: "200m"
    requests.memory: "200Mi"
    limits.cpu: "400m"
    limits.memory: "400Mi"
```

### Checking Quota Tracking:
```bash
$ kubectl describe resourcequota dev-quota -n development
Resource         Used   Hard
--------         ----   ----
limits.cpu       200m   400m
limits.memory    128Mi  400Mi
pods             2      3
requests.cpu     100m   200m
requests.memory  64Mi   200Mi
```

### The Over-Scaling Test (Disapproving the Noise):
We simulated a runaway dev process attempting to scale to 10 pods:
```bash
kubectl scale deployment user-deploy --replicas=10 -n development
```

### The Result:
Kubernetes allowed 1 pod (reaching the hard cap of 3), and **strictly rejected all remaining pods**:
```text
Warning  FailedCreate  replicaset-controller  Error creating: pods "user-deploy-..." is forbidden: exceeded quota: dev-quota, requested: pods=1, used: pods=3, limited: pods=3
```
Production remained **100% unaffected** with all 5 pods running healthy and zero CPU/RAM starvation.

---

## 5. Essential CLI Commands

```bash
# List all namespaces
kubectl get ns

# Inspect resource quota in a namespace
kubectl describe resourcequota dev-quota -n development

# Filter namespaces by label
kubectl get ns -l environment=dev

# List pods inside a specific namespace
kubectl get pods -n development

# List pods across ALL namespaces at once
kubectl get pods -A

# Switch your active terminal namespace (no need to type -n every time!)
kubectl config set-context --current --namespace=development

# Switch back to default
kubectl config set-context --current --namespace=default
```


