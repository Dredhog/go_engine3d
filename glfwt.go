package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"training/glfwt/mesh"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var vertices = []float32{
	//Front face
	1.0, 1.0, 1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, 1.0, 0.0, 0.0, 1.0,
	-1.0, -1.0, 1.0, 0.0, 0.0, 0.0,
	1.0, -1.0, 1.0, 1.0, 0.0, 0.0,

	//Back face
	1.0, 1.0, -1.0, 1.0, 1.0, 1.0,
	-1.0, 1.0, -1.0, 0.0, 1.0, 1.0,
	-1.0, -1.0, -1.0, 0.0, 1.0, 0.0,
	1.0, -1.0, -1.0, 1.0, 1.0, 0.0,
}

var elements = []uint32{
	0, 1, 2, 3,
	7, 6, 5, 4,
	0, 4, 5, 1,
	1, 5, 6, 2,
	4, 0, 3, 7,
	2, 6, 7, 3,
}

const (
	screenWidth  = 1920
	screenHeight = 1080
)

// OpenglWork is responsible for everything drawn in the window context
func openglWork() {

	//Generate the mesh
	vertices, elements, _ = mesh.GeneratePlane(50, 50)

	//Set up glfw
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	//Set up the display window
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(screenWidth, screenHeight, "Opengl", glfw.GetPrimaryMonitor(), nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	//Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	shaderProgram, err := newProgram(vertexSource, fragmentSource)
	if err != nil {
		panic(err)
	}

	gl.UseProgram(shaderProgram)

	//Create the vertex vao
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	//Create the buffer containing vertex information
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	//Create the indices element buffer object
	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(elements)*4, gl.Ptr(elements), gl.STATIC_DRAW)

	//Link the vertex attributes to the variables
	//And vertex data in the vbo
	posAttrib := uint32(gl.GetAttribLocation(shaderProgram, gl.Str("position\x00")))
	gl.VertexAttribPointer(posAttrib, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(posAttrib)
	gl.ClearColor(0.8, 0.8, 1.0, 1.0)

	colAttrib := uint32(gl.GetAttribLocation(shaderProgram, gl.Str("color\x00")))
	gl.VertexAttribPointer(colAttrib, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(colAttrib)

	//Uniform variables
	projectionUniform := gl.GetUniformLocation(shaderProgram, gl.Str("projection\x00"))
	projection := mgl32.Perspective(math.Pi/4, 1.6, 0.0, 10.0)
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	cameraUniform := gl.GetUniformLocation(shaderProgram, gl.Str("camera\x00"))
	modelUniform := gl.GetUniformLocation(shaderProgram, gl.Str("model\x00"))
	timeUniform := gl.GetUniformLocation(shaderProgram, gl.Str("time\x00"))

	//Uniform variables for waves
	amplitudeUniform := gl.GetUniformLocation(shaderProgram, gl.Str("_amplitude\x00"))
	weightXUniform := gl.GetUniformLocation(shaderProgram, gl.Str("_weightX\x00"))
	weightYUniform := gl.GetUniformLocation(shaderProgram, gl.Str("_weightY\x00"))
	weightZUniform := gl.GetUniformLocation(shaderProgram, gl.Str("_weightZ\x00"))
	periodUniform := gl.GetUniformLocation(shaderProgram, gl.Str("_period\x00"))

	//Backing varaibles for wave uniforms
	amplitude := float32(0.1)
	weightX := float32(0)
	weightY := float32(0)
	weightZ := float32(0)
	period := float32(1000000000)

	//Game loop function
	func(window *glfw.Window, shaderProgram uint32, modelUniform int32) {
		/*
			gl.Enable(gl.CULL_FACE)
			gl.CullFace(gl.BACK)
			gl.Enable(gl.DEPTH_TEST)
			gl.DepthFunc(gl.LESS)
		*/
		model := mgl32.Ident4()
		playerPos := mgl32.Vec3{0, 5, 0}
		t0 := time.Now()
		startTime := t0
		frameTime := 8 * time.Millisecond
		frames := 0
		seconds := 0
		for !window.ShouldClose() {
			frames++
			gl.Clear(gl.COLOR_BUFFER_BIT)

			//Input function
			func(window *glfw.Window, pos *mgl32.Vec3) {
				var navigationSpeed float32 = 5.0 / 120
				pressedShift := false

				//Pressing space to exit
				if window.GetKey(glfw.KeySpace) == glfw.Press {
					window.SetShouldClose(true)
				}
				//Pressing enter to exit
				if window.GetKey(glfw.KeyEnter) == glfw.Press {
					window.SetShouldClose(true)
				}
				//First person motion
				if window.GetKey(glfw.KeyW) == glfw.Press {
					pos[2] -= navigationSpeed
				} else if window.GetKey(glfw.KeyS) == glfw.Press {
					pos[2] += navigationSpeed
				}
				if window.GetKey(glfw.KeyA) == glfw.Press {
					pos[0] -= navigationSpeed
				} else if window.GetKey(glfw.KeyD) == glfw.Press {
					pos[0] += navigationSpeed
				}
				//Wave variable editing
				if window.GetKey(glfw.KeyLeftShift) == glfw.Press {
					pressedShift = true
				}
				if window.GetKey(glfw.KeyX) == glfw.Press {
					if pressedShift {
						weightX -= 0.01
					} else {
						weightX += 0.01
					}
				}
				if window.GetKey(glfw.KeyY) == glfw.Press {
					if pressedShift {
						weightY -= 0.01
					} else {
						weightY += 0.01
					}
				}
				if window.GetKey(glfw.KeyZ) == glfw.Press {
					if pressedShift {
						weightZ -= 0.01
					} else {
						weightZ += 0.01
					}
				}
				if window.GetKey(glfw.KeyP) == glfw.Press {
					if pressedShift {
						period *= 1.01
					} else {
						period /= 1.01
					}
				}
				if window.GetKey(glfw.KeyM) == glfw.Press {
					if pressedShift {
						amplitude -= 0.01
					} else {
						amplitude += 0.01
					}
				}
			}(window, &playerPos)

			camera := mgl32.LookAtV(playerPos, playerPos.Add(mgl32.Vec3{0, -1, -3}), mgl32.Vec3{0, 1, 0})

			//Rotate the cube
			/*
			totalTime := float32(time.Since(t0)) / float32(time.Second)
			model = mgl32.HomogRotate3D(totalTime*math.Pi/4, mgl32.Vec3{1, 1, -1}.Normalize())
			*/

			//Update uniform variables
			gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])
			gl.Uniform1f(timeUniform, float32(time.Since(t0)))
			gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
			//Update wave uniforms
			gl.Uniform1f(amplitudeUniform, amplitude)
			gl.Uniform1f(weightXUniform, weightX)
			gl.Uniform1f(weightYUniform, weightY)
			gl.Uniform1f(weightZUniform, weightZ)
			gl.Uniform1f(periodUniform, period)

			gl.DrawElements(gl.TRIANGLES, int32(len(elements)), gl.UNSIGNED_INT, gl.PtrOffset(0))

			time.Sleep(frameTime - time.Since(startTime))
			if int(time.Since(t0)/time.Second) > seconds {
				seconds++
				fmt.Println("FPS:", frames)
				frames = 0
			}
			startTime = time.Now()

			window.SwapBuffers()
			glfw.PollEvents()
		}
	}(window, shaderProgram, modelUniform)
}
