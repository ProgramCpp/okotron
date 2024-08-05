## database for persistence

### requirements
- high update rate for certain commands
    - SWAP
-  access patterns 
    - key value 
    - column updates: SQL
- No transactinal support: atomicity, concurrency
- auto expiry of records
- performance, low latency and hopefully scalability!
- persistence required only for auth tokens. other use cases really only require a cache
- make a quick decision! this is for POC

### Approach 1: SQL
- supports a rich data access pattern. useful to implement all the scenarios
downsides:
- requires a ORM. not required for a POC

### Approach 2: key-value
- simple client API to integrate with DB. very useful for a POC
downsides:
- updates nullify the performance benefits of LSM tree based databases
- multiple databases or key spaces cannot be created. Must be handled in the appliction layer  

### conclucion
- choosing redis. supports expiry, caching and persistence, with simple client interface
- implement update scenario with read-modify-update. there is no concurrency.
- if SWAP command use case latency needs to be optimized, then migrate to an alternate db with microservice pattern maybe.

### Future work
- cockroach db
- yugabyte db
- postgres db
- noSQL with efficient update support
- if it wasnt for swap command, would you have picked a key-value store? maybe. then you dont need the Swap command!🙄 Its a lot of work. Trust me, its totally worth it 😕