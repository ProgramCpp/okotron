## limit order book
there is no need for order match making. this is taken care by okto. every order is independent.

### requirements
- accept limit orders from users
- execute limit order when price match happens

### Approach
A classic case of a task scheduler. instead of time, you act on price!



### References
- https://gist.github.com/halfelf/db1ae032dc34278968f8bf31ee999a25
- https://github.com/i25959341/orderbook