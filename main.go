package main

import (
	"fmt"
	"log"
	"math"
	"runtime"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"

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
	fps          = 120
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

		mesh, skeleton, err0, err1 := collada.ParseMeshSkeleton("data/model/turner.dae")
		if err0 != nil {
			log.Fatalln(err0)
		} else if err1 != nil {
			log.Fatalln(err1)
		}
		model := types.Model{Mesh: mesh}
		model.Animator, err = anim.NewAnimator(skeleton, []anim.Animation{anim.Animation{Duration: keyframe04.SampleTime, Keyframes: []anim.Keyframe{keyframe00, keyframe01, keyframe02, keyframe03, keyframe04}},
			anim.Animation{Duration: keyframe14.SampleTime, Keyframes: []anim.Keyframe{keyframe10, keyframe11, keyframe12, keyframe13, keyframe14}},
			anim.Animation{Duration: keyframe21.SampleTime, Keyframes: []anim.Keyframe{keyframe20, keyframe21}}})
		if err != nil {
			panic(err)
		}
		playerShader, err := shader.NewProgram("player_diffuse_specular")
		if err != nil {
			log.Fatalln(err)
		}
		playerDiffuseTexture, err := texture.NewTexture("rb.png")
		if err != nil {
			log.Fatalln(err)
		}
		playerSpecularTexture, err := texture.NewTexture("pink.png")
		if err != nil {
			log.Fatalln(err)
		}
		model.Mesh.Textures = []types.Texture{{playerSpecularTexture, "specular"}, {playerDiffuseTexture, "diffuse"}}

		light, _, err0, _ := collada.ParseMeshSkeleton("data/model/light.dae")
		if err != nil {
			log.Fatalln(err0)
		}
		lightShader, err := shader.NewProgram("light")
		if err != nil {
			fmt.Print(err)
			log.Fatalln(err)
		}
		//Get uniforms from shader
		vpUniform := gl.GetUniformLocation(playerShader, gl.Str("vp_mat\x00"))
		modelUniform := gl.GetUniformLocation(playerShader, gl.Str("model_mat\x00"))
		modelRotationUniform := gl.GetUniformLocation(playerShader, gl.Str("model_rotation_mat\x00"))
		lightPosUniform := gl.GetUniformLocation(playerShader, gl.Str("light_position\x00"))
		boneUniforms := gl.GetUniformLocation(playerShader, gl.Str("bone_mat\x00"))

		//Decalring gameplay/animation/framerate variables
		toComMatrix := mgl32.Translate3D(0, 0.7, 0)
		toComInvMatrix := toComMatrix.Inv()
		worldGizmo := gizmo{xAxis: mgl32.Vec3{1, 0, 0}, yAxis: mgl32.Vec3{0, 1, 0}, zAxis: mgl32.Vec3{0, 0, 1}}
		camera := newCamera(mgl32.Vec3{0, 1.2, 5}, mgl32.Vec3{0, -1, -6}, &worldGizmo)
		player := player{Dir: worldGizmo.zAxis, Up: worldGizmo.yAxis}
		lightPosition := mgl32.Vec3{0, 1, 2}
		frameTimer := frameTimer{gameLoopStart: float32(glfw.GetTime()), desiredFrameTime: 1 / float32(fps)}
		t := float32(0)
		h := float32(0.5)

		for !window.ShouldClose() {
			//Update the time manager
			frameTimer.Update()

			//Get input
			glfw.PollEvents()
			handleInput(window, &worldGizmo, &frameTimer, &player, &camera, &lightPosition, &t, &h)

			//update variables
			modelRotationMatrix := toComMatrix.Mul4(mgl32.HomogRotate3D(player.TiltAngle, player.TiltAxis).Mul4(mgl32.HomogRotate3DY(player.Angle).Mul4(toComInvMatrix)))
			modelMatrix := mgl32.Translate3D(player.Position[0], player.Position[1], player.Position[2]).Mul4(modelRotationMatrix)
			model.Animator.Update(frameTimer.deltaTime, t, h)

			//FPS display, and debug information
			if frameTimer.isSecondMark {
				fmt.Println("fps:", frameTimer.currentFps)
				fmt.Printf("time: %v;\nt: %v\n", frameTimer.frameStart, t)
				fmt.Printf("speed: %v;\nAcceleration: %v;\n\n", player.Velocity.Len(), player.AccDirection.Len())
			}

			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

			//Update the light shader
			gl.UseProgram(lightShader)
			//Upload unifrom variables
			gl.UniformMatrix4fv(gl.GetUniformLocation(lightShader, gl.Str("vp_mat\x00")), 1, false, &camera.VPMatrix[0])
			gl.Uniform3f(gl.GetUniformLocation(lightShader, gl.Str("light_position\x00")), lightPosition[0], lightPosition[1], lightPosition[2])
			light.Draw(lightShader)

			//Update the player shader
			gl.UseProgram(playerShader)
			//Upload unifrom variables
			gl.UniformMatrix4fv(vpUniform, 1, false, &camera.VPMatrix[0])
			gl.UniformMatrix4fv(modelUniform, 1, false, &modelMatrix[0])
			gl.UniformMatrix4fv(modelRotationUniform, 1, false, &modelRotationMatrix[0])
			gl.Uniform3f(lightPosUniform, lightPosition[0], lightPosition[1], lightPosition[2])
			gl.UniformMatrix4fv(boneUniforms, 15, false, &model.Animator.GlobalPoseMatrices[0][0])
			model.Mesh.Draw(playerShader)
			window.SwapBuffers()
		}
	}(window)
}

