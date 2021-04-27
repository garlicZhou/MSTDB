package mst

import (
	"fmt"
	"testing"
	"time"
)

func TestUpdateIndex(T *testing.T) {
     index := CreateIndex()
	file1 := File{
		Name:    "xiyouji",
		Keys:    []string{"monkey", "pig", "monk", "horse"},
		Times:   []int{2, 4, 3, 5},
		WordSum: 0,
	    }
	/*file2 := File{
		Name:    "xiyouji1",
		Keys:    []string{"monkey", "pig", "monk", "horse"},
		Times:   []int{2, 4, 3, 5},
		WordSum: 0,
	}
	file := File{
		Name:    "xiyouji",
		Keys:    []string{"monkey", "pig", "monk", "horse"},
		Times:   []int{2, 4, 3, 5},
		WordSum: 0,
	}*/
	/*file2 := File{
			Name:    "sanguo",
			Keys:    []string{"horse", "man", "woman"},
			Times:   []int{22, 4, 8},
			WordSum: 0,
		}
		file3 := File{
			Name:    "hlm",
			Keys:    []string{"people", "top", "man", "woman", "friut"},
			Times:   []int{12, 4, 33, 58, 9},
			WordSum: 0,
		}
		file4 := File{
			Name:    "shuihuzhuan",
			Keys:    []string{"monk", "meat", "baijiu", "kill"},
			Times:   []int{12, 74, 63, 115},
			WordSum: 0,
		}
	file5 := File{
			Name:    "labixiaoxi",
			Keys:    []string{"child", "chicken", "monk", "horse"},
			Times:   []int{101, 4, 3, 5},
			WordSum: 0,
		}*/
		index.UpdateIndex(file1,1)

	//index.UpdateIndex(file2,2)
	/*index.UpdateIndex(file3,3)
	index.UpdateIndex(file4,4)
	index.UpdateIndex(file5,5)*/
	fmt.Println(index.list)
	query := File{
		Name:    "wataxi",
		Keys:    []string{"monkey","pig"},
		Times:   []int{2,3},
		WordSum: 0,
	}
	start := time.Now()
	fmt.Println("top-k by taat:", index.TopKbyTAAT(query, 1))
	for i := 0;i < 6000;i++{
		index.TopKbyTAAT(query, 1)
	}
	elapsed := time.Since(start)
	fmt.Println("TAAT running time：", elapsed)

	start2 := time.Now()
	fmt.Println("top-k by daat:", index.TopKbyDAAT(query, 1))
	for i := 0;i < 6000;i++{
		index.TopKbyDAAT(query, 1)
	}
	elapsed2 := time.Since(start2)
	fmt.Println("DAAT running time：", elapsed2)

	start3 := time.Now()
	fmt.Println("top-k by baat:", index.TopKbyBAAT(query, 1,2))
	for i := 0;i <6000;i++{
		index.TopKbyBAAT(query, 1,2)
	}
	elapsed3 := time.Since(start3)
	fmt.Println("BAAT running time：", elapsed3)
}