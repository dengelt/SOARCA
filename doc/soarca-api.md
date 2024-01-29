# SOAR-CA API 
This document describes the SOAR-CA API 


SOAR-CA is developed by TNO Cybersecurity & Technology department (CST)


## Endpoint description

We will use HTTP status codes https://en.wikipedia.org/wiki/List_of_HTTP_status_codes


```plantuml
@startuml
protocol UiEndpoint {
    GET     /workflow
    POST    /workflow
    GET     /workflow/workflow-id
    PUT     /workflow/workflow-id
    DELETE  /workflow/workflow-id


    POST    /trigger/workflow
    POST    /trigger/workflow/id

    GET     /step

    GET     /status
    GET     /status/workflow
    GET     /status/workflow/id
    GET     /status/history

}
@enduml
```




### General messages

#### Error
When an error occurs a 400 status is returned with the following JSON payload, the original call can be omitted in production for security reasons.

responses: 400/Bad request

```plantuml
@startjson
{
    "status": "400",
    "message": "What went wrong.",
    "original-call": "<optional> Request JSON data",
    "downstream-call" : "<optional> downstream call JSON"
}
@endjson
```

#### Unauthorized
When the caller does not have valid authentication 401/unauthorized will be returned.


#### cacao playbook JSON

```plantuml
@startjson
{
            "type": "playbook",
            "spec_version": "1.1",
            "id": "playbook--91220064-3c6f-4b58-99e9-196e64f9bde7",
            "name": "coa flow",
            "description": "This playbook will trigger a specific coa",
            "playbook_types": ["notification"],
            "created_by": "identity--06d8f218-f4e9-4f9f-9108-501de03d419f",
            "created": "2020-03-04T15:56:00.123456Z",
            "modified": "2020-03-04T15:56:00.123456Z",
            "revoked": false,
            "valid_from": "2020-03-04T15:56:00.123456Z",
            "valid_until": "2020-07-31T23:59:59.999999Z",
            "derived_from": [],
            "priority": 1,
            "severity": 1,
            "impact": 1,
            "industry_sectors": ["information-communications-technology", "research", "non-profit"],
            "labels": ["soarca"],
            "external_references": [
                {
                    "name": "TNO SOARCA",
                    "description": "SOARCA Homepage",
                    "source": "TNO - COSSAS - HxxPS://LINK-TO-CODE-REPO.TLD",
                    "url": "HxxPS://LINK-TO-CODE-REPO.TLD",
                    "hash": "00000000000000000000000000000000000000000000000000000000000",
                    "external_id": "TNO/SOARCA 2023.01"
                }
            ],
            "features": {
                "if_logic": true,
                "data_markings": false
            },
            "markings": [],
            "playbook_variables": {
                "$$flow_data_location$$": {
                    "type": "string",
                    "value": "<mongodb_location>",
                    "description": "location of event and flow data",
                    "constant": true
                },
                "$$event_type$$": {
                    "type" : "string",
                    "value": "<event_type_string>",
                    "description": "type of incomming event / trigger",
                    "constant": true	
                }
            },
            "workflow_start": "step--d737c35f-595e-4abf-83ef-d0b6793556b9",
            "workflow_exception": "step--40131926-89e9-44df-a018-5f92f2df7914",
            "workflow": {
                "step--5ea28f63-ac32-4e5e-bd0c-757a50a3a0d7":{
                    "type": "single",
                    "name": "BI for CoAs",
                    "delay": 0,
                    "timeout": 30000,
                    "command": {
                        "type": "http-api",
                        "command": "hxxps://our.bi/key=VALUE"
                    },
                    "on_success": "step--71b15428-275a-49b5-9f09-3944972a0054",
                    "on_failure": "step--71b15428-275a-49b5-9f09-3944972a0054"
                },
                "step--71b15428-275a-49b5-9f09-3944972a0054": {
                    "type": "end",
                    "name": "End Playbook SOARCA Main Flow"
                }
            },
            "targets": { 

            },
            "extension_definitions": { }
        }
@endjson
```
---- 

### /workflow
The workflow endpoinst are used to create workflows in SOAR-CA, new playbook can be added, current ones edited and deleted. 

#### GET `/workflow`
Get all workflow id's that are currently stored in SOAR-CA.

##### Call payload
None

