package response

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/Yalaouf/gostman/pkg/tui/style"
	"github.com/charmbracelet/lipgloss"
)

type NodeType int

const (
	NodeObject NodeType = iota
	NodeArray
	NodeString
	NodeNumber
	NodeBool
	NodeNull
)

type TreeNode struct {
	Key      string
	Value    interface{}
	Type     NodeType
	Children []*TreeNode
	Expanded bool
	Depth    int
}

type JSONTree struct {
	Root     *TreeNode
	cursor   int
	flatList []*TreeNode
	width    int
}

func NewJSONTree(jsonStr string) *JSONTree {
	tree := &JSONTree{}

	var data interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return nil
	}

	tree.Root = tree.buildTree("", data, 0)
	if tree.Root != nil {
		tree.Root.Expanded = true
	}
	tree.rebuildFlatList()

	return tree
}

func (t *JSONTree) buildTree(key string, value interface{}, depth int) *TreeNode {
	node := &TreeNode{
		Key:      key,
		Value:    value,
		Depth:    depth,
		Expanded: depth < 1,
	}

	switch v := value.(type) {
	case map[string]interface{}:
		node.Type = NodeObject
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			child := t.buildTree(k, v[k], depth+1)
			node.Children = append(node.Children, child)
		}
	case []interface{}:
		node.Type = NodeArray
		for i, item := range v {
			child := t.buildTree(fmt.Sprintf("%d", i), item, depth+1)
			node.Children = append(node.Children, child)
		}
	case string:
		node.Type = NodeString
	case float64:
		node.Type = NodeNumber
	case bool:
		node.Type = NodeBool
	case nil:
		node.Type = NodeNull
	}

	return node
}

func (t *JSONTree) rebuildFlatList() {
	t.flatList = nil
	if t.Root != nil {
		t.flattenNode(t.Root)
	}
}

func (t *JSONTree) flattenNode(node *TreeNode) {
	t.flatList = append(t.flatList, node)
	if node.Expanded {
		for _, child := range node.Children {
			t.flattenNode(child)
		}
	}
}

func (t *JSONTree) SetWidth(width int) {
	t.width = width
}

func (t *JSONTree) MoveUp() {
	if t.cursor > 0 {
		t.cursor--
	}
}

func (t *JSONTree) MoveDown() {
	if t.cursor < len(t.flatList)-1 {
		t.cursor++
	}
}

func (t *JSONTree) Toggle() {
	if t.cursor >= 0 && t.cursor < len(t.flatList) {
		node := t.flatList[t.cursor]
		if len(node.Children) > 0 {
			node.Expanded = !node.Expanded
			t.rebuildFlatList()
		}
	}
}

func (t *JSONTree) Expand() {
	if t.cursor >= 0 && t.cursor < len(t.flatList) {
		node := t.flatList[t.cursor]
		if len(node.Children) > 0 && !node.Expanded {
			node.Expanded = true
			t.rebuildFlatList()
		}
	}
}

func (t *JSONTree) Collapse() {
	if t.cursor >= 0 && t.cursor < len(t.flatList) {
		node := t.flatList[t.cursor]
		if len(node.Children) > 0 && node.Expanded {
			node.Expanded = false
			t.rebuildFlatList()
		}
	}
}

func (t *JSONTree) GetCursor() int {
	return t.cursor
}

func (t *JSONTree) ListLength() int {
	return len(t.flatList)
}

func (t *JSONTree) GetSelectedValue() string {
	if t.cursor < 0 || t.cursor >= len(t.flatList) {
		return ""
	}

	node := t.flatList[t.cursor]
	return t.nodeToJSON(node)
}

func (t *JSONTree) nodeToJSON(node *TreeNode) string {
	switch node.Type {
	case NodeObject, NodeArray:
		data, err := json.MarshalIndent(node.Value, "", "  ")
		if err != nil {
			return ""
		}
		return string(data)
	case NodeString:
		return node.Value.(string)
	case NodeNumber:
		return fmt.Sprintf("%v", node.Value)
	case NodeBool:
		return fmt.Sprintf("%v", node.Value)
	case NodeNull:
		return "null"
	}
	return ""
}

func (t *JSONTree) Render() string {
	if t.Root == nil || len(t.flatList) == 0 {
		return ""
	}

	var lines []string
	for i, node := range t.flatList {
		line := t.renderNode(node, i == t.cursor)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func (t *JSONTree) renderNode(node *TreeNode, selected bool) string {
	var sb strings.Builder

	indent := strings.Repeat("  ", node.Depth)
	sb.WriteString(indent)

	if len(node.Children) > 0 {
		if node.Expanded {
			sb.WriteString("▼ ")
		} else {
			sb.WriteString("▶ ")
		}
	} else {
		sb.WriteString("  ")
	}

	keyStyle := lipgloss.NewStyle().Foreground(style.ColorBlue)
	if node.Key != "" {
		sb.WriteString(keyStyle.Render(node.Key))
		sb.WriteString(": ")
	}

	sb.WriteString(t.renderValue(node))

	line := sb.String()

	if selected {
		selectedStyle := lipgloss.NewStyle().
			Background(style.ColorSurface).
			Foreground(style.ColorText)
		if t.width > 0 && len(line) < t.width-4 {
			line = line + strings.Repeat(" ", t.width-4-len(line))
		}
		line = selectedStyle.Render(line)
	}

	return line
}

func (t *JSONTree) renderValue(node *TreeNode) string {
	stringStyle := lipgloss.NewStyle().Foreground(style.ColorGreen)
	numberStyle := lipgloss.NewStyle().Foreground(style.ColorOrange)
	boolStyle := lipgloss.NewStyle().Foreground(style.ColorPurple)
	nullStyle := lipgloss.NewStyle().Foreground(style.ColorGray)
	typeStyle := lipgloss.NewStyle().Foreground(style.ColorGray)

	switch node.Type {
	case NodeObject:
		count := len(node.Children)
		return typeStyle.Render(fmt.Sprintf("{%d}", count))
	case NodeArray:
		count := len(node.Children)
		return typeStyle.Render(fmt.Sprintf("[%d]", count))
	case NodeString:
		val := node.Value.(string)
		if len(val) > 50 {
			val = val[:47] + "..."
		}
		return stringStyle.Render(fmt.Sprintf("\"%s\"", val))
	case NodeNumber:
		return numberStyle.Render(fmt.Sprintf("%v", node.Value))
	case NodeBool:
		return boolStyle.Render(fmt.Sprintf("%v", node.Value))
	case NodeNull:
		return nullStyle.Render("null")
	}
	return ""
}
