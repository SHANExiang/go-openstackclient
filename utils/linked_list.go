package utils


type Node struct {
	data      string
	next      *Node
}

type LinkedList struct {
	head        *Node
}

func (list *LinkedList) addNode(data string) {
	newNode := &Node{data: data}

	if list.head == nil {
		list.head = newNode
	} else {
		current := list.head

		for current.next != nil {
			current = current.next
		}
		current.next = newNode
	}
}
