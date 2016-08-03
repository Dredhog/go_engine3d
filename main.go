package main

import (
	"fmt"
	"log"
	"math"
	"runtime"
	"time"
	"training/engine/anim"
	"training/engine/load/shader"
	"training/engine/load/texture"
	"training/engine/parse/collada"
	"training/engine/types"

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
)

func main() {
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
		log.Fatalln(err)
	}
	window.MakeContextCurrent()

	//Initialize Glow
	if err := gl.Init(); err != nil {
		log.Fatalln(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	model, err := collada.ParseModel("data/model/turner_skinned.dae")
	if err != nil {
		log.Fatalln(err)
	}

	shaderProgram, err := shader.NewProgram("skeleton_diffuse_alfa")
	if err != nil {
		log.Fatalln(err)
	}
	//Game loop function
	func(window *glfw.Window, shaderProgram uint32) {
		gl.Enable(gl.CULL_FACE)
		gl.CullFace(gl.BACK)
		gl.Enable(gl.DEPTH_TEST)
		gl.DepthFunc(gl.LESS)
		gl.ClearColor(0.3, 0.3, 0.4, 1.0)

		rabbitDiffuse, err := texture.NewTexture("baboon.png")
		if err != nil {
			log.Fatalln(err)
		}
		model.Mesh.Textures = []types.Texture{{rabbitDiffuse, "DIFFUSE"}}

		//Get uniforms from shader
		mvpUniform := gl.GetUniformLocation(shaderProgram, gl.Str("mvp_mat\x00"))
		modelUniform := gl.GetUniformLocation(shaderProgram, gl.Str("model_mat\x00"))
		lightPosUniform := gl.GetUniformLocation(shaderProgram, gl.Str("light_position\x00"))
		boneUniforms := gl.GetUniformLocation(shaderProgram, gl.Str("bone_mat\x00"))
		_ = boneUniforms

		modelMatrix := mgl32.Ident4()
		projection := mgl32.Perspective(math.Pi/4, 1.6, 0.1, 100.0)
		playerPos := mgl32.Vec3{0, 1.2, 6}
		lightPos := mgl32.Vec3{0, 1, 2}
		var angle0, angle1, angle2, angle3, angle4, angle5, angle6, angle7 float32
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
			camera := mgl32.LookAtV(playerPos, playerPos.Add(mgl32.Vec3{0, -1, -6}), mgl32.Vec3{0, 1, 0})
			mvpMatrix := projection.Mul4(camera)

			if err := model.Skeleton.CalculateFinalTransformations(
				anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, angle0}},
				anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, angle1}},
				anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, angle2}},
				anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, angle3}},
				anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, angle4}},
				anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, angle5}},
				anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, angle6}},
				anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, angle7}},
			); err != nil {
				panic(err)
			}

			//Upload unifrom variables
			gl.UniformMatrix4fv(mvpUniform, 1, false, &mvpMatrix[0])
			gl.UniformMatrix4fv(modelUniform, 1, false, &modelMatrix[0])
			gl.Uniform3f(lightPosUniform, lightPos[0], lightPos[1], lightPos[2])
			gl.UniformMatrix4fv(boneUniforms, 8, false, &model.Skeleton.FinalTransformations[0][0])

			gl.UseProgram(shaderProgram)

			//FPS: update, maintain, display
			frames++
			time.Sleep(frameTime - time.Since(startTime))
			if int(time.Since(t0)/time.Second) > seconds {
				seconds++
				fmt.Println("FinalTransforms")
				fmt.Println(model.Skeleton.FinalTransformations)
				fmt.Println("fps:", frames)
				frames = 0
			}
			startTime = time.Now()

			//Perform rendering
			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
			model.Mesh.Draw()
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
	//Reset Everything
	if window.GetKey(glfw.KeyR) == glfw.Press {
		*angle0 = 0
		*angle1 = 0
		*angle2 = 0
		*angle3 = 0
		*angle4 = 0
		*angle5 = 0
		*angle6 = 0
		*angle7 = 0
		*angle1 = 0
		*lightPos = mgl32.Vec3{0, 1, 2}
	}
}
