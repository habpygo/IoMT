/* Copyright 2018 Harry Boer - DappDevelopment.com
Licensed under the MIT License, see LICENCE file for details.
*/

/* This package is called to
1. calculate the Merkle Root
2. verify the hashes of all the nodes that make up the root
3. verify the leaf under observation
*/

package merkletree

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
)

type Message interface {
	CheckSum() []byte
	Equals(other Message) bool
}

type MerkleStructMessage interface {
	//PrepareForHashing() []byte
	CalculateHash() []byte
	EqualsStruct(other MerkleStructMessage) bool
}

// Node contains data and its siblings
type Node struct {
	ParentNode, LeftNode, RightNode *Node
	Hash                            []byte // the Node Hash
	M                               Message
	leaf                            bool
}

type MerkleNode struct {
	ParentNode, LeftNode, RightNode *MerkleNode
	Hash                            []byte // the Node Hash
	M                               MerkleStructMessage
	leaf                            bool
}

type MerkleTree struct {
	MerkleRootNode *Node
	MerkleRoot     []byte // the Root Hash
	Leaves         []*Node
}

type MerkleStructTree struct {
	MerkleStructRootNode *MerkleNode
	MerkleStructRoot     []byte // the Root Hash
	Leaves               []*MerkleNode
}

// NewTree builds the tree with the supplied messages using two helper funcs
func NewTree(msg []Message) (*MerkleTree, error) {
	merkleRoot, leaves, err := buildFromMsg(msg)
	if err != nil {
		return nil, fmt.Errorf("Error: couldn't create root and leafs from func buildFromMsg, %v: ", err)
	}

	mt := &MerkleTree{
		MerkleRootNode: merkleRoot,
		MerkleRoot:     merkleRoot.Hash,
		Leaves:         leaves,
	}

	return mt, nil
}

// NewMerkleTree builds the tree with the supplied messages using two helper funcs
func NewMerkleTree(msg []MerkleStructMessage) (*MerkleStructTree, error) {
	merkleStructRoot, leaves, err := buildFromMerkleMsg(msg)
	if err != nil {
		return nil, fmt.Errorf("Error: couldn't create root and leafs from func buildFromMsg, %v: ", err)
	}

	mt := &MerkleStructTree{
		MerkleStructRootNode: merkleStructRoot,
		MerkleStructRoot:     merkleStructRoot.Hash,
		Leaves:               leaves,
	}

	return mt, nil
}

func buildFromMsg(msg []Message) (*Node, []*Node, error) {
	if len(msg) == 0 {
		return nil, nil, errors.New("Error: There are no messages, messsage byte slice is empty")
	}
	var leaves []*Node
	for _, message := range msg {
		leaves = append(leaves, &Node{
			Hash: message.CheckSum(),
			M:    message,
			leaf: true,
		})
	}

	merkleRoot := FindRoot(leaves)
	return merkleRoot, leaves, nil
}

func buildFromMerkleMsg(msg []MerkleStructMessage) (*MerkleNode, []*MerkleNode, error) {
	if len(msg) == 0 {
		return nil, nil, errors.New("Error: There are no messages, messsage byte slice is empty")
	}
	var leaves []*MerkleNode
	for _, message := range msg {
		leaves = append(leaves, &MerkleNode{
			Hash: message.CalculateHash(),
			M:    message,
			leaf: true,
		})
	}

	merkleRoot := FindMerkleRoot(leaves)
	return merkleRoot, leaves, nil
}

// FindRoot takes a node list and retrieves the root node
func FindRoot(nodeList []*Node) *Node {
	if len(nodeList)%2 != 0 {
		nodeList = append(nodeList, nodeList[len(nodeList)-1])
	}
	var nodes []*Node
	for i := 0; i < len(nodeList); i += 2 {
		fmt.Println("i is: ", i, "-------------->")
		h := sha256.New()
		msgHash := append(nodeList[i].Hash, nodeList[i+1].Hash...)
		h.Write(msgHash)
		node := &Node{
			LeftNode:  nodeList[i],
			RightNode: nodeList[i+1],
			Hash:      h.Sum(nil),
		}
		nodes = append(nodes, node)
		nodeList[i].ParentNode = node
		nodeList[i+1].ParentNode = node
		if len(nodeList) == 2 {
			return node
		}
	}
	fmt.Println("******* Nodes is returned and the content is: ", nodes)
	return FindRoot(nodes)
}

func FindMerkleRoot(nodeList []*MerkleNode) *MerkleNode {
	if len(nodeList)%2 != 0 {
		nodeList = append(nodeList, nodeList[len(nodeList)-1])
	}
	var nodes []*MerkleNode
	for i := 0; i < len(nodeList); i += 2 {
		fmt.Println("i is: ", i, "-------------->")
		h := sha256.New()
		msgHash := append(nodeList[i].Hash, nodeList[i+1].Hash...)
		h.Write(msgHash)
		node := &MerkleNode{
			LeftNode:  nodeList[i],
			RightNode: nodeList[i+1],
			Hash:      h.Sum(nil),
		}
		nodes = append(nodes, node)
		nodeList[i].ParentNode = node
		nodeList[i+1].ParentNode = node
		if len(nodeList) == 2 {
			return node
		}
	}
	fmt.Println("******* Nodes is returned and the content is: ", nodes)
	return FindMerkleRoot(nodes)
}

//MerkleRoot returns the unverified Merkle Root Hash.
func (mt *MerkleTree) MerkleRootHash() []byte {
	return mt.MerkleRoot
}

func (mst *MerkleStructTree) MerkleStructRootHash() []byte {
	return mst.MerkleStructRoot
}

