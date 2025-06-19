package types

type MetaData struct{
	UserID string `json:userID`
	FileArray []string `json: fileArray`
	
}

type FileStore struct{
	StoreFile map[string][]string
}

func NewFileStore()*FileStore{
	return &FileStore{
		StoreFile:make(map[string][]string),
	}
}

func (fs *FileStore) AddFile(userID string,fileName string){
	fs.StoreFile[userID]=append(fs.StoreFile[userID], fileName)
}

func (fs *FileStore) GetFile(userID string,fileName string){
	//todo-how?what type to return
}


