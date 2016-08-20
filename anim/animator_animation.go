package anim

import "fmt"

func (a *Animation) SetTotalTicks() {
	a.TotalTicks = a.Keyframes[len(a.Keyframes)-1].Ticks
}

func (a *Animation) SampleLooping(ticks float32, result *Keyframe) error {
	ticks -= a.TotalTicks * float32(int(ticks/a.TotalTicks))
	count := len(a.Keyframes)
	for i := 0; i < count-1; i++ {
		if a.Keyframes[i+1].Ticks > ticks {
			if err := LerpKeyframe(&a.Keyframes[i], &a.Keyframes[i+1], (ticks-a.Keyframes[i].Ticks)/(a.Keyframes[i+1].Ticks-a.Keyframes[i].Ticks), result); err != nil{
				return err
			}
			return nil
		}
	}
	return fmt.Errorf("anim: failed to sample keyframe, sample ticks : %v", ticks)
}

func (a *Animator) BlendAnimations(ticks, t float32) error{
	if err := a.CurrentAnimation.SampleLooping(ticks, &a.CurrentKeyframe); err != nil{
		return fmt.Errorf("blending: %v", err)
	}
	if err := a.UpcomingAnimation.SampleLooping(ticks, &a.UpcomingKeyframe); err != nil{
		return fmt.Errorf("blending: %v", err)
	}
	if err := LerpKeyframe(&a.CurrentKeyframe, &a.UpcomingKeyframe, t, &a.ResultKeyframe); err != nil{
		return err
	}
	return nil
}
