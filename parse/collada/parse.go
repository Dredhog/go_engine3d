package collada

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"training/engine/anim"
	"training/engine/types"

	"github.com/go-gl/mathgl/mgl32"
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

func splitAndRemoveEmpty(str string) []string {
	split := strings.Split(str, " ")
	var result []string
	for _, s := range split {
		if s != "" {
			result = append(result, s)
		}
	}
	return result
}

func parseSkin(skin *skin) ([]float32, []float32, error) {
	var boneDataIndices, influenceCounts []int
	var boneWeights []float32
	indexStride, bonesPerVert := 2, 2

	boneWeights, err := stringToFloatArray(skin.Sources[2].FloatArray.Content)
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

func ParseModel(fileName string) (*types.Model, error) {
	collada, err := Parse(fileName)
	if err != nil {
		return nil, err
	}
	if collada.LibraryGeometries == nil {
		return nil, fmt.Errorf("collada to mesh: no geometry data found in %v\n", fileName)
	}
	model := types.Model{}

	meshCollada := &collada.LibraryGeometries.Geometries[0].Meshes[0]
	controllerCollada := &collada.LibraryControllers.Controllers[0]
	model.Mesh, err = extractMesh(meshCollada, controllerCollada)
	if err != nil {
		return nil, fmt.Errorf("collada mode: error extracting mesh : %v", err)
	}
	model.Skeleton, err = extractSkeleton(&controllerCollada.Skin, collada.LibraryVisualScenes)
	if err != nil {
		return nil, fmt.Errorf("collada mode: error extracting skeleton: %v", err)
	}

	return &model, nil
}

func extractMesh(meshCollada *mesh, controllerCollada *controller) (*types.Mesh, error) {
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
		attribs, _ := stringToFloatArray(meshCollada.Sources[a].FloatArray.Content)
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
	if controllerCollada != nil {
		boneIndexAttribs, boneWeightAttribs, err = parseSkin(&controllerCollada.Skin)
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
		if err != nil {
			panic(err)
		}
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

func extractSkeleton(skin *skin, libraryVisualScenes *libraryVisualScenes) (*anim.Skeleton, error) {
	sidToInvBindMat := make(map[string][]float32)
	sidToIndex := make(map[string]int)
	boneNames := splitAndRemoveEmpty(skin.Sources[0].NameArray.Content)
	invBindMatriceFloats, err := stringToFloatArray(skin.Sources[1].FloatArray.Content)
	if err != nil {
		return nil, fmt.Errorf("collada: skeleton: getting inv bind matrices %v", err)
	}
	bindShapeFloats, err := stringToFloatArray(skin.BindShapeMatrix)
	if err != nil {
		return nil, fmt.Errorf("collada: skeleton: getting bind shape matrix %v", err)
	}
	bindShapeArray := [16]float32{}
	copy(bindShapeArray[:], bindShapeFloats)
	bindShapeMatrix := mgl32.Mat4(bindShapeArray)

	for i, name := range boneNames {
		sidToInvBindMat[name] = invBindMatriceFloats[i*16 : (i+1)*16]
		sidToIndex[name] = i
	}

	var armature node
	for _, node := range libraryVisualScenes.VisualScene.Nodes {
		if node.Name == "Armature" {
			armature = node
		}
	}
	for _, node := range armature.Nodes {
		if node.Type == "JOINT" {
			bones := make([]anim.Bone, len(boneNames))
			registerBoneAndItsChildren(bones, &node, sidToIndex[node.Name], bindShapeMatrix, sidToIndex, sidToInvBindMat)

			skeleton := anim.NewSkeleton(bones, bindShapeMatrix, sidToIndex[node.Sid])
			skeleton.RootIndex = sidToIndex[node.Name]
			return skeleton, nil
		}
	}
	return nil, fmt.Errorf("collada: no root node found in skeleton data")
}

func registerBoneAndItsChildren(bones []anim.Bone, currentNode *node, parentIndex int, bindShape mgl32.Mat4, sidToIndex map[string]int, sidToInvBindMat map[string][]float32) {
	currentIndex := sidToIndex[currentNode.Sid]
	bones[currentIndex].Name = currentNode.Name
	bones[currentIndex].Index = sidToIndex[currentNode.Sid]
	bones[currentIndex].ParentIndex = parentIndex

	copy(bones[currentIndex].InverseBindPose[:], sidToInvBindMat[currentNode.Sid])
	bones[currentIndex].InverseBindPose = bones[currentIndex].InverseBindPose.Transpose().Mul4(bindShape.Transpose())
	bones[currentIndex].BindPose = bones[currentIndex].InverseBindPose.Inv()
	for _, node := range currentNode.Nodes {
		registerBoneAndItsChildren(bones, &node, sidToIndex[currentNode.Sid], bindShape, sidToIndex, sidToInvBindMat)
	}
}
