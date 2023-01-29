## Introduction

Welcome to the Nauticus tutorial page! Here, you will learn about the powerful features of Nauticus, a Kubernetes controller that simplifies the management of spaces within your cluster. Spaces are fully-managed namespaces that integrate with RBAC, network policies, resources, and quotas. With Nauticus, you can easily create, update, and delete spaces, as well as manage additional resources such as service accounts and limit ranges. This tutorial will guide you through the process of using Nauticus to manage your spaces and resources, and provide examples of how to utilize its features to improve your workflow. Let's get started!

## Specify Space's owner(s)

In the Nauticus controller, you have the option to specify the owner(s) of a space. This allows you to assign specific users or teams to be responsible for managing the resources within a particular space. This feature is particularly useful in a multi-tenant environment, where different teams may be responsible for different parts of the application.

To specify the owner(s) of a space, you can add the "owners" field in the space specification. The field should contain a list of email addresses of the users or teams that will be designated as the owner(s) of the space. For example, if you want to assign the space to the team `dev-team@example.com`, the space specification would look like this:

```yaml title="config/samples/space_with_owners.yaml"
{% include "../config/samples/space_with_owners.yaml" %}
```

```bash title="Create Space with Resource Quota"
kubectl apply -f config/samples/space_with_owners.yaml
```

## Space with Resource Quota

This section is dedicated to showcasing how Nauticus can be used to assign resource quotas to underlying namespaces. This feature is particularly useful for ensuring that resources are being utilized efficiently within a namespace, and can also be used to prevent overconsumption of resources.

To begin, you will need to create a new Space resource and specify the desired resource quotas within the resource definition. For example, you may want to set a limit on the number of pods that can be created within the namespace, or limit the amount of CPU and memory that can be consumed.

Here his an example of a Space with resource quota specification:

```yaml title="config/samples/space_with_resource_quotas.yaml"
{% include "../config/samples/space_with_resource_quotas.yaml" %}
```


Once you have defined the resource quotas, you can use the kubectl apply command to create the Space and apply the resource quotas to the underlying namespace.

```bash title="Create Space with Resource Quota"
kubectl apply -f config/samples/space_with_resource_quotas.yaml
```

You can also update or remove the resource quotas by updating or delete the space resource. The update will be automatically propagate to the underlying namespace.


## Space with Network Policies

__Network policies__ are an important aspect of securing and isolating resources within a Kubernetes cluster. Nauticus allows you to easily create and manage network policies for your spaces.

When you create a space, you have the option to specify a `enableDefaultStrictMode` parameter. When set to true, Nauticus will create a default network policy that restricts ingress communication from other spaces or namespaces, except for namespaces that have the label `nauticus.io/role: system`. This default policy will help secure your space and prevent unauthorized access to your resources.

To create a network policy for your space, you can use the kubectl command line tool and specify the desired rules in a manifest file. Here is an example of a manifest file that creates a network policy in underlying space's namespace.

```yaml title="config/samples/space_with_network_policy.yaml"
{% include "../config/samples/space_with_network_policy.yaml" %}
```

```bash title="Create Space with Network Policy"
kubectl apply -f config/samples/space_with_network_policy.yaml
```

## Assign Additional RoleBindings
One of the features of Nauticus is the ability to assign additional role bindings to a space. This allows you to grant specific roles and permissions to users or service accounts within the space.

To assign additional role bindings, you can include the `additionalRoleBindings` field in the space specification. This field should contain a list of objects, each with `subjects` and `roleRef` field.

```yaml title="config/samples/space_with_additional_rolebindings.yaml"
{% include "../config/samples/space_with_additional_rolebindings.yaml" %}
```

```bash title="Create Space with Network Policy"
kubectl apply -f config/samples/space_with_additional_rolebindings.yaml
```

## Assign Limit Ranges

The limit range feature in Nauticus allows users to set constraints on the resources that can be requested and consumed by the containers in a namespace. This feature provides an additional layer of control over resource utilization in the cluster, ensuring that high resource-intensive workloads do not affect the performance of other services. To utilize this feature, users can specify limit ranges in their `Space` configuration and Nauticus will enforce these limits at runtime. This provides a simple and effective way to manage the resource consumption in a multi-tenant cluster environment.

```yaml title="config/samples/space_with_limit_ranges.yaml"
{% include "../config/samples/space_with_limit_ranges.yaml" %}
```

```bash title="Create Space with Network Policy"
kubectl apply -f config/samples/space_with_limit_ranges.yaml
```

## Space Example with All Features Combined Together

In this example, we will demonstrate how to create a Space in Nauticus that combines all the features discussed so far: resource quota, network policy, space owners, and additional role bindings.

```yaml title="config/samples/space_with_all.yaml"
{% include "../config/samples/space_with_all.yaml" %}
```

```bash title="Create Space with all features"
kubectl apply -f config/samples/space_with_all.yaml
```


After creating the Space resource, you can verify the resource quota and network policy by checking the respective resources in the namespace created for this Space. You can also check the role bindings and owners by using the 

```bash title="Retreive Space informations"
kubectl get rolebindings,networkpolicies,resourcequotas,limitranges -n space-all-in-one 
kubectl describe space space-all-in-one 
```

This Space example demonstrates how Nauticus enables you to easily manage and control your Kubernetes resources by combining various features in a single resource.