<p align="center">
<img
    src="https://raw.githubusercontent.com/schollz/snaptext/master/static/favicon/android-icon-144x144.png?token=AGPyE68M8fOMP_cj87oSPy7gbOR2yVimks5alHtwwA%3D%3D"
    width="144px" border="0" alt="snaptext">
<br>
<a href="https://travis-ci.org/schollz/snaptext"><img src="https://travis-ci.org/schollz/snaptext.svg?branch=master" alt="Build Status"></a>
<a href="https://github.com/schollz/snaptext/releases/latest"><img src="https://img.shields.io/badge/version-0.1.0-brightgreen.svg?style=flat-square" alt="Version"></a>
<a href="https://goreportcard.com/report/github.com/schollz/croc"><img src="https://goreportcard.com/badge/github.com/schollz/croc" alt="Go Report Card"></a>
<a href="https://www.paypal.me/ZackScholl/5.00"><img src="https://img.shields.io/badge/donate-$5-brown.svg" alt="Donate"></a>
</p>

<p align="center">Like snapchat, but just for text.</p>

Leave a text/html message for anyone on the internet. For example, I could write a message for "Zack" which will be available at [`snaptext.live/zack`](https://snaptext.live/zack). Once the message is opened, it is deleted and never shown again. Messages are updated in realtime, so you can just leave your browser open to gaurantee recieving the message.

# Why?

# Usage

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



# Run your own server

The easiest way is using Go:

```
go get github.com/schollz/snaptext
```

or you can download a release for your system on the releases.

# License

MIT
