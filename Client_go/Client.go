package main

import (
	"flag"
	"fmt"
	"net"
)

type Clinet struct{
	ServerIp	string
	ServerPort 	int
	Name		string
	conn 		net.Conn
	flag		int //模式

}

func NewClint(serverIp string,serverPort int) *Clinet {
	//创建客户端对象
	client := &Clinet{
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
func (client *Clinet)menu() bool {
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
func (client *Clinet)Run()  {
	for client.flag !=0{
		for client.menu()!=true{
		}
		switch client.flag {
		case 1://群聊模式
			fmt.Println("使用群聊模式。。。。。")


			break

		case 2://私聊模式
			fmt.Println("使用私聊模式。。。。。")
			break

		case 3://rename username
			fmt.Println("修改用户名。。。。。。")
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
	fmt.Println("conection success\n")

	//启动服务端业务
	client.Run()


}