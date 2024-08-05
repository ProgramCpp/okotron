## database for persistence

### requirements
- high update rate for certain commands
    - ex: SWAP, LIMIT_ORDER
-  access patterns 
    - key value 
    - column updates: SQL
- No transactinal support: atomicity, concurrency
- auto expiry of records
- performance, low latency and hopefully scalability!
- persistence required only for auth tokens. other use cases really only require a cache
- make a quick decision! this is for a POC

### Approach 1: SQL
- supports a rich data access pattern. useful to implement all the scenarios
downsides:
- requires a ORM. not required for a POC

### Approach 2: key-value
- simple client API to integrate with DB. very useful for a POC
downsides:
- updates nullify the performance benefits of LSM tree based databases
- multiple databases or key spaces cannot be created. Must be handled in the application layer  

### conclucion
- choosing redis. supports,
    - key-value store
    - simple client interface
    - partial updates with HMSET
    - updates with json merging/ partial json updates
    - expiry
    - caching
    - persistence
