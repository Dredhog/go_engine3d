package main

import (
	"fmt"
	"log"
	"math"
	"time"
	"training/engine/anim"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	screenWidth  = 768
	screenHeight = 768
	fps          = 1200
)

var vertices = []float32{
	0.1, 0, 0, //position
	0, 0.1, 0, //rotation
	1, 0, 0, //color
	0, 0, //bones
	0.2, 0.8, //weights

	-0.1, 0, 0, //position
	0, 0, 0, //rotation
	1, 0, 0, //color
	0, 0, //bones
	0.2, 0.8, //weights

	0.1, 0.5, 0, //position
	0, 0, 0, //rotation
	0, 1, 0, //color
	1, 0, //bones
	0.2, 0.8, //weights

	-0.1, 0.5, 0, //position
	0, 0, 0, //rotation
	0, 1, 0, //color
	1, 0, //bones
	0.2, 0.8, //weights

	0.1, 0.8, 0, //position
	0, 0, 0, //rotation
	0, 0, 1, //color
	2, 2, //bones
	0.2, 0.8, //weights

	-0.1, 0.8, 0, //position
	0, 0, 0, //rotation
	0, 0, 1, //color
	2, 2, //bones
	0.2, 0.8, //weights
}
var elements []uint32 = []uint32{
	0, 1,
	3, 5,
	4, 2,
	0,
}

// OpenglWork is responsible for everything drawn in the window context
func runEngine() {

	//---------------------------------ANIMATION---------------------------------
	crane := anim.NewSkeleton(3)
	crane.Bones[0] = anim.Bone{Name: "root", ToRoot: anim.Transform{mgl32.Vec3{1, 1, 1}, [3]float32{0, 0, 0}, [3]float32{0, 0, 0}}}
	crane.Bones[1] = anim.Bone{Name: "neck", ToRoot: anim.Transform{mgl32.Vec3{1, 1, 1}, mgl32.Vec3{0, -0.5, 0}, mgl32.Vec3{}}, Index: 1, ParentIndex: 0}
	crane.Bones[2] = anim.Bone{Name: "head", ToRoot: anim.Transform{mgl32.Vec3{1, 1, 1}, mgl32.Vec3{0, -0.8, 0}, mgl32.Vec3{}}, Index: 2, ParentIndex: 1}
	//---------------------------------------------------------------------------

	//vertices, elements := obj_reader.ParseFile("dust2x2/dust2x2.obj", false, true)

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

	//Create the vertex array object
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
	gl.VertexAttribPointer(posAttrib, 3, gl.FLOAT, false, 13*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(posAttrib)
	gl.ClearColor(0.8, 0.8, 1.0, 1.0)

	normalAttrib := uint32(gl.GetAttribLocation(shaderProgram, gl.Str("normal\x00")))
	gl.VertexAttribPointer(normalAttrib, 3, gl.FLOAT, false, 13*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(normalAttrib)

	boneIndiceAttrib := uint32(gl.GetAttribLocation(shaderProgram, gl.Str("bones\x00")))
	gl.VertexAttribPointer(boneIndiceAttrib, 2, gl.FLOAT, false, 13*4, gl.PtrOffset(9*4))
	gl.EnableVertexAttribArray(boneIndiceAttrib)

	boneWeightAttrib := uint32(gl.GetAttribLocation(shaderProgram, gl.Str("weights\x00")))
	gl.VertexAttribPointer(boneWeightAttrib, 2, gl.FLOAT, false, 13*4, gl.PtrOffset(11*4))
	gl.EnableVertexAttribArray(boneWeightAttrib)

	//Game loop function
	func(window *glfw.Window, shaderProgram uint32) {
		//gl.Enable(gl.CULL_FACE)
		//gl.CullFace(gl.BACK)
		gl.Enable(gl.DEPTH_TEST)
		gl.DepthFunc(gl.LESS)

		//Declare uniform variables
		mvpUniform := gl.GetUniformLocation(shaderProgram, gl.Str("mvp_mat\x00"))
		modelUniform := gl.GetUniformLocation(shaderProgram, gl.Str("model_mat\x00"))
		lightPosUniform := gl.GetUniformLocation(shaderProgram, gl.Str("light_position\x00"))
		boneUniforms := gl.GetUniformLocation(shaderProgram, gl.Str("bone_mat\x00"))

		projection := mgl32.Perspective(math.Pi/4, 1.6, 0.1, 100.0)
		playerPos := mgl32.Vec3{0, 1, 2}
		lightPos := mgl32.Vec3{0, 10, 0}
		angle0 := float32(0)
		angle1 := float32(0)
		angle2 := float32(0)
		t0 := time.Now()
		startTime := t0
		frameTime := time.Second / fps
		frames := 0
		seconds := 0

		for !window.ShouldClose() {
			frames++

			handleInput(window, &playerPos, &lightPos, &angle0, &angle1, &angle2)

			camera := mgl32.LookAtV(playerPos, playerPos.Add(mgl32.Vec3{0, -0.5, -3}), mgl32.Vec3{0, 1, 0})

			//Update the main viewport matrix
			mvp := projection.Mul4(camera)

			//Rotate the cube
			totalTime := float32(time.Since(t0)) / float32(time.Second)
			model := mgl32.HomogRotate3D(totalTime*math.Pi/4, mgl32.Vec3{0, 1, 0}.Normalize())
			crane.CalculateFinalTransformations(anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, angle0}}, anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, angle1}}, anim.Transform{mgl32.Vec3{1, 1, 1}, mgl32.Vec3{}, mgl32.Vec3{0, 0, angle2}})

			//Upload unifrom variables
			gl.UniformMatrix4fv(mvpUniform, 1, false, &mvp[0])
			gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
			gl.Uniform3f(lightPosUniform, lightPos[0], lightPos[1], lightPos[2])
			gl.UniformMatrix4fv(boneUniforms, 3, false, &crane.FinalTransforms[0][0])

			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
			//gl.DrawArrays(gl.TRIANGLES, 0, int32(len(vertices)))
			gl.DrawElements(gl.LINE_STRIP, int32(len(elements)), gl.UNSIGNED_INT, gl.PtrOffset(0))

			time.Sleep(frameTime - time.Since(startTime))
			if int(time.Since(t0)/time.Second) > seconds {
				seconds++
				fmt.Println("fps:", frames)
				fmt.Println("Angle1: ", angle1)
				fmt.Println("Angle2: ", angle2)
				fmt.Println("LocalTransforms")
				fmt.Println(crane.LocalTransforms)
				fmt.Println("FinalTransforms")
				fmt.Println(crane.FinalTransforms)
				frames = 0
			}
			startTime = time.Now()

			window.SwapBuffers()
			glfw.PollEvents()
		}
	}(window, shaderProgram)
}

