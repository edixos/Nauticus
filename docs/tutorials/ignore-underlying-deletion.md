# Ignore Underling Deletion

In Nauticus, the `nauticus.io/ignore-underlying-deletion annotation allows you to delete a Space without deleting the underlying resources that were created by the Nauticus controller.

## Overview

By default, when a Space is deleted in Nauticus, all the underlying resources (including namespaces, quotas, and other objects created by Nauticus controller) are also deleted. This is usually the desired behavior to avoid orphan resources which can consume system resources.

However, there may be cases where you want to delete a Space, but retain the underlying resources. This can be useful, for example, for debugging purposes, or when you want to move resources from one Space to another.

## How to Use

To delete a Space without deleting the underlying resources, you simply add the `nauticus.io/ignore-underlying-deletion` annotation to the Space manifest, and set its value to `true`.

```yaml title="config/samples/space_with_annotation.yaml"
{% include "../../config/samples/space_with_annotation.yaml" %}
```

```bash title="Create the Space with the annotation"
kubectl apply -f config/samples/space_with_annotation.yaml
```

With this annotation in place, when you delete the Space, Nauticus controller will ignore the underlying resources, leaving them intact.

!!! Warning "Important Note"
    Please note that using this annotation will leave behind resources that are no longer managed by a Space. These resources will continue to consume cluster resources and may need to be manually cleaned up when no longer needed. Use this annotation with caution and understand the implications on your resource usage and management.