# Drone Github Extensions

An extension to improve Github/Drone integration, creating a temporary per-pipeline access tokens with access scoped to
the owner of the build job and injecting it as environment variables into the build.

![Example pipeline output](./doc/example-job-output.png)

## Installation

### Create a shared secret:

```bash
openssl rand -hex 16
```

### Download and run the plugin via docker:

```console
$ docker run -d \
  -p 3000:3000 \
  --env=DRONE_DEBUG=true \
  --env=DRONE_SECRET=<shared-secret> \
  --restart=always \
  --name=drone-github-extensions \
  ghcr.io/davidhiendl/drone-github-extensions:master
```

### Deploy the plugin to Kubernetes via Helm:

See folder [./charts](./charts)

### Multiple Plugins at once

The individual plugins can be used concurrently. Simply add all the relevant environment variables
to the agent.

## Environment Plugin Configuration

Automatically injects environment variables into the build. Requires less boilerplate than the secret plugin alternative
but creates a token and injects the environment variables into every build regardless if it is needed.

**Configuration Options**
| Key | Description | Default |
| --- | --- |  --- |
| EMULATE_CI_PREFIXED_ENV_VARS | Generate various commonly used CI_ environment variables | true |
| ENV_ADD_TAG_SEMVER | Parse tags as semver and add SEMVER_ prefix variables for convenience. | true |

Update your runner configuration to include the plugin address and the shared secret as environment variable:

```bash
DRONE_ENV_PLUGIN_ENDPOINT=http://1.2.3.4:3000/env
DRONE_ENV_PLUGIN_TOKEN==<shared-secret>
```

## Convert Plugin Configuration

Implements new directives for .drone-ci.yaml transformation

Update your drone server (IMPORTANT: this has to be added to the server, not the runners unlike the other plugins) configuration to include the plugin address and the shared secret as environment variable:

```bash
DRONE_CONVERT_PLUGIN_ENDPOINT=http://1.2.3.4:3000/convert
DRONE_CONVERT_PLUGIN_SECRET=<shared-secret>
```

### Directive _include

Allows including a remote YAML at the location of the directive.

Limitations: May only be at line start, no nested YAML keys are supported. This is due to how the directive is processed which is by text substituation in order to support all YAML features including anchors which would be difficult when parsing the YAML.

**Configuration**

| Config | Value | Default |
|-----------------------|--|---------|
| DRONE_CONFIG_INCLUDE_MAX | Number of include directives allow when processing a yaml file (including recursive includes) | 20 |

```bash
DRONE_CONFIG_INCLUDE_MAX=20
```

Use in pipelines:

**yaml to be included with _include**
```yaml
.StepTemplate: &StepTemplate
  image: alpine
  commands:
    - echo "do something"
```

**Project drone-ci.yaml**
```yaml
_include: https://yourdomain.tld/example.yaml

kind: pipeline
name: default

steps:
  - name: test
    <<: *StepTemplate
    image: ubuntu # overwrite
```

**Resulting merged yaml**
```yaml
# DIRECTIVE_START _include: https://yourdomain.tld/example.yaml
.StepTemplate: &StepTemplate
  image: alpine
  commands:
    - echo "do something"
_include: https://yourdomain.tld/example.yaml
  # DIRECTIVE_END _include: https://yourdomain.tld/example.yaml

kind: pipeline
name: default

steps:
  - name: test
    <<: *StepTemplate
    image: ubuntu # overwrite
```
