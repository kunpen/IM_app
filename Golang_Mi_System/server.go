package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip		string
	Port 	int
	//添加user表,读写锁
	OnlineMap map[string]*User
	mapLock	  sync.RWMutex

	//添加进行消息广播的msg管道
	Message chan string


}
//create a server intserface


func NewServer(ip string,port int) *Server {
	server := &Server{
		Ip: ip,
		Port: port,
		OnlineMap: make(map[string]*User),
		Message: make(chan string),
	}
	return server
}
//监听广播消息,一旦有消息就发送给在线user
func (this *Server)ListenMessage()  {
	for  {
		msg := <- this.Message
		//发送给在线user
		this.mapLock.Lock()
		for _,cli := range this.OnlineMap{
			cli.C<-msg
		}
		this.mapLock.Unlock()
	}
}


//广播方法
func (this *Server)BordCast(user *User,msg string)  {
	sendMsg:= "["+user.Addr+"]"+user.Name+":"+msg
	this.Message<-sendMsg
}


func (this *Server)Handler(conn net.Conn)  {
	//当前连接的业务
	fmt.Println("connect success!,new user online!!")

	//用户上线了，将用户加入online map中
	user := NewUser(conn,this)
	user.Online()



	//监听用户是否活跃的channel
	isLive:=make(chan bool)


	//接收客户端发送的消息
	go func() {
		buf:= make([]byte,4096)
		for  {
			n,err := conn.Read(buf)
			if n==0 {
				user.Offline()
				return
			}
			if err!=nil&& err!=io.EOF{
				fmt.Println("conn read err:",err)
			}
			msg:= string(buf[:n-1])
			//广播消息
			//针对用户进行消息处理
			user.DoMessage(msg)
			//活跃用户
			isLive<-true


		}
	}()

	for  {
		select {
			//定时器
			case <-isLive:
				//活跃用户，重置定时器
				//不需要做操作 会执行下面条件判断的代码，重置定时器
				//激活select
				//但是不会执行 判断后的代码

			case <-time.After(time.Second * 60):
				//case 中有数据，说明已经超时
				//关闭user 连接

				user.SendMessage("you are not activite ,good by\n")
				close(user.C)
				conn.Close()

				this.mapLock.Lock()
				delete(this.OnlineMap,user.Name)
				this.mapLock.Unlock()

				//退出handler
				return





		}

	}




	//广播用户上线消息



	
}




//start server function

func (this Server)Start()  {
	//TODO socket listen
	listener,err :=net.Listen("tcp",fmt.Sprintf("%s:%d",this.Ip,this.Port))
	if err !=nil {
		fmt.Println("net.listen error:",err)
		return
	}
	// close listen socket
	defer listener.Close()
	//启动监听message的goroutine
	go this.ListenMessage()


	for {
		//accept
		conn,err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept error:",err)
			continue
		}


		//do handler
		go this.Handler(conn)


	}



}
