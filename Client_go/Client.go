package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct{
	ServerIp	string
	ServerPort 	int
	Name		string
	conn 		net.Conn
	flag		int //模式

}

func NewClint(serverIp string,serverPort int) *Client {
	//创建客户端对象
	client := &Client{
		ServerIp: serverIp,
		ServerPort: serverPort,
		flag: 999,
	}

	//连接server
	conn,err :=net.Dial("tcp",fmt.Sprintf("%s:%d",serverIp,serverPort))
	if err!=nil {
		fmt.Println("net.dial err:",err)
		return nil
	}
	client.conn=conn
	return client



	//返回对象


}
func (client *Client)menu() bool {
	var flag int
	fmt.Println("1.群聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)

	if flag>=0 && flag <=3{
		client.flag=flag
		return true


	}else {
		fmt.Println(">>>给出合法范围的数字<<<")
		return false
	}

}
func (client *Client)UpdateName()bool{
	fmt.Println(">>输入用户名<<")
	fmt.Scanln(&client.Name)
	sendMsg := "rename|"+client.Name+"\n"
	_,err := client.conn.Write([]byte(sendMsg))
	if err!=nil{
		fmt.Println("conn.write error",err)
		return false

	}
	return true
	
}
func (client *Client)PublicChat()  {
	var chatMsg string
	fmt.Println(">>请输入聊天内容,exit退出<<")
	fmt.Scanln(&chatMsg)
	for chatMsg !="exit"{
		//发送给服务器消息
		//不为空则发送
		if len(chatMsg)!=0{
			sendMsg := chatMsg+"\n\r"
			_,err:=client.conn.Write([]byte(sendMsg))
			if err!=nil {
				fmt.Println("conn write err:",err)
				break
			}

		}
		chatMsg=""
		fmt.Println(">>请输入聊天内容,exit退出<<")
		fmt.Scanln(&chatMsg)

	}

}
//查询在线用户
func (client *Client)SelectUsers()  {
	sendMsg:="who\n"
	_,err:=client.conn.Write([]byte(sendMsg))
	if err!=nil{
		fmt.Println("conn write err",err)
		return
	}



}
func (client *Client)PrivateChat()  {
	var RemoteName string
	var ChatMsg string
	//查询在线用户
	//展示在线用户
	client.SelectUsers()

	//提示用户从中挑选目标用户
	fmt.Println(">>请输入聊天对象[用户名],exit 退出<<")
	fmt.Scanln(&RemoteName)

	if RemoteName!="exit" {

		fmt.Println("请输入消息内容.exit退出")
		fmt.Scanln(&ChatMsg)

		for ChatMsg!="exit" {
			if len(ChatMsg)!=0{
				sendMsg := "to|"+RemoteName+"|"+ChatMsg+"\n\r"
				_,err:=client.conn.Write([]byte(sendMsg))
				if err!=nil {
					fmt.Println("conn write err:",err)
					break
				}

			}
			ChatMsg=""
			fmt.Println(">>请输入聊天内容,exit退出<<")
			fmt.Scanln(&ChatMsg)

		}
		client.SelectUsers()
		fmt.Println(">>请输入聊天对象[用户名],exit 退出<<")
		fmt.Scanln(&RemoteName)

	}



}



//处理回应消息，显示到标准输出即可
func (client *Client)DealRepsone()  {

	//for  {
	//	buff := make()
	//	client.conn.Read(buff)
	//	fmt.Println(buff)
	//}
	//等效io.copy


	io.Copy(os.Stdout,client.conn)
	//一旦io.conn有数据就拷贝到 os.stdout
	//永久阻塞监听

}
func (client *Client)Run()  {
	for client.flag !=0{
		for client.menu()!=true{
		}
		switch client.flag {
		case 1://群聊模式
			fmt.Println("使用群聊模式。。。。。")
			client.PublicChat()


			break

		case 2://私聊模式
			fmt.Println("使用私聊模式。。。。。")
			client.PrivateChat()
			break

		case 3://rename username
			fmt.Println("修改用户名。。。。。。")
			client.UpdateName()
			break

		}
	}
	
}

var serverIp string
var serverPort int
//client -ip 127.0.0.1 -port
func init(){
	flag.StringVar(&serverIp,"ip","127.0.0.1","设置服务器地址，默认为 127.0.0.1")
	flag.IntVar(&serverPort,"port",8888,"设置默认连接端口，默认为8888")



}
func main(){
	//命令行解析
	flag.Parse()
	client:=NewClint(serverIp,serverPort)
	if client==nil{
		fmt.Println("connection fault\n")
		return
	}


	//单独启动监听go程
	go client.DealRepsone()

	fmt.Println("conection success\n")

	//启动服务端业务
	client.Run()


}