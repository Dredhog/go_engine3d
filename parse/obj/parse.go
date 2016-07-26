package obj

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func deletEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func ParseFile(fileName string, useT, useN bool) (floats []float32, indices []uint32) {
	relativePath, err := os.Getwd()
	check(err)

	file, err := os.Open(relativePath + "/data/models/" + fileName)
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	//the vertex variables
	var positions []float32
	var normals []float32
	var uvs []float32
	var objIndices []uint32

	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Split(line, " ")
		words = deletEmpty(words)
		if len(words) == 0 {
			continue
		}

		switch words[0] {
		case "v":
			x, _ := strconv.ParseFloat(words[1], 32)
			y, _ := strconv.ParseFloat(words[2], 32)
			z, _ := strconv.ParseFloat(words[3], 32)
			positions = append(positions, float32(x), float32(y), float32(z))
		case "vn":
			x, _ := strconv.ParseFloat(words[1], 32)
			y, _ := strconv.ParseFloat(words[2], 32)
			z, _ := strconv.ParseFloat(words[3], 32)
			normals = append(normals, float32(x), float32(y), float32(z))
		case "vt":
			u, _ := strconv.ParseFloat(words[1], 32)
			v, _ := strconv.ParseFloat(words[2], 32)
			uvs = append(uvs, float32(u), float32(v))
		case "f":
			vert1 := strings.Split(words[1], "/")
			vert2 := strings.Split(words[2], "/")
			vert3 := strings.Split(words[3], "/")
			//first vertex
			v0, _ := strconv.ParseUint(vert1[0], 10, 32)
			vt0, _ := strconv.ParseUint(vert1[1], 10, 32)
			vn0, _ := strconv.ParseUint(vert1[2], 10, 32)

			//second vertex
			v1, _ := strconv.ParseUint(vert2[0], 10, 32)
			vt1, _ := strconv.ParseUint(vert2[1], 10, 32)
			vn1, _ := strconv.ParseUint(vert2[2], 10, 32)

			//thirds vertex
			v2, _ := strconv.ParseUint(vert3[0], 10, 32)
			vt2, _ := strconv.ParseUint(vert3[1], 10, 32)
			vn2, _ := strconv.ParseUint(vert3[2], 10, 32)

			switch {
			case !useT && !useN:
				objIndices = append(objIndices, uint32(v0)-1, uint32(v1)-1, uint32(v2)-1)
			case !useT:
				objIndices = append(objIndices, uint32(v0-1), uint32(vn0-1), uint32(v1-1), uint32(vn1-1), uint32(v2-1), uint32(vn2-1))
			case !useN:
				objIndices = append(objIndices, uint32(v0-1), uint32(vt0-1), uint32(v1-1), uint32(vt1-1), uint32(v2-1), uint32(vt2-1))
			default:
				objIndices = append(objIndices, uint32(v0-1), uint32(vt0-1), uint32(vn0-1), uint32(v1-1), uint32(vt1-1), uint32(vn1-1), uint32(v2-1), uint32(vt2-1), uint32(vn2-1))
			}
		}
	}

	floatsPerVertex := 3
	attributesPerIndex := 1
	if useT {
		floatsPerVertex += 2
		attributesPerIndex++
	}
	if useN {
		floatsPerVertex += 3
		attributesPerIndex++
	}

	floats = make([]float32, len(objIndices)/attributesPerIndex*floatsPerVertex)
	indices = make([]uint32, len(objIndices)/attributesPerIndex)

	for i := 0; i < len(objIndices)/attributesPerIndex; i++ {
		//Add the position
		floats[i*floatsPerVertex] = positions[objIndices[i*attributesPerIndex]*3]
		floats[i*floatsPerVertex+1] = positions[objIndices[i*attributesPerIndex]*3+1]
		floats[i*floatsPerVertex+2] = positions[objIndices[i*attributesPerIndex]*3+2]

		//Add uvs
		if useT {
			floats[i*floatsPerVertex+3] = uvs[objIndices[i*attributesPerIndex+1]*2]
			floats[i*floatsPerVertex+4] = uvs[objIndices[i*attributesPerIndex+1]*2+1]
		}
		//add normals
		if offsetF, offsetI := 0, 0; useN {
			if useT {
				offsetF = 2
				offsetI = 1
			}
			floats[i*floatsPerVertex+3+offsetF] = normals[objIndices[i*attributesPerIndex+offsetI+1]*3]
			floats[i*floatsPerVertex+4+offsetF] = normals[objIndices[i*attributesPerIndex+offsetI+1]*3+1]
			floats[i*floatsPerVertex+5+offsetF] = normals[objIndices[i*attributesPerIndex+offsetI+1]*3+2]
		}
	}
	for i := 0; i < len(indices); i++ {
		indices[i] = uint32(i)
	}
	return
}
