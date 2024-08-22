## copy trade
user must be able to mimic/ copy another user's trades

### requirements
- user must see top traders and their metrics
- user can select another user to copy trade
- how does users identity other users

### implementation
- identify users by telegram identity. since the trades are all tracked within the bot and there is no support from okto api's
    - for milestone 1, cache telegram username at the time of login. does not handle chainging usernames
        - save user id for username
        - the list of followers is mainted for each followee
            - when another user enters the user name, lookup the id and save the follower
        - every trade request should have the user id
        - every time a trade executes, lookup followers by id and execute copy trades

    - for milestone 2, resolve telegram user id for username with telegram api

portfolio
    - for milestone 1, user will be able to enter another user's telegram name manually to copy her trades
    - for milestone 2, build a user portfolio and user can select another user to copy trade

is chat id the best user identifier? does telegram have a user id?
do not use update.message.Chat.Id. instead use, update.message.From.Id

### References
- https://core.telegram.org/api
- https://core.telegram.org/method/contacts.resolveUsername
- https://github.com/gotd/td
- https://github.com/danog/MadelineProto/blob/7629e33f857bc35a5b13043e73f0e2702519caae/src/MTProtoTools/PeerDatabase.php#L336