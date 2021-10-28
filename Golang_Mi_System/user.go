package main

import (
	"net"
)

type User struct{
	Name 	string
	Addr 	string
	C 		chan string
	conn 	net.Conn
	server  *Server
}

//create user function

func NewUser(conn net.Conn,server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:	  make(chan string),
		conn: conn,
		server: server,
	}
	//启动监听当前user 的channel 的goroutinue
	go user.ListMessage()
	return user
}
func (this *User)Online()  {
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name]=this
	this.server.mapLock.Unlock()
	//广播上线消息
	this.server.BordCast(this,"online!!")
}
func (this *User)Offline()  {
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap,this.Name)
	this.server.mapLock.Unlock()
	//广播上线消息
	this.server.BordCast(this,"offline!!")
	
}
//给当前user发送消息
func (this User)SendMessage(msg string)  {
	this.conn.Write([]byte(msg))
}


func (this *User)DoMessage(msg string)  {
	if msg == "who" {
		//查询当前在线用户
		this.server.mapLock.Lock()
		for _,user := range this.server.OnlineMap{
			onlineMsg := "["+user.Addr+"]"+user.Name+" online\n"
			this.SendMessage(onlineMsg)

		}

		this.server.mapLock.Unlock()

		
	}else{
		this.server.BordCast(this,msg)
	}

	
}

//监听User channel 的方法，一旦有消息就发送到客户端
func (this *User)ListMessage()  {
	for  {
		msg := <-this.C
		this.conn.Write([]byte(msg+"\n\r\n\r"))
	}
}

