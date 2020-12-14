package parser

type Parser struct {
	Link            string              `json:"link"`
	ChildURLIsFull  bool                `json:"url_is_full"`
	Attributes      []Attribute         `json:"attributes"`
	ChildAttributes ChildPagesAttribute `json:"child_attributes"`
}

type Attribute struct {
	DivClass string `json:"div_class"`
	AClass   string `json:"a_class"`
}

type ChildPagesAttribute struct {
	ChildDivClass    string `json:"child_div_class"`
	ClassTitle       string `json:"class_title"`
	ClassDescription string `json:"class_description"`
}
