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
4. Modify /etc/hosts on the host to direct traffic to the kind cluster’s ingress controller. We’ll need to get the IP address of our kind node’s Docker container first by running:
```
docker container inspect kind-control-plane \
              --f '{{ .NetworkSettings.Networks.kind.IPAddress }}'
```
Then add an entry to /etc/hosts with the IP address found that looks like:
`172.18.0.2 evally.com`

5. Create the ingress yaml
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: example-ingress
spec:
  rules:
    - host: evally.com
      http:
        paths:
          - pathType: Prefix
            path: "/status"
            backend:
              service:
                name: foo-service
                port:
                  number: 5678
          - pathType: Prefix
            path: "/offer"
            backend:
              service:
                name: bar-service
                port:
                  number: 5678
```
6. Finally, we can curl evally.com:
`curl evally.com/status`
`curl evally.com/offer`
