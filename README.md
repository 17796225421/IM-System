IM即时通讯

Server结构，

NewServer函数，构造Server结构

Start函数，启动服务器，listen后，死循环不断accept到一个客户端连接，就交给goroutine Handle函数。

Handler函数，将客户端加入在线用户表中，在线用户表有锁。要广播给所有User说明有一个客户端上线，将消息交给Server的channel，需要有一个gotoutine将Server的channel消息转发给所有的User，也就是Server的listenMessage函数。

Handler函数需要有一个goroutine来接收自己的负责的客户端发送过来的消息，将收到的数据进行判断，可以对数据进行广播，也就是将消息放到Server的channel，会有一个goroutine将Server的channel消息转发给所有User，也就是Server的listenMessage函数，可以对数据进行判断。

如果数据是who，也就是客户端发送过来的消息是who，将在线用户表的所有用户信息发送给客户端。

如果数据是rename|zhangsan，把当前客户端用户名改为zhangsan。修改完将修改成功消息发送给客户端。

如果数据是to|zhangsan|nihao，私聊，根据用户名，获取User，对该User客户端发送消息。

Handler函数需要实现超时踢出当前客户端用户的功能，可以用一个管道，每次，需要计时器，只要客户端发送过来任何消息，就需要一个channel，往channel输入true表示客户端活着。Handler函数可以用select的一个case来从channel取出数据，一旦取出数据，就重置计时器。select的另一个case来每10s触发一次超时踢出。

listenMessage函数，将Server的channel消息转发给所有的User，在在向用户表中所有的User。

main函数，对server构造，启动。

User结构，专门负责与一个客户端进行交互。

NewUser函数，构造User结构。

ListenMessage函数，如果User的channel有数据，就转发给客户端。

Client结构，代表一个客户端

newClient函数，构造，连接服务器

init函数，解析命令行 ./client -ip 127.0.0.1 -port 8888

main函数，构造客户端连接服务器后，开始业务。

Run函数执行业务，显示菜单，执行对应业务公聊、私聊、更新用户名

UpdateName函数更新用户名，组装rename|zhouzihong，发送给服务器。

DealResponse函数，负责处理服务器消息回应的goroutine，直接打印到屏幕。

PublicChat函数公聊，死循环接收标准输入，只要不是exit，就发送给服务器。

PrivateChat函数私聊，先组装who\n发送给服务器查询所有在线用户，再死循环组装to|zhouzihong|nihao\n发送给服务器
