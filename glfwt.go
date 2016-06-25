package main

import (
	"fmt"
	"log"
	"math"
	"time"

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
	screenWidth  = 640
	screenHeight = 480
)

// OpenglWork is responsible for everything drawn in the window context
func openglWork() {

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
	window, err := glfw.CreateWindow(screenWidth, screenHeight, "Opengl", nil, nil)
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
	projection := mgl32.Perspective(math.Pi/4, 1.2, 0.0, 10.0)
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	cameraUniform := gl.GetUniformLocation(shaderProgram, gl.Str("camera\x00"))
	camera := mgl32.LookAtV(mgl32.Vec3{-3, 3, -3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	modelUniform := gl.GetUniformLocation(shaderProgram, gl.Str("model\x00"))

	//Game loop function
	func(window *glfw.Window, shaderProgram uint32, modelUniform int32) {
		gl.Enable(gl.CULL_FACE)
		gl.CullFace(gl.BACK)
		model := mgl32.Ident4()
		playerPos := mgl32.Vec3{0, -2, 5}
		t0 := time.Now()
		startTime := t0
		frameTime := 16 * time.Millisecond
		frames := 0
		seconds := 0
		for !window.ShouldClose() {
			frames++

			//Input function
			func(window *glfw.Window, pos *mgl32.Vec3) {
				var navigationSpeed float32 = 5.0 / 60

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
			}(window, &playerPos)

			//Alter the spectator position
			camera := mgl32.LookAtV(playerPos, playerPos.Add(mgl32.Vec3{0, 1, -3}), mgl32.Vec3{0, 1, 0})
			gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

			//Rotate the cube
			totalTime := float32(time.Since(t0)) / float32(time.Second)
			model = mgl32.HomogRotate3D(totalTime*math.Pi/4, mgl32.Vec3{1, 1, 1}.Normalize())

			gl.Clear(gl.COLOR_BUFFER_BIT)
			gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
			gl.DrawElements(gl.TRIANGLE_FAN, 4, gl.UNSIGNED_INT, gl.PtrOffset(0))
			gl.DrawElements(gl.TRIANGLE_FAN, 4, gl.UNSIGNED_INT, gl.PtrOffset(4*4))
			gl.DrawElements(gl.TRIANGLE_FAN, 4, gl.UNSIGNED_INT, gl.PtrOffset(8*4))
			gl.DrawElements(gl.TRIANGLE_FAN, 4, gl.UNSIGNED_INT, gl.PtrOffset(12*4))
			gl.DrawElements(gl.TRIANGLE_FAN, 4, gl.UNSIGNED_INT, gl.PtrOffset(16*4))
			gl.DrawElements(gl.TRIANGLE_FAN, 4, gl.UNSIGNED_INT, gl.PtrOffset(20*4))
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
