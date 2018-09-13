package mynet

import (
	"fmt"
	"log"
	"net"
	"shoppingzone/conf"
	inner "shoppingzone/mylib/myinnercommand"
	"shoppingzone/mylib/mylog"
	"shoppingzone/myutil"
	"strconv"
	"time"
)

//Server : struct
type Server struct{}

//Client : struct
type Client struct{}

var conns = make(chan net.Conn, conf.Config.NetMaxConns)

/*
Start :
@Param {int} p  123
@Param {string}  p ":123"
*/
func (s *Server) Start(p interface{}) net.Listener {
	var port string
	if p == nil || p == "" {
		port = conf.Config.NetPort
	} else if v, ok := p.(string); ok {
		port = v
		mylog.Tf("[Info]", "Server", "Start", "Success manually set the port %s ", port)
	} else if i, ok := p.(int); ok {
		port = `:` + strconv.Itoa(i)
		mylog.Tf("[Info]", "Server", "Start", "Success manually set the port %s ", port)
	}
	listenner, err := net.Listen("tcp", port)
	if err != nil {
		mylog.Tf("[Error]", "Server", "Start", "Fail to start net server listenning on port %s %s", port, err.Error())
		return nil
	}
	mylog.Tf("[Info]", "Server", "Start", "Success to start net server listenning on port %s ", port)
	return listenner
}

/*
GetConn :
*/
func (c *Client) GetConn(addr string) {
	go func() {
		fmt.Println("lc :", len(conns))
		var port string
		if addr == "" {
			port = conf.Config.NetPort
		} else {
			port = addr
			mylog.Tf("[Info]", "Client", "GetConn", "Success manually set the port %s ", port)
		}
		Conn, err := net.Dial("tcp", port)
		if err != nil {
			mylog.Tf("[Error]", "Client", "GetConn", "Fail to conncent on port %s %s", port, err.Error())
		} else {
			mylog.Tf("[Info]", "Client", "GetConn", "Success to connect port %s ", port)
			conns <- Conn
			fmt.Println("lc :", len(conns))
		}
	}()
}

/*
Send :
*/
func (c *Client) Send(m string) {
	c.GetConn(conf.Config.NetPort)
	conn := <-conns
	defer conn.Close()
	buff := []byte(m)
	n, _ := conn.Write(buff)
	log.Println(n)
}

/////////////////////////////////
////////////////////////////////
///////////////////////////////

//MyNet :
type MyNet struct {
	p *Protocol
}

//Message :
type Message struct {
	c net.Conn
	m string
}

//ConnHandel : handel of server listenner
func (n *MyNet) ConnHandel(conn net.Conn) {
	tmpBuffer := make([]byte, 0)
	buff := make([]byte, 2048)
	for {
		i, err := conn.Read(buff)
		if err != nil {
			mylog.Tf("[Error]", "MyNet", "ConnHandel", "failed to read from conn! %s", err.Error())
			return
		}
		hb := make(chan byte)
		go n.HeartBeat(conn, hb, 10)
		go byte2chl(buff[:i], hb)
		tmpBuffer = n.p.DePackage(append(tmpBuffer, buff[:i]...))
		pms := n.ProcessMessage(Message{conn, string(tmpBuffer)})
		wms := n.WriteMessage(pms)
		for m := range wms {
			mylog.Tf("[Info]", "MyNet", "ConnHandel", "%s", m)
		}
	}
}

//ProcessMessage :
func (n *MyNet) ProcessMessage(message Message) <-chan Message {
	out := make(chan Message)
	go func() {
		mylog.Tf("[Info]", "MyNet", "ProcessMessage", "processe massage : %s", message.m)
		s2j := myutil.String2JSON(message.m)
		for m := range n.Exec(s2j) {
			mylog.Tf("[Info]", "MyNet", "ProcessMessage", "afert pms : %s", message.m)
			out <- Message{message.c, m}
		}
		close(out)
	}()
	return out
}

//WriteMessage :
func (n *MyNet) WriteMessage(messages <-chan Message) <-chan string {
	out := make(chan string)
	go func() {
		for mess := range messages {
			_, err := mess.c.Write([]byte(mess.m))
			if err != nil {
				mylog.Tf("[Error]", "MyNet", "WriteMessage", "failed to write message to conn. %s", err.Error())
				out <- "no"
			}
			out <- "yes"
		}
		close(out)
	}()
	return out
}

//Exec : executive the message command
func (n *MyNet) Exec(m <-chan myutil.Message) <-chan string {
	out := make(chan string)
	go func() {
		for cds := range m {
			switch cds.Mothed {
			case "checkUser":
				mylog.Tf("[Info]", "MyNet", "Exec", "%s", cds.Data)
				t := cds.Data.(map[string]interface{})
				chkN, ok1 := <-inner.CheckUserNameIsNull(t["Username"].(string))
				chkM, ok2 := <-inner.CheckUserMobileIsNull(t["Mobilenumber"].(string))
				for ok1 || ok2 {
					if !chkN {
						mylog.Tf("[Info]", "MyNet", "Exec", "Failed to use this username , uername already used")
						out <- `{Mothed:"checkUser",Data:{Status:false,Message:"can't register"}}`
						break
					}
					if !chkM {
						mylog.Tf("[Info]", "MyNet", "Exec", "Failed to use this mobile , mobile already used")
						out <- `{Mothed:"checkUser",Data:{Status:false,Message:"can't register"}}`
						break
					}
					if ok1 && ok2 && chkN && chkM {
						out <- `{Mothed:"checkUser",Data:{Status:true,Message:"can register"}}`
						break
					}
				}
			default:
				mylog.Tf("[Info]", "MyNet", "Exec", "Haven't this inner command")
			}
		}
		close(out)
	}()
	return out
}

//HeartBeat : test connect
func (n *MyNet) HeartBeat(conn net.Conn, rchl chan byte, timeout int) {
	select {
	case <-rchl:
		mylog.Tf("[Info]", "MyNet", "HeartBeat", "%s  message recived!", conn.RemoteAddr().String())
		conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Minute))
		break
	case <-time.After(time.Second * 5):
		mylog.Tf("[Info]", "MyNet", "HeartBeat", "%s  timeout to close!", conn.RemoteAddr().String())
		conn.Close()
	}
}

func byte2chl(b []byte, chl chan byte) {
	for _, v := range b {
		chl <- v
	}
	close(chl)
}
