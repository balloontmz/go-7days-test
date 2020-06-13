//Package gee trie for route match
package gee

import (
	"strings"
)

type node struct {
	pattern  string  // 待匹配路由
	part     string  //路由中的一部分
	children []*node // 子节点 例如 [doc, tutorial, intro]
	isWild   bool    //是否精确匹配, part 含有 : 或 * 时为 true
}

//第一个匹配成功的节点,用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

//所有匹配成功的节点,用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern // 只有在匹配到最底层才添加 path,如果底层是 wild,则多次赋值也没关系
		return
	}
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

//这里其实没有判断 wild 和 named 共同存在的情况 -- child 是有序的切片,所以如果按照顺序,按理来说没问题??? 路由命名注意先具名再切片!!!
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") { // *的情况
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		// fmt.Print("当前遍历的 child 为:", child, "\n")
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
