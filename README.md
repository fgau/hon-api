# API for hot or not website

**HTTP Request**

```Get https://localhost:8090/getperson/{gender}```

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

Curl Example:

curl -i -H "Accept: application/json" -H "Content-Type: application/json" -X GET http://localhost:8090/getperson/female

---

**HTTP Post Request**

```Post https://localhost:8090/voteperson```

Json Payload:

<pre>

{
    "ID":"dp4noljbokrm9u3t5dq067jn52",
    "Nickname":"julileein",
    "Vote":"0"
}

</pre>

Curl Example:

curl -X POST -H "Content-Type: application/json" -d '{"ID":"0815","Nickname":"BLNMausi","Vote":"0"}' --insecure  https://localhost:8090/voteperson

**HTTP Response**
<pre>

{"result":"received data succesfully"}

</pre>
