package main

import (
	"fmt"
	"net"
	"sync"
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
	user := NewUser(conn)
	this.mapLock.Lock()
	this.OnlineMap[user.Name]=user
	this.mapLock.Unlock()

	this.BordCast(user,"online!!")
	select {

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
