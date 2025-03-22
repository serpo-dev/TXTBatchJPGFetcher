package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	fileList      []string
	downloadQueue map[string]map[int]string
	timeout       int = 100
	startTime     time.Time
	totalFiles    int
	downloaded    int
)

func main() {
	fmt.Println("Photo Downloader")

	reader := bufio.NewReader(os.Stdin)

	for {
		clearScreen()
		printMenu()

		choiceStr, _ := reader.ReadString('\n')
		choiceStr = strings.TrimSpace(choiceStr)

		choice, err := strconv.Atoi(choiceStr)
		if err != nil {
			fmt.Println("Invalid choice. Please try again.")
			time.Sleep(1 * time.Second)
			continue
		}

		switch choice {
		case 1:
			addTxtFilesOrDirectories(reader)
		case 2:
			clearFiles()
		case 3:
			setOutputFolder(reader)
		case 4:
			setTimeout(reader)
		case 5:
			startDownload(reader)
		case 6:
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
			time.Sleep(1 * time.Second)
		}
	}
}

func clearScreen() {
	fmt.Print("\033[H\033[2J") // Очистка экрана
}

func printMenu() {
	fmt.Println("Menu:")
	fmt.Println("1. Add TXT Files or Directories")
	fmt.Println("2. Clear All Files")
	fmt.Println("3. Set Output Folder")
	fmt.Println("4. Set Timeout (ms)")
	fmt.Println("5. Start Download")
	fmt.Println("6. Exit")
	fmt.Print("Choose an option: ")
}

func addTxtFilesOrDirectories(reader *bufio.Reader) {
	fmt.Print("Enter the path to the TXT file or directory: ")
	path, _ := reader.ReadString('\n')
	path = strings.TrimSpace(path)
	path = filepath.Clean(path)

	info, err := os.Stat(path)
	if err != nil {
		fmt.Println("Path not found:", path)
		time.Sleep(1 * time.Second)
		return
	}

	if info.IsDir() {
		err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(info.Name(), ".txt") {
				fileList = append(fileList, filePath)
				fmt.Println("File added:", filePath)
			}
			return nil
		})
		if err != nil {
			fmt.Println("Error walking directory:", err)
		}
	} else {
		if strings.HasSuffix(info.Name(), ".txt") {
			fileList = append(fileList, path)
			fmt.Println("File added:", path)
		} else {
			fmt.Println("The file is not a .txt file:", path)
		}
	}
	time.Sleep(1 * time.Second)
}

func clearFiles() {
	fileList = nil
	downloadQueue = make(map[string]map[int]string)
	fmt.Println("All files cleared.")
	time.Sleep(1 * time.Second)
}

func setOutputFolder(reader *bufio.Reader) {
	fmt.Print("Enter the output folder path: ")
	folderPath, _ := reader.ReadString('\n')
	folderPath = strings.TrimSpace(folderPath)
	folderPath = filepath.Clean(folderPath)

	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		fmt.Println("Folder does not exist. Creating...")
		err := os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			fmt.Println("Failed to create folder:", err)
			time.Sleep(1 * time.Second)
			return
		}
	}

	fmt.Println("Output folder set to:", folderPath)
	time.Sleep(1 * time.Second)
}

func setTimeout(reader *bufio.Reader) {
	fmt.Print("Enter the timeout in milliseconds: ")
	timeoutStr, _ := reader.ReadString('\n')
	timeoutStr = strings.TrimSpace(timeoutStr)

	if t, err := strconv.Atoi(timeoutStr); err == nil {
		timeout = t
		fmt.Println("Timeout set to:", timeout, "ms")
	} else {
		fmt.Println("Invalid timeout value.")
	}
	time.Sleep(1 * time.Second)
}

func startDownload(reader *bufio.Reader) {
	if len(fileList) == 0 {
		fmt.Println("No files selected. Please add TXT files first.")
		time.Sleep(1 * time.Second)
		return
	}

	fmt.Print("Enter the output folder path: ")
	outputFolder, _ := reader.ReadString('\n')
	outputFolder = strings.TrimSpace(outputFolder)
	outputFolder = filepath.Clean(outputFolder)

	if _, err := os.Stat(outputFolder); os.IsNotExist(err) {
		fmt.Println("Folder does not exist. Creating...")
		err := os.MkdirAll(outputFolder, os.ModePerm)
		if err != nil {
			fmt.Println("Failed to create folder:", err)
			time.Sleep(1 * time.Second)
			return
		}
	}

	startTime = time.Now()
	totalFiles = 0
	downloaded = 0
	downloadQueue = make(map[string]map[int]string)

	for _, filePath := range fileList {
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println("Failed to open file:", filePath, err)
			continue
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		folderName := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
		folderPath := filepath.Join(outputFolder, folderName)
		os.MkdirAll(folderPath, os.ModePerm)

		index := 1
		downloadQueue[folderPath] = make(map[int]string)
		for scanner.Scan() {
			url := scanner.Text()
			if url != "" {
				downloadQueue[folderPath][index] = url
				totalFiles++
				index++
			}
		}
	}

	fmt.Println("Starting download...")
	fmt.Printf("Total files to download: %d\n", totalFiles)

	for len(downloadQueue) > 0 {
		for folderPath, files := range downloadQueue {
			for index, url := range files {
				if downloadImage(index, url, folderPath) {
					delete(downloadQueue[folderPath], index)
					downloaded++
					printProgress(downloaded, totalFiles)
				} else {
					fmt.Printf("\rFailed to download: %d. Retrying...", index)
				}

				time.Sleep(time.Duration(timeout) * time.Millisecond)
			}

			if len(downloadQueue[folderPath]) == 0 {
				delete(downloadQueue, folderPath)
			}
		}
	}

	fmt.Println("\nDownload completed.")
	time.Sleep(2 * time.Second)
}

func downloadImage(index int, url, folderPath string) bool {
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false
	}

	fileName := filepath.Join(folderPath, fmt.Sprintf("%d.jpg", index))
	file, err := os.Create(fileName)
	if err != nil {
		return false
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err == nil
}

func printProgress(downloaded, totalFiles int) {
	elapsed := time.Since(startTime)
	speed := float64(downloaded) / elapsed.Seconds()
	fmt.Printf("\rDownloaded: %d/%d | Time elapsed: %02d:%02d:%02d | Speed: %.2f photos/sec", downloaded, totalFiles, int(elapsed.Hours()), int(elapsed.Minutes())%60, int(elapsed.Seconds())%60, speed)
}
