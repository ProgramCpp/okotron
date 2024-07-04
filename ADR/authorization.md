## okto auhthorization flow.

okto account is managed by user's google account. oktron must request user authorization to access okto API's.

### requirements
1. do not expose telegram chat_id. it is not a secret key and not managed by the bot. redact chat_id.
2. user is already authenticated to telegram. telegram bot is an extension to telegram. thus, telegram user id can be used as the user identifier in the bot
3. there is no need for explicit authentication to the bot. If the bot requires authentication, How do you map the telegram user identity with the bot user identity?

## User flow

### Approach 1: browser flow
downside: security concerns since this is not the standard oauth 2.0 flow. enforce best security practices. bot session passed between user client's - bot authorization token is curated in the redirection url. malicious actor can gain control over okto with access to this url. The short lived token can protect the user but it is not a bulletproof solution. Never roll your own security algorithm!

1. user initiates oktron authorization from telegram
2. a short lived session id is created for the chat_id. the chat_id for the session_id is persisted with a short expiry - time enough to complete authorization.
3. the bot responds with a link to the bot authorization page along with the session_id
4. user opens the authorization page, which sets the session_id cookie
5. user completes the google oauth authorization flow and is redirected back to bot authorization finish page
6. the user submits both the session id and the google auth code to finish authorization
7. the user is redrected back to telegram bot chat

#### Future Work

- how to pass session information from an app to the browser?


### Approach 2: [device authorization flow](https://www.oauth.com/oauth2-servers/device-flow/user-flow/)
downside: not the best UX

### Conclusion: Approach 2 is the reccomended flow for clients with limited input functionality


