package main

import "net"

type User struct{
	Name 	string
	Addr 	string
	C 		chan string
	conn 	net.Conn
}

//create user function

func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:	  make(chan string),
		conn: conn,
	}
	//启动监听当前user 的channel 的goroutinue
	go user.ListMessage()
	return user
}

//监听User channel 的方法，一旦有消息就发送到客户端
func (this *User)ListMessage()  {
	for  {
		msg := <-this.C
		this.conn.Write([]byte(msg+"\n\r\n\r"))
	}
}

