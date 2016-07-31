package collada

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"training/engine/types"
)

func Parse(relativePath string) (collada, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return collada{}, fmt.Errorf("parse collada: absolute path error: %v\n", err)
	}
	colladaBytes, err := ioutil.ReadFile(workingDir + "/" + relativePath)
	if err != nil {
		return collada{}, fmt.Errorf("parse collada: %v\n", err)
	}
	var coll collada
	if err := xml.Unmarshal(colladaBytes, &coll); err != nil {
		return collada{}, fmt.Errorf("parse collada: %v\n", err)
	}
	return coll, nil
}

func stringToFLoatArray(str string) ([]float32, error) {
	split := strings.Split(str, " ")
	var result []float32
	for _, s := range split {
		if s != "" {
			f, err := strconv.ParseFloat(s, 32)
			if err != nil {
				return nil, fmt.Errorf("Collada: string to float conv error: %v, for number: %v\n", err, f)
			}
			result = append(result, float32(f))
		}
	}
	return result, nil
}

func stringToIntArray(str string) ([]int, error) {
	split := strings.Split(str, " ")
	var result []int
	for _, s := range split {
		if s != "" {
			i, err := strconv.Atoi(s)
			if err != nil {
				return nil, fmt.Errorf("collada: string to int conv error: %v, for number: %v\n", err, i)
			}
			result = append(result, i)
		}
	}
	return result, nil
}

func ParseToMesh(fileName string) (*types.Mesh, error) {
	collada, err := Parse(fileName)
	if err != nil {
		return nil, err
	}

	if len(collada.LibraryGeometries.Geometries) == 0 {
		return nil, fmt.Errorf("collada to mesh: no geometry data found in %v\n", fileName)
	}

	meshCollada := collada.LibraryGeometries.Geometries[0].Meshes[0]
	indicesCollada, err := stringToIntArray(meshCollada.Polylist.P)
	if err != nil {
		return nil, err
	}

	floatsPerVert := 0
	for _, source := range meshCollada.Sources {
		tempInt, _ := strconv.Atoi(source.TechniqueCommon.Accessor.Stride)
		floatsPerVert += tempInt
	}

	attrMask := uint32(0)
	attrOffsets := [5]int{}
	indexStride := len(meshCollada.Polylist.Inputs)
	floats := make([]float32, floatsPerVert*len(indicesCollada)/indexStride)

	for a, floatCount := 0, 0; a < indexStride; a++ {
		attribute := meshCollada.Polylist.Inputs[a]
		switch attribute.Semantic {
		case "VERTEX":
			attrMask += types.USE_POSITIONS
			attrOffsets[0] = floatCount
		case "NORMAL":
			attrMask += types.USE_NORMALS
			attrOffsets[1] = floatCount
		case "TEXCOORD":
			attrMask += types.USE_TEXCOORDS
			attrOffsets[2] = floatCount
		case "COLOR":
			attrMask += types.USE_COLORS
			attrOffsets[3] = floatCount
		}
		attribs, _ := stringToFLoatArray(meshCollada.Sources[a].FloatArray.Floats)
		floatStride, _ := strconv.Atoi(meshCollada.Sources[a].TechniqueCommon.Accessor.Stride)
		indexOffset, _ := strconv.Atoi(meshCollada.Polylist.Inputs[a].Offset)
		for b := indexOffset; b < len(indicesCollada); b += indexStride {
			for c := 0; c < floatStride; c++ {
				floats[floatCount] = attribs[floatStride*indicesCollada[b]+c]
				floatCount++
			}
		}
	}
	indices := make([]uint32, len(floats)/floatsPerVert)
	for i, _ := range indices {
		indices[i] = uint32(i)
	}
	mesh := types.Mesh{}
	mesh.Init(floats, indices, attrMask, attrOffsets, nil)

	return &mesh, nil
}
