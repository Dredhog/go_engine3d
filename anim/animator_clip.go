package anim

import "fmt"

func (c *Clip) SetDuration() {
	c.Duration = c.Keyframes[len(c.Keyframes)-1].SampleTime
}

func (c *Clip) SampleLooping(time float32, result *Keyframe) error {
	time -= c.Duration * float32(int(time/c.Duration))
	count := len(c.Keyframes)
	for i := 0; i < count-1; i++ {
		if c.Keyframes[i+1].SampleTime > time {
			if err := LerpKeyframe(&c.Keyframes[i], &c.Keyframes[i+1], (time-c.Keyframes[i].SampleTime)/(c.Keyframes[i+1].SampleTime-c.Keyframes[i].SampleTime), result); err != nil {
				return err
			}
			//*result = Keyframe{}
			//*result = c.Keyframes[i]
			return nil
		}
	}
	return fmt.Errorf("anim: failed to sample keyframe, sample time : %v", time)
}

func (a *Animator) BlendClips(time, t float32) error {
	if err := a.CurrentClip.SampleLooping(time, &a.CurrentKeyframe); err != nil {
		return fmt.Errorf("blending: %v", err)
	}
	if err := a.UpcomingClip.SampleLooping(time, &a.UpcomingKeyframe); err != nil {
		return fmt.Errorf("blending: %v", err)
	}
	if err := LerpKeyframe(&a.CurrentKeyframe, &a.UpcomingKeyframe, t, &a.ResultKeyframe); err != nil {
		return err
	}
	return nil
}
