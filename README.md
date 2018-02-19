# snaptext

This is like a voice mail machine, except you can leave text messages for anyone. A message to someone named "X" will be available at `snaptext.schollz/X`. Once someone looks at the message, it will be deleted. Messages are updated in realtime, so you can just leave your browser open (like a very silly isntant messenger).

## API

The API is incredibly simple. There is only one endpoint, to post a message.

**POST /** - send a messsage

Use the following payload to set the recipient (`to`), the sender name (`from`), the message (`message`), and the number of seconds before the message is deleted after being shown (`display`).

Example JSON:

```json
{
    "to":"Zack",
    "from":"Your friend",
    "message":"Hi Zack!"
}
```

