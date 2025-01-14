---
title: Executer Modules
weight: 6
categories: [architecture]
tags: [components]
description: >
    Native executer modules 
---

Executer modules are part of the SOARCA core. Executer modules perform the actual commands in CACAO playbook steps.

## Native modules in SOARCA
The following capability modules are currently defined in SOARCA:
 
- ssh
- http-api
- openc2-http

The capability will be selected based on the agent in the CACAO playbook step. The agent should be of type `soarca` and have a name corresponding to `soarca-[capability name]`. 

### SSH capability

This capability executes [SSH Commands](https://docs.oasis-open.org/cacao/security-playbooks/v2.0/cs01/security-playbooks-v2.0-cs01.html#_Toc152256500) on the specified targets.

This capability support [User Authentication](https://docs.oasis-open.org/cacao/security-playbooks/v2.0/cs01/security-playbooks-v2.0-cs01.html#_Toc152256508) using the `user-auth` type. For SSH authentication username/password is authentication supported.

#### Success and failure

The SSH step is considered successful if a proper connection to each target can be initialized, the supplied command executes without error, and returns with zero exit status.

In every other circumstance the step is considered to have failed.

#### Variables

This module does not define specific variables as input, but variable interpolation is supported in the command and target definitions. It has the following output variables:

```json
{
    "__soarca_ssh_result__": {
        "type": "string",
        "value": "<stdout of the last command>"
    }
}
```

#### Example

```json
{
    "workflow": {
        "action--7777c6b6-e275-434e-9e0b-d68f72e691c1": {
            "type": "action",
            "agent": "soarca--00010001-1000-1000-a000-000100010001",
            "targets": ["linux--c7e6af1b-9e5a-4055-adeb-26b97e1c4db7"],
            "commands": [
                {
                    "type": "ssh",
                    "command": "ls -la"
                }
            ]
        }
    },
    "agent_definitions": {
        "soarca--00010001-1000-1000-a000-000100010001": {
            "type": "soarca",
            "name": "soarca-ssh"
        }
    },
    "target_definitions": {
        "linux--c7e6af1b-9e5a-4055-adeb-26b97e1c4db7": {
            "type": "linux",
            "name": "target",
            "address": { "ipv4": ["10.0.0.1"] }
        }
    }
}
```


### HTTP-API capability

This capability implements the [HTTP API Command](https://docs.oasis-open.org/cacao/security-playbooks/v2.0/cs01/security-playbooks-v2.0-cs01.html#_Toc152256495).

Both [HTTP Basic Authentication](https://docs.oasis-open.org/cacao/security-playbooks/v2.0/cs01/security-playbooks-v2.0-cs01.html#_Toc152256506) with user_id/password and token based [OAuth2 Authentication](https://docs.oasis-open.org/cacao/security-playbooks/v2.0/cs01/security-playbooks-v2.0-cs01.html#_Toc152256507) are supported.

At this time, redirects are not supported.

#### Success and failure

The command is considered to have successfully completed if a successful HTTP response is returned from each target. An HTTP response is successful if it's response code is in the range 200-299.

#### Variables

This capability supports variable interpolation in the command, port, authentication info, and target definitions.

The result of the step is stored in the following output variables:

```json
{
    "__soarca_http_api_result__": {
        "type": "string",
        "value": "<http response body>"
    }
}
```

#### Example

```json
{
    "workflow": {
        "action--8baa7c78-751b-4de9-81d4-775806cee0fb": {
            "type": "action",
            "agent": "soarca--00020001-1000-1000-a000-000100010001",
            "targets": ["http-api--4ebae9c3-9454-4e28-b25b-0f43cd97f9e0"],
            "commands": [
                {
                    "type": "http-api",
                    "command": "GET /overview HTTP/1.1",
                    "port": "8080"
                }
            ]
        }
    },
    "agent_definitions": {
        "soarca--00020001-1000-1000-a000-000100010001": {
            "type": "soarca",
            "name": "soarca-http-api"
        }
    },
    "target_definitions": {
        "http-api--4ebae9c3-9454-4e28-b25b-0f43cd97f9e0": {
            "type": "http-api",
            "name": "target",
            "address": { "dname": ["my.server.com"] }
        }
    }
}
```

### OpenC2 capability

This capability implements the [OpenC2 HTTP Command](https://docs.oasis-open.org/cacao/security-playbooks/v2.0/cs01/security-playbooks-v2.0-cs01.html#_Toc152256498), by sending [OpenC2 messages](https://docs.oasis-open.org/openc2/oc2ls/v1.0/cs01/oc2ls-v1.0-cs01.html)  using the [HTTPS transport method](https://docs.oasis-open.org/openc2/open-impl-https/v1.0/open-impl-https-v1.0.html).

It supports the same authentication mechanisms as the HTTP-API capability.

#### Success and failure

Any successful HTTP response from an OpenC2 compliant endpoint (with a status code in the range 200-299) is considered a success. Connection failures and HTTP responses outside the 200-299 range are considered a failure.

#### Variables

It supports variable interpolation in the command, headers, and target definitions.

The result of the step is stored in the following output variables:

```json
{
    "__soarca_openc2_http_result__": {
        "type": "string",
        "value": "<openc2-http response body>"
    }
}
```

#### Example

```json
{
    "workflow": {
        "action--aa1470d8-57cc-4164-ae07-05745bef24f4": {
            "type": "action",
            "agent": "soarca--00030001-1000-1000-a000-000100010001",
            "targets": ["http-api--5a274b6d-dc65-41f7-987e-9717a7941876"],
            "commands": [{
                "type": "openc2-http",
                "command": "POST /openc2-api/ HTTP/1.1",
                "content_b64": "ewogICJoZWFkZXJzIjogewogICAgInJlcXVlc3RfaWQiOiAiZDFhYzA0ODktZWQ1MS00MzQ1LTkxNzUtZjMwNzhmMzBhZmU1IiwKICAgICJjcmVhdGVkIjogMTU0NTI1NzcwMDAwMCwKICAgICJmcm9tIjogInNvYXJjYS5ydW5uZXIubmV0IiwKICAgICJ0byI6IFsKICAgICAgImZpcmV3YWxsLmFwaS5jb20iCiAgICBdCiAgfSwKICAiYm9keSI6IHsKICAgICJvcGVuYzIiOiB7CiAgICAgICJyZXF1ZXN0IjogewogICAgICAgICJhY3Rpb24iOiAiZGVueSIsCiAgICAgICAgInRhcmdldCI6IHsKICAgICAgICAgICJmaWxlIjogewogICAgICAgICAgICAiaGFzaGVzIjogewogICAgICAgICAgICAgICJzaGEyNTYiOiAiMjJmZTcyYTM0ZjAwNmVhNjdkMjZiYjcwMDRlMmI2OTQxYjVjMzk1M2Q0M2FlN2VjMjRkNDFiMWE5MjhhNjk3MyIKICAgICAgICAgICAgfQogICAgICAgICAgfQogICAgICAgIH0KICAgICAgfQogICAgfQogIH0KfQ==",
                "headers": {
                    "Content-Type": ["application/openc2+json;version=1.0"]
                }
            }]
        }
    },
    "agent_definitions": {
        "soarca--00030001-1000-1000-a000-000100010001": {
            "type": "soarca",
            "name": "soarca-openc2-http"
        }
    },
    "target_definitions": {
        "http-api--5a274b6d-dc65-41f7-987e-9717a7941876": {
            "type": "http-api",
            "name": "openc2-compliant actuator",
            "address": { "ipv4": ["187.0.2.12"] }
        }
    }
}
```

---

## MQTT fin module
This module is used by SOARCA to communicate with fins (capabilities) see [fin documentation](/docs/soarca-extensions/) for more information
