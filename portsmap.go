package main

import (
	"fmt"
	"net"
	"sync"

	"gopkg.in/ini.v1"
)

// read缓存容量
const Buffersize = 2048

// 映射类型
type PortMap struct {
	Name   string `ini:"name"`
	Addr   string `ini:"addr"`
	Map    string `ini:"map"`
	Active bool   `ini:"active"`
}

var Portmaplist = []PortMap{}                // 映射列表
var Portsmap = make(map[string]string)       // 查重后 // addr查重
var Portsmapvals = make(map[string]struct{}) // mapaddr查重

// 初始化，解析配置文件
func init() {
	fmt.Println("reading portsmap.ini ...")
	cfg, err := ini.Load("portsmap.ini")
	if err != nil {
		fmt.Println("read portsmap.ini err", err)
		return
	}
	// fmt.Println("config:", cfg.SectionStrings())
	// 收集Portmaplist
	for _, key := range cfg.SectionStrings() {
		if key == "DEFAULT" {
			addrs := cfg.Section(key).KeyStrings()
			for _, addr := range addrs {
				var portmap = PortMap{
					Name:   "",
					Addr:   addr,
					Map:    cfg.Section(key).Key(addr).String(),
					Active: true,
				}
				Portmaplist = append(Portmaplist, portmap)
			}
		} else {
			var portmap = &PortMap{
				Name:   "",
				Addr:   "",
				Map:    "",
				Active: true,
			}
			err := cfg.Section(key).MapTo(portmap)
			if err != nil {
				fmt.Println("mapto portmap err", key, err)
			} else {
				Portmaplist = append(Portmaplist, *portmap)
			}
		}
	}
	// fmt.Println("Portmaplist", Portmaplist)
	for _, portmap := range Portmaplist {
		if portmap.Active {
			if _, ok := Portsmap[portmap.Addr]; !ok {
				if _, ok := Portsmapvals[portmap.Map]; !ok {
					Portsmap[portmap.Addr] = portmap.Map
					Portsmapvals[portmap.Map] = struct{}{}
				} else {
					fmt.Println("映射地址重复: ", portmap.Name, portmap.Addr, "=>", portmap.Map)
				}
			} else {
				fmt.Println("原始地址重复: ", portmap.Name, portmap.Addr, "=>", portmap.Map)
			}
		}
	}
	// fmt.Println("Portsmap", Portsmap)
}

func main() {
	fmt.Println("\nstart portsmap ...")
	var wg = sync.WaitGroup{} // 计数信号量，用来记录并维护运行的 goroutine
	for addr, mapaddr := range Portsmap {
		// 监听mapaddr
		maplisten, err := net.Listen("tcp", mapaddr)
		if err != nil {
			fmt.Println("映射地址监听失败: ", mapaddr, err)
			continue
		}
		go func(maplisten net.Listener, addr string, mapaddr string) {
			defer maplisten.Close()
			defer wg.Done()
			fmt.Println("地址映射建立: ", addr, "=>", mapaddr)
			for {
				// 等待连接mapaddr
				mapconn, err := maplisten.Accept()
				if err != nil {
					fmt.Println("映射连接建立失败: ", mapaddr, err)
					break
				}
				// 主动连接addr
				addrconn, err := net.Dial("tcp", addr)
				if err != nil {
					fmt.Println("原始地址连接失败: ", addr, err)
					mapconn.Close()
					break
				}
				// 建立工作连接
				workbind(mapconn, addrconn)
			}
			fmt.Println("映射地址监听结束: ", mapaddr)
		}(maplisten, addr, mapaddr)
		wg.Add(1)
	}
	// select {} // 据说这种阻塞方式不好
	// wg.Add(len(Portsmap))
	wg.Wait()
	fmt.Println("\nquit portsmap ...")
}

// 建立工作连接
func workbind(conn1, conn2 net.Conn) {
	var bindconn = func(readconn, writeconn net.Conn, middleware func([]byte) []byte) {
		defer readconn.Close()
		defer writeconn.Close()
		for {
			// read
			var buffer [Buffersize]byte
			bufsize, readerr := readconn.Read(buffer[:])
			if readerr != nil {
				fmt.Println("read err:", readerr)
				break
			}
			// middleware
			msg := middleware(buffer[:bufsize])
			// write
			_, writeerr := writeconn.Write(msg)
			if writeerr != nil {
				fmt.Println("write err:", writeerr)
				break
			}
		}
	}
	go bindconn(conn1, conn2, decode)
	go bindconn(conn2, conn1, encode)
}

// 加密
func encode(msg []byte) []byte {
	return msg
}

// 解密
func decode(msg []byte) []byte {
	return msg
}
