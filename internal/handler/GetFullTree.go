package handler

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/starboy011/Family-Tree-Backend/internal/db"
)

type TreeResult struct {
	ID       int    `json:"id"`
	ParentID int    `json:"parentid"`
	Name     string `json:"name"`
	Img      string `json:"img"`
}

type TreeNode struct {
	ID       int        `json:"id"`
	Label    string     `json:"label"`
	Img      string     `json:"img"`
	Children []TreeNode `json:"children,omitempty"`
}

// ConvertToTree converts the flat list of nodes into a hierarchical structure
func ConvertToTree(nodes []TreeResult) TreeNode {
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
	return buildTree(root, nodeMap)
}

// buildTree recursively constructs the tree from the node map
func buildTree(current TreeResult, nodeMap map[int][]TreeResult) TreeNode {
	treeNode := TreeNode{
		ID:    current.ID,
		Label: current.Name,
		Img:   current.Img,
	}
	for _, child := range nodeMap[current.ID] {
		treeNode.Children = append(treeNode.Children, buildTree(child, nodeMap))
	}
	return treeNode
}

func GetFamilyTreeData(w http.ResponseWriter, r *http.Request) {
	imageDir := "images"
	defaultImage := "Default.jpg"

	db, err := db.InitDb(w, r)
	if err != nil {
		http.Error(w, "Error in connecting to db", http.StatusInternalServerError)
		log.Fatalf("Error connecting to db: %v", err)
		return
	}
	defer db.Close()

	query := `SELECT "ID", "Name", "ParentID" FROM mulvansham`

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
	tree := ConvertToTree(results)

	// Convert tree to JSON
	jsonData, err := json.Marshal(tree)
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