##### Response
200/OK with payload:

```plantuml
@startjson
{
    "workflows": [
        {
            "workflow-id": "xxxxxxxx-xxxx-Mxxx-Nxxx-xxxxxxxxxxxx",
            "workflow-name": "name",
            "workflow-description": "description"
        }
    ]
}
@endjson
```

##### Error
400/BAD REQUEST with payload:
General error


#### POST `/workflow`
Create a new workflow that and store it in SOAR-CA. The format is 


##### Payload

```plantuml
@startjson
{
    "workflow": "<cacao-playbook> (json)"
}
@endjson
```



##### Response
201/CREATED

```plantuml
@startjson
{
    "workflow-id": "xxxxxxxx-xxxx-Mxxx-Nxxx-xxxxxxxxxxxx"
}
@endjson
```

##### Error
400/BAD REQUEST with payload: General error, 409/CONFLICT if the entry already exists


#### GET `/workflow/{workflow-id}`
Get workflow details

##### Call payload
None

##### Response
200/OK with payload:

```plantuml
@startjson
{
    "workflow-id": "xxxxxxxx-xxxx-Mxxx-Nxxx-xxxxxxxxxxxx",
    "workflow": "<cacao-playbook> (json)"
    
}
@endjson
```
##### Error
400/BAD REQUEST

----

#### PUT `/workflow/{workflow-id}``
An existing workflow can be updated with PUT. 

##### Call payload
A playbook like <<cacao playbook JSON>>


##### Response
200/OK with the edited playbook <<cacao playbook JSON>>

##### Error
400/BAD REQUEST for malformed request

When updated it will return 200/OK or General error in case of an error.

----


#### DELETE `/workflow/{workflow-id}`
An existing workflow can be deleted with DELETE. When removed it will return 200/OK or general error in case of an error.

##### Call payload
None

##### Response
200/OK if deleted

##### Error
400/BAD REQUEST if resource does not exist

---

#### POST `/trigger/workflow/xxxxxxxx-xxxx-Mxxx-Nxxx-xxxxxxxxxxxx` 
Execute workflow with specific id

##### Call payload
None

##### Response
Will return 200/OK when finished with workflow playbook.

```plantuml
@startjson
{
    "execution-id": "xxxxxxxx-xxxx-Mxxx-Nxxx-xxxxxxxxxxxx",
    "playbook-id": "xxxxxxxx-xxxx-Mxxx-Nxxx-xxxxxxxxxxxx"
}
@endjson
```

##### Error
400/BAD REQUEST general error on error.

---

#### POST `/trigger/workflow`
Execute an adhoc playbook

##### Call payload
A playbook like <<cacao playbook JSON>>

##### Response
Will return 200/OK when finished with playbook.

```plantuml
@startjson
{
    "execution-id": "xxxxxxxx-xxxx-Mxxx-Nxxx-xxxxxxxxxxxx",
    "playbook-id": "xxxxxxxx-xxxx-Mxxx-Nxxx-xxxxxxxxxxxx"
}
@endjson
```

##### Error
400/BAD REQUEST general error on error.

----

### /step
Get capable steps for SOARCA to allow a coa builder to generate or build valid coa's

#### GET `/step`
Get all available steps for SOARCA. 

##### Call payload
None

##### Response
200/OK


```plantuml
@startjson
{
    
    "steps": [{
        "module": "executor-module",
        "category" : "analyses",
        "context" : "external",
        "step--5ea28f63-ac32-4e5e-bd0c-757a50a3a0d7":{
                    "type": "single",
                    "name": "BI for CoAs",
                    "delay": 0,
                    "timeout": 30000,
                    "command": {
                        "type": "http-api",
                        "command": "hxxps://our.bi/key=VALUE"
                    },
                    "on_success": "step--71b15428-275a-49b5-9f09-3944972a0054",
                    "on_failure": "step--71b15428-275a-49b5-9f09-3944972a0054"
                }}]
}
@endjson
```

Module is the executing module name that will do the executer call.

Category defines what kind of step is executed:
```plantuml
@startuml
enum workflowType {
    analyses
    action
    asset-look-up
    etc...
}
@enduml
```
Context will define whether the call is internal or external:

```plantuml
@startuml
enum workflowType {
    internal
    external
}
@enduml
```

