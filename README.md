# cilium-nodeport-deleter

> [!WARNING]
> This repository has been archived

As of cilium 1.18, the API this project used to delete nodeports is [no longer available](https://github.com/cilium/cilium/issues/12781#issuecomment-3180128433). Following this upstream change, together with us adopting a different workaround (see below), we have made the decision to stop maintaining this project.

## Post-archival workaround

We are currently using [Mutating Admission Policies](https://kubernetes.io/docs/reference/access-authn-authz/mutating-admission-policy/) to enforce `allocateLoadBalancerNodePorts=true` on LoadBalancer type services, which achieves the same goal as deleting them after creation:

```yaml
---
apiVersion: admissionregistration.k8s.io/v1alpha1
kind: MutatingAdmissionPolicyBinding
metadata:
  name: disable-allocate-lb-nodeports
spec:
  policyName: disable-allocate-lb-nodeports
---
# This policy sets spec.allocateLoadBalancerNodePorts to false for services of type LoadBalancer
# Ref: https://kubernetes.io/docs/concepts/services-networking/service/#load-balancer-nodeport-allocation
# Ref: https://kubernetes.io/docs/reference/access-authn-authz/mutating-admission-policy/
apiVersion: admissionregistration.k8s.io/v1alpha1
kind: MutatingAdmissionPolicy
metadata:
  name: disable-allocate-lb-nodeports
spec:
  matchConstraints:
    resourceRules:
    - apiGroups:   [""]
      apiVersions: ["v1"]
      operations:  ["CREATE"]
      resources:   ["services"]
  matchConditions:
    - name: type-load-balancer
      expression: |-
        object.spec.type == "LoadBalancer"
    - name: not-allow-annotation
      # svc.tenshi.es/allow-nodeports
      # There might be a way to do this without escaping, but I haven't found it yet.
      expression: |-
        !has(object.metadata.annotations.svc__dot__tenshi__dot__es__slash__allow__dash__nodeports)
  failurePolicy: Fail
  reinvocationPolicy: IfNeeded
  mutations:
    - patchType: "ApplyConfiguration"
      applyConfiguration:
        expression: >
          Object{
            spec: Object.spec{
              allocateLoadBalancerNodePorts: false
            }
          }
```

CEL expressions for admission policies are, at the time of writing (k8s 1.33), in alpha state, and can be enabled by:
1. Setting the `MutatingAdmissionPolicy` feature gate to `true` in the API Server and Controller Manager (e.g. via `--feature-gates`)
1. Adding `admissionregistration.k8s.io/v1alpha1=true` to the API Server's `--runtime-config`.

## Original README

As of today a CiliumClusterwideNetworkPolicy (CCNP) blocks all traffic that reached a node
which is not managed by Kubernetes but all Kubernetes managed traffic flows (like LoadBalancers
and NodePorts) remain open.

By legacy a service which type is LoadBalancer will always open a NodePort. This is an expected
behaviour but (from my perspective) Cilium should provide a way to block or filter this traffic
the same as you can filter or block all traffic from any other Kubernetes managed endpoint.

There is an issue to follow this feature: https://github.com/cilium/cilium/issues/12781
