# SOARCA controller

The SOARCA controller will control the incoming calls and the execution of steps. 

```plantuml
interface IWorkflow{
    void Get()
    void Get(PlaybookId id)
    void Add(Playbook workflow)
    void Update(Playbook workflow)
    void Remove(Playbook workflow)
}

interface ICoa{
    void GetAll()
    void Get(CoaId id)
    void Add(Playbook coa)
    void Update(Playbook coa)
    void Remove(Playbook coa)
    
}

interface IStatus{

}

interface ITrigger{
    void LoadAndStartPlaybook(PlaybookId id)
}

Interface IWorkflowDatabase

Interface IDecomposer
Interface IExecuter

class Controller
class Decomposer

IWorkflow <|.. Controller
ICoa <|.. Controller
ITrigger <|.. Controller
IStatus <|.. Controller

Controller -> IWorkflowDatabase
IWorkflowDatabase <|.. WorkflowDatabase


IDecomposer <- Controller
IDecomposer <|.. Decomposer
IExecuter -> Decomposer

```

## Main application flow
