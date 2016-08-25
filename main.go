package main

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"log"
	"math"
	"runtime"
	"time"

	"training/engine/anim"
	"training/engine/load/shader"
	"training/engine/load/texture"
	"training/engine/parse/collada"
	"training/engine/types"
)

func init() {
	runtime.LockOSThread()
}

const (
	screenWidth  = 1920
	screenHeight = 1080
	fps          = 244
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
	window, err := glfw.CreateWindow(screenWidth, screenHeight, "Opengl", glfw.GetPrimaryMonitor(), nil)
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

		model, err := collada.ParseModel("data/model/turner_skinned_fixed_skin.dae")
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
		modelRotationUniform := gl.GetUniformLocation(shaderProgram, gl.Str("model_rotation_mat\x00"))
		lightPosUniform := gl.GetUniformLocation(shaderProgram, gl.Str("light_position\x00"))
		boneUniforms := gl.GetUniformLocation(shaderProgram, gl.Str("bone_mat\x00"))

		//Decalring gameplay/animation/framerate variables
		modelMatrix := mgl32.Ident4()
		modelRotationMatrix := mgl32.Ident4()
		toComMatrix := mgl32.Translate3D(0, 0.7, 0)
		toComInvMatrix := toComMatrix.Inv()
		projection := mgl32.Perspective(math.Pi/4, 1.6, 0.1, 100.0)
		cameraPosition := mgl32.Vec3{0, 1.2, 2.5}
		cameraDirection := mgl32.Vec3{0, -1, -6}
		camera := mgl32.LookAtV(cameraPosition, cameraPosition.Add(cameraDirection), mgl32.Vec3{0, 1, 0})
		mvpMatrix := projection.Mul4(camera)
		player := player{Dir: mgl32.Vec3{0, 0, 1}, Angle: 0}
		lightPosition := mgl32.Vec3{0, 1, 2}
		forward := mgl32.Vec3{0, 0, -1}
		left := mgl32.Vec3{-1, 0, 0}
		NegYAxis := mgl32.Vec3{0, -1, 0}
		zAxis := mgl32.Vec3{0, 0, 1}

		frames := 0
		seconds := 0
		walkAnimation := anim.Animation{Keyframes: []anim.Keyframe{keyframe00, keyframe01, keyframe02, keyframe03, keyframe04}}
		walkAnimation.SetTotalTicks()
		runAnimation := anim.Animation{Keyframes: []anim.Keyframe{keyframe10, keyframe11, keyframe12, keyframe13, keyframe14}}
		runAnimation.SetTotalTicks()
		model.Animator = &anim.Animator{CurrentAnimation: &walkAnimation, UpcomingAnimation: &runAnimation, CurrentKeyframe: anim.Keyframe{Transforms: make([]anim.Transform, len(keyframe00.Transforms))}, UpcomingKeyframe: anim.Keyframe{Transforms: make([]anim.Transform, len(keyframe00.Transforms))}, ResultKeyframe: anim.Keyframe{Transforms: make([]anim.Transform, len(keyframe00.Transforms))}}
		ticks := float32(0)
		frameStart := glfw.GetTime()
		frameTime := 1 / float32(fps)
		deltaTime := frameTime
		t := float32(0)

		for !window.ShouldClose() {
			//Get input
			glfw.PollEvents()
			handleInput(window, deltaTime, &player, &lightPosition, &forward, &left, &zAxis, &t)

			deltaTime = float32(glfw.GetTime() - frameStart)
			if deltaTime < frameTime {
				time.Sleep(time.Duration(int64((frameTime - deltaTime) * 1e9)))
			}
			deltaTime = frameTime
			frameStart = glfw.GetTime()

			//update variables
			tiltAxis := player.Dir.Cross(NegYAxis)
			modelRotationMatrix = toComMatrix.Mul4(mgl32.HomogRotate3D(player.TiltAngle, tiltAxis).Mul4(mgl32.HomogRotate3DY(player.Angle).Mul4(toComInvMatrix)))
			modelMatrix = mgl32.Translate3D(player.Position[0], player.Position[1], player.Position[2]).Mul4(modelRotationMatrix)
			ticks += 250 * t  * deltaTime
			if t < 0 {
				t = 0
			} else if t > 1 {
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
			gl.UniformMatrix4fv(modelRotationUniform, 1, false, &modelRotationMatrix[0])
			gl.Uniform3f(lightPosUniform, lightPosition[0], lightPosition[1], lightPosition[2])
			gl.UniformMatrix4fv(boneUniforms, 15, false, &model.Skeleton.FinalMatrices[0][0])

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

type player struct {
	Position     mgl32.Vec3
	Velocity     mgl32.Vec3
	AccDirection mgl32.Vec3
	Dir          mgl32.Vec3
	TiltAngle    float32
	DestAngle    float32
	Angle        float32
}

//Input function
func handleInput(window *glfw.Window, deltaTime float32, player *player, lightPosition, forward, left, zAxis *mgl32.Vec3, t *float32) {
	var lightSpeed float32 = 5 * deltaTime
	var maxTiltAngle float32 = 0.25
	var tiltSpeed float32 = 0.6
	var maxSpeed float32 = 8
	var angularVelocity float32 = 15
	var acc float32 = 30
	var deacc float32 = 20

	//Pressing space to exit
	if window.GetKey(glfw.KeySpace) == glfw.Press {
		window.SetShouldClose(true)
	}

	//PLAYER MOTION
	player.AccDirection = mgl32.Vec3{}
	if window.GetKey(glfw.KeyW) == glfw.Press {
		player.AccDirection = player.AccDirection.Add(*forward)
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		player.AccDirection = player.AccDirection.Add(forward.Mul(-1))
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		player.AccDirection = player.AccDirection.Add(*left)
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		player.AccDirection = player.AccDirection.Add(left.Mul(-1))
	}
	if accDirLen := player.AccDirection.Len(); accDirLen != 0 {
		player.AccDirection = player.AccDirection.Mul(1/accDirLen)
		player.Velocity = player.Velocity.Add(player.AccDirection.Mul(deltaTime * acc))
		player.Dir = player.Dir.Add(player.AccDirection.Mul(deltaTime * acc * 2))
	} else if speed := player.Velocity.Len(); speed >= deltaTime*deacc {
		player.Velocity = player.Velocity.Sub(player.Velocity.Mul((1/speed) * deltaTime * deacc))
		player.TiltAngle -= tiltSpeed * deltaTime
	} else {
		player.Velocity = mgl32.Vec3{}
		player.TiltAngle -= tiltSpeed * deltaTime
	}
	//Limit the player's velocity
	if speed := player.Velocity.Len(); speed > maxSpeed {
		player.Velocity = player.Velocity.Mul((1 / speed) * maxSpeed)
	}
	if speed := player.Velocity.Len(); speed != 0 {
		player.Dir = player.Velocity.Mul(1/speed)
		player.TiltAngle += player.Dir.Dot(player.AccDirection) * tiltSpeed * deltaTime
	}
	//Update the player's position
	player.Position = player.Position.Add(player.Velocity.Mul(deltaTime))

	//Limit the player's tilt angle
	if player.TiltAngle > maxTiltAngle{
		player.TiltAngle = maxTiltAngle
	} else if player.TiltAngle < 0{
		player.TiltAngle = 0
	}

	//Determine the player's orientation
	dtr := float32(math.Pi / 180)
	player.DestAngle = float32(math.Acos(float64(player.Dir.Dot(*zAxis))))
	if player.Dir[0] < 0 {
		player.DestAngle = (-1)*player.DestAngle + 360*dtr
	}
	if player.Angle < 0 {
		player.Angle += 360 * dtr
	} else if player.Angle > 360*dtr {
		player.Angle -= 360 * dtr
	}
	if delta := player.DestAngle - player.Angle; delta > 0 {
		switch {
		case delta <= 180*dtr:
			if player.Angle+angularVelocity*deltaTime < player.DestAngle {
				player.Angle += angularVelocity * deltaTime
			}
		case delta < 360*dtr:
			player.Angle -= angularVelocity * deltaTime
		}
	} else {
		switch {
		case -180*dtr < delta:
			if player.Angle-angularVelocity*deltaTime > player.DestAngle {
				player.Angle -= angularVelocity * deltaTime
			}
		case -360*dtr < delta:
			player.Angle += angularVelocity * deltaTime
		}
	}
	//LIGHT MOTION
	if window.GetKey(glfw.KeyUp) == glfw.Press {
		lightPosition[2] -= 5 * lightSpeed
	} else if window.GetKey(glfw.KeyDown) == glfw.Press {
		lightPosition[2] += 5 * lightSpeed
	}
	if window.GetKey(glfw.KeyLeft) == glfw.Press {
		lightPosition[0] -= 5 * lightSpeed
	} else if window.GetKey(glfw.KeyRight) == glfw.Press {
		lightPosition[0] += 5 * lightSpeed
	}
	//RESET BUTTON
	if window.GetKey(glfw.KeyR) == glfw.Press {
		lightPosition[0] = 0
		lightPosition[1] = 1
		lightPosition[2] = 2
		player.Position[0] = 0
		player.Position[1] = 0
		player.Position[2] = 0
		player.Velocity[0] = 0
		player.Velocity[1] = 0
		player.Velocity[2] = 0
		player.Dir[0] = 0
		player.Dir[1] = 0
		player.Dir[2] = 1
	}
	//ANIMATION BLENDING
	*t = player.Velocity.Len()/maxSpeed
	if window.GetKey(glfw.Key0) == glfw.Press {
		*t -= 2 * deltaTime
	}
	if window.GetKey(glfw.Key1) == glfw.Press {
		*t += 2 * deltaTime
	}
}

var keyframe10 = anim.Keyframe{Ticks: 0, Transforms: []anim.Transform{
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, -10}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, -30}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{90, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 10}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{-20, 0, 30}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{20, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{-10, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{20, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{-90, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}}}}
var keyframe11 = anim.Keyframe{Ticks: 30, Transforms: []anim.Transform{
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 20, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, -10}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{45, 0, -30}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{130, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, -20, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{-50, -20, 30}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{60, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, -30, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{-40, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{-20, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{-20, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{80, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{-40, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{-50, 0, 0}}}}
var keyframe12 = anim.Keyframe{Ticks: 60, Transforms: []anim.Transform{
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, -10}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{-20, 0, -30}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{20, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 10}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 30}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{90, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{20, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{-90, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{-10, 0, 0}}}}
var keyframe13 = anim.Keyframe{Ticks: 90, Transforms: []anim.Transform{
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, -20, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 20, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{-50, 20, -30}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{60, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 10}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{45, 0, 30}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{130, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 30, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{80, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{-40, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{-50, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{-40, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{-20, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{-20, 0, 0}}}}
var keyframe14 = anim.Keyframe{Ticks: 120, Transforms: []anim.Transform{
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, -10}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, -30}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{90, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 10}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{-20, 0, 30}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{20, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{-10, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{20, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{-90, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}}}}
var keyframe00 = anim.Keyframe{Ticks: 0, Transforms: []anim.Transform{
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, -10}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, -30}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 10}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 30}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}}}}
var keyframe01 = anim.Keyframe{Ticks: 30, Transforms: []anim.Transform{
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, -10}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, -30}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 10}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 30}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, -5, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{-20, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{15, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{-10, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{15, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{15, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{-30, 0, 0}}}}
var keyframe02 = anim.Keyframe{Ticks: 60, Transforms: []anim.Transform{
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, -10}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, -30}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 10}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 30}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}}}}
var keyframe03 = anim.Keyframe{Ticks: 90, Transforms: []anim.Transform{
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, -10}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, -30}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 10}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 30}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 5, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{15, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{15, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{-30, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{-20, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{15, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{-10, 0, 0}}}}
var keyframe04 = anim.Keyframe{Ticks: 120, Transforms: []anim.Transform{
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, -10}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, -30}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 10}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 30}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}}}}
