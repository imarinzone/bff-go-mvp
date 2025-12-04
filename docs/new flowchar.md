flowchart TD
  A[Client / UI / Mobile] -->|HTTP/gRPC| B(Go API Router)
  B -->|Start Workflow / Signal| C[Temporal Orchestration -Java]
  C --> D{Decision / Activity}
  D -->|call| E[Downstream Service 1 - HTTP/gRPC]
  D -->|call| F[Downstream Service 2 - DB / Cache]
  D -->|call| G[Third-party API]
  E --> H[Activity Result]
  F --> H
  G --> H
  H --> C
  C -->|Workflow Completed / Result| B
  B -->|HTTP Response| A

  subgraph Observability
    I[Logs / Traces]---C
    I---B
    J[Metrics & Alerts]---C
    J---B
  end

  subgraph Temporal
    C
  end


sequenceDiagram
    participant Client
    participant GoAPI as Go API Router
    participant Temporal as Temporal (Java)
    participant Activity1 as DownstreamSvc1
    participant Activity2 as DownstreamSvc2

    Client->>GoAPI: POST /process {payload}
    GoAPI->>Temporal: StartWorkflow(processRequest)
    Note right of Temporal: Workflow receives request
    Temporal->>Activity1: ExecuteActivity(callService1)
    Activity1-->>Temporal: 200 {data}
    Temporal->>Activity2: ExecuteActivity(queryDB)
    Activity2-->>Temporal: 200 {data}
    Temporal-->>GoAPI: WorkflowCompleted({aggregateResult})
    GoAPI-->>Client: 200 {aggregateResult}