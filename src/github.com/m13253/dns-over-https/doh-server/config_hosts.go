// 文件名: config_hosts.go.go
//
// 作者: 张立丹 (zhangld@enlink.cn)
// 描述:
//
// Copyright 2018 (C) 2018 Enlink, Inc. All Right Reserved.
package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

type NameMap struct {
	Name  string              // 域名中每一级的名字, 如果为空则表明为根节点(.com 之前的节点).
	Addrs  []net.IP              // IP 地址(IPv4 或者 IPv6), .
	Items map[string]*NameMap // 所有下一级域名映射表.
	Cursor int                // 域名对应的ip地址的序号
}

// loadConfigHosts() - 加载域名 hosts 文件, 生成域名查询树.
//
// @path - 配置文件路径.
//
// 返回域名查询树及错误信息.
func (this *NameMap) loadConfigHosts(path string) (*NameMap, error) {
	this.Items = make(map[string]*NameMap) // 初始化根节点.

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader, lineNo := bufio.NewReader(file), 0
	for {
		lineNo += 1

		if line, err := reader.ReadString('\n'); len(line) != 0 {
			line = strings.TrimSpace(line)

			// 支持以 "#" 开头的单行注释.
			// 同时忽略所有空行.
			if strings.HasPrefix(line, "#") || len(line) == 0 {
				continue
			}

			if err := this.insertLine(line, lineNo); err != nil {
				return nil, err
			}
		} else {
			if err == io.EOF {
				return this, nil
			}
			return nil, err
		}
	}

	return this, nil
}

// insertLine() - 向查询树中插入一行记录, 字条串格式原始配置值.
//
// @line - 字符串配置信息(原始值)
// @lineNo - 行号(调试或日志使用).
//
// 返回错误信息.
func (this *NameMap) insertLine(line string, lineNo int) error {
	list := strings.Fields(line)

	if len(list) != 2 {
		msg := fmt.Sprintf("Invalid line(%d): %s", lineNo, line)
		return errors.New(msg)
	}

	ips, domain := strings.TrimSpace(list[0]), strings.TrimSpace(list[1])
	if len(ips) == 0 || len(domain) == 0 {
		msg := fmt.Sprintf("Invalid line(%d): %s", lineNo, line)
		return errors.New(msg)
	}

	strArrs := strings.Split(ips, ",")
	var addrs []net.IP
	for i := 0; i < len(strArrs); i++ {
		ipAddr := net.ParseIP(strArrs[i])
		if ipAddr == nil {
			msg := fmt.Sprintf("Invalid IP address in line(%d): %s", lineNo, strArrs[i])
			return errors.New(msg)
		}
		addrs = append(addrs, ipAddr)
	}


	tmpList, listName := strings.Split(domain, "."), make([]string, 0)
	for cc := len(tmpList); cc > 0; cc -- {
		name := strings.TrimSpace(tmpList[cc-1])
		if len(name) == 0 || (cc != 1 && name == "*") {
			msg := fmt.Sprintf("Invalid domain in line(%d): %s", lineNo, domain)
			return errors.New(msg)
		}

		listName = append(listName, name)
	}

	this.insertRecord(addrs, listName, -1)
	return nil
}

// insertRecord() - 向查询树中插入一条记录(解析过的信息)
//
// @addr  - 解析后的 IP 地址, 可以是 IPv4 也可以是 IPv6.
// @list  - 分级且倒序排列的域名.
// @level - 当前域名级别.
//
// 无返回值.
func (this *NameMap) insertRecord(addrs []net.IP, list []string, level int) {
	if level >= len(list) {
		return
	}

	// 非根节点, 记录域名.
	if level != -1 {
		this.Name = list[level]
	}

	// 叶子节点记录 IP 地址.
	if level+1 >= len(list) {
		this.Addrs = addrs
		this.Cursor = 0
		return
	}

	// 非叶子节点则增加下一级节点.
	next := list[level+1]
	if _, ok := this.Items[next]; !ok {
		this.Items[next] = this.createNameMap()
	}
	this.Items[next].insertRecord(addrs, list, level+1)

}

// findAddr() - 向查询树中插入一条记录(解析过的信息)
//
// @addr  - 解析后的 IP 地址, 可以是 IPv4 也可以是 IPv6.
// @list  - 分级且倒序排列的域名.
// @level - 当前域名级别.
//
// 无返回值.
func (this *NameMap) findAddr(domain string) net.IP {
	domain = strings.TrimSpace(domain)
	if strings.HasSuffix(domain, ".") {
		domain = strings.TrimRight(domain, ".")
	}

	// 分割并倒排序.
	tmpList, listName := strings.Split(domain, "."), make([]string, 0)
	for cc := len(tmpList); cc > 0; cc -- {
		listName = append(listName, strings.TrimSpace(tmpList[cc-1]))
	}

	// 根节点特殊处理.
	name := listName[0]
	if nameMap, ok := this.Items[name]; ok {
		if addr := nameMap.findAddrInner(listName, 0); addr != nil {
			return addr
		}
	}

	// 未查找到对应记录.
	// 如果存在 "*" 记录, 则返回对应 IP 地址.
	if nameMap, ok := this.Items["*"]; ok {
		cursor := nameMap.Cursor
		nameMap.Cursor++
		if cursor + 1 > len(nameMap.Addrs) {
			nameMap.Cursor = 0
		}
		return nameMap.Addrs[cursor]
	}

	return nil
}

func (this *NameMap) findAddrInner(list []string, level int) net.IP {
	if level >= len(list) {
		return nil
	}

	name := list[level]
	if this.Name != name {
		return nil // 当前级别域名不匹配.
	}

	// 完全匹配, 即要查询的域名匹配结束.
	if level+1 >= len(list) {
		cursor := this.Cursor
		this.Cursor++
		if cursor + 1 > len(this.Addrs) {
			this.Cursor = 0
		}
		return this.Addrs[cursor]
	}

	// 优先查询子节点.
	// 如果查询成功, 则直接返回对应地址, 查询结束.
	if nameMap, ok := this.Items[list[level+1]]; ok {
		if addr := nameMap.findAddrInner(list, level+1); addr != nil {
			return addr
		}
	}

	// 当前域名子域名如果存在 * 记录, 则返回对应的 IP 地址.
	if nameMap, ok := this.Items["*"]; ok {
		cursor := nameMap.Cursor
		nameMap.Cursor++
		if cursor + 1 > len(nameMap.Addrs) {
			nameMap.Cursor = 0
		}
		return nameMap.Addrs[cursor]
	}

	return nil
}

func (this *NameMap) printTree(level int) {
	spaces := ""
	for cc := 0; cc < level; cc++ {
		spaces += "  "
	}

	fmt.Printf("%sname: %s, addr: %s\n", spaces, this.Name, this.Addrs)
	for _, item := range this.Items {
		item.printTree(level + 1)
	}
}

func (this *NameMap) createNameMap() *NameMap {
	return &NameMap{Items: make(map[string]*NameMap)}
}
