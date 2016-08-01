package collada

import "encoding/xml"

type collada struct {
	XMLName             xml.Name             `xml:"COLLADA"`
	LibraryGeometries   *libraryGeometries   `xml:"library_geometries"`
	LibraryControllers  *libraryControllers  `xml:"library_controllers"`
	LibraryVisualScenes *libraryVisualScenes `xml:"library_visual_scenes"`
}

//------library_geometries-----------

type libraryGeometries struct {
	Geometries []geometry `xml:"geometry"`
}

type geometry struct {
	Meshes []mesh `xml:"mesh"`
}

type mesh struct {
	Sources  []source `xml:"source"`
	Polylist polylist `xml:"polylist"`
}

type source struct {
	Id              string          `xml:"id,attr"`
	NameArray       *nameArray      `xml:"Name_array,omitempty"`
	FloatArray      *floatArray     `xml:"float_array,omitempty"`
	TechniqueCommon techniqueCommon `xml:"technique_common"`
}

type polylist struct {
	Matterial string  `xml:"material,attr"`
	Count     string  `xml:"count,attr"`
	Inputs    []input `xml:"input"`
	P         string  `xml:"p"`
}

type techniqueCommon struct {
	Accessor accessor `xml:"accessor"`
}

type accessor struct {
	Source string  `xml:"source,attr"`
	Count  string  `xml:"count,attr"`
	Stride string  `xml:"stride,attr"`
	Params []param `xml:"param,omitempty"`
}

type param struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
}

type input struct {
	Semantic string `xml:"semantic,attr"`
	Source   string `xml:"source,attr"`
	Offset   string `xml:"offset,attr"`
}

type nameArray struct {
	Id      string `xml:"id,attr"`
	Count   string `xml:"count,attr"`
	Content string `xml:",chardata"`
}

type floatArray struct {
	Id      string `xml:"id,attr"`
	Count   string `xml:"count,attr"`
	Content string `xml:",chardata"`
}

//------library_visual_scenes-----------

type libraryVisualScenes struct {
	VisualScene visualScene `xml:"visual_scene"`
}

type visualScene struct {
	Nodes []node `xml:"node"`
}

type node struct {
	Id     string `xml:"id,attr"`
	Sid    string `xml:"sid,attr"`
	Name   string `xml:"name,attr"`
	Type   string `xml:"type,attr"`
	Matrix matrix `xml:"matrix"`
	Nodes  []node `xml:"node"`
}

type matrix struct {
	Sid     string `xml:"sid,attr"`
	Content string `xml:",chardata"`
}

//--------library_controllers-------------

type libraryControllers struct {
	Controllers []controller `xml:"controller"`
}

type controller struct {
	Id   string `xml:"id,attr"`
	Name string `xml:"name,attr"`
	Skin skin   `xml:"skin"`
}

type skin struct {
	Source        string        `xml:"source,attr"`
	Sources       []source      `xml:"source"`
	Joints        joints        `xml:"joints"`
	VertexWeights vertexWeights `xml:"vertex_weights"`
}

type joints struct {
	Inputs []input `xml:"input"`
}

type vertexWeights struct {
	Inputs []input `xml:"input"`
	VCount string  `xml:"vcount"`
	V      string  `xml:"v"`
}
