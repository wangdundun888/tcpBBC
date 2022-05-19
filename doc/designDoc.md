# tcpBBC
基于tcp的在线聊天系统

### 1.概述
基于tcp+个人设计简单协议通信,分为客户端和服务端两部分,客户端根据用户的输入处理消息、发送消息、接收消息和显示消息,
服务端负责处理各个客户端的消息.

### 2.客户端模块详细设计
#### 2.1 接收模块
&emsp;开启模块所需参数:  
&emsp;&emsp;(1) tcp连接:用以接收消息  
&emsp;&emsp;(2) process通道,将接收到的消息交与处理模块  
&emsp;&emsp;(3) context,全局控制退出  
&emsp;先读取两个字节确认报文的长度,再读取剩余报文,确认报文没有丢失后交由处理模块处理.  
#### 2.2 发送模块
&emsp;开启模块所需参数:  
&emsp;&emsp;(1) tcp连接:用以发送消息  
&emsp;&emsp;(2) send通道,将接收到的消息发送出去  
&emsp;&emsp;(3) context,全局控制退出  
&emsp;直接将从通道接收到的消息发送出去.  
#### 2.3 处理模块
&emsp;开启模块所需参数:  
&emsp;&emsp;(1) recv通道:用以和接收模块通信    
&emsp;&emsp;(2) send通道,用以和发送模块通信    
&emsp;&emsp;(3) context,全局控制退出  
&emsp;先开启一个协程,主要负责处理控制台输入;然后开启一个for循环,负责处理用户的输入以及接收模块发来的消息;  
&emsp;&emsp;用户输入:  
&emsp;&emsp;&emsp;根据不同状态处理用户输入,主要有创建用户,创建聊天室,请求聊天室列表,进入和退出聊天室,在聊天室聊天等操作  
&emsp;&emsp;接收模块发来消息:  
&emsp;&emsp;&emsp;根据协议,先判断消息类型,消息类型包括创建用户确认、请求聊天室列表确认、进入聊天室确认、退出聊天室确认、
发送消息确认、创建聊天室确认、接收其他人消息和其他如进入聊天室提示.
#### 2.4 主模块
&emsp;首先,从控制台读取ip参数;然后根据ip参数获取一个tcp链接,然后初始化创建一些模块间通信所用通道;最后依次开启接收、发送和处理模块.  
### 3.服务端模块详细设计
#### 3.1 接收模块
&emsp;开启模块所需参数:  
&emsp;&emsp;(1) tcp连接:该tcp连接为一个客户端连接  
&emsp;&emsp;(2) process通道,将接收到的消息交与处理模块  
&emsp;&emsp;(3) context,全局控制退出  
&emsp;先读取两个字节确认报文的长度,再读取剩余报文,确认报文没有丢失后交由处理模块处理,每到达一个新的客户端,都会启动一个该模块的协程去处理.    
#### 3.2 发送模块
&emsp;开启模块所需参数:  
&emsp;&emsp;(1) send通道,将接收到的消息发送出去  
&emsp;&emsp;(2) context,全局控制退出  
&emsp;直接将从通道接收到的消息发送出去,通道出来的信息包含了一个指定的tcp链接.  
#### 3.3 处理模块
&emsp;&emsp;(1) recv通道:用以和接收模块通信    
&emsp;&emsp;(2) send通道,用以和发送模块通信    
&emsp;&emsp;(3) context,全局控制退出  
&emsp;创建四个map,负责管理客户端链接、用户名唯一校验、聊天室名唯一校验以及同聊天室用户,从接收模块接收消息,
消息类型分别有链接、断开链接、创建用户、创建聊天室、请求聊天室列表、进入某一聊天室、退出某一聊天室以及发送消息,
根据不同消息类型,处理消息,并给出回复,交由发送模块发出.  
#### 3.4 主模块
&emsp;监听tcp某一端口,然后创建发送,接收等通道等初始资源,开启一个for循环,每到一个新的客户端,go一个接收模块协程处理该客户端通信.  

