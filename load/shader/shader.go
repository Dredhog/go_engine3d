package shader

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
)

func NewProgram(fileName string) (uint32, error) {
	relativePath, err := os.Getwd()
	if err != nil {
		return 0, fmt.Errorf("shader: realative path read error: %v", err)
	}

	//Read vertex shader
	vertexSourceBytes, err := ioutil.ReadFile(relativePath + "/shaders/" + fileName + ".vert")
	if err != nil {
		return 0, fmt.Errorf("shader: vertex file read error: %v", err)
	}
	vertexShader, err := compileShader(string(vertexSourceBytes)+"\x00", gl.VERTEX_SHADER)
	if err != nil {
		return 0, fmt.Errorf("shader: vertex compilation error: %v", err)
	}

	//Read fragment shader
	fragmentSourceBytes, err := ioutil.ReadFile(relativePath + "/shaders/" + fileName + ".frag")
	if err != nil {
		return 0, fmt.Errorf("shader: fragment file read error: %v", err)
	}
	fragmentShader, err := compileShader(string(fragmentSourceBytes)+"\x00", gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, fmt.Errorf("shader: fragment compilation error: %v", err)
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}
func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