//Input function
func handleInput(window *glfw.Window, playerPos *mgl32.Vec3, lightPos *mgl32.Vec3, angle0, angle1, angle2 *float32) {
	var navigationSpeed float32 = 5.0 / fps

	//Pressing space to exit
	if window.GetKey(glfw.KeySpace) == glfw.Press ||
		window.GetKey(glfw.KeyEnter) == glfw.Press {
		window.SetShouldClose(true)
	}
	//First person motion
	if window.GetKey(glfw.KeyW) == glfw.Press {
		playerPos[2] -= navigationSpeed
	} else if window.GetKey(glfw.KeyS) == glfw.Press {
		playerPos[2] += navigationSpeed
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		playerPos[0] -= navigationSpeed
	} else if window.GetKey(glfw.KeyD) == glfw.Press {
		playerPos[0] += navigationSpeed
	}
	//light motion
	if window.GetKey(glfw.KeyUp) == glfw.Press {
		lightPos[2] -= 2 * navigationSpeed
		*angle1 += 50 * navigationSpeed
	} else if window.GetKey(glfw.KeyDown) == glfw.Press {
		lightPos[2] += 2 * navigationSpeed
		*angle1 -= 50 * navigationSpeed
	}
	if window.GetKey(glfw.KeyLeft) == glfw.Press {
		lightPos[0] -= 2 * navigationSpeed
		*angle2 -= 50 * navigationSpeed
	} else if window.GetKey(glfw.KeyRight) == glfw.Press {
		lightPos[0] += 2 * navigationSpeed
		*angle2 += 50 * navigationSpeed
	}
	if window.GetKey(glfw.Key0) == glfw.Press {
		*angle0 -= 50 * navigationSpeed
	} else if window.GetKey(glfw.Key1) == glfw.Press {
		*angle0 += 50 * navigationSpeed
	}
}
