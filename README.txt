# Large File Upload Demo

A demonstration of large file upload with chunk upload and resume capability.

## Features

- Chunk-based file upload
- Upload progress tracking
- Resume interrupted uploads
- MD5 hash verification
- Automatic chunk merging
- Modern and responsive UI

## Tech Stack

Backend:
- Go 1.20+
- Gin web framework
- File system storage

Frontend:
- HTML5
- Alpine.js for reactivity
- Axios for HTTP requests
- Tailwind CSS for styling
- SparkMD5 for file hashing

## Project Structure

```
chunk-upload/
├── backend/
│   ├── go.mod         # Go module definition
│   └── main.go        # Backend server implementation
├── frontend/
│   └── index.html     # Frontend UI and logic
└── uploads/
    ├── temp/          # Temporary chunk storage
    └── complete/      # Completed file storage
```

## How It Works

1. File Selection:
   - User selects a large file
   - Frontend calculates file MD5 hash
   - File is split into 2MB chunks

2. Upload Process:
   - Frontend checks existing chunks on server
   - Only missing chunks are uploaded
   - Real-time progress tracking
   - Automatic resume capability

3. File Merging:
   - Server automatically merges chunks when all are received
   - Original file is reconstructed in the complete directory
   - Temporary chunks are cleaned up

## Setup and Run

1. Start the backend server:
   ```bash
   cd backend
   go mod tidy
   go run main.go
   ```

2. Open frontend:
   - Open frontend/index.html in a web browser
   - Server runs on http://localhost:8080

## API Endpoints

1. GET /upload/chunk/check
   - Check existing chunks
   - Query params: hash

2. POST /upload/chunk/add
   - Upload a single chunk
   - Form params: chunk, index, hash, name, total

## Features in Detail

1. Chunk Upload:
   - Fixed 2MB chunk size
   - Parallel upload support
   - Automatic chunk management

2. Resume Support:
   - Tracks uploaded chunks
   - Continues from last successful chunk
   - No duplicate uploads

3. Progress Tracking:
   - Real-time upload progress
   - Individual chunk progress
   - Overall file progress

4. Error Handling:
   - Network error recovery
   - Invalid file handling
   - Upload cancellation support

## Security Considerations

- File hash verification
- Safe file path handling
- No original filename in temporary storage
- Automatic cleanup of temporary files
