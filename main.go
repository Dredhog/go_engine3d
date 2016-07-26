package main

import (
	"fmt"
	"log"
	"math"
	"runtime"
	"time"
	"training/engine/anim"
	"training/engine/parse/obj"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}
const (
	screenWidth  = 1920
	screenHeight = 1080
	fps          = 122
	w = 25
	h = 1
)

var vertices = []float32{
	//Root Vertex
	0, 0, 0, //position
	1, 0, 0, //normal
	0, 0, //bones
	1.0, 0, //weights

	//Neck
	0, 0.1, 0, //position
	1, 0, 0, //normal
	1, 0, //bones
	0.7, 0.3, //weights

	//Left shoulder
	-0.2, 0.05, 0, //position
	1, 0, 1, //normal
	2, 4, //bones
	0.6, 4.0, //weights

	//Right shoulder
	0.2, 0.05, 0, //position
	1, 0, 1, //normal
	3, 5, //bones
	0.6, 4.0, //weights

	//Left elbow
	-0.4, 0.05, 0, //position
	1, 0, 0, //normal
	4, 6, //bones
	0.6, 4.0, //weights

	//Right elbow
	0.4, 0.05, 0, //position
	0, 0, 1, //normal
	5, 7, //bones
	0.6, 4.0, //weights

	//Left wrist
	-0.6, 0.05, 0, //position
	0, 1, 1, //normal
	6, 6, //bones
	0.6, 4.0, //weights

	//Right wrist
	0.6, 0.05, 0, //position
	0, 1, 1, //normal
	7, 7, //bones
	0.6, 4.0, //weights
	//-----------Bottom of arm----------------------
	//Left shoulder
	-0.2, -0.05, 0, //position
	1, 0, 0, //normal
	2, 4, //bones
	0.8, 0.2, //weights

	//Right shoulder
	0.2, -0.05, 0, //position
	1, 0, 0, //normal
	3, 5, //bones
	0.8, 0.2, //weights

	//Left elbow
	-0.4, -0.05, 0, //position
	1, 0, 0, //normal
	4, 6, //bones
	0.8, 0.2, //weights

	//Right elbow
	0.4, -0.05, 0, //position
	1, 0, 0, //normal
	5, 7, //bones
	0.8, 0.2, //weights

	//Left wrist
	-0.6, -0.05, 0, //position
	1, 0, 0, //normal
	6, 6, //bones
	0.8, 0.2, //weights

	//Right wrist
	0.6, -0.05, 0, //position
	1, 0, 0, //normal
	7, 7, //bones
	0.8, 0.2, //weights
}

var elements []uint32 = []uint32{
	0, 1,
	0, 2,
	0, 3,
	2, 4,
	3, 5,
	4, 6,
	5, 7,
	8, 10,
	10, 12,
	9, 11,
	11, 13,
}

