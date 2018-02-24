<p align="center">
<img
    src="https://raw.githubusercontent.com/schollz/snaptext/master/static/favicon/logo.png?token=AGPyE4FL_L452-C_VhQ1bi8WiJhpB6ALks5alK3HwA%3D%3D"
    width="260px" border="0" alt="snaptext">
<br>
<a href="https://travis-ci.org/schollz/snaptext"><img src="https://travis-ci.org/schollz/snaptext.svg?branch=master" alt="Build Status"></a>
<a href="https://github.com/schollz/snaptext/releases/latest"><img src="https://img.shields.io/badge/version-0.1.0-brightgreen.svg?style=flat-square" alt="Version"></a>
<a href="https://goreportcard.com/report/github.com/schollz/snaptext"><img src="https://goreportcard.com/badge/github.com/schollz/snaptext" alt="Go Report Card"></a>
<a href="https://www.paypal.me/ZackScholl/5.00"><img src="https://img.shields.io/badge/donate-$5-brown.svg" alt="Donate"></a>
</p>

<p align="center">Like snapchat, but for text.</p>

*snaptext* is a web app (and API) that lets you easily send and receive self-destructing messages in real-time. For example, you can go to [`snaptext.live/?to=schollz`](https://snaptext.live/?to=schollz) and write me a message. The message will be stored in a queue for me (`schollz`) and it will be destroyed when a browser is opened at [`snaptext.live/schollz`](https://snaptext.live/schollz) which pops the first message. 

Messaging occurs in real-time using websockets, so to guarantee that you receive the message its best to have the browser open or use a obfuscated ID. Messages are queued for each ID, so you can send multiple messages and they will be read in order (FIFO).

# Why?

I recently made [a "turnkey" solution for the Raspberry Pi](https://github.com/schollz/raspberry-pi-turnkey) to easily assign the Pi WiFi credentials without using SSH or writing to the boot (useful for shipping to customers). The turnkey image Pi starts up a temporary WiFI access point and the user enters their home WiFi credentials. The Pi then restarts and connects to the new WiFi. At this point, it needs a way to communicate to the user that it is connected and provide its LAN IP. Email is not an option here because I cannot ship a Pi using my own SMTP credentials. Thus, I made *snaptext* so that the Pi sends the user the message through the temporary webpage, like `snaptext.live/abc234basd3b`, which tells the user that it is online and gives its IP address.

There may be other uses for *snaptext* - it is basically a simple, transient way of sending short messages once a URL is shared between the parties.

# Usage

*snaptext* only supports doing two things: writing or reading messages.

## Writing messages

You can write messages online. Goto [`snaptext.live`](https://snaptext.live) to write messages. The message can be text or HTML, though a limited number of HTML tags are allowed (to prevent XSS attacks). 

You can also write messages from other programs. The API is incredibly simple. There is only one endpoint, to post a message: **POST /**. Use the following payload to set the recipient (`to`), the sender name (`from`), and the message (`message`).

```json
{
    "to":"snaptext",
    "from":"schollz",
    "message":"Just a test"
}
```

The recipient controls where the message can be seen (this particular message will be seen at `snaptext.live/snapchat`). The `from` just tells who is sending the message. An example CURL:

```bash
curl  -d '{"to":"snaptext","from":"schollz","message":"Just a test"}' -X POST https://snaptext.live
```

## Reading messages

Goto [`snaptext.live/snapchat`](https://snaptext.live/ID) to read messages that have been written to `snapchat`. Once a message is read, it is destroyed. There is no check on who reads a message - it is first come first serve. However, anyone with a browser currently connected can read an incoming message.

# Run your own server

The easiest way is using Go (requires Go 1.9+):

```
$ go install -v github.com/schollz/snaptext
$ snaptext
```

# License

MIT
