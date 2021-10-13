## Externally Expose application services in kubernetes cluster using ingress

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

```go
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

`Note :` Create your deployment and ingress all in ingress-controller namespace. This namespace is created while creating the controller in step 2.

## Test role and rolebinding

Set the yaml for role and rolebinding with a user.

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  namespace: default
  name: pod-reader
rules:
  - verbs: ["get", "watch" , "list"]
    resources: ["pods"]
    apiGroups: [""] # "" indicates the core API group
---
apiVersion: rbac.authorization.k8s.io/v1
# This role binding allows "raihan@appscode.com" to read pods in the "default" namespace.
# You need to already have a Role named "pod-reader" in that namespace.
kind: RoleBinding
metadata:
  name: read-pods
  namespace: default
subjects:
  # You can specify more than one "subject"
  - kind: User
    name: raihan # "name" is case sensitive
    apiGroup: rbac.authorization.k8s.io
roleRef:
  # "roleRef" specifies the binding to a Role / ClusterRole
  kind: Role #this must be Role or ClusterRole
  name: pod-reader # this must match the name of the Role or ClusterRole you wish to bind to
  apiGroup: rbac.authorization.k8s.io
```

Now test if you can get(list) the pods in default namespace using --as flag

```azure
kubectl get pods --as raihan
```


## Use client-go library to develop kubernetes native app

Retrieve the kubeconfig file path from local host in order to build config using clientcmd package.

```go
var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
  config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
  ```

if you are running your application in kubernetes cluster, Then build the config like this -

```go
config, err = rest.InClusterConfig()
```
// InClusterConfig returns a config object which uses the service account
// kubernetes gives to pods. It's intended for clients that expect to be
// running inside a pod running on kubernetes. It will return ErrNotInCluster
// if called from a process not running in a kubernetes environment.

Build your dynamic client using this config

```go
dynamicClient, err := dynamic.NewForConfig(config)
```

NewForConfig creates a new dynamic client or returns an error.

Create deployment resource
```go
depResource := schema.GroupVersionResource{
		Group:    "apps",
		Version:  "v1",
		Resource: "deployments",
	}
```
// GroupVersionResource unambiguously identifies a resource.  It doesn't anonymously include GroupVersion
// to avoid automatic coercion.  It doesn't use a GroupVersion to avoid custom marshalling
type GroupVersionResource struct {
Group    string
Version  string
Resource string
}
```go

```

create a deployment using client-go dynamic package. 
	
// Unstructured allows objects that do not have Golang structs registered to be manipulated
// generically. This can be used to deal with the API objects from a plug-in. Unstructured
// objects still have functioning TypeMeta features-- kind, version, etc.

```go
deployment := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name": "apiserver",
			},
			"spec": map[string]interface{}{
				"replicas": 2,
				"selector": map[string]interface{}{
					"matchLabels": map[string]interface{}{
						"app": "server",
					},
				},
				"template": map[string]interface{}{
					"metadata": map[string]interface{}{
						"labels": map[string]interface{}{
							"app": "server",
						},
					},
					"spec": map[string]interface{}{
						"containers": []map[string]interface{}{
							{
								"name":  "ecommerce",
								"image": "raihankhanraka/ecommerce-api:v1.1",
								"ports": []map[string]interface{}{
									{
										"name":          "http",
										"protocol":      "TCP",
										"containerPort": 8080,
									},
								},
							},
						},
					},
				},
			},
		},
	}
  ```

Now create deployment.

```go
dep, err := dynamicClient.Resource(depResource).Namespace("default").Create(context.TODO(), deployment, v1.CreateOptions{})
```

```go
svcResource := schema.GroupVersionResource{
//Group:    "",
Version:  "v1",
Resource: "services",
}

service := &unstructured.Unstructured{
Object: map[string]interface{}{
"apiVersion": "v1",
"kind":       "Service",
"metadata": map[string]interface{}{
"name": "server-svc",
},
"spec": map[string]interface{}{
"selector": map[string]interface{}{
"app": "server",
},
"ports": []map[string]interface{}{
{
"protocol":   "TCP",
"targetPort": 8080,
"port":       8080,
},
},
},
},
}

fmt.Printf("creating service %s\n", service.GetName())
svc, err := dynamicClient.Resource(svcResource).Namespace("default").Create(context.TODO(), service, v1.CreateOptions{})
if err != nil {
panic(fmt.Errorf("failed to create service -- %s\n", err.Error()))
}
```

## kubernetes under the hood

 - kubectl translates your imperative command into a declarative Kubernetes Deployment object. A Deployment is a higher-level API that allows rolling updates (see below).
 - kubectl sends the Deployment to the Kubernetes API server, kube-apiserver, which runs in-cluster.
 - kube-apiserver saves the Deployment to etcd (a distributed key-value store), which also runs in-cluster, and responds to kubectl.
 - Asynchronously, the Kubernetes controller manager, kube-controller-manager, which watches for Deployment events (among others), creates a ReplicaSet from the Deployment and sends it to kube-apiserver. A ReplicaSet is a version of a Deployment. During a rolling update, a new ReplicaSet will be created and progressively scaled out to the desired number of replicas, while the old one is scaled in to zero.
 - kube-apiserver saves the ReplicaSet to etcd.
 - Asynchronously, kube-controller-manager, creates two Pods (or more if we scale out) from the ReplicaSet and sends them to kube-apiserver. Pods are the basic unit of Kubernetes. They represent one or several containers sharing a Linux cgroup and namespaces.
 - kube-apiserver saves the Pods to etcd.
 - Asynchronously, the Kubernetes scheduler, kube-scheduler, which watches for Pod events, updates each Pod to assign it to a Node and sends them back to kube-apiserver.
 - kube-apiserver saves the Pods to etcd.
 - Finally, the kubelet that runs on the assigned Node, always watching, actually starts the container.

Note: the controller, scheduler and kubelet also send status information back to the API server.

In summary, Kubernetes is a CRUD API with control loops.

