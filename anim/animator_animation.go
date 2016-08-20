package anim

import "fmt"

func (a *Animation) SetTotalTicks(){
	a.TotalTicks = a.Keyframes[len(a.Keyframes)-1].Ticks
}

func (a *Animation) LoopedLinearSample(ticks float32) (*Keyframe, error){
	ticks -= a.TotalTicks*float32(int(ticks/a.TotalTicks))
	count := len(a.Keyframes)
	for i := 0; i < count-1; i++{
		if a.Keyframes[i+1].Ticks > ticks{
			return LerpKeyframe(&a.Keyframes[i], &a.Keyframes[i+1], (ticks - a.Keyframes[i].Ticks)/(a.Keyframes[i+1].Ticks - a.Keyframes[i].Ticks)), nil
		}
	}
	return nil, fmt.Errorf("anim: failed to sample keyframe, sample ticks : %v", ticks)
}

func (a *Animator) StartAnimation(animationIndex int){
}

func (a *Animator) StepAnimations(deltaTicks int) *Keyframe {
	return nil
}
