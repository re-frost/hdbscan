package edgeDetection

type Data struct {
	img *ImageCV

	depthFile []float64

	coordXYZ []Vector3
	coordUV  []Vector2

	points [][]float64

	indexXYZ [][3]int
	indexUV  [][3]int

	clusterNumber []int
	mean          []Vector3
}

type Vector3 struct {
	X, Y, Z float64
}

type Vector2 struct {
	X, Y float64
}

func Detection(argument string) [][]float64 {

	verticesXYZ, verticesUV, xyzs, uvs := ReadObjFile(argument + "/mesh.obj")
	depth := ReadDepthFile(argument + "/depth.txt")
	d := &Data{
		img:       ImageControler(argument + "/texture.png"),
		depthFile: depth,
		coordXYZ:  verticesXYZ,
		coordUV:   verticesUV,
		indexXYZ:  xyzs,
		indexUV:   uvs,
	}
	whitePoints := d.findCorrespondendingPoints()
	whitePoints.showImg("whitePoints")

	return d.points
}
