package mst

import (
	"fmt"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"math"
	"strings"
)

type key_file struct {
	Key       string
	File_list []string
	Score     []float64
	keySum    int
}

type Inverted_list struct {
	k         int
	list      []key_file
	docSum    int
	keyToKey_file map[string]int
	keyToTimes   map[string]int
	file_list    []File_list
	nameToFileList map[string]File_list
	blockToFileLists  map[uint][]File_list
	blockAvgScore  map[uint]map[string]block_key
 	Db   ethdb.Database
}

type block_key struct {
	time int
	score float64
}

type key_file_list []key_file

type score_rank struct {
	docName       []string
	//docScore      []float64
	nameToScore   map[string]float64
}

var kfToScore map[string]float64

func (sr *score_rank) sort() {
	for i := 0; i < len(sr.docName); i++ {
		for j := 1; j < len(sr.docName) - i; j++ {
			if sr.nameToScore[sr.docName[j]] < sr.nameToScore[sr.docName[j - 1]] {
				sr.docName[j], sr.docName[j - 1] = sr.docName[j - 1], sr.docName[j]
			}
		}
	}
}

/*func (index *Inverted_list) swap(i, j int) *Inverted_list {
	key_file1 := index.list[i]
	index.list[i] = index.list[j]
	index.list[j] = key_file1
	return index
}*/

func (index *Inverted_list) insert(file key_file) {
	flag := 0
	var p int
	if len(index.list) == 0 {
		index.list = append(index.list, file)
		index.keyToTimes = make(map[string]int)
		index.keyToTimes[file.Key] ++
		index.keyToKey_file = make(map[string]int)
		index.keyToKey_file[file.Key] = 0
	} else {
		for i := 0; i < len(index.list); i++ {
			if strings.Compare(file.Key, index.list[i].Key) == 0 {
				flag = 1
				index.list[i].File_list = append(index.list[i].File_list, file.File_list[0])
				index.list[i].Score = append(index.list[i].Score, file.Score[0])
				m := len(index.list[i].File_list)
				for j := m; j > 1; j-- {
					if index.list[i].Score[j - 1] > index.list[i].Score[j - 2]{
						index.list[i].File_list[j - 1], index.list[i].File_list[j - 2] = index.list[i].File_list[j - 2], index.list[i].File_list[j - 1]
						index.list[i].Score[j - 1], index.list[i].Score[j - 2] = index.list[i].Score[j - 2], index.list[i].Score[j - 1]
					} else {
						break
					}
				}
				p = i
			}
		}
		if flag == 0 {
			index.list = append(index.list, file)
			index.keyToTimes[file.Key] ++
			index.keyToKey_file[file.Key] = p
			p = len(index.list) - 1
		}
	}
	index.k = len(index.list)
//	list_data,_ := rlp.EncodeToBytes(index.list[p])
//	_ = index.Db.Put([]byte(file.Key), list_data)
}

func (index *Inverted_list) list_sort() *Inverted_list {
	for j := index.k; j > 0; j--{
		for i := 0; i < j - 1; i++ {
			if len(index.list[i].File_list) < len(index.list[i+1].File_list){
				index.list[i].File_list, index.list[i + 1].File_list = index.list[i + 1].File_list,index.list[i].File_list
			}
		}
	}
	return index
}

/*func (index *Inverted_list) RenewList() {
	iter := index.Db.NewIterator()
	for iter.Next() {
		files := []string{}
		_ = rlp.DecodeBytes(iter.Value(), &files)
		key_file1 := key_file{Key:string(iter.Key()), File_list: files}
		index.list = append(index.list, key_file1)
		index.k = len(index.list)
	}
	iter.Release()
	index.list_sort()
}*/

func (index *Inverted_list) searchKey(key string) key_file {
	var keyfile1 key_file
	var files []string
	data1, _ := index.Db.Get([]byte(key))
	err :=rlp.DecodeBytes(data1, &files)
	if err != nil{
		fmt.Println(err)
	}
	keyfile1.File_list = files
	return keyfile1
}

