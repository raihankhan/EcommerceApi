## Ingress

Ingress exposes HTTP and HTTPS routes from outside the cluster to services within the cluster. Traffic routing is controlled by rules defined on the Ingress resource.You must have an Ingress controller to satisfy an Ingress. Only creating an Ingress resource has no effect.

You may need to deploy an Ingress controller such as ingress-nginx. Follow the steps - 

1. Create a cluster. Must use extraPortMappings and node lables in cluster congiguration.
```yaml
apiVersion: kind.x-k8s.io/v1alpha4
kind: Cluster
nodes:
  - role: control-plane
    kubeadmConfigPatches:
      - |
        kind: InitConfiguration
        nodeRegistration:
          kubeletExtraArgs:
            node-labels: "ingress-ready=true"
    extraPortMappings:
    - containerPort: 80
      hostPort: 80
      protocol: TCP
    - containerPort: 443
      hostPort: 443
      protocol: TCP
```
This configuration will expose port 80 and 443 on the host. It’ll also add a node label so that the nginx-controller may use a node selector to target only this node. If a kind configuration has multiple nodes, it’s essential to only bind ports 80 and 443 on the host for one node because port collision will occur otherwise.
Then create a kind cluster using this config file via:

`kind create cluster --config cluster-extraportmapping.yaml`

2. Create ingress-nginx-controller and other required resources by executing this command

  `kubectl apply --filename https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/static/provider/kind/deploy.yaml`
  
3. Deploy necessary pods and services
In this case, I am deploying following deployment and service.
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api-server
  labels:
    app: server
spec:
  replicas: 2
  selector:
    matchLabels:
      app: server
  template:
    metadata:
      labels:
        app: server
    spec:
      containers:
      - name: ecommerce
        image: raihankhanraka/ecommerce-api:v1.0
        ports:
        - containerPort: 8080
---
apiVersion: v1
kind: Service
metadata:
  name: server-svc
spec:
  selector:
    app: server
  ports:
    - protocol: TCP
      port: 8080

---

```

4. Modify /etc/hosts on the host to direct traffic to the kind cluster’s ingress controller. We’ll need to get the IP address of our kind node’s Docker container first by running:

```
docker container inspect kind-control-plane \
              --f '{{ .NetworkSettings.Networks.kind.IPAddress }}'
```

Then add an entry to /etc/hosts with the IP address found that looks like:

`172.18.0.2 e-sell.com`

5. Create ingress with the yaml

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: ingress
spec:
  rules:
    - host: e-sell.com
      http:
        paths:
          - pathType: Prefix
            path: "/login"
            backend:
              service:
                name: server-svc
                port:
                  number: 8080

          - pathType: Prefix
            path: "/products"
            backend:
              service:
                name: server-svc
                port:
                  number: 8080
```

6. Finally, Go to Postman and send these queries to test that we have been able to successfully expose our application in kubernetes cluster using ingress -

`POST` `http://e-sell.com/login`

`GET` `http://e-sell.com/products`

`GET` `http://e-sell.com/products/LT01`

`Note :` Use create your deployment and ingress all in ingress-controller namespace. This namespace is created while creating the controller in step 2.