##### Error
400/BAD REQUEST general error on error.

----

### /status
The status endpoints are used to get various statuses. 

#### GET `/status`
Call this endpoint to see if SOAR-CA is up and ready. This call has no payload body.

##### Call payload
None

##### Response
200/OK

```plantuml
@startjson
{
    "version": "1.0.0",
    "components": [
        {
            "name": "Component name",
            "status": "ready/running/failed/stopped/...",
            "message": "Some message",
            "version": "semver verison: 1.0.0"
        }
    ]
}
@endjson
```

##### Error
5XX/Internal error, 500/503/504 message.

---

#### GET /status/workflow

##### Call payload

##### Response
200/OK

```plantuml
@startjson
{
            "workflows": [
                {"type": "playbook",
                "spec_version": "1.1",
                "id": "playbook--91220064-3c6f-4b58-99e9-196e64f9bde7",
                "name": "SOARCA Main Flow",
                "description": "This playbook will run for each trigger event in SOARCA",
                "playbook_types": ["notification"],
                "created_by": "identity--06d8f218-f4e9-4f9f-9108-501de03d419f",
                "created": "2020-03-04T15:56:00.123456Z",
                "modified": "2020-03-04T15:56:00.123456Z",
                "revoked": false,
                "valid_from": "2020-03-04T15:56:00.123456Z",
                "valid_until": "2020-07-31T23:59:59.999999Z",
                "derived_from": [],
                "priority": 1,
                "severity": 1,
                "impact": 1}
            ]

}
@endjson
```

##### Error
400/BAD REQUEST general error on error.

---- 

#### GET `/status/workflow/xxxxxxxx-xxxx-Mxxx-Nxxx-xxxxxxxxxxxx (workflow-id)`
Get workflow details which is running which will return cacao playbook JSON

##### Call payload
None

##### Response
200/OK

See <<cacao playbook JSON>>
Empty payload if no workflows are running

##### Error
400/BAD REQUEST general error on error.

----

#### GET `/status/workflow/xxxxxxxx-xxxx-Mxxx-Nxxx-xxxxxxxxxxxx/coas`
Get workflow details which is running which will return cacao playbook JSON

##### Call payload
None

##### Response
200/OK

list of caoids

##### Error
400/BAD REQUEST general error on error.

---

#### GET `/status/workflow/{playbook-id}`
Get coarse of action list for coa awaiting action.

##### Call payload
None

##### Response
200/OK

```plantuml
@startjson
{
            "actions": [
                {
                    "playbook_id": "playbook--91220064-3c6f-4b58-99e9-196e64f9bde7",
                    "status": "running/finished/failed/stopped/paused"
                }
            ]

}
@endjson
```

##### Error
400/BAD REQUEST general error on error.

---

#### GET /status/history
Get all workflow and coarse of action id's and statuses that have been run excluded those that are running or paused.

##### Call payload
None

##### Response
200/OK

```plantuml
@startjson
{
            "actions": [
                {
                    "playbook_id": "playbook--91220064-3c6f-4b58-99e9-196e64f9bde7",
                    "status": "running/finished/failed/stopped/paused"
                }
            ]

}
@endjson
```

##### Error
400/BAD REQUEST general error on error.


## Usage example flow

### Stand alone

```plantuml
@startuml
participant "SWAGGER" as gui
control "SOAR-CA API" as api
control "controller" as controller
control "Executor" as exe
control "SSH-module" as ssh


gui -> api : /trigger/workflow with playbook body
api -> controller : execute playbook playload
controller -> exe : execute playbook
exe -> ssh : get url from log
exe <-- ssh : return result
controller <-- exe : results
api <-- controller: results

@enduml
```

### Database load and execution

```plantuml
@startuml
participant "SWAGGER" as gui
control "SOAR-CA API" as api
control "controller" as controller
database "Mongo" as db
control "Executor" as exe
control "SSH-module" as ssh


gui -> api : /trigger/workflow/playbook--91220064-3c6f-4b58-99e9-196e64f9bde7
api -> controller : load playbook from database
controller -> db: retreive playbook
controller <-- db: playbook json
controller -> controller: validate playbook
controller -> exe : execute playbook
exe -> ssh : get url from log
exe <-- ssh : return result
controller <-- exe : results
api <-- controller: results

@enduml
```