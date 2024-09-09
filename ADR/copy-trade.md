## copy trade
user must be able to mimic/ copy another user's trades

### requirements
- user must see top traders and their metrics
- user can select another user to copy trade
- how does users identity other users? 
    - address would be ideal to identify other users

### implementation
Approach 1:
- identify users by telegram identity. since the trades are all tracked within the bot and there is no support from okto api's
    - for milestone 1, cache telegram username at the time of login. does not handle chainging usernames
        - save user id for username
        - the list of followers is mainted for each followee
            - when another user enters the user name, lookup the id and save the follower
        - every trade request should have the user id
        - every time a trade executes, lookup followers by id and execute copy trades

    - for milestone 2, resolve telegram user id for username with telegram api

is chat id the best user identifier? does telegram have a user id?
do not use update.message.Chat.Id. instead use, update.message.From.Id

Approach 2:
- simply accept address to copy trades. 
- whenever a trade occurs from the given network, copy the trade 

portfolio
    - for milestone 1, user will be able to enter another user's telegram name manually to copy her trades
    - for milestone 2, build a user portfolio and user can select another user to copy trade

### Future work
- list created copy orders
- cancel copying trades
- unified copy trades across networks? user has to create separate copy trades for each network, by their network address

### References
- https://core.telegram.org/api
- https://core.telegram.org/method/contacts.resolveUsername
- https://github.com/gotd/td
- https://github.com/danog/MadelineProto/blob/7629e33f857bc35a5b13043e73f0e2702519caae/src/MTProtoTools/PeerDatabase.php#L336