func main() {

	//---------------------------------ANIMATION---------------------------------
	arm := anim.NewSkeleton(8)
	arm.Bones[0] = anim.Bone{Name: "root          ", Index: 0, ParentIndex: 0, ToRoot: anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{}}}
	arm.Bones[1] = anim.Bone{Name: "neck          ", Index: 1, ParentIndex: 0, ToRoot: anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{}}}
	arm.Bones[2] = anim.Bone{Name: "left_clavicle ", Index: 2, ParentIndex: 0, ToRoot: anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{}}}
	arm.Bones[3] = anim.Bone{Name: "right_clavicle", Index: 3, ParentIndex: 0, ToRoot: anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{}}}
	arm.Bones[4] = anim.Bone{Name: "left_shoulder ", Index: 4, ParentIndex: 2, ToRoot: anim.Transform{[3]float32{1, 1, 1}, [3]float32{0.2, 0, 0}, [3]float32{}}}
	arm.Bones[5] = anim.Bone{Name: "right_shoulder", Index: 5, ParentIndex: 3, ToRoot: anim.Transform{[3]float32{1, 1, 1}, [3]float32{-0.2, 0, 0}, [3]float32{}}}
	arm.Bones[6] = anim.Bone{Name: "left_elbow    ", Index: 6, ParentIndex: 4, ToRoot: anim.Transform{[3]float32{1, 1, 1}, [3]float32{0.4, 0, 0}, [3]float32{}}}
	arm.Bones[7] = anim.Bone{Name: "right_elbow   ", Index: 7, ParentIndex: 5, ToRoot: anim.Transform{[3]float32{1, 1, 1}, [3]float32{-0.4, 0, 0}, [3]float32{}}}
	//---------------------------------------------------------------------------

	vertices, elements := obj.ParseFile("rabbot.obj", true, true)

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

	shaderProgram, err := newProgram("diffuse_alfa")
	if err != nil {
		panic(err)
	}
	//Game loop function
	func(window *glfw.Window, shaderProgram uint32) {
		gl.Enable(gl.CULL_FACE)
		gl.CullFace(gl.BACK)
		gl.Enable(gl.DEPTH_TEST)
		gl.DepthFunc(gl.LESS)
		gl.ClearColor(0.3, 0.3, 0.4, 1.0)

		rabbitDiffuse, err := newTexture("rabbit.png")
		if err != nil {
			panic(err)
		}

		var personMesh Mesh
		personMesh.Init(vertices, elements, []Texture{{rabbitDiffuse, "DIFFUSE"}})
		//Get uniforms from shader
		mvpUniform := gl.GetUniformLocation(shaderProgram, gl.Str("mvp_mat\x00"))
		modelUniform := gl.GetUniformLocation(shaderProgram, gl.Str("model_mat\x00"))
		lightPosUniform := gl.GetUniformLocation(shaderProgram, gl.Str("light_position\x00"))
		boneUniforms := gl.GetUniformLocation(shaderProgram, gl.Str("bone_mat\x00"))

		model := mgl32.Ident4()
		projection := mgl32.Perspective(math.Pi/4, 1.6, 0.1, 100.0)
		playerPos := mgl32.Vec3{0, 1, 2}
		lightPos := mgl32.Vec3{0, 2, 0}
		angle0 := float32(0)
		angle1 := float32(0)
		angle2 := float32(0)
		angle3 := float32(0)
		angle4 := float32(0)
		angle5 := float32(0)
		angle6 := float32(0)
		angle7 := float32(0)
		t0 := time.Now()
		startTime := t0
		frameTime := time.Second / fps
		frames := 0
		seconds := 0

		for !window.ShouldClose() {
			//Get input
			glfw.PollEvents()
			handleInput(window, &playerPos, &lightPos, &angle0, &angle1, &angle2, &angle3, &angle4, &angle5, &angle6, &angle7)

			//UPDATE VARIABLES
			camera := mgl32.LookAtV(playerPos, playerPos.Add(mgl32.Vec3{0, -0.5, -3}), mgl32.Vec3{0, 1, 0})
			mvp := projection.Mul4(camera)
			arm.CalculateFinalTransformations(
				anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, angle0}},
				anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, angle1}},
				anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, angle2}},
				anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, angle3}},
				anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, angle4}},
				anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, angle5}},
				anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, angle6}},
				anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, angle7}},
			)

			//Upload unifrom variables
			gl.UniformMatrix4fv(mvpUniform, 1, false, &mvp[0])
			gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
			gl.Uniform3f(lightPosUniform, lightPos[0], lightPos[1], lightPos[2])
			gl.UniformMatrix4fv(boneUniforms, 8, false, &arm.FinalTransforms[0][0])

			gl.UseProgram(shaderProgram)

			//FPS: update, maintain, display
			frames++
			time.Sleep(frameTime - time.Since(startTime))
			if int(time.Since(t0)/time.Second) > seconds {
				seconds++
				fmt.Println("LocalTransforms")
				fmt.Println(arm.LocalTransforms)
				fmt.Println("FinalTransforms")
				fmt.Println(arm.FinalTransforms)
				fmt.Println("fps:", frames)
				frames = 0
			}
			startTime = time.Now()

			//Perform rendering
			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
			for i := 0; i < w; i++ {
				for j := 0; j < h; j++{
					model = mgl32.Translate3D(float32(i)*1.5, 0, float32(j)*1.5)
					gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
					personMesh.Draw()
				}
			}
			window.SwapBuffers()
		}
	}(window, shaderProgram)
}

//Input function
func handleInput(window *glfw.Window, playerPos *mgl32.Vec3, lightPos *mgl32.Vec3, angle0, angle1, angle2, angle3, angle4, angle5, angle6, angle7 *float32) {
	var navigationSpeed float32 = 5.0 / fps
	pressedShift := false

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
	if window.GetKey(glfw.KeyLeftShift) == glfw.Press {
		pressedShift = true
	}
	//light motion
	if window.GetKey(glfw.KeyUp) == glfw.Press {
		lightPos[2] -= 5 * navigationSpeed
	} else if window.GetKey(glfw.KeyDown) == glfw.Press {
		lightPos[2] += 5 * navigationSpeed
	}
	if window.GetKey(glfw.KeyLeft) == glfw.Press {
		lightPos[0] -= 5 * navigationSpeed
	} else if window.GetKey(glfw.KeyRight) == glfw.Press {
		lightPos[0] += 5 * navigationSpeed
	}
	//Joint motion
	if pressedShift {
		navigationSpeed *= -1
	}
	if window.GetKey(glfw.Key0) == glfw.Press {
		*angle0 += 50 * navigationSpeed
	}
	if window.GetKey(glfw.Key1) == glfw.Press {
		*angle1 += 50 * navigationSpeed
	}
	if window.GetKey(glfw.Key2) == glfw.Press {
		*angle2 += 50 * navigationSpeed
	}
	if window.GetKey(glfw.Key3) == glfw.Press {
		*angle3 += 50 * navigationSpeed
	}
	if window.GetKey(glfw.Key4) == glfw.Press {
		*angle4 += 50 * navigationSpeed
	}
	if window.GetKey(glfw.Key5) == glfw.Press {
		*angle5 += 50 * navigationSpeed
	}
	if window.GetKey(glfw.Key6) == glfw.Press {
		*angle6 += 50 * navigationSpeed
	}
	if window.GetKey(glfw.Key7) == glfw.Press {
		*angle7 += 50 * navigationSpeed
	}
	if pressedShift {
		navigationSpeed *= -1
	}
}
