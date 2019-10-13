`grafana-annotation-resource`
-----------------------------

A resource for adding annotations to Grafana dashboards.

The resource represents a single Grafana (`version >= 6.3`) and this resource
can be used to add and update annotations using tags.

Usage
-----

In the following example, imagine we are running smoke tests from Concourse,
and we want to annotate our Grafana dashboard with a region during which the
smoke tests ran.

First define `grafana-annotation` as a resource type:

```
resource_types:
  - name: grafana-annotation
    type: docker-image
    source:
      repository: tlwr/grafana-annotation-resource
      tag: latest
```

Then declare a `grafana-annotation` resource with a descriptive name:

```
resources:
  - name: run-smoke-tests-annotation
    type: grafana-annotation
    source:
      url: http://grafana:3000
      username: admin
      password: admin
    tags:
      - run-from-concourse
      - smoke-tests
```

Then use the resource:

```
jobs:
  - name: run-smoke-tests
    plan:
      - put: run-smoke-tests-annotation
        params:
          tags:
            - started

      - task: run-smoke-tests
        config:
          platform: linux
          image_resource:
            type: docker-image
            source:
              repository: my-smoke-tests
              tag: latest
          run:
            path: /bin/smoke

      - put: run-smoke-tests-annotation
        params:
          tags:
            - finished

          # The path is the resource name, so the ID is identical to above.
          # This tells the resource that we are updating an existing annotation

          path: run-smoke-tests-annotation
```

Note: the tags in source and params are merged.

Note: when an existing annotation is updated, the tags from creation are
overwritten with the tags from the update.
