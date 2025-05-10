package gee

import (
	"net/http"
	"strings"
)

// router 是路由器的核心结构，包含路由树(roots) 和路由处理函数表(handlers).
type router struct {
	roots    map[string]*node       // 每种 HTTP 方法对应一颗 Trie 树.
	handlers map[string]HandlerFunc // 路由处理函数，键是 "GET-/p/:lang" 的形式.
}

// newRouter 创建并初始化一个 router 实例.
func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// parsePattern 将 pattern 按照 '/' 分割成片段数组，处理动态参数和通配符
// 最多允许一个 '*' 结尾.
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/") // 例如 "/p/:lang" -> ["", "p",":lang"]

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			// 如果是通配符，后续则无需再分片.
			if item[0] == '*' {
				break
			}
		}
	}

	return parts
}

// addRoute 添加一个新路由规则，
// method 是请求方法(GET/POST 等)，
// pattern 是 URL 路径，
// handler 是处理函数.
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern) // 解析路径为切片
	key := method + "-" + pattern  // 用于 handlers 映射，例如 "GET-/p/:lang"

	// 若该 method 尚未创建根节点，初始化一颗新的 Trie 树.
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}

	// 将路由路径插入到对应方法的 Trie 树中.
	r.roots[method].insert(pattern, parts, 0)

	// 保存对应的处理函数.
	r.handlers[key] = handler
}

// getRoute 根据请求方法和 URL 路径查找路由节点，返回匹配的 node 和参数映射.
func (r *router) getRoute(method, path string) (*node, map[string]string) {
	searchParts := parsePattern(path) // 将实际请求路径解析为切片.
	params := make(map[string]string) // 用于保存动态参数，如 {lang: "go"}

	root, ok := r.roots[method]
	if !ok {
		// 若不存在对应方法的 Trie 树，直接返回空.
		return nil, nil
	}

	// 在 Trie 树中查找匹配路径的节点.
	n := root.search(searchParts, 0)
	if n != nil {
		// 找到匹配的节点后，需要进一步提取动态参数.
		parts := parsePattern(n.pattern) // 取出注册时的原始切片.
		for index, part := range parts {
			if part[0] == ':' {
				// 如果是动态参数，记录为键值对，例如 ":lang" -> {"lang": "go"}.
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				// 如果是通配符，截取剩余路径，保存为参数.
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

// getRoutes 获取所有注册的路径节点(pattern 不为空)，用于调试或打印所有路径.
func (r *router) getRoutes(method string) []*node {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}
	nodes := make([]*node, 0)
	root.travel(&nodes) // 遍历整棵 Trie 树，收集所有注册了 pattern 的节点.
	return nodes
}

// handle 接收一个 Context，将请求交给对应的路由处理函数，或返回 404.
func (r *router) handle(c *Context) {
	// 根据请求方法和路径，查询匹配的路由节点及提取出的参数.
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		// 找到节点：将解析出的参数保存到 Context 中.
		c.Params = params
		// 构造 handler 在 handlers 表中的键：形如 "GET-/p/:lang"
		key := c.Method + "-" + n.pattern
		// 调用对应的处理函数，传入 Context.
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		// 没有匹配到任何路由：返回 404 错误.
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}
