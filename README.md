# GoBo
Twitch bot written in go.

warning: this bot is still in development 

## Builtin Commands
### For any user
``` bach
< !echo text
> text

< !me
> [caller username]

< !insert text0 %% text1 %% text3 << data 
> text0 data text1 data text3

< !repeat 4 text
> text text text text

< !+ 1.0 1.0
> 2.00

< !- 1 2.0 -1
> 0.00

< !* 2 2.0 2.0 -1.0
> -8.00

< !/ 1.0 0
> inf

< !% 3.0 2
> 1.00
```

### Broadcaster only
``` bach
// just adds the new command to the Bot and DB
< !addcmd [permission] [cmd name] [cmd body]
> 

// just adds the new recurrent command to the Bot and DB
< !addrcmd [cmd name] [duration]
>
```

## The composition operator '$'
You can compose commands using '$', for example,
``` bach
< !insert 2 + 3 * 3 = %% << $ !+ 2 $ !* 3 3
> 2 + 3 * 3 = 11.00
```