func (index *Inverted_list) TopKbyTAAT (file File, k int) []string {
	score_rank1 := score_rank{ nil, nil}
	score_rank1.nameToScore = make(map[string]float64)
	for i,j := range file.Keys {
		tf := float64(file.Times[i]) / float64(file.WordSum)
		idf := math.Log2(float64( index.docSum) / float64(index.keyToTimes[j] + 1))
		score := tf * idf
		p := index.keyToKey_file[j]
		for m , n := range index.list[p].File_list {
			if score_rank1.nameToScore[n] == 0 {
				score_rank1.docName = append(score_rank1.docName, n)
				score_rank1.nameToScore[j] = index.list[p].Score[m] * score
			} else {
				score_rank1.nameToScore[j] += index.list[p].Score[m] * score
			}
		}
		score_rank1.sort()
	}
	return score_rank1.docName
}

func (index *Inverted_list) TopKbyDAAT (file File, k int) []File_list {
	kfToScore = make(map[string]float64)
	var key_file_list1  key_file_list
	file_name_list := []string{file.Name}
	score_rank1 := []File_list{}
	file.calSum()
	for i,j := range file.Keys {
		tf := float64(file.Times[i]) / float64(file.WordSum)
		idf := math.Log2(float64( index.docSum) / float64(index.keyToTimes[j] + 1))
		score := tf * idf
		key_file1 := key_file{j, file_name_list, []float64{score}, 0}
		key_file_list1 = append(key_file_list1, key_file1)
	}
	for _, b := range index.file_list {
		kfToScore[b.Name] = 0
	}
	for _, n := range key_file_list1 {
		for _, m := range index.list[index.keyToKey_file[n.Key]].File_list {
			if kfToScore[m] == 0  {
				file_list := index.nameToFileList[m]
				for _, s := range key_file_list1{
					file_list.Score += s.Score[0] * file_list.keyToScore[s.Key]
				}
				//fmt.Println(file_list.Score)
				kfToScore[m] = file_list.Score
				score_rank1 = append(score_rank1, file_list)
			}
		}
	}
	for i := 0; i < len(score_rank1); i++ {
		for j := 1; j < len(score_rank1)-i; j++ {
			if score_rank1[j].Score > score_rank1[j-1].Score {
				score_rank1[j], score_rank1[j-1] = score_rank1[j-1], score_rank1[j]
			}
		}
	}
	return score_rank1
}

 func (index *Inverted_list) TopKbyBAAT (file File, k int, blockNumber uint) []File_list {
	var key_file_list1  key_file_list
	file_name_list := []string{file.Name}
	score_rank1 := []File_list{}
	file_list1 := File_list{
		Name:       "",
		Keys:       nil,
		keyToScore: nil,
		Score:      0,
	}
	 score_rank1 = append(score_rank1, file_list1)
	var maxScore float64
	blockMaxScore := float64(0)
	file.calSum()
	for i,j := range file.Keys {
		tf := float64(file.Times[i]) / float64(file.WordSum)
		idf := math.Log2(float64( index.docSum) / float64(index.keyToTimes[j] + 1))
		score := tf * idf
		key_file1 := key_file{j, file_name_list, []float64{score}, 0}
		key_file_list1 = append(key_file_list1, key_file1)
	}
    for p := blockNumber ; p > 0 ; p-- {
    	for _,n := range key_file_list1 {
    		blockMaxScore += n.Score[0] * index.blockAvgScore[blockNumber][n.Key].score
		}

		if score_rank1 != nil {
			maxScore = score_rank1[0].Score
		} else {
			maxScore = 0
		}
		if blockMaxScore < maxScore {
			continue
		} else {
			for _,s := range index.blockToFileLists[blockNumber] {
				score1 := float64(0)
				for _, t := range key_file_list1 {
					score1 += t.Score[0] * s.keyToScore[t.Key]
				}
				if score1 > score_rank1[0].Score && len(score_rank1) == k {
					score_rank1[0] = s
					score_rank1[0].Score = score1
				} else {
					score_rank1 = append(score_rank1, s)
				}
				/*for v := 0; v < k - 1; v++ {
					if score_rank1[v].Score <  score_rank1[v + 1].Score {
						score_rank1[v], score_rank1[v + 1] = score_rank1[v + 1], score_rank1[v]
					}
				}*/
			}
		}
	}
	return score_rank1
}













