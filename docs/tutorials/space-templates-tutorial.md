## Introduction

The current process for creating Spaces with Nauticus, while straightforward, often involves repetitive configuration of the same parameters. This redundancy can hinder efficiency and result in time-consuming manual work. To address these challenges, the Space Templates feature has been developed. It allows administrators to create predefined templates that encompass common configurations, making Space creation more efficient and standardized.

## Space Templates

Space Templates are pre-configured templates that administrators can create to define common settings and configurations. These templates can include resource quotas, network policies, role bindings, and more. By referencing a Space Template during Space creation, users can take advantage of these predefined settings.

## Usage


### SpaceTemplate Reference example
In this example, we will demonstrate how to create a SpaceTemplate in Nauticus that combines all the features so far: resource quota, network policy, Limit ranges and additional role bindings
```yaml title="config/samples/space_template_with_all.yaml"
{% include "../../config/samples/space_template_with_all.yaml" %}
```

To utilize Space Templates, users can reference a Space Template in their Space resource during creation. Any specifications provided in the Space resource will override the corresponding parameters in the Space Template. Here's an example of how this works:

```yaml title="config/samples/space_with_template_ref.yaml"
{% include "../../config/samples/space_with_template_ref.yaml" %}
```

### Overrides
One of the key features of Space Templates is the ability to override specific configurations when referencing a template in a Space. Users can selectively modify parameters to match the requirements of their Space while still benefiting from the predefined template.

By understanding the concept of overrides, administrators can take full advantage of Space Templates and achieve a balance between standardization and flexibility when configuring Spaces.
#### Example: Overriding ResourceQuotas
In this example, a Space named `space-tpl-ref-override` references a SpaceTemplate called `space-tpl-sample`. While leveraging predefined configurations from the template, it overrides the default resource quotas with custom values.

* The `templateRef` references the `space-tpl-sample` SpaceTemplate.
* This demonstrates the flexibility of SpaceTemplates, allowing Spaces to maintain standard configurations while adjusting specific settings as needed.

    ```yaml title="config/samples/space_wth_tpl_ref_overrides.yaml"
    {% include "../../config/samples/space_wth_tpl_ref_overrides.yaml" %}
    ```
#### Example: Merging Additional Role Bindings

In this example, a Space named `space-tpl-ref-merge` references a Space Template called `space-tpl-sample`. It merges additional role bindings with predefined configurations from the template. This demonstrates the capability to combine and customize various settings while maintaining consistency.

* The templateRef references the `space-tpl-sample` SpaceTemplate.
* The Space includes additional role bindings for both `viewer` and `editor` roles, with specific subjects.
* The merged role bindings enrich the Space's access control settings, ensuring flexibility and control.
    ```yaml title="config/samples/space_wth_tpl_merge.yaml"
    {% include "../../config/samples/space_wth_tpl_merge.yaml" %}
    ```