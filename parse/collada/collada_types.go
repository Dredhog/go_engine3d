package collada

import "encoding/xml"

type collada struct {
	XMLName           xml.Name          `xml:"COLLADA"`
	LibraryGeometries libraryGeometries `xml:"library_geometries"`
}

type libraryGeometries struct {
	Geometries []geometry `xml:"geometry"`
}

type geometry struct {
	Meshes []mesh `xml:"mesh"`
}

type mesh struct {
	Sources []source `xml:"source"`
	Indices string   `xml:"polylist>p"`
}

type source struct {
	Id          string `xml:"id,attr"`
	FloatArrays string `xml:"float_array"`
}

type floatArray struct {
	Id     string `xml:"id,attr"`
	Count  string `xml:"count,attr"`
	Floats string `xml:",chardata"`
}
