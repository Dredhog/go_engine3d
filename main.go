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
	"training/engine/collision"
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
	window, err := glfw.CreateWindow(screenWidth, screenHeight, "Engine", glfw.GetPrimaryMonitor(), nil)
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
		gl.Enable(gl.BLEND)
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
		gl.ClearColor(0.2, 0.3, 0.5, 1.0)
		window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

		//Load player data
		mesh, skeleton, err0, err1 := collada.ParseMeshSkeleton("data/model/turner1.dae")
		if err0 != nil {
			log.Fatalln(err0)
		} else if err1 != nil {
			log.Fatalln(err1)
		}
		model := types.Model{Mesh: mesh}
		model.Animator, err = anim.NewAnimator(skeleton, []anim.Animation{anim.Animation{Duration: keyframe04.SampleTime, Keyframes: []anim.Keyframe{keyframe00, keyframe01, keyframe02, keyframe03, keyframe04}},
			anim.Animation{Duration: keyframe14.SampleTime, Keyframes: []anim.Keyframe{keyframe10, keyframe11, keyframe12, keyframe13, keyframe14}},
			anim.Animation{Duration: keyframe21.SampleTime, Keyframes: []anim.Keyframe{keyframe20, keyframe21}},
			anim.Animation{Duration: keyframe32.SampleTime, Keyframes: []anim.Keyframe{keyframe30, keyframe31, keyframe32}},
			anim.Animation{Duration: keyframe41.SampleTime, Keyframes: []anim.Keyframe{keyframe40, keyframe41}}})
		model.Animator.SetPlaybackRate(0, 1.3)
		model.Animator.SetPlaybackRate(1, 1.3)
		model.Animator.SetLooping(2, false)
		if err != nil {
			panic(err)
		}
		playerShader, err := shader.NewProgram("player_diffuse_specular")
		if err != nil {
			log.Fatalln(err)
		}
		playerDiffuseTexture, err := texture.NewTexture("rb.png", gl.CLAMP_TO_EDGE)
		if err != nil {
			log.Fatalln(err)
		}
		model.Mesh.Textures = []types.Texture{{playerDiffuseTexture, "diffuse"}}

		//Load collision data
		colliderMesh, _, err0, _ := collada.ParseMeshSkeleton("data/model/rock.dae")
		if err0 != nil {
			log.Fatalln(err0)
		}
		shapeA := collision.GenerateCollisionPointsFromConvexMesh(colliderMesh)
		dynamicMesh, _, err0, _ := collada.ParseMeshSkeleton("data/model/64verts.dae")
		if err0 != nil {
			log.Fatalln(err0)
		}
		colliderShader, err0 := shader.NewProgram("diffuse_variable_color")
		if err != nil {
			log.Fatalln(err)
		}
		debugShader, err := shader.NewProgram("debug_collision")
		if err != nil {
			log.Fatalln(err)
		}
		simplexShader, err := shader.NewProgram("simplex")
		if err != nil {
			log.Fatalln(err)
		}

		//Load level data
		level, _, err0, _ := collada.ParseMeshSkeleton("data/model/dust2x2_scaled_UVs.dae")
		if err0 != nil {
			log.Fatalln(err0)
		}
		squareTexture, err := texture.NewTexture("squares_small.png", gl.REPEAT)
		if err != nil {
			log.Fatalln(err)
		}
		level.Textures = []types.Texture{{squareTexture, "diffuse"}}
		shaderDiffuseTexture, err := shader.NewProgram("diffuse_texture")
		if err != nil {
			log.Fatalln(err)
		}
		shaderPointLitTexture, err := shader.NewProgram("diffuse_texture_point_lit")
		if err != nil {
			log.Fatalln(err)
		}
		shaderDiffuseTextureWaving, err := shader.NewProgram("diffuse_texture_waving")
		if err != nil {
			log.Fatalln(err)
		}

		//Decalring gameplay/animation/framerate variables
		toComMatrix := mgl32.Translate3D(0, 0.7, 0)
		toComInvMatrix := toComMatrix.Inv()
		worldGizmo := gizmo{xAxis: mgl32.Vec3{1, 0, 0}, yAxis: mgl32.Vec3{0, 1, 0}, zAxis: mgl32.Vec3{0, 0, 1}}
		player := player{Dir: worldGizmo.zAxis, Up: worldGizmo.yAxis}
		camera := newCamera(mgl32.Vec3{0, 1.5, 5}, &player.Position, &worldGizmo)
		lightPosition := mgl32.Vec3{0, 0, 0}
		colliderPosition := mgl32.Vec3{0, 0, 0}
		colliderRotation := float32(0)
		colliderMat := mgl32.Ident4()
		identityMat := mgl32.Ident4()
		vecZero := mgl32.Vec3{}
		frameTimer := frameTimer{gameLoopStart: float32(glfw.GetTime()), desiredFrameTime: 1 / float32(fps)}
		editBone := int32(-1)
		pressedN := false
		speed := float32(0)
		head := float32(0.5)
		height := float32(0)
		red := mgl32.Vec4{1, 0, 0, 1}
		purple := mgl32.Vec4{1, 0, 1, 1}
		green := mgl32.Vec4{0, 1, 0, 1}
		blue := mgl32.Vec4{0, 0, 1, 1}
		white := mgl32.Vec4{1, 1, 1, 1}
		//black := mgl32.Vec4{0, 0, 0, 1}
		//cyan := mgl32.Vec4{0, 1, 1, 1}
		environmentShader := shaderDiffuseTexture
		shapeB := make([]mgl32.Vec3, len(shapeA))
		CSO := make([]mgl32.Vec3, len(shapeA)*len(shapeB))
		collisionSteps := 15

		for !window.ShouldClose() {
			//Update the time manager
			frameTimer.Update()

			//Get input
			glfw.PollEvents()
			handleInput(window, &worldGizmo, &frameTimer, &player, &camera, &colliderPosition, &lightPosition, &colliderRotation, &speed, &head, &height, &editBone, &collisionSteps, &pressedN, &environmentShader, shaderDiffuseTexture, shaderPointLitTexture, shaderDiffuseTextureWaving)

			//update variables
			colliderMat = mgl32.HomogRotate3DY(colliderRotation)
			colliderMat = mgl32.Translate3D(colliderPosition[0], colliderPosition[1], colliderPosition[2]).Mul4(colliderMat)
			transformShape(shapeA, shapeB, colliderPosition, colliderRotation)
			if err := calcMinkowskiDiff(shapeA, shapeB, CSO); err != nil {
				panic(err)
			}

			x, y := window.GetCursorPos()
			camera.Update(x, y, speed, &player.Dir)
			modelRotationMatrix := toComMatrix.Mul4(mgl32.HomogRotate3D(player.TiltAngle, player.TiltAxis).Mul4(mgl32.HomogRotate3DY(player.Angle).Mul4(toComInvMatrix)))
			modelMatrix := mgl32.Translate3D(player.Position[0], player.Position[1], player.Position[2]).Mul4(modelRotationMatrix)
			model.Animator.Update(frameTimer.deltaTime, func() {
				a := model.Animator
				a.SampleAtGlobalTime(0, 0)     //Sample Walk
				a.SampleAtGlobalTime(1, 1)     //Sample Run
				a.LinearBlend(0, 1, speed, 1)  //LERP(Walk, Run) => move
				a.SampleAtGlobalTime(3, 0)     //Sample Idle
				a.LinearBlend(0, 1, speed, 0)  //LERP(move, Idle) => ground
				a.SampleLinear(4, height, 1)   //Sample Jump
				a.LinearBlend(0, 1, height, 0) //LERP(ground, Jump) => body
				a.SampleLinear(2, head, 1)     //Sample head Turn
				a.AdditiveBlend(0, 1, 1.0, 0)  //AdditiveBlend(body, head) => final
			})

			//FPS display, and debug information
			if frameTimer.isSecondMark {
				fmt.Println("fps:", frameTimer.currentFps)
				fmt.Printf("time: %v;\n", frameTimer.frameStart)
				fmt.Printf("speed: %v;\nAcceleration: %v;\n\n", player.Velocity.Len(), player.AccDirection.Len())
				fmt.Printf("collider position %v;\t rotation %v\n", colliderPosition, colliderRotation)
			}

			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

			//Update the enviromnet shader
			gl.UseProgram(environmentShader)
			gl.UniformMatrix4fv(gl.GetUniformLocation(environmentShader, gl.Str("vp_mat\x00")), 1, false, &camera.VPMatrix[0])
			gl.Uniform3f(gl.GetUniformLocation(environmentShader, gl.Str("light_position\x00")), lightPosition[0], lightPosition[1], lightPosition[2])
			gl.Uniform1f(gl.GetUniformLocation(environmentShader, gl.Str("time\x00")), frameTimer.frameStart)
			level.Draw(environmentShader, gl.TRIANGLES)

			//Update the player shader
			gl.UseProgram(playerShader)
			gl.UniformMatrix4fv(gl.GetUniformLocation(playerShader, gl.Str("vp_mat\x00")), 1, false, &camera.VPMatrix[0])
			gl.UniformMatrix4fv(gl.GetUniformLocation(playerShader, gl.Str("model_mat\x00")), 1, false, &modelMatrix[0])
			gl.UniformMatrix4fv(gl.GetUniformLocation(playerShader, gl.Str("model_rotation_mat\x00")), 1, false, &modelRotationMatrix[0])
			gl.Uniform3f(gl.GetUniformLocation(playerShader, gl.Str("light_position\x00")), lightPosition[0], lightPosition[1], lightPosition[2])
			gl.UniformMatrix4fv(gl.GetUniformLocation(playerShader, gl.Str("bone_mat\x00")), 15, false, &model.Animator.GlobalPoseMatrices[0][0])
			gl.Uniform1i(gl.GetUniformLocation(playerShader, gl.Str("edit_bone\x00")), editBone)
			model.Mesh.Draw(playerShader, gl.TRIANGLES)

			//DRAW POINT MESHES
			gl.PointSize(8)
			gl.UseProgram(debugShader)
			gl.UniformMatrix4fv(gl.GetUniformLocation(debugShader, gl.Str("vp_mat\x00")), 1, false, &camera.VPMatrix[0])
			//Draw the origin
			gl.Uniform3fv(gl.GetUniformLocation(debugShader, gl.Str("var_positions\x00")), 1, &vecZero[0])
			gl.Uniform1i(gl.GetUniformLocation(debugShader, gl.Str("var_count\x00")), 1)
			gl.Uniform4fv(gl.GetUniformLocation(debugShader, gl.Str("var_color\x00")), 1, &white[0])
			dynamicMesh.Draw(debugShader, gl.POINTS)

			//Draw shape a
			gl.Uniform3fv(gl.GetUniformLocation(debugShader, gl.Str("var_positions\x00")), int32(len(shapeA)), &shapeA[0][0])
			gl.Uniform1i(gl.GetUniformLocation(debugShader, gl.Str("var_count\x00")), int32(len(shapeA)))
			gl.Uniform4fv(gl.GetUniformLocation(debugShader, gl.Str("var_color\x00")), 1, &green[0])
			dynamicMesh.Draw(debugShader, gl.POINTS)
			//Draw shape b
			gl.Uniform3fv(gl.GetUniformLocation(debugShader, gl.Str("var_positions\x00")), int32(len(shapeB)), &shapeB[0][0])
			gl.Uniform1i(gl.GetUniformLocation(debugShader, gl.Str("var_count\x00")), int32(len(shapeB)))
			gl.Uniform4fv(gl.GetUniformLocation(debugShader, gl.Str("var_color\x00")), 1, &blue[0])
			dynamicMesh.Draw(debugShader, gl.POINTS)
			//Draw the minokwski difference
			//gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
			gl.Uniform3fv(gl.GetUniformLocation(debugShader, gl.Str("var_positions\x00")), int32(len(CSO)), &CSO[0][0])
			gl.Uniform1i(gl.GetUniformLocation(debugShader, gl.Str("var_count\x00")), int32(len(CSO)))
			gl.Uniform4fv(gl.GetUniformLocation(debugShader, gl.Str("var_color\x00")), 1, &red[0])
			dynamicMesh.Draw(debugShader, gl.POINTS)
			//gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)

			//SYMPLEX RENDERING
			simplex, order, collided := collision.BGJK(shapeA, shapeB, collisionSteps)
			gl.UseProgram(simplexShader)
			gl.UniformMatrix4fv(gl.GetUniformLocation(simplexShader, gl.Str("vp_mat\x00")), 1, false, &camera.VPMatrix[0])
			gl.Uniform3fv(gl.GetUniformLocation(simplexShader, gl.Str("var_positions\x00")), int32(len(simplex)), &simplex[0][0])
			gl.Uniform1i(gl.GetUniformLocation(simplexShader, gl.Str("var_count\x00")), int32(order+1))
			gl.Uniform4fv(gl.GetUniformLocation(simplexShader, gl.Str("start_color\x00")), 1, &green[0])
			gl.Uniform4fv(gl.GetUniformLocation(simplexShader, gl.Str("end_color\x00")), 1, &purple[0])
			dynamicMesh.Draw(simplexShader, gl.LINE_STRIP)

			color := green
			if collided {
				color = blue
			}

			//Green collider
			gl.UseProgram(colliderShader)
			gl.UniformMatrix4fv(gl.GetUniformLocation(colliderShader, gl.Str("vp_mat\x00")), 1, false, &camera.VPMatrix[0])
			gl.UniformMatrix4fv(gl.GetUniformLocation(colliderShader, gl.Str("model_mat\x00")), 1, false, &colliderMat[0])
			gl.Uniform3fv(gl.GetUniformLocation(colliderShader, gl.Str("var_color\x00")), 1, &color[0])
			colliderMesh.Draw(colliderShader, gl.TRIANGLES)

			//Red collider
			gl.UseProgram(colliderShader)
			gl.UniformMatrix4fv(gl.GetUniformLocation(colliderShader, gl.Str("vp_mat\x00")), 1, false, &camera.VPMatrix[0])
			gl.UniformMatrix4fv(gl.GetUniformLocation(colliderShader, gl.Str("model_mat\x00")), 1, false, &identityMat[0])
			gl.Uniform3fv(gl.GetUniformLocation(colliderShader, gl.Str("var_color\x00")), 1, &red[0])
			colliderMesh.Draw(colliderShader, gl.TRIANGLES)

			window.SwapBuffers()
		}
	}(window)
}

