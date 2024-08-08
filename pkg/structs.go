package pkg

type TreeResult struct {
	ID       int    `json:"id"`
	ParentID int    `json:"parentid"`
	Name     string `json:"name"`
	Img      string `json:"img"`
}

type TreeNode struct {
	ID       int        `json:"id"`
	Label    string     `json:"label"`
	Children []TreeNode `json:"children,omitempty"`
}
