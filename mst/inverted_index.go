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
 	Db   ethdb.Database
}

type key_file_list []key_file

type score_rank struct {
	docName       []string
	//docScore      []float64
	nameToScore   map[string]float64
}

func (sr *score_rank) sort() {
	for i := 0; i < len(sr.docName); i++ {
		for j := 1; j < len(sr.docName) - i; j++ {
			if sr.nameToScore[sr.docName[j]] < sr.nameToScore[sr.docName[j - 1]] {
				sr.docName[j], sr.docName[j - 1] = sr.docName[j - 1], sr.docName[j]
			}
		}
	}
}

func (index *Inverted_list) swap(i, j int) *Inverted_list {
	key_file1 := index.list[i]
	index.list[i] = index.list[j]
	index.list[j] = key_file1
	return index
}

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
				m := len(index.list[i].Score)
				for j := m; j > 1; j-- {
					if index.list[i].Score[j - 1] > index.list[i].Score[j - 2]{
						oldDoc := index.list[i].File_list[j - 1]
						index.list[i].File_list[j - 1] = index.list[i].File_list[j - 2]
						index.list[i].File_list[j - 2] = oldDoc
						oldScore := index.list[i].Score[j - 1]
						index.list[i].Score[j - 1] = index.list[i].Score[j - 2]
						index.list[i].Score[j - 2] = oldScore
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
	list_data,_ := rlp.EncodeToBytes(index.list[p])
	_ = index.Db.Put([]byte(file.Key), list_data)
}

func (index *Inverted_list) list_sort() *Inverted_list {
	for j := index.k; j > 0; j--{
		for i := 0; i < j - 1; i++ {
			if len(index.list[i].File_list) < len(index.list[i+1].File_list){
				index.swap(i, i + 1)
			}
		}
	}
	return index
}

func (index *Inverted_list) RenewList() {
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
}

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

func (index *Inverted_list) TopKbyTAAT (file *File, k int) []string {
	var key_file_list1  key_file_list
	file_name_list := []string{file.Name}
	score_rank1 := score_rank{ nil, nil}
	score_rank1.nameToScore = make(map[string]float64)
	for i,j := range file.Keys {
		tf := float64(file.Times[i] / file.WordSum)
		idf := math.Log2(float64( index.docSum / index.keyToTimes[j] + 1))
		score := tf * idf
		key_file1 := key_file{j, file_name_list, []float64{score}, 0}
		key_file_list1 = append(key_file_list1, key_file1)
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
	return score_rank1.docName[ : k]
}

/*func (index *Inverted_list) TopKbyDAAT (file *File, k int) []string {
	var key_file_list1  key_file_list
	file_name_list := []string{file.Name}
	for i,j := range file.Keys {
		tf := float64(file.Times[i] / file.WordSum)
		idf := math.Log2(float64( index.docSum / index.keyToTimes[j] + 1))
		score := tf * idf
		key_file1 := key_file{j, file_name_list, []float64{score}, 0}
		key_file_list1 = append(key_file_list1, key_file1)
	}




}

func (index *Inverted_list) TopKbyBAAT (file *File, k int) []string {
	var key_file_list1  key_file_list
	file_name_list := []string{file.Name}
	for i,j := range file.Keys {
		tf := float64(file.Times[i] / file.WordSum)
		idf := math.Log2(float64( index.docSum / index.keyToTimes[j] + 1))
		score := tf * idf
		key_file1 := key_file{j, file_name_list, []float64{score}, 0}
		key_file_list1 = append(key_file_list1, key_file1)
	}




}*/













