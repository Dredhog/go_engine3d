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
	Sources   []source   `xml:"source"`
	Polylists []polylist `xml:"polylist"`
}

type source struct {
	Id         string     `xml:"id,attr"`
	FloatArray floatArray `xml:"float_array"`
}

type polylist struct {
	Matterial string  `xml:"material,attr"`
	count     string  `xml:"count,attr"`
	Inputs    []input `xml:"input"`
	P         string  `xml:"p"`
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
