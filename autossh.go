package main

import (
	"os"
	"fmt"
	"net"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	if len(os.Args) == 5 {
		config := &ssh.ClientConfig{
		    User: os.Args[3],
		    Auth: []ssh.AuthMethod{
		        ssh.Password(os.Args[4]),
		    },
		    HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
		}
		addr := os.Args[1] + ":" + os.Args[2]
		client, err := ssh.Dial("tcp", addr, config)
		if err != nil {
			fmt.Println("建立连接出错")
			return
		}
		defer client.Close()

		session, err := client.NewSession()
		if err != nil {
			fmt.Println("创建Session出错")
			return
		}
		defer session.Close()
		
		fd := int(os.Stdin.Fd())
		oldState, err := terminal.MakeRaw(fd)
		if err != nil {
			fmt.Println("创建文件描述符出错")
			return
		}

		session.Stdout = os.Stdout
		session.Stderr = os.Stderr
		session.Stdin = os.Stdin

		termWidth, termHeight, err := terminal.GetSize(fd)
		if err != nil {
			fmt.Println("获取窗口宽高出错")
			return
		}

		defer terminal.Restore(fd, oldState)

		modes := ssh.TerminalModes{
			ssh.ECHO:          1,
			ssh.TTY_OP_ISPEED: 14400,
			ssh.TTY_OP_OSPEED: 14400,
		}

		if err := session.RequestPty("xterm-256color", termHeight, termWidth, modes); err != nil {
			fmt.Println("创建终端出错")
			return
		}

		err = session.Shell()
		if err != nil {
			fmt.Println("执行Shell出错")
			return
		}

		err = session.Wait()
		if err != nil {
			// fmt.Println("执行Wait出错")
			return
		}
	} else {
		fmt.Printf("参数个数不对,请依次输入ip port username password")
	}
}
