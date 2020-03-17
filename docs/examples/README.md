# Examples

This directory has example manifests which will help you install and use the CCM with your cluster.

- [cloud-controller-manager.yml](cloud-controller-manager.yml) - All resources required for the Vultr Cloud Controller Manager
- [load-balancer-https.yml](load-balancer-https.yml) -  Creates a Vultr LoadBalancer that listens on 80 and 443. Note this will require you to create your own [kubernetes tls secret](https://kubernetes.io/docs/concepts/services-networking/ingress/#tls)