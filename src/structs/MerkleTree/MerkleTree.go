package MerkleTree

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"os"
)

type Merkle struct {

	Tree [][] string
}


func (m Merkle) CreateTree(filename string, numOfChunks int64, chunkSize int64) {


	leaves := make([][]byte, numOfChunks)

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	fStat, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(fStat.Size())


	for i:=0; i<int(numOfChunks) ; i++ {

		buffer := make([]byte, chunkSize)

		if i == int(numOfChunks) - 1 {

			buffer = make([]byte, fStat.Size() % chunkSize)

		}

		fmt.Println(len(buffer))

		bytesRead, err := file.Read(buffer)

		if err != nil {
			fmt.Println(err)
			log.Fatal(err)
		}

		//fmt.Printf("Buffer : %+v\n", buffer)

		leaves[i] = buffer[:bytesRead]

		//fmt.Printf("Leave : %+v\n", leaves[i])

	}

	//m.Tree = make([]string, math.Pow(2, float64(numOfChunks)) - 1)
	hasher := md5.New()

	hashedLeaves := make([]string, len(leaves))

	for i:=0; i<len(leaves) ;i++  {
		hasher.Write(leaves[i])
		hashedLeaves[i] = hex.EncodeToString(hasher.Sum(nil))
	}

	currentLevel := new([]string)
	currentLevel = &hashedLeaves

	nextLevel := new([]string)

	m.Tree = append(m.Tree, *currentLevel)


	for ; ;  {

		for i:=0;i<len(*currentLevel) ;i = i + 2 {

			if i+1 < len(*currentLevel) {
				tmpStr := (*currentLevel)[i] + (*currentLevel)[i+1] + string(i)
				hasher.Write([]byte(tmpStr))
				*nextLevel = append(*nextLevel, hex.EncodeToString(hasher.Sum(nil)))
			} else {
				tmpStr := (*currentLevel)[i] + (*currentLevel)[i] + string(i)
				hasher.Write([]byte(tmpStr))
				*nextLevel = append(*nextLevel, hex.EncodeToString(hasher.Sum(nil)))

			}

		}

		m.Tree = append(m.Tree, *nextLevel)
		currentLevel = nextLevel

		if len(*nextLevel) == 1 {
			break
		}

		nextLevel = new([]string)

	}

	fmt.Printf("Merkle tree : %+v\n", m)

}
