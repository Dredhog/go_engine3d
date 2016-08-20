package main

import (
	"fmt"
	"log"
	"math"
	"time"
	"runtime"
	"training/engine/anim"
	"training/engine/load/shader"
	"training/engine/load/texture"
	"training/engine/parse/collada"
	"training/engine/types"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

func init() {
	runtime.LockOSThread()
}

const (
	screenWidth  = 1920
	screenHeight = 1080
	fps          = 60
)

func main() {
	//Set up glfw
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	//Set up the display window
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
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

	//Game loop function
	func(window *glfw.Window) {
		gl.Enable(gl.CULL_FACE)
		gl.CullFace(gl.BACK)
		gl.Enable(gl.DEPTH_TEST)
		gl.DepthFunc(gl.LESS)
		gl.ClearColor(0.2, 0.3, 0.5, 1.0)

		model, err := collada.ParseModel("data/model/turner_skinned.dae")
		if err != nil {
			log.Fatalln(err)
		}
		shaderProgram, err := shader.NewProgram("skeleton_diffuse_alfa")
		if err != nil {
			log.Fatalln(err)
		}
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

		modelMatrix := mgl32.Ident4()
		projection := mgl32.Perspective(math.Pi/4, 1.6, 0.1, 100.0)
		playerPos := mgl32.Vec3{0, 1.2, 6}
		lightPos := mgl32.Vec3{0, 1, 2}

		frames := 0
		seconds := 0

		keyframe00 := anim.Keyframe{Ticks: 0, Transforms: []anim.Transform{
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 20}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 30}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 25}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}}}}
		keyframe01 := anim.Keyframe{Ticks: 30, Transforms: []anim.Transform{
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 25}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 45}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 65}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}}}}
		keyframe02 := anim.Keyframe{Ticks: 60, Transforms: []anim.Transform{
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 20}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 30}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 25}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}}}}
		walkAnimation := anim.Animation{Keyframes: []anim.Keyframe{keyframe00, keyframe01, keyframe02}}
		walkAnimation.SetTotalTicks()
		keyframe10 := anim.Keyframe{Ticks: 0, Transforms: []anim.Transform{
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
			anim.Transform{[3]float32{1.5, 1.5, 1.5}, [3]float32{}, [3]float32{0, 0, 0}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}}}}
		keyframe11 := anim.Keyframe{Ticks: 30, Transforms: []anim.Transform{
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, -20}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}}}}
		keyframe12 := anim.Keyframe{Ticks: 60, Transforms: []anim.Transform{
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
			anim.Transform{[3]float32{1.5, 1.5, 1.5}, [3]float32{}, [3]float32{0, 0, 0}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}}}}
		keyframe13 := anim.Keyframe{Ticks: 90, Transforms: []anim.Transform{
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 20}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}}}}
		keyframe14 := anim.Keyframe{Ticks: 120, Transforms: []anim.Transform{
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
			anim.Transform{[3]float32{1.5, 1.5, 1.5}, [3]float32{}, [3]float32{0, 0, 0}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
			anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}}}}
		runAnimation := anim.Animation{Keyframes: []anim.Keyframe{keyframe10, keyframe11, keyframe12, keyframe13, keyframe14}}
		runAnimation.SetTotalTicks()
		model.Animator = &anim.Animator{CurrentAnimation: &walkAnimation, UpcomingAnimation: &runAnimation, CurrentKeyframe: anim.Keyframe{Transforms: make([]anim.Transform, len(keyframe00.Transforms))}, UpcomingKeyframe: anim.Keyframe{Transforms: make([]anim.Transform, len(keyframe00.Transforms))}, ResultKeyframe: anim.Keyframe{Transforms: make([]anim.Transform, len(keyframe00.Transforms))}}
		ticks := float32(0)
		frameStart := glfw.GetTime()
		frameTime := 1/float32(fps)
		deltaTime := frameTime
		t := float32(0)

		for !window.ShouldClose() {
			//Get input
			glfw.PollEvents()
			handleInput(window, deltaTime, &playerPos, &lightPos, &t)

			deltaTime = float32(glfw.GetTime() - frameStart)
			if deltaTime < frameTime{
				time.Sleep(time.Duration(int64((frameTime-deltaTime)*1e9)))
			}
			deltaTime = frameTime
			frameStart = glfw.GetTime()

			//update variables
			ticks += 50 * deltaTime
			camera := mgl32.LookAtV(playerPos, playerPos.Add(mgl32.Vec3{0, -1, -6}), mgl32.Vec3{0, 1, 0})
			mvpMatrix := projection.Mul4(camera)
			if t < 0{
				t = 0
			} else if t > 1{
				t = 1
			}
			if err := model.Animator.BlendAnimations(ticks, t); err != nil {
				panic(err)
			}
			if err = model.Skeleton.ApplyKeyframe(&model.Animator.ResultKeyframe); err != nil {
				panic(err)
			}

			//Upload unifrom variables
			gl.UniformMatrix4fv(mvpUniform, 1, false, &mvpMatrix[0])
			gl.UniformMatrix4fv(modelUniform, 1, false, &modelMatrix[0])
			gl.Uniform3f(lightPosUniform, lightPos[0], lightPos[1], lightPos[2])
			gl.UniformMatrix4fv(boneUniforms, 8, false, &model.Skeleton.FinalMatrices[0][0])

			//FPS: update, maintain, display
			frames++
			if int(glfw.GetTime()) > seconds {
				seconds++
				fmt.Println("fps:", frames)
				fmt.Printf("ticks:\t%v; t\t%v\n", ticks, t)
				frames = 0
			}

			//Perform rendering
			gl.UseProgram(shaderProgram)
			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
			model.Mesh.Draw()
			window.SwapBuffers()
		}
	}(window)
}

//Input function
func handleInput(window *glfw.Window, deltaTime float32, playerPos *mgl32.Vec3, lightPos *mgl32.Vec3, t *float32) {
	var navigationSpeed float32 = 5 * deltaTime

	//Pressing space/enter to exit
	if window.GetKey(glfw.KeySpace) == glfw.Press{
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
		lightPos[2] -= 5 * navigationSpeed
	} else if window.GetKey(glfw.KeyDown) == glfw.Press {
		lightPos[2] += 5 * navigationSpeed
	}
	if window.GetKey(glfw.KeyLeft) == glfw.Press {
		lightPos[0] -= 5 * navigationSpeed
	} else if window.GetKey(glfw.KeyRight) == glfw.Press {
		lightPos[0] += 5 * navigationSpeed
	}
	 if window.GetKey(glfw.Key0) == glfw.Press{
		*t -= 2 * deltaTime
	}
	if window.GetKey(glfw.Key1) == glfw.Press{
		*t += 2 * deltaTime
	}
}
