package main

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"time"

	"training/engine/anim"
)

type gizmo struct {
	xAxis mgl32.Vec3
	yAxis mgl32.Vec3
	zAxis mgl32.Vec3
}

type player struct {
	Position     mgl32.Vec3
	Velocity     mgl32.Vec3
	AccDirection mgl32.Vec3
	Up           mgl32.Vec3
	TiltAxis     mgl32.Vec3
	Dir          mgl32.Vec3
	TiltAngle    float32
	DestAngle    float32
	Angle        float32
	InAir        bool
	LookAtLight  bool
}

type camera struct {
	Offset           mgl32.Vec3
	Position         mgl32.Vec3
	Target           *mgl32.Vec3
	Forward          mgl32.Vec3
	Left             mgl32.Vec3
	ViewMatrix       mgl32.Mat4
	ProjectionMatrix mgl32.Mat4
	VPMatrix         mgl32.Mat4
	angleX           float32
	angleY           float32
	mouseX           float32
	mouseY           float32
}

type frameTimer struct {
	gameLoopStart    float32
	frameStart       float32
	deltaTime        float32
	desiredFrameTime float32
	longestFrame     float32
	shortestFrame    float32
	frames           int
	seconds          int
	currentFps       int
	isSecondMark     bool
}

func newCamera(offset mgl32.Vec3, target *mgl32.Vec3, world *gizmo) camera {
	cam := camera{Offset: offset, Target: target}
	cam.Left = world.yAxis.Cross(cam.Offset.Mul(-1)).Normalize()
	cam.Forward = cam.Left.Cross(world.yAxis).Normalize()
	cam.ViewMatrix = mgl32.LookAtV(cam.Target.Add(cam.Offset), *cam.Target, world.yAxis)
	cam.ProjectionMatrix = mgl32.Perspective(3.14159/4, 1.6, 0.1, 1000)
	cam.VPMatrix = cam.ProjectionMatrix.Mul4(cam.ViewMatrix)
	return cam
}

func (c *camera) Update(x, y float64, speedPercent float32, playerDir *mgl32.Vec3) {
	c.angleX -= (float32(x) - c.mouseX) / screenWidth
	c.angleY -= (float32(y) - c.mouseY) / screenHeight
	c.mouseX = float32(x)
	c.mouseY = float32(y)
	c.angleY = clamp(-1, c.angleY, 1)
	rotX := mgl32.Rotate3DX(c.angleY)
	rotY := mgl32.Rotate3DY(c.angleX)
	c.Position = rotY.Mul3x1(rotX.Mul3x1(c.Offset))
	c.Forward = rotY.Mul3x1(mgl32.Vec3{0, 0, -1})
	c.Left = rotY.Mul3x1(mgl32.Vec3{-1, 0, 0})
	//Field of view for speed effect
	towardsSpeed := clamp(0, playerDir.Dot(c.Forward), 1)
	c.ProjectionMatrix = mgl32.Perspective(3.14159/(4-speedPercent*towardsSpeed*0.4), 1.6, 0.1, 1000)
	c.ViewMatrix = mgl32.LookAtV(c.Target.Add(c.Position), c.Target.Add(mgl32.Vec3{0, 1, 0}), mgl32.Vec3{0, 1, 0})
	c.VPMatrix = c.ProjectionMatrix.Mul4(c.ViewMatrix)
}

func (f *frameTimer) Update() {
	currentTime := float32(glfw.GetTime())
	f.deltaTime = currentTime - f.frameStart
	if excessFrameTime := f.desiredFrameTime - f.deltaTime; excessFrameTime > 0 {
		time.Sleep(time.Duration(int64(excessFrameTime * 1e9)))
		f.deltaTime = f.desiredFrameTime
	}
	f.frames++
	if f.seconds < int(currentTime-f.gameLoopStart) {
		f.seconds = int(currentTime)
		f.isSecondMark = true
		f.currentFps = f.frames
		f.frames = 0
	} else {
		f.isSecondMark = false
	}
	f.frameStart = float32(glfw.GetTime())
}

var keyframe00 = anim.Keyframe{SampleTime: 0, Transforms: []anim.Transform{
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
var keyframe01 = anim.Keyframe{SampleTime: 0.1, Transforms: []anim.Transform{
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{0, 0.05, 0}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, -10}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, -30}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 10}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{-20, 0, 30}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{10, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, -5, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{-20, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{15, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{-10, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{15, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{15, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{-30, 0, 0}}}}
var keyframe02 = anim.Keyframe{SampleTime: 0.2, Transforms: []anim.Transform{
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
var keyframe03 = anim.Keyframe{SampleTime: 0.3, Transforms: []anim.Transform{
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{0, 0.05, 0}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, -10}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{-20, 0, -30}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{10, 0, 0}},
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
var keyframe04 = anim.Keyframe{SampleTime: 0.4, Transforms: []anim.Transform{
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
var keyframe10 = anim.Keyframe{SampleTime: 0, Transforms: []anim.Transform{
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
var keyframe11 = anim.Keyframe{SampleTime: 0.1, Transforms: []anim.Transform{
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{0, 0.025, 0}, [3]float32{0, 0, 0}},
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
var keyframe12 = anim.Keyframe{SampleTime: 0.2, Transforms: []anim.Transform{
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
var keyframe13 = anim.Keyframe{SampleTime: 0.3, Transforms: []anim.Transform{
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{0, 0.025, 0}, [3]float32{0, 0, 0}},
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
var keyframe14 = anim.Keyframe{SampleTime: 0.4, Transforms: []anim.Transform{
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
var keyframe20 = anim.Keyframe{SampleTime: -1, Transforms: []anim.Transform{
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, -60, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, -45, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}}}}
var keyframe21 = anim.Keyframe{SampleTime: 1, Transforms: []anim.Transform{
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 60, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 45, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}}}}
var keyframe30 = anim.Keyframe{SampleTime: 0, Transforms: []anim.Transform{
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, -10}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, -40}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{10, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 10}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 40}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{10, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{0, -0.005, 0}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}}}}
var keyframe31 = anim.Keyframe{SampleTime: 1, Transforms: []anim.Transform{
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, -10}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, -40}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{10, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 10}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 40}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{10, 0, 0}},
	anim.Transform{[3]float32{1.005, 1.005, 1.005}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}}}}
var keyframe32 = anim.Keyframe{SampleTime: 2, Transforms: []anim.Transform{
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, -10}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, -40}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{10, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 10}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 40}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{10, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{0, -0.005, 0}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}}}}
var keyframe40 = anim.Keyframe{SampleTime: 0, Transforms: []anim.Transform{
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, -10}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, -40}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{10, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 10}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 40}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{10, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}}}}
var keyframe41 = anim.Keyframe{SampleTime: 1, Transforms: []anim.Transform{
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 20}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, -20}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{100, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{-80, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{-20, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}}}}
