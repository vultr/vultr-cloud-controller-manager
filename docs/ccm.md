# Cloud Controller Manager

The Vultr Cloud Controller manager implements the provided interfaces laid out by the Kubernetes CCM guidelines. The CCM allows kubernetes to communicate with Vultr as a first class citizen. Here are a few highlights of what the CCM manages.

- Node resources are assigned their respective Vultr instance hostnames, Region, PlanID and public/private IPs.
- Node resources get put into their proper state if they are shutdown or removed. This allows for Kubernetes to properly reschedule pods
- Vultr LoadBalancers are automatically deployed when a LoadBalancer service is deployed.

More information about the cloud controller manager can be found here
- [Concepts Underlying the Cloud Controller Manager](https://kubernetes.io/docs/concepts/architecture/cloud-controller/)
- [Developing Cloud Controller Manager](https://kubernetes.io/docs/tasks/administer-cluster/developing-cloud-controller-manager/)