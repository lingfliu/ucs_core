package utils

//TODO: Not used, remove on release
type TreeNode struct {
	Parent *TreeNode
	Childs []*TreeNode
	Value  any
}

func (node *TreeNode) AddChild(child *TreeNode) {
	child.Parent = node
	node.Childs = append(node.Childs, child)
}

func (node *TreeNode) AddChildWithValue(value any) {
	child := &TreeNode{Value: value}
	child.Parent = node
	node.Childs = append(node.Childs, child)
}

func (node *TreeNode) AddChildWithValueAndReturn(value any) *TreeNode {
	child := &TreeNode{Value: value}
	child.Parent = node
	node.Childs = append(node.Childs, child)
	return child
}

/**
 * depth first search
 */
func (node *TreeNode) SearchChildWithValue(value any) *TreeNode {
	for _, child := range node.Childs {
		if child.Value == value {
			return child
		}
		if child.SearchChildWithValue(value) != nil {
			return child.SearchChildWithValue(value)
		}
	}
	return nil
}

/**
 * search the end node of the tree
 */
func (node *TreeNode) SearchEndNodeWithValue(value any) *TreeNode {
	if len(node.Childs) == 0 {
		return node
	}
	return node.Childs[0].SearchEndNodeWithValue(value)
}
