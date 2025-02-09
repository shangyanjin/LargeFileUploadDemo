package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strings"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
)

// ChunkInfo contains information about a file chunk
type ChunkInfo struct {
    Index    int    `json:"index"`    // Index of the chunk
    Hash     string `json:"hash"`     // File hash for identification
    Name     string `json:"name"`     // Original file name
    Total    int    `json:"total"`    // Total number of chunks
}

func main() {
    r := gin.Default()

    // Configure CORS
    config := cors.DefaultConfig()
    config.AllowAllOrigins = true
    config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type"}
    r.Use(cors.New(config))

    // Create upload directories
    os.MkdirAll("uploads/temp", 0755)
    os.MkdirAll("uploads/complete", 0755)

    // Check which chunks already exist
    r.GET("/upload/chunk/check", func(c *gin.Context) {
        hash := c.Query("hash")
        
        chunks := make([]int, 0)
        pattern := fmt.Sprintf("uploads/temp/%s-*", hash)
        
        // Find all chunks for this file
        files, _ := filepath.Glob(pattern)
        for _, file := range files {
            // Extract chunk index from filename
            parts := strings.Split(file, "-")
            if len(parts) > 0 {
                if index, err := strconv.Atoi(parts[len(parts)-1]); err == nil {
                    chunks = append(chunks, index)
                }
            }
        }

        c.JSON(200, gin.H{"chunks": chunks})
    })

    // Handle chunk upload
    r.POST("/upload/chunk/add", func(c *gin.Context) {
        file, _ := c.FormFile("chunk")
        index, _ := strconv.Atoi(c.PostForm("index"))
        hash := c.PostForm("hash")
        name := c.PostForm("name")
        total, _ := strconv.Atoi(c.PostForm("total"))

        // Save the chunk to temp directory
        tempPath := fmt.Sprintf("uploads/temp/%s-%d", hash, index)
        if err := c.SaveUploadedFile(file, tempPath); err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }

        // Check if all chunks have been uploaded
        pattern := fmt.Sprintf("uploads/temp/%s-*", hash)
        files, _ := filepath.Glob(pattern)
        
        if len(files) == total {
            // All chunks uploaded, start merging
            finalPath := filepath.Join("uploads/complete", name)
            finalFile, err := os.Create(finalPath)
            if err != nil {
                c.JSON(500, gin.H{"error": err.Error()})
                return
            }
            defer finalFile.Close()

            // Merge all chunks in order
            for i := 0; i < total; i++ {
                chunkPath := fmt.Sprintf("uploads/temp/%s-%d", hash, i)
                chunkFile, err := os.Open(chunkPath)
                if err != nil {
                    c.JSON(500, gin.H{"error": fmt.Sprintf("Error opening chunk %d: %v", i, err)})
                    return
                }
                
                // Copy chunk data to final file
                _, err = io.Copy(finalFile, chunkFile)
                chunkFile.Close()
                if err != nil {
                    c.JSON(500, gin.H{"error": fmt.Sprintf("Error copying chunk %d: %v", i, err)})
                    return
                }
                
                // Clean up: remove merged chunk
                os.Remove(chunkPath)
            }

            c.JSON(200, gin.H{
                "message": "All chunks uploaded and merged successfully",
                "merged": true,
            })
            return
        }

        c.JSON(200, gin.H{
            "message": "Chunk uploaded successfully",
            "merged": false,
        })
    })

    r.Run(":8080")
}