// CheckNode verifies each leaf node and returns the hash of the leaf
func (node *Node) checkSingleNode() []byte {
	if node.leaf {
		return node.M.CheckSum()
	}
	h := sha256.New()
	h.Write(append(node.LeftNode.checkSingleNode(), node.RightNode.checkSingleNode()...))
	return h.Sum(nil)
}

// CheckNode verifies each leaf node and returns the hash of the leaf
func (mnode *MerkleNode) checkSingleStructNode() []byte {
	if mnode.leaf {
		return mnode.M.CalculateHash()
	}
	h := sha256.New()
	h.Write(append(mnode.LeftNode.checkSingleStructNode(), mnode.RightNode.checkSingleStructNode()...))
	return h.Sum(nil)
}

func (node *Node) calculateNodeHash() []byte {
	if node.leaf {
		return node.M.CheckSum()
	}
	h := sha256.New()
	h.Write(append(node.LeftNode.Hash, node.RightNode.Hash...))
	return h.Sum(nil)
}

func (mnode *MerkleNode) calculateStructNodeHash() []byte {
	if mnode.leaf {
		return mnode.M.CalculateHash()
	}
	h := sha256.New()
	h.Write(append(mnode.LeftNode.Hash, mnode.RightNode.Hash...))
	return h.Sum(nil)
}

func (mt *MerkleTree) CheckMerkleTree() bool {
	calculatedMerkleRootHash := mt.MerkleRootNode.checkSingleNode()
	if bytes.Compare(mt.MerkleRoot, calculatedMerkleRootHash) == 0 {
		return true
	}
	return false
}

func (mst *MerkleStructTree) CheckMerkleStructTree() bool {
	calculatedMerkleStructRootHash := mst.MerkleStructRootNode.checkSingleStructNode()
	if bytes.Compare(mst.MerkleStructRoot, calculatedMerkleStructRootHash) == 0 {
		return true
	}
	return false
}

func (mt *MerkleTree) CheckIfContentIsAuthentic(authenticMerkleRoot []byte, message Message) bool {
	for _, leaf := range mt.Leaves {
		if leaf.M.Equals(message) {
			fmt.Println("messages are equal, leaf.ParentNode is: ", leaf.ParentNode)
			leafParent := leaf.ParentNode
			for leafParent != nil {
				h := sha256.New()
				if leafParent.LeftNode.leaf && leafParent.RightNode.leaf {
					fmt.Println("Parent.LeftNode & RightNode exist. ")
					fmt.Println("leafParent.LeftNode.leaf is: ", leafParent.LeftNode.leaf)
					fmt.Println("leafParent.RightNode.leaf is: ", leafParent.RightNode.leaf)
					h.Write(append(leafParent.LeftNode.calculateNodeHash(), leafParent.RightNode.calculateNodeHash()...))
					if bytes.Compare(h.Sum(nil), leafParent.Hash) != 0 {
						fmt.Println("bytes.Compare is false")
						return false
					}
					leafParent = leafParent.ParentNode
				} else {
					h.Write(append(leafParent.LeftNode.calculateNodeHash(), leafParent.RightNode.calculateNodeHash()...))
					if bytes.Compare(h.Sum(nil), leafParent.Hash) != 0 {
						fmt.Println("bytes.Compare from merklewebmamgiota is: false")
						return false
					}
					leafParent = leafParent.ParentNode
				}
			}
			return true
		}

	}
	return false
}

func (mst *MerkleStructTree) CheckIfStructContentIsAuthentic(authenticMerkleRoot []byte, merkleSM MerkleStructMessage) bool {
	for _, leaf := range mst.Leaves {
		if leaf.M.EqualsStruct(merkleSM) {
			leafParent := leaf.ParentNode
			for leafParent != nil {
				h := sha256.New()
				if leafParent.LeftNode.leaf && leafParent.RightNode.leaf {
					h.Write(append(leafParent.LeftNode.calculateStructNodeHash(), leafParent.RightNode.calculateStructNodeHash()...))
					if bytes.Compare(h.Sum(nil), leafParent.Hash) != 0 {
						fmt.Println("bytes.Compare is false")
						return false
					}
					leafParent = leafParent.ParentNode
				} else {
					h.Write(append(leafParent.LeftNode.calculateStructNodeHash(), leafParent.RightNode.calculateStructNodeHash()...))
					if bytes.Compare(h.Sum(nil), leafParent.Hash) != 0 {
						fmt.Println("bytes.Compare from merklewebmamgiota is: false")
						return false
					}
					leafParent = leafParent.ParentNode
				}
			}
			return true
		}

	}
	return false
}

/**************************** Utils for merkle tree ***********************/
// IsHashListEmpty checks whether there is a Merkle tree to begin with
func HashListEmpty(leaflist []*Node) (bool, error) {
	if len(leaflist) == 0 {
		return true, errors.New("Error: cannot construct tree with no content.")
	}
	return false, nil
}

// FixNumberOfHashesOdd checks whether there is a leaf missing,
// if not return original list else duplicate last hash
func FixNumberOfOddLeaves(leaflist []*Node) ([]*Node, error) {

	if len(leaflist)%2 == 0 {
		return leaflist, nil
	}
	duplicateLeaf := leaflist[len(leaflist)-1]
	leaflist = append(leaflist, duplicateLeaf) // e.g. leaflilst now becomes 1 2 3 4 5 5, an odd number of leafs as binary trees should be

	return leaflist, nil
}

type Content struct {
	c string
}

// CheckSum hashes the values of the message content
func (content Content) CheckSum() []byte {
	h := sha256.New()
	h.Write([]byte(content.c))
	return h.Sum(nil)
}

//Equals tests for equality of two Contents
func (content Content) Equals(other Content) bool {
	return content.c == other.c
}
