# WiseLabz — v2 Backlog

Features explicitly deferred from v1 during frontend planning. Each replaces or
extends a v1 decision (see the frontend plan, §8).

## Deferred to v2

| Feature                                             | Context                                                                                                                                                                 |
|-----------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| SSE endpoint for AI suggestions                     | Alternative to WebSocket streaming if the `/ws` channel becomes too complex to multiplex. Evaluate after the v1 AI module is stable (decision §8.4).                    |

> **Shipped since:** doc edit presence and soft locks ("X is currently editing")
> were delivered under #92.
