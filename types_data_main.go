package main

import (
	"time"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/glfw/v3.2/glfw"
	"math"

	"training/engine/anim"
)

type gizmo struct{
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
}

type camera struct {
	Position         mgl32.Vec3
	Direction        mgl32.Vec3
	Forward          mgl32.Vec3
	Left             mgl32.Vec3
	ViewMatrix       mgl32.Mat4
	ProjectionMatrix mgl32.Mat4
	VPMatrix	 mgl32.Mat4
}

type frameTimer struct {
	gameLoopStart    float32
	frameStart       float32
	deltaTime        float32
	desiredFrameTime float32
	frames           int
	seconds          int
	currentFps	 int
	isSecondMark     bool
}

func newCamera(position, direction mgl32.Vec3, world *gizmo) camera {
	cam := camera{Position: position, Direction: direction}
	globalUp := mgl32.Vec3{0, 1, 0}
	cam.Left = world.yAxis.Cross(cam.Direction).Normalize()
	cam.Forward = cam.Left.Cross(world.yAxis).Normalize()
	cam.ViewMatrix = mgl32.LookAtV(cam.Position, cam.Position.Add(cam.Direction), globalUp)
	cam.ProjectionMatrix = mgl32.Perspective(math.Pi/4, 1.6, 0.1, 100.0)
	cam.VPMatrix = cam.ProjectionMatrix.Mul4(cam.ViewMatrix)
	return cam
}


func (f *frameTimer) Update() {
	f.frames++
	currentTime := float32(glfw.GetTime())
	f.deltaTime = currentTime - f.frameStart
	if excessFrameTime := f.desiredFrameTime - f.deltaTime; excessFrameTime  > 0{
		time.Sleep(time.Duration(int64(excessFrameTime * 1e9)))
		f.deltaTime = f.desiredFrameTime
	}
	if f.seconds < int(currentTime - f.gameLoopStart){
		f.seconds = int(currentTime)
		f.isSecondMark = true
		f.currentFps = f.frames
		f.frames = 0
	} else{
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
var keyframe01 = anim.Keyframe{SampleTime: 0.15, Transforms: []anim.Transform{
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
var keyframe02 = anim.Keyframe{SampleTime: 0.3, Transforms: []anim.Transform{
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
var keyframe03 = anim.Keyframe{SampleTime: 0.45, Transforms: []anim.Transform{
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
var keyframe04 = anim.Keyframe{SampleTime: 0.6, Transforms: []anim.Transform{
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
var keyframe11 = anim.Keyframe{SampleTime: 0.15, Transforms: []anim.Transform{
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
var keyframe12 = anim.Keyframe{SampleTime: 0.3, Transforms: []anim.Transform{
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
var keyframe13 = anim.Keyframe{SampleTime: 0.45, Transforms: []anim.Transform{
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
var keyframe14 = anim.Keyframe{SampleTime: 0.6, Transforms: []anim.Transform{
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
var keyframe20 = anim.Keyframe{SampleTime: 0, Transforms: []anim.Transform{
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, -50, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, -40, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}}}}
var keyframe21 = anim.Keyframe{SampleTime: 1, Transforms: []anim.Transform{
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 50, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 40, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	anim.Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}}}}
