# textmailmachine

This is like a voice mail machine, except you can leave text messages for anyone. A message to X will be available at `textmailmachine.schollz/X` for the specified amount of time and then it is deleted off the server forever.

## API

**POST /** - send a messsage

Use the following payload to set the recipient (`to`), the sender name (`from`), the message (`message`), and the number of seconds before the message is deleted after being shown (`display`).

Example JSON:

```json
{
    "to":"Zack",
    "from":"Your friend",
    "message":"Hi Zack!",
    "display":3
}
```

