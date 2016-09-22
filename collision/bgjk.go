package collision

import (
	//"fmt"
	"github.com/go-gl/mathgl/mgl32"
)

func BGJK(shapeA, shapeB []mgl32.Vec3, stepCount int) ([]mgl32.Vec3, int, bool) {
	_ = stepCount
	var simplexStack [4]mgl32.Vec3
	simplex := simplexStack[:]
	simplex[0] = shapeA[0].Sub(shapeB[0])
	order := 0

	dir := simplex[0].Mul(-1)
	for i := 0; ; i++ {
		if i == stepCount{
			break
		}
		//fmt.Printf("-------------Step %v---------------\n", i)
		a := support(shapeA, shapeB, dir)
		if a.Dot(dir) < 0 {
			return simplex, order, false
		}
		order++
		simplex[order] = a
		i++
		if i == stepCount{
			break
		}
		if order == 3 && doSimplex3(simplex, &dir, &order) {
			return simplex, order, true
		} else if order == 2 {
			doSimplex2(simplex, &dir, &order)
		} else if order == 1 {
			doSimplex1(simplex, &dir)
		}
	}
	return simplex, order, false
}

func doSimplex1(simplex []mgl32.Vec3, dir *mgl32.Vec3) {
	ab := simplex[0].Sub(simplex[1])
	*dir = ab.Cross(simplex[1].Mul(-1)).Cross(ab)
}

func doSimplex2(simplex []mgl32.Vec3, dir *mgl32.Vec3, order *int) {
	ao := simplex[2].Mul(-1)
	ab := simplex[1].Sub(simplex[2])
	ac := simplex[0].Sub(simplex[2])
	abc := ab.Cross(ac)
	if abc.Cross(ac).Dot(ao) > 0 { //ac edge
		*order = 1
		simplex[0] = simplex[1]
		simplex[1] = simplex[2]
		*dir = ac.Cross(ao).Cross(ac)
	} else if ab.Cross(abc).Dot(ao) > 0 { //ab edge
		*order = 1
		simplex[0] = simplex[2]
		simplex[1] = simplex[1]
		*dir = ab.Cross(ao).Cross(ab)
	} else {
		if abc.Dot(ao) > 0 { //towards triangle normal
			*dir = abc
			return
		} else {
			*dir = abc.Mul(-1)
			//reverse triangle winding
			temp := simplex[0]
			simplex[0] = simplex[2]
			simplex[2] = temp
			//fmt.Println("reverse order 2")
		}
	}
}

func doSimplex3(simplex []mgl32.Vec3, dir *mgl32.Vec3, order *int) bool {
	ao := simplex[3].Mul(-1)
	ab := simplex[2].Sub(simplex[3])
	ac := simplex[1].Sub(simplex[3])
	ad := simplex[0].Sub(simplex[3])
	abc := ab.Cross(ac)
	adb := ad.Cross(ab)
	acd := ac.Cross(ad)
	if abc.Dot(ao) > 0 {
		if abc.Cross(ac).Dot(ao) > 0 { //counter-clockwise of abc region
			if ac.Cross(acd).Dot(ao) > 0 { //ac region
				simplex[0] = simplex[3] //assign a to [0]
				simplex[1] = simplex[1] //assign c to [1]
				*dir = ac.Cross(ao).Cross(ac)
				*order = 1
				//fmt.Println("ac 3")
				return false
			}
			//acd region
			simplex[0] = simplex[0]	//assign d to [0]
			simplex[1] = simplex[1] //assign c to [1]
			simplex[2] = simplex[3] //assign a to [2]
			*dir = acd
			*order = 2
			//fmt.Println("acd 3")
			return false
		} else if ab.Cross(abc).Dot(ao) > 0 { //clockwise of abc region
			if adb.Cross(ab).Dot(ao) > 0 { //ab region
				simplex[0] = simplex[3] //assign a to [0]
				simplex[1] = simplex[2] //assign b to [1]
				*dir = ab.Cross(ao).Cross(ab)
				*order = 1
				//fmt.Println("ab 3")
				return false
			}
			//adb region
			simplex[0] = simplex[0]	//assign d to [0]
			simplex[1] = simplex[3] //assign b to [1]
			simplex[2] = simplex[2] //assign a to [2]
			*dir = adb
			*order = 2
			//fmt.Println("adb 3")
			return false
		}
		//abc region
		simplex[0] = simplex[1]	//assign c to [0]
		simplex[1] = simplex[2] //assign b to [1]
		simplex[2] = simplex[3] //assign a to [2]
		*dir = abc
		*order = 2
		//fmt.Println("abc 3")
		return false
	} else if adb.Dot(ao) > 0 {
		if adb.Cross(ab).Dot(ao) > 0 { //region ab
			simplex[0] = simplex[3] //assign a to [0]
			simplex[1] = simplex[2] //assign b to [1]
			*dir = ab.Cross(ao).Cross(ab)
			*order = 1
			//fmt.Println("ab 3")
			return false
		} else if ad.Cross(adb).Dot(ao) > 0 { //clockwise of adb
			if acd.Cross(ad).Dot(ao) > 0 { //ad region
				simplex[0] = simplex[0] //assign d to [0]
				simplex[1] = simplex[3] //assign a to [1]
				*dir = ad.Cross(ao).Cross(ad)
				*order = 1
				//fmt.Println("ad 3")
				return false
			}
			//acd region
			simplex[0] = simplex[0]	//assign d to [0]
			simplex[1] = simplex[1] //assign c to [1]
			simplex[2] = simplex[3] //assign a to [2]
			*dir = acd
			*order = 2
			//fmt.Println("acd 3")
			return false
		}
		//adb region
		simplex[0] = simplex[0]	//assign d to [0]
		simplex[1] = simplex[3] //assign a to [1]
		simplex[2] = simplex[2] //assign b to [2]
		*dir = adb
		*order = 2
		//fmt.Println("adb 3")
		return false
	} else if acd.Dot(ao) > 0 { //acd face
		if acd.Cross(ad).Dot(ao) > 0 { //ad region
			simplex[0] = simplex[3] //assign a to [0]
			simplex[1] = simplex[0] //assign d to [1]
			*dir = ad.Cross(ao).Cross(ad)
			*order = 1
			//fmt.Println("ad 3")
			return false
		} else if ac.Cross(acd).Dot(ao) > 0 { //ac region
			simplex[0] = simplex[3] //assign a to [0]
			simplex[1] = simplex[1] //assign c to [1]
			*dir = ac.Cross(ao).Cross(ac)
			*order = 1
			//fmt.Println("ac 3")
			return false
		}
		//acd region
		simplex[0] = simplex[0]	//assign d to [0]
		simplex[1] = simplex[1] //assign c to [1]
		simplex[2] = simplex[3] //assign a to [2]
		*dir = acd
		*order = 2
		//fmt.Println("acd 3")
		return false
	}
	return true
}

func support(shapeA, shapeB []mgl32.Vec3, dir mgl32.Vec3) mgl32.Vec3 {
	maxA, minB := shapeA[0].Dot(dir), shapeB[0].Dot(dir)
	indA, indB := 0, 0
	for i := 1; i < len(shapeA); i++ {
		if temp := shapeA[i].Dot(dir); temp > maxA {
			maxA = temp
			indA = i
		}
	}
	for i := 1; i < len(shapeB); i++ {
		if temp := shapeB[i].Dot(dir); temp < minB {
			minB = temp
			indB = i
		}
	}
	return shapeA[indA].Sub(shapeB[indB])
}
