function! rpc#validate#Name(msg, rpc_name)
    if a:msg["rpc"] == a:rpc_name
        return v:true
    endif
    return v:false
endfun
