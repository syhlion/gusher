# Gusher

[Pusher](https://pusher.com/docs) implment server-side by golang

## Requirements

* sqlite 3

## Usage

Install from source

```
$ go get github.com/syhlion/gusher
$ gusher -h

```

Install from binary

## Command

* start

Param | Type | Default|Dircetions
---|---|---|----
--config, -c | string |./default.json| config Optinal

## Config  

```
{
        "auth_account":"account", //Super Admin 
        "auth_password":"password",
        "environment":"DEBUG", // DEBUG || DEVELOPMENT || PRODUCATION
        "logfile":"./", // "console" or "/homs/user/x.log"
        "listen":":8001",
        "sqlfile":"./appdata.sqlite",
        "max_wait_hook":10
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
auth_account | string | app admin basic auth account. Require
auth_password | string | app admin basic auth password. Require
connect_hook | string | You can take this verification [webhook](https://github.com/syhlion/gopusher/#web-hook). Optinal

* 200 status Response:  

```
{
    "app_name":"test",
    "auth_account":"app_admin",
    "auth_password":"password",
    "app_key":"abcdefghijklmnop",
    "connect_hook":"http://localhost/vaildlogin"
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

`[GET] /api/app-list/{limit}/{page}`  

* 200 status Response:

```
{
    "limit":2,
    "offset":1,
    "total":2,
    "data":
    [
        {
            "app_key":"db15759925b279b4b037d7a4e1f92b0f",
            "app_name":"test",
            "auth_account":"scott",
            "auth_password":"760804",
            "connect_hook":"http://localhost/vaildlogin"
            "request_ip":"127.0.0.1:50040",
            "date":"2015/08/25 11:39:12",
            "timestamp":"1440473952104"
        },
        {
            "app_key":"56f08a519be877060fb4a4ea2a75aad8",
            "app_name":"test2",
            "auth_account":"scott",
            "auth_password":"760804",
            "connect_hook":"http://localhost/vaildlogin"
            "request_ip":"127.0.0.1:55567",
            "date":"2015/08/26 12:10:54",
            "timestamp":"1440562254686"
        }
    ]
}

```


#### Push:  

`[POST] /api/push/{app_key}`  

* Param:  

Name|Type|Dircetions
---|---|---
content| string | the message you send
user_tag | string | the message who will receive. Support Regex. Optinal

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

`[GET] /api/{app_key}/listonlineuser/{limit}/{page}`  

* 200 status Response:

```
{
    "app_key":"abcdefghijklmnop",
    "total":1,
    "limit":1,
    "page":1
    "user_tags":["A:1"]
}
```

#### Client Connect:  

If you have fill in Register api connect_hook, use this:  

`[GET] ws://localhost/ws/{app_key}/{user_tag}?token={token}`  


If it is not,use this:  

`[GET] ws://localhost/ws/{app_key}/{user_tag}`


## WEB hook

#### Connect Hook

If you use this, you need to prepare a hook api, to accept http request POST "token"  

If successful at http body output "ok"   

* Param:  

Name|Type|Dircetions
---|---|---
token| string | Your server api  will get POST from GuPusher
