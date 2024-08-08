package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/starboy011/Family-Tree-Backend/internal/db"
	"github.com/starboy011/Family-Tree-Backend/pkg"
)

// GetIdFromName fetches the ID from the database using the given name
func GetIdFromName(w http.ResponseWriter, r *http.Request, name string) int {
	db, err := db.InitDb(w, r)
	if err != nil {
		http.Error(w, "Error connecting to database", http.StatusInternalServerError)
		log.Printf("Error connecting to db: %v", err)
		return 0
	}
	defer db.Close()

	query := `SELECT "ID" FROM mulvansham WHERE "Relationship" = 1 AND "Name" = $1;`
	var id int
	err = db.QueryRow(query, name).Scan(&id)
	if err != nil {
		http.Error(w, "Error executing query", http.StatusInternalServerError)
		log.Printf("Error executing query: %v", err)
		return 0
	}

	return id
}

// ExtractTreeData retrieves the subtree rooted at a specific ID
func ExtractTreeData(id int) (*pkg.TreeNode, error) {
	db, err := db.InitDb(nil, nil)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}
	defer db.Close()

	query := `SELECT "ID", "Name", "ParentID" FROM mulvansham WHERE "Relationship" IN (1, 2)`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	var results []pkg.TreeResult

	for rows.Next() {
		var result pkg.TreeResult
		if err := rows.Scan(&result.ID, &result.Name, &result.ParentID); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		results = append(results, result)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through rows: %v", err)
	}

	tree := pkg.ConvertToTreefromID(results)
	subtree := pkg.ExtractSubtreeWithPath(tree, id)

	if subtree == nil {
		return nil, fmt.Errorf("node with ID %d not found", id)
	}

	return subtree, nil
}

// mergeTrees merges two trees into one
func mergeTrees(tree1, tree2 *pkg.TreeNode) *pkg.TreeNode {
	nodeMap := make(map[int]*pkg.TreeNode)

	// Add nodes from both trees to the map
	addToMap(tree1, nodeMap)
	addToMap(tree2, nodeMap)

	// Create a new root node for the merged tree
	mergedTree := &pkg.TreeNode{
		ID:       0,
		Label:    "Merged Tree",
		Children: []pkg.TreeNode{},
	}

	// Add all nodes from the map to the merged tree
	for _, node := range nodeMap {
		if node.ID == 0 {
			continue
		}
		mergedTree.Children = append(mergedTree.Children, *node)
	}

	return mergedTree
}

// addToMap adds nodes to the map, merging children if necessary
func addToMap(node *pkg.TreeNode, nodeMap map[int]*pkg.TreeNode) {
	if node == nil {
		return
	}

	if existingNode, found := nodeMap[node.ID]; found {
		existingNode.Children = mergeChildren(existingNode.Children, node.Children)
	} else {
		nodeMap[node.ID] = node
	}

	for _, child := range node.Children {
		addToMap(&child, nodeMap)
	}
}

// mergeChildren merges children of two nodes
func mergeChildren(children1, children2 []pkg.TreeNode) []pkg.TreeNode {
	childMap := make(map[int]*pkg.TreeNode)

	for _, child := range children1 {
		childMap[child.ID] = &child
	}

	for _, child := range children2 {
		if existingChild, found := childMap[child.ID]; found {
			existingChild.Children = mergeChildren(existingChild.Children, child.Children)
		} else {
			childMap[child.ID] = &child
		}
	}

	var mergedChildren []pkg.TreeNode
	for _, child := range childMap {
		mergedChildren = append(mergedChildren, *child)
	}

	return mergedChildren
}

// GetRelationshipFromName fetches and merges trees for two names and returns as JSON
func GetRelationshipFromName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	firstName := vars["firstname"]
	secondName := vars["secondname"]

	firstId := GetIdFromName(w, r, firstName)
	secondId := GetIdFromName(w, r, secondName)

	if firstId == 0 || secondId == 0 {
		http.Error(w, "One or both names not found", http.StatusNotFound)
		return
	}

	firstTreeData, err := ExtractTreeData(firstId)
	if err != nil {
		http.Error(w, "Error fetching subtree for first ID", http.StatusInternalServerError)
		log.Printf("Error fetching subtree for first ID: %v", err)
		return
	}

	secondTreeData, err := ExtractTreeData(secondId)
	if err != nil {
		http.Error(w, "Error fetching subtree for second ID", http.StatusInternalServerError)
		log.Printf("Error fetching subtree for second ID: %v", err)
		return
	}

	mergedTree := mergeTrees(firstTreeData, secondTreeData)

	jsonData, err := json.Marshal(mergedTree)
	if err != nil {
		http.Error(w, "Error marshalling JSON", http.StatusInternalServerError)
		log.Printf("Error marshalling JSON: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
