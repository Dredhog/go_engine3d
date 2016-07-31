package collada

import "encoding/xml"

type collada struct {
	XMLName             xml.Name            `xml:"COLLADA"`
	LibraryGeometries   libraryGeometries   `xml:"library_geometries"`
	LibraryVisualScenes libraryVisualScenes `xml:"library_visual_scenes"`
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
	FloatArray      floatArray      `xml:"float_array"`
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
	Source string `xml:"source,attr"`
	Count  string `xml:"count,attr"`
	Stride string `xml:"stride,attr"`
}

type input struct {
	Semantic string `xml:"semantic,attr"`
	Source   string `xml:"source,attr"`
	Offset   string `xml:"offset,attr"`
}

type floatArray struct {
	Id     string `xml:"id,attr"`
	Count  string `xml:"count,attr"`
	Floats string `xml:",chardata"`
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
	Name   string `xml:"name,attr"`
	Type   string `xml:"type,attr"`
	Matrix matrix `xml:"matrix"`
	Nodes  []node `xml:"node"`
}

type matrix struct {
	Sid     string `xml:"sid,attr"`
	Content string `xml:",chardata"`
}
