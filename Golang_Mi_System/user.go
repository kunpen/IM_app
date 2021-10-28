package main

import (
	"net"
	"strings"
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

		
	}else if len(msg)>7&&msg[:7] == "rename|"  {
		//修改username 指令为 rename|张三
		newName := strings.Split(msg,"|")[1]
		//判断用户名是否被占用
		_,ok:=this.server.OnlineMap[newName]
		if ok {
			this.SendMessage("username is used\n")
		}else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap,this.Name)
			this.server.OnlineMap[newName]=this
			this.server.mapLock.Unlock()
			this.Name=newName
			this.SendMessage("username changed\n")

		}


	}else if len(msg)>4 &&msg[:3] == "to|"{
		//格式 to|张三|消息内容
		//获取username
		remoteName := strings.Split(msg,"|")[1]
		if remoteName==""{
			this.SendMessage("message correct!use  to|username|message\n")
			return

		}

		//根据username 获取user对象
		remoteUser,ok := this.server.OnlineMap[remoteName]
		if !ok {
			this.SendMessage(" invide user\n")
			return
		}
		//获取消息内容，根据user对象将消息发送出去
		content := strings.Split(msg,"|")[2]
		if content=="" {
			this.SendMessage(" empty message\n")
			return
		}
		remoteUser.SendMessage(this.Name+" said to u:"+content)



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

