## limit order book
there is no need for order match making. this is taken care by okto. every order is independent.

### requirements
- accept limit orders from users
- execute limit order when price match happens
- handle price volatality: slippage handled by okto. not configurable by user, atleast at this point

### Approach
A classic case of a task scheduler. instead of time, you act on price! even simpler, a callback to execute ALL trade transactions at a given limit price.

- save a list of orders for limit prices in redis
- get live price feeds
- process orders concurrently. 

the order volume is still not high! Handle them as volume grows
    - sequentially process all orders, process concurrently
    - slippage not too much


request payload:
to simplify and to be deterministic, accept from-network from user. no need to deduce the balances of tokens from all networks.
in a conventional limit order, lets say stocks, there is only one accont. in this case there are number of networks, pick the source and target network

- from network : the network where the tokens are, that is used to pay
- from token : the token used to pay
- to network : the network, where the bought tokens should be
- to token : the token to purchase
- limit price: the target market price of the to-token at which, the order must execute
- quantity: the amount of to-token to buy

### References
- https://gist.github.com/halfelf/db1ae032dc34278968f8bf31ee999a25
- https://github.com/i25959341/orderbook