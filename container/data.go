package container

var (
	SearchMap = map[int]string{
		5: "wu",
		6: "lu",
		11: "shiyi",
		8: "ba",
		13: "shisan",
	}
	SearchMapLength = 5
	SearchMapNotKeys = []int{9, 4, 1, -1}
)

func init() {
	SearchMap[11] = "shiyi2"
}
