package types



type FileStore struct {
	StoreFile map[string][]string
}

func NewFileStore() *FileStore {
	return &FileStore{
		StoreFile: make(map[string][]string),
	}
}

func (fs *FileStore) AddFile(userID string, fileName string) {
	fs.StoreFile[userID] = append(fs.StoreFile[userID], fileName)

}

func (fs *FileStore) GetFile(userID string, fileName string) (bool){
	files,ok:=fs.StoreFile[userID]
	 if !ok {
		return false
	 }

	 for _,file:=range(files){
		if file==fileName{
			return true
		}
	 }
	 return false


}
