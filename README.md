# c-hdfs (Custom Distributed File System)

A simple distributed file system written in Go, inspired by HDFS. This project demonstrates file encryption, chunking, and distributed storage across multiple nodes.

## Features

- **File Upload:** Uploads files, encrypts them, splits them into chunks, and distributes chunks to multiple nodes.
- **File Download:** Retrieves and reassembles file chunks from nodes, then decrypts the file.
- **Thread-Safe Metadata:** Tracks uploaded files per user using a thread-safe in-memory store.
- **Basic Replication:** Each chunk is uploaded to two nodes for redundancy.
- **AES Encryption:** Files are encrypted before storage for security.

## Project Structure

```
go.mod
go.sum
node1/
  chunking.go         # File splitting and merging logic
  main.go             # Main server and API endpoints
  setNodeInfo.go      # Node configuration and initialization
  uploading.go        # Chunk upload/download logic
  decryptedFiles/     # (output) Decrypted files
  encryptedFiles/     # (output) Encrypted files
  encryption/
    encryption.go     # AES encryption/decryption utilities
  temp/               # (output) Temporary files
  types/
    types.go          # Thread-safe file metadata store
```

## Getting Started

### Prerequisites

- Go 1.18 or newer

### Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/mirdha8846/c-hdfs.git
    cd c-hdfs/node1
    ```

2. Download dependencies:
    ```sh
    go mod tidy
    ```

### Running the Node

Start the node server:
```sh
go run main.go
```
The server will start on port `8082`.

### API Endpoints

#### Upload File

- **POST** `/api/fileUpload`
- **Form Data:**
  - `userID`: User identifier
  - `file`: File to upload

#### Download File

- **POST** `/api/getFiles`
- **Form Data:**
  - `userID`: User identifier
  - `fileName`: Name of the file to retrieve

## How It Works

1. **Upload:**
   - The file is saved to a temp directory.
   - It is encrypted using AES.
   - The encrypted file is split into 3 chunks.
   - Each chunk is uploaded to two nodes for redundancy.
   - Metadata is stored in a thread-safe in-memory map.

2. **Download:**
   - Checks if the file exists for the user.
   - Downloads all chunks from the nodes.
   - Merges the chunks and decrypts the file.
   - Returns the decrypted file.

## Limitations

----

## License

MIT

---

**For educational and demonstration purposes