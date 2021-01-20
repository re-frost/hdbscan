package kmean

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type Vector3 struct {
	X, Y, Z float64
}

type DistanceFunction func(first, second Vector3) (float64, error)

func Kmeans(distFunc DistanceFunction, data []Vector3, k, threshold int) ([]int, []Vector3, error) {

	seeds := seed(data, k, distFunc)

	for i, o := range seeds {
		fmt.Println("Cluster: ", i, "Mean Init: ", o)
	}

	clusterNumber, mean, err := kmeans(data, seeds, distFunc, threshold)

	return clusterNumber, mean, err
}

// Instead of initializing randomly the seeds, make a sound decision of initializing
func seed(data []Vector3, k int, distanceFunction DistanceFunction) []Vector3 {

	time := time.Now().Nanosecond()
	rand.Seed(int64(time))

	s := make([]Vector3, k)
	s[0] = data[rand.Intn(len(data))]
	d2 := make([]float64, len(data))
	for i := 1; i < k; i++ {
		var sum float64
		for j, p := range data {
			_, dMin := near(p, s[:i], distanceFunction)
			d2[j] = dMin * dMin
			sum += d2[j]
		}
		target := rand.Float64() * float64(sum)
		j := 0
		for sum = d2[0]; sum < float64(target); sum += d2[j] {
			j++
		}
		s[i] = data[i]
	}
	return s
}

func near(p Vector3, mean []Vector3, distanceFunction DistanceFunction) (int, float64) {
	indexOfCluster := 0
	minSquaredDistance, _ := distanceFunction(p, mean[0])
	for i := 1; i < len(mean); i++ {
		squaredDistance, _ := distanceFunction(p, mean[i])
		if squaredDistance < minSquaredDistance {
			minSquaredDistance = squaredDistance
			indexOfCluster = i
		}
	}
	return indexOfCluster, math.Sqrt(minSquaredDistance)
}

// K-Means Algorithm
func kmeans(data []Vector3, mean []Vector3, distanceFunction DistanceFunction, threshold int) ([]int, []Vector3, error) {

	counter := 0
	clusterNumber := make([]int, len(data))
	for i, vec := range data {
		closestCluster, _ := near(vec, mean, distanceFunction)
		clusterNumber[i] = closestCluster
	}

	mLen := make([]int, len(mean))
	// fmt.Println("len(mLen)", len(mLen))
	for {
		for j := range mean {
			mean[j] = Vector3{0, 0, 0}
			mLen[j] = 0
		}
		for i, p := range data {
			mean[clusterNumber[i]] = mean[clusterNumber[i]].Add3(p)
			mLen[clusterNumber[i]]++
		}

		for ii := range mean {
			div := 1. / float64(mLen[ii])
			mean[ii] = mean[ii].MultiplyByScalar3(div)
		}

		var changes int
		for ii, p := range data {
			if closestCluster, _ := near(p, mean, distanceFunction); closestCluster != clusterNumber[ii] {
				changes++
				clusterNumber[ii] = closestCluster
			}
		}
		counter++

		// fmt.Println("Mean Vectors: ", mean)
		if changes == 0 || counter > threshold {
			return clusterNumber, mean, nil
		}
	}
}

func (a Vector3) MultiplyByScalar3(s float64) Vector3 {
	return Vector3{
		X: a.X * s,
		Y: a.Y * s,
		Z: a.Z * s,
	}
}

func (a Vector3) Add3(b Vector3) Vector3 {
	return Vector3{
		X: a.X + b.X,
		Y: a.Y + b.Y,
		Z: a.Z + b.Z,
	}
}
