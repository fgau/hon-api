# API for hot or not website

**HTTP Request**

```Get https://example.com/getperson/{gender}```

Parameter **gender**:

* **female**
* **male**

**HTTP Response**
<pre>

{
    "id":"dp4noljbokrm9u3t5dq067jn52",
    "nickname":"julileein",
    "gender":"female",
    "pixurl":"http://pix.hotornot.de/p/f/9/4/f94fe06f6d00.jpg"
}

</pre>

---

**HTTP Post Request**

```Post https://example.com/voteperson```

Json Payload:

<pre>

{
    "ID":"dp4noljbokrm9u3t5dq067jn52",
    "Nickname":"julileein",
    "Vote":"0"
}

</pre>

Curl Example:

curl -X POST -H "Content-Type: application/json" -d '{"ID":"0815","Nickname":"BLNMausi","Vote":"0"}' --insecure  https://example.com/voteperson

**HTTP Response**
<pre>

{"result":"received data succesfully"}

</pre>