### 简单通讯协议设计
第 0,1 个byte:  
&emsp;数据报长度,大端编码  
第 2 个byte:  
&emsp;消息类型:  当消息类型不同时,后面定义也有所不同    
&emsp;&emsp;0,请求聊天室列表  
&emsp;&emsp;1,请求聊天室列表确认  
&emsp;&emsp;&emsp;第 3 个byte:字符串分隔符    
&emsp;&emsp;&emsp;第 4 个及以后byte:聊天室列表切片,形式为聊天室id+聊天室名称    
&emsp;&emsp;2,进入某一聊天室  
&emsp;&emsp;&emsp;第 3,4 个byte:聊天室id,大端编码  
&emsp;&emsp;3,进入某一聊天室确认    
&emsp;&emsp;&emsp;第 3 个byte:0为success,1为fail  
&emsp;&emsp;&emsp;如为success，第 4,5 个byte则为聊天室id,大端编码  
&emsp;&emsp;&emsp;如为fail,第 4个及以后byte则为简单的fail reason  
&emsp;&emsp;4,退出某一聊天室  
&emsp;&emsp;&emsp;第 3,4 个byte:聊天室id,大端编码  
&emsp;&emsp;5,退出某一聊天室确认  
&emsp;&emsp;&emsp;第 3 个byte:0为success,1为fail  
&emsp;&emsp;&emsp;如为success,第 4,5 个byte则为聊天室id,大端编码  
&emsp;&emsp;&emsp;如为fail,第 4个及以后byte则为简单的fail reason  
&emsp;&emsp;6,创建某一聊天室  
&emsp;&emsp;&emsp;第 3个及以后byte:聊天室名称,不能超过十个中文字符,即不能超过40个byte  
&emsp;&emsp;7,创建某一聊天室确认  
&emsp;&emsp;&emsp;第 3 个byte:0为success,1为fail  
&emsp;&emsp;&emsp;如为success,第 4,5 个byte则为聊天室id,大端编码  
&emsp;&emsp;&emsp;如为fail,第 4个及以后byte则为简单的fail reason  
&emsp;&emsp;8,在某一聊天室发送消息  
&emsp;&emsp;&emsp;第 3,4 个byte:聊天室id,大端编码  
&emsp;&emsp;&emsp;第 5,6 个byte:个人id,大端编码  
&emsp;&emsp;&emsp;第 7,8,9 个byte:留待扩展  
&emsp;&emsp;&emsp;第 10个及以后byte:消息  
&emsp;&emsp;9,在某一聊天室发送消息确认  
&emsp;&emsp;&emsp;第 3,4 个byte:聊天室id,大端编码  
&emsp;&emsp;&emsp;第 5,6 个byte:个人id,大端编码  
&emsp;&emsp;&emsp;第 7,8,9 个byte:留待扩展  
&emsp;&emsp;&emsp;第 10个byte:0为success,1为fail  
&emsp;&emsp;&emsp;第 11个及以后byte:如果fail,则为简单的fail reason  
&emsp;&emsp;10,创建临时聊天用户名  
&emsp;&emsp;&emsp;第 3个及以后byte:用户名,不能超过十个中文字符,即不能超过40个byte  
&emsp;&emsp;11,创建临时聊天用户名确认  
&emsp;&emsp;&emsp;第 3 个byte:0为success,1为fail   
&emsp;&emsp;&emsp;如为success,第 4,5 个byte则为个人id,大端编码  
&emsp;&emsp;&emsp;如为fail,第 4个及以后byte则为简单的fail reason  
&emsp;&emsp;12,接收其他人的消息  
&emsp;&emsp;&emsp;第 3,4 个byte:他人昵称长度,大端编码  
&emsp;&emsp;&emsp;第 5个至长度byte：他人昵称  
&emsp;&emsp;&emsp;后面附加他人发送消息  
&emsp;&emsp;13,他人加入聊天室  
&emsp;&emsp;&emsp;第 3个及以后byte:他人用户名  
&emsp;&emsp;其余待用.  
    


