let g:VGRPC_router = {
      \ "Ping": function("handlers#ping#Ping"),
      \ "GetEnv": function("handlers#env#GetEnv"),
      \ "RegisterCommand": function("handlers#commands#RegisterCommand")
      \ }
