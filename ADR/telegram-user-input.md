## telegram user input
- when a command is incomplete and expects user input, it sets the next command/ sub command in message_<message-id>: <next-command>
- if the telegram callback handler finds a next-command, calls the corresponding handler
- the user can respond in two ways
	- by replying to a bot message. this is identified by `update.Message.ReplyToMessage` 
	- by interacting with an inline keyboard. this is identified by `update.CallbackQuery`

### Future Work:
why do we need two different ways to handle user input? consolidate into one. the inline keyboard is more intuitive and readable
migrate PIN flow to inline keyboard flow.