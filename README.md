# GoPusher

[Pusher](https://pusher.com/) implment server-side by golang

## Requirements

* Linux 64 bit
* sqlite 3

## Usage

Install from source

```
$ go get github.com/syhlion/gopusher
$ gopusher -h

```

Install from binary

## Command

Param | Type | Default|Dircetions
---|---|---|----
--addr, -a | string |:8001| Input listent ip port
--env, -e | string | PRODUCTION|Log Level PRODUCTION, DEVELOPMENT, DEBUG
--log, -l | string |console| Input console or /home/user/gusher/gusher.log


## API

#### Register:

`[POST] /api/register`  

* Param:  

Name|Type|Directions
---|---|---
app_name | string | a app_name

* 200 status Response:  

```
{
    "app_name":"test",
    "app_key":"abcdefghijklmnop",
    "request_ip":"127.0.0.1:77777"
}
```

#### Unregister:  

`[DELETE] /api/{app_key}/unregister`  

* 200 status Response:

```
{
    "message":"Scuess"
}
```


#### List All App:  

`[GET] /api/app-list`  

* 200 status Response:

```
[
    {
        "app_key":"abcdefghijklmnop",
        "app_name":"test",
        "request_ip":"127.0.0.1:77777",
        "date":"2015/08/04 11:22:33",
        "timestamp":"1440149593490"
    }
]
```


#### Push:  

`[POST] /api/push/{app_key}`  

* Param:  

Name|Type|Dircetions
---|---|---
content| string | the message you send
user_tag | string | the message who will receive. Support Regex. OPTION

* 200 status Response:  

```
{
    "app_key":"abcdefghijklmnop",
    "content":"hello world",
    "user_tag":"A:1"
    "total":"1"   // receive this message client total
}
```

#### List Online User:  

`[GET] /api/{app_key}/listonlineuser`  

* 200 status Response:

```
{
    "app_key":"abcdefghijklmnop",
    "total_online_user":"1",
    "online_user:["A:1"]
}
```

#### Client Connect:  

`[GET] ws://localhost/ws/{app_key}/{user_tag}`



