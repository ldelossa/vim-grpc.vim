let g:vgrpc_environment = { 
                \ "appName": "vim",
                \ "appRoot": system('git rev-parse --show-toplevel'),
                \ "language": v:lang,
                \ "machineId": system('hostname'),
                \ "remoteName": "",
                \ "sessionId": "session-" . rand(),
                \ "shell": $SHELL,
                \}

function! handlers#env#GetEnv(channel, envelope) abort
    if !rpc#validate#Name(a:envelope, "GetEnv")
        return
    endif
    let a:envelope["body"] = g:vgrpc_environment
    call ch_sendexpr(a:channel, a:envelope) 
endfunc