func clamp(a float32, i *float32, b float32) {
	if *i < a {
		*i = a
	} else if *i > b {
		*i = b
	}
}

func max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

//Input function
func handleInput(window *glfw.Window, world *gizmo, frameTimer *frameTimer, player *player, camera *camera, lightPosition *mgl32.Vec3, s, h *float32) {
	var maxTiltAngle float32 = 0.25
	var lightSpeed float32 = 10
	var maxSpeed float32 = 10
	//var headTurnSpeed float32 = 5
	var jumpVerticalSpeed float32 = 5
	var angularVelocity float32 = 15
	var acc float32 = 30
	var deacc float32 = 15

	//Pressing Esc to exit
	if window.GetKey(glfw.KeyEscape) == glfw.Press ||
		window.GetKey(glfw.KeyEnter) == glfw.Press {
		window.SetShouldClose(true)
	}
	if window.GetKey(glfw.KeyLeftShift) == glfw.Press {
		frameTimer.deltaTime /= 8
	}
	deltaTime := frameTimer.deltaTime
	if !player.InAir && window.GetKey(glfw.KeySpace) == glfw.Press {
		player.InAir = true
		player.Velocity = player.Velocity.Add(world.yAxis.Mul(jumpVerticalSpeed))
	}

	//PLAYER MOTION
	player.AccDirection = mgl32.Vec3{}
	if window.GetKey(glfw.KeyW) == glfw.Press {
		player.AccDirection = player.AccDirection.Add(camera.Forward)
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		player.AccDirection = player.AccDirection.Add(camera.Forward.Mul(-1))
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		player.AccDirection = player.AccDirection.Add(camera.Left)
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		player.AccDirection = player.AccDirection.Add(camera.Left.Mul(-1))
	}
	if !player.InAir {
		if accDirLen := player.AccDirection.Len(); accDirLen != 0 {
			player.AccDirection = player.AccDirection.Mul(1 / accDirLen)
			if temp := player.Up.Normalize().Add(player.AccDirection.Mul(3 * deltaTime)); math.Acos(float64(temp.Normalize().Dot(world.yAxis))) < float64(maxTiltAngle) {
				player.Up = temp.Mul(1 / player.Up[1])
			}
			player.Velocity = player.Velocity.Add(player.AccDirection.Mul(acc * deltaTime))
			player.Dir = player.Dir.Add(player.AccDirection.Mul(deltaTime * acc * 2))
		} else if speed := player.Velocity.Len(); speed >= deltaTime*deacc {
			player.Velocity = player.Velocity.Sub(player.Velocity.Mul((1 / speed) * deltaTime * deacc))
		} else {
			player.Velocity = mgl32.Vec3{}
		}
		//Limit the player's velocity
		if speed := player.Velocity.Len(); speed != 0 {
			player.Dir = player.Velocity.Mul(1 / speed)
			if speed > maxSpeed {
				player.Velocity = player.Dir.Mul(maxSpeed)
			}
		}
	} else {
		player.Velocity = player.Velocity.Add(world.yAxis.Mul(-9.81 * deltaTime))
	}
	//basic falling collision
	if player.Position[1] < 0 {
		player.InAir = false
		player.Position[1] = 0
		player.Velocity[1] = 0
	}

	//Determine the player's tilt
	if tiltAxis := world.yAxis.Cross(player.Up.Normalize()); tiltAxis.Len() != 0 {
		player.TiltAxis = tiltAxis
		player.TiltAngle = float32(math.Asin(float64(player.TiltAxis.Len())))
		player.TiltAxis = player.TiltAxis.Normalize()
	}

	//Update the player's position
	player.Position = player.Position.Add(player.Velocity.Mul(deltaTime))

	//Determine the player's rotation around the y axis
	dtr := float32(math.Pi / 180)
	player.DestAngle = float32(math.Acos(float64(player.Dir.Dot(world.zAxis))))
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
			} else {
				player.Angle = player.DestAngle
			}
		case delta < 360*dtr:
			player.Angle -= angularVelocity * deltaTime
		}
	} else {
		switch {
		case -180*dtr <= delta:
			if player.Angle-angularVelocity*deltaTime > player.DestAngle {
				player.Angle -= angularVelocity * deltaTime
			} else {
				player.Angle = player.DestAngle
			}
		case -360*dtr < delta:
			player.Angle += angularVelocity * deltaTime
		}
	}
	//LIGHT MOTION
	if window.GetKey(glfw.KeyUp) == glfw.Press {
		lightPosition[2] -= lightSpeed * deltaTime
	}
	if window.GetKey(glfw.KeyDown) == glfw.Press {
		lightPosition[2] += lightSpeed * deltaTime
	}
	if window.GetKey(glfw.KeyLeft) == glfw.Press {
		lightPosition[0] -= lightSpeed * deltaTime
	}
	if window.GetKey(glfw.KeyRight) == glfw.Press {
		lightPosition[0] += lightSpeed * deltaTime
	}
	//HeadRotation
	toLight := lightPosition.Sub(player.Position)
	toLight[1] = 0
	toLight = toLight.Normalize()
	t := toLight.Dot(player.Dir)
	t += -1
	clamp(-1, &t, 0)
	//right hemisphere, where 0's at 12 o'clock
	if toLight.Cross(player.Dir)[1] < 0 {
		t *= -1
	}
	fmt.Println(t)
	*h += 4 * (t - (*h)) * deltaTime
	//ANIMATION BLENDING
	*s = player.Velocity.Len() / maxSpeed
	clamp(-1, h, 1)
	clamp(0, s, 1)
	//RESET BUTTON
	if window.GetKey(glfw.KeyR) == glfw.Press {
		*lightPosition = mgl32.Vec3{0, 1, 2}
		player.Position = mgl32.Vec3{0, 0, 0}
		player.Velocity = mgl32.Vec3{0, 0, 0}
		player.Dir = mgl32.Vec3{0, 0, 1}
		player.Up = mgl32.Vec3{0, 1, 0}
		player.TiltAngle = 0
		*h = 0
	}
}
