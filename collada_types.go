package collada

import "encoding/xml"

type Collada struct{
	XMLName xml.Name `xml:"COLLADA"`
	Geometry
}
