function! handlers#commands#RegisterCommand(channel, envelope) abort
    if !rpc#validate#Name(a:envelope, "RegisterCommand")
        return
    endif
    let title = a:envelope["body"]["title"]
    let cmd = a:envelope["body"]["command"]

    exec 'command! ' . title . ' ' . 'call DoCommand("' . cmd . '")'

    let registration = {
                \ "registered": v:true,
                \ "reason": ""
                \}
    let a:envelope["body"] =  registration
    call ch_sendexpr(a:channel, a:envelope)
endfunc

function! DoCommand(command)
    let envelope = {
                \ "mailbox": 0,
                \ "rpc": "CommandIssued",
                \ "body": { "command": a:command }
                \}
    call ch_sendexpr(g:vgrpc_channel, envelope)
endfunc
