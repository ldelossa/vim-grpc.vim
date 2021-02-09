function! handlers#ping#Ping(channel, msg)
    if !rpc#validate#Name(a:msg, "Ping")
        return
    endif
    let a:msg["rpc"] = "Pong"
    call ch_sendexpr(a:channel, a:msg)
endfun
