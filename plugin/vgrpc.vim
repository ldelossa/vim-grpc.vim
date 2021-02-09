let g:vgrpc_channel = ""

function! VGRPC_route_rpc(channel, msg) 
  echom "got rpc message " . a:msg["rpc"]
  call g:VGRPC_router[a:msg["rpc"]](a:channel, a:msg)
endfun

function! s:VGRPC_start() 
    let g:vgrpc_channel = ch_open("localhost:7999", {
          \ "waittime": 0,
          \ "callback": "VGRPC_route_rpc"
          \})
endfun

function! s:VGRPC_stop() 
  call ch_close(g:vgrpc_channel)
endfun

command! -nargs=* VGRPCStart call s:VGRPC_start()
command! -nargs=* VGRPCStop  call s:VGRPC_stop()


