package handler

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/starboy011/Family-Tree-Backend/internal/db"
)

// ConvertToTree converts the flat list of nodes into a hierarchical structure
func ConvertToTreefromID(nodes []TreeResult) TreeNode {
	nodeMap := make(map[int][]TreeResult)
	var root TreeResult

	// Create a map of parentID to child nodes and identify the root node
	for _, node := range nodes {
		if node.ParentID == 0 {
			root = node
		} else {
			nodeMap[node.ParentID] = append(nodeMap[node.ParentID], node)
		}
	}

	// Recursively build the tree
	return buildTreeFromId(root, nodeMap)
}

// buildTree recursively constructs the tree from the node map
func buildTreeFromId(current TreeResult, nodeMap map[int][]TreeResult) TreeNode {
	treeNode := TreeNode{
		ID:    current.ID,
		Label: current.Name,
		Img:   current.Img,
	}
	for _, child := range nodeMap[current.ID] {
		treeNode.Children = append(treeNode.Children, buildTreeFromId(child, nodeMap))
	}
	return treeNode
}

// ExtractSubtreeWithPath finds the subtree rooted at the specified node ID and ensures the path from the root to the specified node is included
func ExtractSubtreeWithPath(root TreeNode, id int) *TreeNode {
	if root.ID == id {
		return &root
	}
	for _, child := range root.Children {
		subtree := ExtractSubtreeWithPath(child, id)
		if subtree != nil {
			// Return the current root with the subtree attached to it
			return &TreeNode{
				ID:       root.ID,
				Label:    root.Label,
				Img:      root.Img,
				Children: []TreeNode{*subtree},
			}
		}
	}
	return nil
}

func GetTreeDataFromId(w http.ResponseWriter, r *http.Request) {
	imageDir := "images"
	defaultImage := "Default.jpg"
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	db, err := db.InitDb(w, r)
	if err != nil {
		http.Error(w, "Error in connecting to db", http.StatusInternalServerError)
		log.Fatalf("Error connecting to db: %v", err)
		return
	}
	defer db.Close()

	query := `SELECT "ID", "Name", "ParentID" FROM mulvansham WHERE "Relationship" IN (1, 2)`

	// Execute query
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Error executing query", http.StatusInternalServerError)
		log.Fatalf("Error executing query: %v", err)
		return
	}
	defer rows.Close()

	var results []TreeResult

	// Iterate through rows and populate results
	for rows.Next() {
		var result TreeResult
		if err := rows.Scan(&result.ID, &result.Name, &result.ParentID); err != nil {
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			log.Fatalf("Error scanning row: %v", err)
			return
		}
		idStr := strconv.Itoa(result.ID)
		imageName := idStr + ".jpg"
		imagePath := path.Join(imageDir, imageName)

		// Check if the image file exists
		if _, err := os.Stat(imagePath); os.IsNotExist(err) {
			// Use default image if the image file does not exist
			imagePath = path.Join(imageDir, defaultImage)
		}

		imageBytes, err := os.ReadFile(imagePath)
		if err != nil {
			http.Error(w, "Error reading image file", http.StatusInternalServerError)
			log.Fatalf("Error reading image file: %v", err)
			return
		}

		// Encode image bytes to base64
		imageBase64 := base64.StdEncoding.EncodeToString(imageBytes)
		result.Img = "data:image/jpg;base64," + imageBase64 // Adjust according to your image type

		results = append(results, result)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "Error iterating through rows", http.StatusInternalServerError)
		log.Fatalf("Error iterating through rows: %v", err)
		return
	}

	// Convert results slice to hierarchical tree structure
	tree := ConvertToTreefromID(results)

	// Extract the subtree rooted at the specified node, including the path from root to the node
	subtree := ExtractSubtreeWithPath(tree, id)
	if subtree == nil {
		http.Error(w, "Node not found", http.StatusNotFound)
		return
	}

	// Convert subtree to JSON
	jsonData, err := json.Marshal(subtree)
	if err != nil {
		http.Error(w, "Error marshalling JSON", http.StatusInternalServerError)
		log.Fatalf("Error marshalling JSON: %v", err)
		return
	}

	// Set Content-Type header and write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
