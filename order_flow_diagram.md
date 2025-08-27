# Order Flow Diagrams

## Current Order Processing Flow

```mermaid
flowchart TD
    A[Client] --> B[Create Order]
    B --> C[Save to DB with NEW status]
    C --> D[Return Order ID]
    
    E[Client] --> F[Get Orders]
    F --> G[Query DB]
    G --> H[Return Orders List]
    
    I[Client] --> J[Get Balance]
    J --> K[Query DB for Balance]
    K --> L[Return Balance]
    
    M[Client] --> N[Withdraw]
    N --> O[Check Balance]
    O --> P{Sufficient Balance?}
    P -->|Yes| Q[Save Withdrawal]
    P -->|No| R[Return Error]
    Q --> S[Update Balance]
```

## Order Processing with Goroutine Pools

```mermaid
flowchart TD
    subgraph "Order Checker"
        A[Check DB for NEW orders]
        B[Mark as PROCESSING]
        C[Add to channel]
    end
    
    subgraph "Database"
        DB[(Orders Table)]
    end
    
    subgraph "Channel"
        CH[Order Channel]
    end
    
    subgraph "Accrual Poller"
        D[Get order from channel]
        E[Poll Accrual Service]
        F[Update final status]
    end
    
    subgraph "External"
        ACC[Accrual Service]
    end
    
    A -->|Find NEW| DB
    B -->|Mark PROCESSING| DB
    C -->|Add to queue| CH
    
    CH -->|Get order| D
    D --> E
    E -->|Max 5 min| ACC
    ACC -->|Response| E
    E --> F
    F -->|Update status| DB
    
    subgraph "Recovery"
        R[Service Restart]
        S[Find PROCESSING orders]
        T[Re-add to channel]
    end
    
    R --> S
    S -->|Scan DB| DB
    T --> CH
    
    style A fill:#e1f5fe
    style D fill:#f3e5f5
    style CH fill:#fff3e0
    style ACC fill:#e8f5e8
    style R fill:#ffebee
```

## Key Features

1. **Order Checker**:
    - Checks DB for `NEW` orders
    - Marks them as `PROCESSING`
    - Adds to channel queue

2. **Accrual Poller**:
    - Takes orders from channel
    - Polls accrual service (max 5 min)
    - Updates final status

3. **Recovery**:
    - On restart, finds `PROCESSING` orders
    - Re-adds them to channel
