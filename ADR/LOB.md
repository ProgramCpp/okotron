## limit order book
there is no need for order match making. this is taken care by okto. every order is independent.

### requirements
- accept limit orders from users
- execute limit order when price match happens
- handle price volatality: slippage handled by okto. not configurable by user, atleast at this point
- handle multiple payment tokens. 
    - the limit price is not constant, if denominated in the source tokens.
    - it is also hard for the user to determine the value of the limit order.
- the order volume is still not high! Handle them as volume grows
    - sequentially process all orders, process concurrently
    - slippage not too much, since the order volume is low and orders execute almost instantaneously
- handle buy and sell limit orders

### Approach
Handling token denomination:
- to simplify, denominate the target price of the token in USD/ stable coins
- if the user is paying with another token, the direction of the from-token price doesnt matter. 
// Wow, this brings out a unique use case where user can trade an unfavorable token for a favorable token at a favorable price, both at once! But only downside of no limits for the payable token. better yet, user can place a buy limit order and a sell limit order
- user can always choose to use USDT/ USDC for trades for a stable conventional limit order trades

Implementation: A classic case of a task scheduler. instead of time, you act on price! even simpler, a callback to execute ALL trade transactions at a given limit price.
- save a list of orders for limit prices in redis
- get live price feeds
- process orders concurrently. 
  - keep polling for token price changes
  - trace a list of orders with the changed token price
  - process only those orders

request payload:
to simplify and to be deterministic, accept from-network from user. no need to deduce the balances of tokens from all networks.
in a conventional limit order, lets say stocks, there is only one account. in this case there are number of networks, pick the source and target network

- buy-or-sell : buy order or sell order
- from network : the network where the tokens are, that is used to pay/ sell
- from token : the token used to pay/ sell
- to network : the network, where the bought tokens/returns should be trasferred to
- to token : the token to purchase/ get paid
- limit price: the target market price of the to-token(for buy orders)/ from-token(for sell orders) in USD at which, the order must execute
- quantity: the amount of to-token to buy/ from-token to sell

### Future Work
- should you have a separate map for each token to-token/ from-token. improve this as volume grows

### References
- https://gist.github.com/halfelf/db1ae032dc34278968f8bf31ee999a25
- https://github.com/i25959341/orderbook