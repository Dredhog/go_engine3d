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

func stringToFloatArray(str string) ([]float32, error) {
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

func extractSkinData(libraryControllers *libraryControllers) ([]float32, []float32, error) {
	skin := libraryControllers.Controllers[0].Skin
	var boneDataIndices, influenceCounts []int
	var boneWeights []float32
	indexStride, bonesPerVert := 2, 2

	boneWeights, err := stringToFloatArray(skin.Sources[2].FloatArray.Floats)
	if err != nil {
		return nil, nil, fmt.Errorf("collada: skin error getting bone weights %v", err)
	}
	influenceCounts, err = stringToIntArray(skin.VertexWeights.VCount)
	if err != nil {
		return nil, nil, fmt.Errorf("collada: skin: %v", err)
	}
	boneDataIndices, err = stringToIntArray(skin.VertexWeights.V)
	if err != nil {
		return nil, nil, fmt.Errorf("collada: skin: %v", err)
	}

	indices := make([]float32, bonesPerVert*len(influenceCounts))
	weights := make([]float32, bonesPerVert*len(influenceCounts))
	currentIndex := 0
	for i := 0; i < len(influenceCounts); i++ {
		switch influenceCounts[i] {
		case 0:
			indices[bonesPerVert*i] = 0
			indices[bonesPerVert*i+1] = 0
			weights[bonesPerVert*i] = 0
			weights[bonesPerVert*i+1] = 0
		case 1:
			indices[bonesPerVert*i] = float32(boneDataIndices[currentIndex])
			weights[bonesPerVert*i] = boneWeights[boneDataIndices[currentIndex+1]]
			indices[bonesPerVert*i+1] = 0
			weights[bonesPerVert*i+1] = 0
		case 2:
			indices[bonesPerVert*i] = float32(boneDataIndices[currentIndex])
			weights[bonesPerVert*i] = boneWeights[boneDataIndices[currentIndex+1]]
			indices[bonesPerVert*i+1] = float32(boneDataIndices[currentIndex+2])
			weights[bonesPerVert*i+1] = boneWeights[boneDataIndices[currentIndex+3]]
		default:
			indices[bonesPerVert*i] = float32(boneDataIndices[currentIndex])
			weights[bonesPerVert*i] = boneWeights[boneDataIndices[currentIndex+1]]
			indices[bonesPerVert*i+1] = float32(boneDataIndices[currentIndex+2])
			weights[bonesPerVert*i+1] = boneWeights[boneDataIndices[currentIndex+3]]
		}
		currentIndex += influenceCounts[i] * indexStride
	}
	return indices, weights, nil
}

func ParseToMesh(fileName string) (*types.Mesh, error) {
	collada, err := Parse(fileName)
	if err != nil {
		return nil, err
	}

	if collada.LibraryGeometries == nil {
		return nil, fmt.Errorf("collada to mesh: no geometry data found in %v\n", fileName)
	}
	meshCollada := collada.LibraryGeometries.Geometries[0].Meshes[0]
	indices, err := stringToIntArray(meshCollada.Polylist.P)
	if err != nil {
		return nil, err
	}

	floatsPerVert := 0
	for _, source := range meshCollada.Sources {
		tempInt, _ := strconv.Atoi(source.TechniqueCommon.Accessor.Stride)
		floatsPerVert += tempInt
	}

	indexStride := len(meshCollada.Polylist.Inputs)
	attrMask := uint32(0)
	attrOffsets := [6]int{}
	floats := make([]float32, floatsPerVert*len(indices)/indexStride)

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
		attribs, _ := stringToFloatArray(meshCollada.Sources[a].FloatArray.Floats)
		floatStride, _ := strconv.Atoi(meshCollada.Sources[a].TechniqueCommon.Accessor.Stride)
		indexOffset, _ := strconv.Atoi(meshCollada.Polylist.Inputs[a].Offset)
		for b := indexOffset; b < len(indices); b += indexStride {
			for c := 0; c < floatStride; c++ {
				floats[floatCount] = attribs[floatStride*indices[b]+c]
				floatCount++
			}
		}
	}

	//Bones
	var boneIndexAttribs, boneWeightAttribs []float32
	if collada.LibraryControllers != nil {
		boneIndexAttribs, boneWeightAttribs, err = extractSkinData(collada.LibraryControllers)
		if err != nil {
			return nil, fmt.Errorf("collada to mesh: skin data extraction error: %v", err)
		}
		attrMask += types.USE_BONES
		bonesPerVert := 2
		weightAttribOffset := bonesPerVert * len(indices) / indexStride
		boneFloats := make([]float32, 2*bonesPerVert*len(indices)/indexStride)
		for i := 0; i < len(indices)/indexStride; i++ {
			for j := 0; j < bonesPerVert; j++ {
				boneFloats[bonesPerVert*i+j] = boneIndexAttribs[bonesPerVert*indices[indexStride*i]+j]
				boneFloats[weightAttribOffset+bonesPerVert*i+j] = boneWeightAttribs[bonesPerVert*indices[indexStride*i]+j]
			}
		}
		attrOffsets[4] = len(floats)
		attrOffsets[5] = attrOffsets[4] + weightAttribOffset
		floats = append(floats, boneFloats...)
	}
	//Indices ("0, 1, 2, ..., n")
	finalIndices := make([]uint32, len(indices)/indexStride)
	for i := range finalIndices {
		finalIndices[i] = uint32(i)
	}
	mesh := types.Mesh{}
	mesh.Init(floats, finalIndices, attrMask, attrOffsets, nil)

	return &mesh, nil
}
