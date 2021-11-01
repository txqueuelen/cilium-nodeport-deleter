# cilium-nodeport-deleter

As of today a CiliumClusterwideNetworkPolicy (CCNP) blocks all traffic that reached a node
which is not managed by Kubernetes but all Kubernetes managed traffic flows (like LoadBalancers
and NodePorts) remain open.

By legacy a service which type is LoadBalancer will always open a NodePort. This is an expected
behaviour but (from my perspective) Cilium should provide a way to block or filter this traffic
the same as you can filter or block all traffic from any other Kubernetes managed endpoint.

There is an issue to follow this feature: https://github.com/cilium/cilium/issues/12781
