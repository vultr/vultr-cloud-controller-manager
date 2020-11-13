# Load Balancers

## Overview

Kubernetes Services type `LoadBalancer` will be deployed through [Vultr Load Balancers](https://www.vultr.com/products/load-balancers/). This is provisioned by the Cloud Controller Manager. For generic info and faq please visit the [Vultr LoadBalancer Doc](https://www.vultr.com/docs/vultr-load-balancers).

Examples of `LoadBalancer` resources can be found [here](examples) 

## Annotations

The Vultr CCM allows you to configure your `LoadBalancer` resource to be deployed in a specific way through annotations.

All of the annotations are listed below. Please note that all annotations listed below **must** be prepended with `service.beta.kubernetes.io/vultr-loadbalancer-` and are case sensitive.

Annotation (Suffix) | Values | Default | Description
---|---|---|---
`protocol` | `tcp`, `http` | `tcp` | This is used to specify the protocol to be used for your LoadBalancer protocol.
`https-ports` | string | | Defines which ports should be used for HTTPS. You can pass in a comma separated list: 443,8443
`ssl` | string | | The string you provide should be the name of a Kubernetes TLS Secret which store your cert + key
`ssl-pass-through` | `true`, `false` | `false` | If you want SSL termination to happen on your `pods` or `ingress` then this must be enabled. This is to be used with the `https-ports` annotation
`healthcheck-protocol` | `tcp` `http` | `tcp` | The protocol to be used for your LoadBalancer HealthCheck
`healthcheck-path` | string | `/` | The URL path to check on the back-end during health checks
`healthcheck-port` | int | `defaults to what kubernetes defines` | The port that should be called for health checks
`healthcheck-check-interval` | int | `15` | Interval between health checks (in seconds)
`healthcheck-response-timeout` | int | `5` | Response timeout (in seconds)
`healthcheck-unhealthy-threshold` | int | `5` | Number of unhealthy requests before a back-end is removed
`healthcheck-healthy-threshold` | int | `5` | Number of healthy requests before a back-end is added back in
`algorithm` | `least_connections`, `roundrobin` | `roundrobin` | Balancing algorithm 
`ssl-redirect` | `true`, `false`| `false` | Force HTTP to HTTPS
`sticky-session-enabled` | `on`, `off`| `off` | Enables Sticky Sessions. If enabled you must provide `sticky-session-cookie-name`
`sticky-session-cookie-name"` | string |  | Name of sticky session