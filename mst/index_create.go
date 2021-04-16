package mst

import (
	"github.com/ethereum/go-ethereum/ethdb"
	"math"
)

type File struct {
	Name   string
	Keys   []string
	Times  []int
	WordSum  int
}

type File_list struct {
	Name          string
	Keys          []string
	keyToScore    map[string]float64
	Score         float64
}

type Block_file_list struct {

}

func CreateIndex() *Inverted_list {
	index := &Inverted_list{0, nil, 0, nil, nil, nil, nil, nil}
	index.nameToFileList = make(map[string]File_list)
	return index
}

func (in *Inverted_list) putDb(db ethdb.Database) {
	in.Db = db
}

//TODO
/*func UpdateIndex(f1 File, blockNumber uint, in *Inverted_list, t *MST) *MST {
	in.docSum ++
	index := Inverted_list{Db: in.Db}
	index.RenewList()
	if f1.WordSum == 0 {
		f1.calSum()
	}
	key_file_list1 := f1.fileToKey(in)
	for _, j:= range key_file_list1 {
		index.insert(j)
	}
	index.list_sort()
	//f1.keysSort(in.Db)
	t.root_insert(index_info{key: f1.Keys, pos: blockNumber}, in.Db)
	return t
}*/

func UpdateIndex(f1 File, blockNumber uint, in *Inverted_list) {
	in.docSum ++
	index := Inverted_list{Db: in.Db}
	index.RenewList()
	if f1.WordSum == 0 {
		f1.calSum()
	}
	key_file_list1 := f1.fileToKey(in)
	for _, j:= range key_file_list1 {
		index.insert(j)
	}
	index.list_sort()
}


func (file *File) fileToKey(in *Inverted_list) key_file_list {
	var key_file_list1  key_file_list
	file_list1 := File_list{
		Name:       file.Name,
		Keys:       file.Keys,
		keyToScore: make(map[string]float64),
		Score:      0,
	}
	for i,j := range file.Keys {
		in.keyToTimes[j]++
		file_name_list := []string{file.Name}
		tf := float64(file.Times[i] / file.WordSum)
		idf := math.Log2(float64( in.docSum / in.keyToTimes[j] + 1))
		score := tf * idf
		file_list1.keyToScore[j] = score
		key_file1 := key_file{j, file_name_list, []float64{score}, 0}
		key_file_list1 = append(key_file_list1, key_file1)
	}
	in.nameToFileList[file.Name] = file_list1
	in.file_list = append(in.file_list, file_list1)
	return key_file_list1
}

/*func (file1 *File) keysSort(db ethdb.Database) {
	var key_file_pre key_file
	var key_file_next key_file
	for j := len(file1.Keys); j > 0; j-- {
		for i := 0;i < j - 1;i++ {
			data, _ := db.Get([]byte(file1.Keys[i]))
            key_file_pre.Key = file1.Keys[i]
			rlp.DecodeBytes(data,&key_file_pre.File_list)
			key_file_next.Key = file1.Keys[i]
			rlp.DecodeBytes(data,&key_file_next.File_list)
            if len(key_file_pre.Key) < len(key_file_next.Key) {
            	file1.swap(i, i + 1)
			}
		}
	}
}*/

func (file *File) swap(i, j int) {
	flag := file.Keys[i]
	file.Keys[j] = file.Keys[i]
	file.Keys[i] = flag
}


func (file *File) calSum() {
	for _,j := range file.Times {
		file.WordSum += j
	}
}