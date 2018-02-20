<p align="center">
<img
    src="https://raw.githubusercontent.com/schollz/snaptext/master/static/favicon/logo.png?token=AGPyE4FL_L452-C_VhQ1bi8WiJhpB6ALks5alK3HwA%3D%3D"
    width="260px" border="0" alt="snaptext">
<br>
<a href="https://travis-ci.org/schollz/snaptext"><img src="https://travis-ci.org/schollz/snaptext.svg?branch=master" alt="Build Status"></a>
<a href="https://github.com/schollz/snaptext/releases/latest"><img src="https://img.shields.io/badge/version-0.1.0-brightgreen.svg?style=flat-square" alt="Version"></a>
<a href="https://goreportcard.com/report/github.com/schollz/croc"><img src="https://goreportcard.com/badge/github.com/schollz/croc" alt="Go Report Card"></a>
<a href="https://www.paypal.me/ZackScholl/5.00"><img src="https://img.shields.io/badge/donate-$5-brown.svg" alt="Donate"></a>
</p>

<p align="center">Like snapchat, but for text.</p>

*snaptext* is a web app (and API) that lets you easily send and receive self-destructing messages in real-time. For example, you can go to [`snaptext.live/?to=schollz`](https://snaptext.live/?to=schollz) and write me a message. The message will be stored in a queue for me (`schollz`) and it will be destroyed when a browser is opened at [`snaptext.live/schollz`](https://snaptext.live/schollz) which pops the latest message. 

Messaging occurs in real-time using websockets, so to guarantee that you receive the message its best to have the browser open or use a obfuscated ID. Messages are queued for each ID, so you can send multiple messages and they will be read in order.

# Why?

I recently made [a "turnkey" solution for the Raspberry Pi](https://github.com/schollz/raspberry-pi-turnkey) to solve a problem with starting up a Pi without using SSH. The turnkey image Pi starts up a temporary WiFI access point and the user enters their home WiFi credentials. The Pi then restarts and connects to the new WiFi, and needs a way to tell the user is connected. Email is not an option because I can't distribute Pi's with my own SMTP credentials. Thus, I made *snaptext* so that the user can login to a temporary webpage, like `snaptext.live/abc234basd3b` where the Pi will send a message once it becomes online.

There may be other uses for *snaptext* - it is a very transient way of sending short messages once a URL is shared between the parties.

# Usage

## Writing messages online

Goto [`snaptext.live`](https://snaptext.live) to write messages.

## Reading messages online

Goto [`snaptext.live/maru`](https://snaptext.live/ID) to read messages that have been written to `maru`. Once a message is read, it is destroyed. There is no check on who reads a message - it is first come first serve.

## API

The API is incredibly simple. There is only one endpoint, to post a message.

**POST /** - send a messsage

Use the following payload to set the recipient (`to`), the sender name (`from`), and the message (`message`).

```json
{
    "to":"maru",
    "from":"schollz",
    "message":"Hi Maru!"
}
```

The recipient controls where the message can be seen (it will be seen at `snaptext.live/X`) where `X` is the recipient.

The `from` just tells who is sending the message. The message can be text or HTML, though a limited number of HTML tags are allowed (to prevent XSS attacks).

# Run your own server

The easiest way is using Go:

```
go get github.com/schollz/snaptext
```

or you can download a release for your system on the releases.

# License

MIT
