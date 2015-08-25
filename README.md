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

* start

Param | Type | Default|Dircetions
---|---|---|----
--config, -c | string |./default.json| config

## Config  

```
{
        "auth_account":"account", //Super Admin 
        "auth_password":"password",
        "environment":"DEBUG",
        "logdir":"./", // "console" or "/homs/user/x.log"
        "listen":":8001",
        "sqldir":"./appdata.sqlite"
}
                        
```


## API

All api need http basic Auth, Super Admin can access all api

#### Register:

`[POST] /api/register`  

* Param:  

Name|Type|Directions
---|---|---
app_name | string | a app_name
auth_account | string | app admin basic auth account
auth_password | string | app admin basic auth password

* 200 status Response:  

```
{
    "app_name":"test",
    "auth_account":"app_admin",
    "auth_password":"password",
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



