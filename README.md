# Kubernetes Cluster Setup

***(single control-plane cluster***)

## Install kubeadm (all nodes)

```bash
user@master:~$ sudo apt-get update && sudo apt-get install -y apt-transport-https curl
user@master:~$ curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
user@master:~$ cat <<EOF | sudo tee /etc/apt/sources.list.d/kubernetes.list
deb https://apt.kubernetes.io/ kubernetes-xenial main
EOF
user@master:~$ sudo apt-get update
user@master:~$ sudo apt-get install -y kubelet kubeadm kubectl
user@master:~$ sudo apt-mark hold kubelet kubeadm kubectl
```

## Initialise master node

```bash
user@master:~$ sudo kubeadm init --pod-network-cidr=192.168.0.0/16
user@master:~$ mkdir -p $HOME/.kube
user@master:~$ sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
user@master:~$ sudo chown $(id -u):$(id -g) $HOME/.kube/config
```

## Install Calico (master only)

```bash
user@master:~$ kubectl apply -f https://docs.projectcalico.org/v3.11/manifests/calico.yaml
```

## Join master node

```bash
user@master:~$ sudo kubeadm token list  # use this token
user@master:~$ openssl x509 -pubkey -in /etc/kubernetes/pki/ca.crt | openssl rsa -pubin -outform der 2>/dev/null | \ # get hash
   openssl dgst -sha256 -hex | sed 's/^.* //'
   
user@slave:~$ sudo kubeadm join {{MASTER_IP}}:6443 --token {{TOKEN}} \
    --discovery-token-ca-cert-hash sha256:{{HASH}}s
```



### Kubernetes cluster is up and running 👑

![Alt Text](https://media.giphy.com/media/vFKqnCdLPNOKc/giphy.gif)





## Configure kubectl on slave nodes

```bash
user@slave:~$ mkdir -p $HOME/.kube
user@slave:~$ scp user@{MASTER_IP}:/home/user/.kube/config .kube/config
user@slave:~$ sudo chown $(id -u):$(id -g) $HOME/.kube/config
```



## Install Web UI

![Alt Text](https://octant.dev/docs/master/octant-demo.gif)

I suggest using Octane, because default Kubernetes Dashboard is a bit hard to configure correctly

```bash
user@master:~$ wget https://github.com/vmware-tanzu/octant/releases/download/v0.12.1/octant_0.12.1_Linux-64bit.deb
user@master:~$ dpkg -i octant_0.12.1_Linux-64bit.deb
```



Use OCTANT_LISTENER_ADDR=0.0.0.0:7777 to allow connections from any host

```bash
OCTANT_LISTENER_ADDR=0.0.0.0:7777 octant
```

If you are in production, set OCTANT_ACCEPTED_HOSTS to the value of the host you are trying to connect to Octant from, and  OCTANT_LISTENER_ADDR to host IP

```bash
OCTANT_ACCEPTED_HOSTS={EXTERNAL_IP} OCTANT_LISTENER_ADDR={SERVER_IP}:80 octant
```



## Install Helm

```bash
user@master:~$ cd /tmp
user@master:~$ curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3

user@master:~$ chmod u+x get_helm-helm.sh
user@master:~$ ./get_helm.sh
user@master:~$ helm repo add stable https://kubernetes-charts.storage.googleapis.com/
user@master:~$ helm repo update
```