func calcMinkowskiDiff(shapeA, shapeB, CSO []mgl32.Vec3) error {
	if len(CSO) != len(shapeA)*len(shapeB) {
		return fmt.Errorf("Mikowski: number of points in cso does not match input shapes.")
	}
	vertInd := 0
	for i := 0; i < len(shapeA); i++ {
		for j := 0; j < len(shapeB); j++ {
			CSO[vertInd] = shapeA[i].Sub(shapeB[j])
			vertInd++
		}
	}
	return nil
}

func transformShape(input, output []mgl32.Vec3, translation mgl32.Vec3, rotationAngle float32) {
	for i := 0; i < len(input); i++ {
		output[i] = (mgl32.Rotate3DY(rotationAngle)).Mul3x1(input[i]).Add(translation)
	}
}

func clamp(a, i, b float32) float32 {
	if i < a {
		return a
	} else if i > b {
		return b
	}
	return i
}
func clampInt(a, i, b int) int {
	if i < a {
		return a
	} else if i > b {
		return b
	}
	return i
}

//Input function
func handleInput(window *glfw.Window, world *gizmo, frameTimer *frameTimer, player *player, camera *camera, colliderPosition, lightPosition *mgl32.Vec3, colliderRotation, speed, head, height *float32, editBone *int32, collisionSteps *int, pressedN *bool, envShader *uint32, firstShader, secondShader, thirdShader uint32) {
	var maxTiltAngle float32 = 0.25
	var lightSpeed float32 = 3
	var maxSpeed float32 = 10
	var initialJumpSpeed float32 = 5
	var angularVelocity float32 = 15
	var acc float32 = 30
	var deacc float32 = 15

	//Pressing Esc to exit
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}
	//Shader switching
	if window.GetKey(glfw.Key1) == glfw.Press {
		*envShader = firstShader
	}
	if window.GetKey(glfw.Key2) == glfw.Press {
		*envShader = secondShader
	}
	if window.GetKey(glfw.Key3) == glfw.Press {
		*envShader = thirdShader
	}
	pressedShift := false
	//Time slowdown
	if window.GetKey(glfw.KeyLeftShift) == glfw.Press {
		frameTimer.deltaTime /= 8
		pressedShift = true
	}
	deltaTime := frameTimer.deltaTime
	if !player.InAir && window.GetKey(glfw.KeySpace) == glfw.Press {
		player.InAir = true
		player.Velocity = player.Velocity.Add(world.yAxis.Mul(initialJumpSpeed))
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
	//LIGHT/COLLIDER MOTION
	if window.GetKey(glfw.KeyUp) == glfw.Press {
		*lightPosition = lightPosition.Add(camera.Forward.Mul(lightSpeed * deltaTime))
		*colliderPosition = colliderPosition.Add(camera.Forward.Mul(lightSpeed * deltaTime))
	}
	if window.GetKey(glfw.KeyDown) == glfw.Press {
		*lightPosition = lightPosition.Add(camera.Forward.Mul(-lightSpeed * deltaTime))
		*colliderPosition = colliderPosition.Add(camera.Forward.Mul(-lightSpeed * deltaTime))
	}
	if window.GetKey(glfw.KeyLeft) == glfw.Press {
		*lightPosition = lightPosition.Add(camera.Left.Mul(lightSpeed * deltaTime))
		*colliderPosition = colliderPosition.Add(camera.Left.Mul(lightSpeed * deltaTime))
	}
	if window.GetKey(glfw.KeyRight) == glfw.Press {
		*lightPosition = lightPosition.Add(camera.Left.Mul(-lightSpeed * deltaTime))
		*colliderPosition = colliderPosition.Add(camera.Left.Mul(-lightSpeed * deltaTime))
	}
	//COLLISION VISUALIZATION
	if window.GetKey(glfw.KeyJ) == glfw.Press {
		*colliderRotation += lightSpeed * deltaTime
	}
	if window.GetKey(glfw.KeyK) == glfw.Press {
		*colliderRotation -= lightSpeed * deltaTime
	}
	//POSE EDITING/COLLISION STEPING
	if window.GetKey(glfw.KeyN) == glfw.Press {
		*pressedN = true
	}
	if *pressedN && window.GetKey(glfw.KeyN) == glfw.Release {
		//*editBone++;
		*collisionSteps++
		if pressedShift {
			*collisionSteps -= 2
		}
		if *editBone >= 0 {
			*editBone %= 15
		}
		*collisionSteps = clampInt(0, *collisionSteps, 20)
		*pressedN = false
	}
	//ANIMATION BLENDING
	if window.GetKey(glfw.KeyL) == glfw.Press {
		player.LookAtLight = true
	} else if window.GetKey(glfw.KeyC) == glfw.Press {
		player.LookAtLight = false
	}
	lookDir := camera.Forward
	if player.LookAtLight {
		lookDir = lightPosition.Sub(player.Position)
	}
	lookDir[1] = 0
	lookDir = lookDir.Normalize()
	t := lookDir.Dot(player.Dir)
	t += -1
	t = clamp(-1, t, 0)
	if lookDir.Cross(player.Dir)[1] < 0 {
		t *= -1
	}
	*head += 8 * (t - (*head)) * deltaTime
	*speed = float32(math.Sqrt(float64((player.Velocity[0]*player.Velocity[0])+(player.Velocity[2]*player.Velocity[2])))) / maxSpeed
	*head = clamp(-1, *head, 1)
	*speed = clamp(0, *speed, 1)
	//jump blend
	maxFallTime := initialJumpSpeed / 9.81
	maxHeight := maxFallTime * maxFallTime * 9.81 / 2
	if player.Velocity[1] > 0 {
		*height = clamp(0.2, 2*player.Position[1]/maxHeight, 1)
	} else {
		*height = player.Position[1] / maxHeight
	}
	//RESET BUTTON
	if window.GetKey(glfw.KeyR) == glfw.Press {
		*lightPosition = mgl32.Vec3{}
		player.Position = mgl32.Vec3{}
		player.Velocity = mgl32.Vec3{}
		player.Dir = mgl32.Vec3{0, 0, 1}
		player.Up = mgl32.Vec3{0, 1, 0}
		player.TiltAngle = 0
		*head = 0
	}
}
