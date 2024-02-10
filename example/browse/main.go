package main

import (
	"fmt"

	"github.com/huskar-t/opcda"
)

func main() {
	opcda.Initialize()
	defer opcda.Uninitialize()
	host := "localhost"
	progID := "Matrikon.OPC.Simulation.1"
	server, err := opcda.Connect(progID, host)
	if err != nil {
		panic(err)
	}
	browser, err := server.CreateBrowser()
	if err != nil {
		panic(err)
	}
	browser.MoveToRoot()
	root := &Tree{"root", nil, []*Tree{}, []Leaf{}}
	buildTree(browser, root)
	PrettyPrint(root)
}

type Tree struct {
	Name     string
	Parent   *Tree
	Branches []*Tree
	Leaves   []Leaf
}

type Leaf struct {
	Name string
	Tag  string
}

func buildTree(browser *opcda.OPCBrowser, branch *Tree) {
	err := browser.ShowLeafs(false)
	if err != nil {
		panic(err)
	}
	count := browser.GetCount()
	for i := 0; i < count; i++ {
		item, err := browser.Item(i)
		if err != nil {
			panic(err)
		}
		itemID, err := browser.GetItemID(item)
		if err != nil {
			panic(err)
		}
		l := Leaf{Name: item, Tag: itemID}
		branch.Leaves = append(branch.Leaves, l)
	}
	err = browser.ShowBranches()
	if err != nil {
		panic(err)
	}
	count = browser.GetCount()
	for i := 0; i < count; i++ {
		nextName, err := browser.Item(i)
		if err != nil {
			panic(err)
		}
		err = browser.MoveDown(nextName)
		if err != nil {
			panic(err)
		}
		nextBranch := &Tree{nextName, branch, []*Tree{}, []Leaf{}}
		branch.Branches = append(branch.Branches, nextBranch)
		buildTree(browser, nextBranch)
		err = browser.MoveUp()
		if err != nil {
			panic(err)
		}
		err = browser.ShowBranches()
		if err != nil {
			panic(err)
		}
	}
}

func PrettyPrint(tree *Tree) {
	fmt.Println(tree.Name)
	printSubtree(tree, 1)
}

func printSubtree(tree *Tree, level int) {
	space := ""
	for i := 0; i < level; i++ {
		space += "  "
	}
	for _, l := range tree.Leaves {
		fmt.Println(space, "-", l.Tag)
	}
	for _, b := range tree.Branches {
		fmt.Println(space, "+", b.Name)
		printSubtree(b, level+1)
	}
}
