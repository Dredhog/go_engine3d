package collision

import "github.com/go-gl/mathgl/mgl32"

func BGJK(shapeA, shapeB []mgl32.Vec3) bool {
	var simplexStack [4]mgl32.Vec3
	simplex := simplexStack[:]
	simplex[0] = shapeA[0].Sub(shapeB[0])
	order := 0

	dir := simplex[0].Mul(-1)
	for {
		a := support(shapeA, shapeB, dir)
		if a.Dot(dir) < 0 {
			return false
		}
		order++
		simplex[order] = a
		if order == 3 && DoSimplex3(simplex, &dir, &order) {
			return true
		} else if order == 2 {
			doSimplex2(simplex, &dir, &order)
		} else if order == 1 {
			doSimplex1(simplex, &dir)
		}
	}
	return false
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
		*order--
		simplex[1] = simplex[2] //assign a in place of b
		*dir = ac.Cross(ao.Cross(ac))
	} else if (ab.Cross(abc)).Dot(ao) > 0 { //ab edge
		*order--
		simplex[0] = simplex[1] //assign b in place of c
		simplex[1] = simplex[2] //make a the second vertex
		*dir = ab.Cross(ao.Cross(ab))
	} else {
		if abc.Dot(ao) > 0 { //towards triangle normal
			*dir = abc.Cross(ao).Cross(abc)
		} else {
			*dir = abc.Cross(ao).Cross(abc).Mul(-1)
			//reverse triangle winding
			ao = simplex[0]
			simplex[0] = simplex[1]
			simplex[1] = ao
		}
	}
}

func DoSimplex3(simplex []mgl32.Vec3, dir *mgl32.Vec3, order *int) bool {
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
				return false
			}
			//acd region
			return false
		} else if ab.Cross(abc).Dot(ao) > 0 { //clockwise of abc region
			if adb.Cross(ab).Dot(ao) > 0 { //ab region
				return false
			}
			//adb region
			return false
		}
		//abc region
		return false
	} else if adb.Dot(ao) > 0 {
		if adb.Cross(ab).Dot(ao) > 0 { //region ab
		} else if ad.Cross(adb).Dot(ao) > 0 { //clockwise of adb
			if acd.Cross(ad).Dot(ao) > 0 { //ad region
				return false
			}
			//acd region
			return false
		}
		//adb region
		return false
	} else if acd.Dot(ao) > 0 { //acd face
		if acd.Cross(ad).Dot(ao) > 0 { //ad region
			return false
		} else if ac.Cross(acd).Dot(ao) > 0 { //ac region
			return false
		}
		//acd region
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
