package gee

import (
	"fmt"
	"strings"
)

// node 表示 Trie 的节点.
type node struct {
	pattern  string  // 待匹配路由，例如 "/p/:lang"，只有在叶子节点进行设置
	part     string  // 路由中的一部分，例如 "p"、":lang"、"*filepath"
	children []*node // 子节点集合，例如 [doc, tutorial, intor]、用于向下一层继续匹配
	isWild   bool    // 是否为动态路由
}

// String 实现了 fmt.Stringer 接口，方便打印调试节点信息.
func (n *node) String() string {
	// 返回节点的关键信息：完整路由、当前片段和值路由类型标志.
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}

// insert 负责将一个完整路由 pattern 安装 parts 切分后插入到匹配树中.
func (n *node) insert(pattern string, parts []string, height int) {
	// 如果已经递归到最后一个切片，则将当前节点标记为路由终点，保存完整 pattern.
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	// 否则，取当前层需要插入的片段.
	part := parts[height]
	// 在已有子节点中查找能匹配该片段的一个子节点（精确或模糊）.
	child := n.matchChild(part)
	if child == nil {
		// 如果不存在，就新建一个节点，并根据是否动态或通配符设置 isWild.
		child = &node{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*',
		}
		// 将新节点追加到 children 列表.
		n.children = append(n.children, child)
	}
	// 递归处理下一层片段.
	child.insert(pattern, parts, height+1)
}

// search 根据请求切片 parts 在匹配树中查找对应节点，height 为当前递归深度.
func (n *node) search(parts []string, height int) *node {
	// 如果到达切片末尾，或者当前节点是通配符(*)节点.
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		// 如果当该节点存储了完整的 pattern 时，才视为匹配成功.
		if n.pattern == "" {
			return nil
		}
		return n
	}

	// 获取当前要匹配的片段.
	part := parts[height]
	// 找出所有可匹配该片段的子节点（精确或模糊）.
	children := n.matchChildren(part)

	// 遍历每个匹配的子节点，递归继续向下匹配.
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	// 所有分支都未匹配成功，返回 nil.
	return nil
}

// travel 遍历以当前节点为根的整个子树，收集所有设置了 pattern 的节点（即路由叶子节点）
func (n *node) travel(list *([]*node)) {
	// 如果当前节点存储了完整路由，则加入结果列表.
	if n.pattern != "" {
		*list = append(*list, n)
	}

	// 递归遍历所有子节点.
	for _, child := range n.children {
		child.travel(list)
	}
}

// matchChild 在 children 中返回第一个能匹配 part 的子节点（精确或模糊）.
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		// 如果子节点 part 与目标相同，或子节点为动态/通配符(isWild)，即可匹配.
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// matchChildren 在 children 中返回所有能匹配 part 的子节点列表.
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		// 与 matchChild 相同的匹配条件，当返回
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}
