# TXTBatchJPGFetcher

### Batch .JPG images fetcher from .TXT files

Photo Downloader is a command-line tool written in Go that allows you to download images from URLs listed in `.txt` files. The program supports batch processing, allowing you to add multiple `.txt` files or directories containing `.txt` files. It also provides options to set an output folder, configure a timeout between downloads, and track download progress.

## Features

- **Add TXT Files or Directories**: Add individual `.txt` files or directories containing `.txt` files.
- **Clear All Files**: Clear the list of added files.
- **Set Output Folder**: Specify the folder where downloaded images will be saved.
- **Set Timeout**: Configure the delay (in milliseconds) between each download.
- **Start Download**: Begin downloading images from the URLs listed in the `.txt` files.
- **Progress Tracking**: View real-time progress, including the number of files downloaded, elapsed time, and download speed.

## Installation

1. Ensure you have [Go](https://golang.org/dl/) installed on your system.
2. Clone this repository or download the source code.
1. Go to rhe directory and run the command:
   ```bash
   go run main.go
   ```

### Build
1. Navigate to the project directory and build the program:
   ```bash
   go build -o TXTBatchJPGFetcher.exe
   ```
2. Run the program:
    ```bash
    ./TXTBatchJPGFetcher.exe
    ```

## Usage

```bash
Menu:
1. Add TXT Files or Directories
2. Clear All Files
3. Set Output Folder
4. Set Timeout (ms)
5. Start Download
6. Exit
```

1. Add TXT Files or Directories:
- Enter the path to a .txt file or a directory containing .txt files.
- The program will scan the files and add them to the download queue.
- You can repeat it many times to add more files from different directories.
2. Clear All Files:
- Select option 2 to clear the list of added files.
3. Set Output Folder:
- If the folder does not exist, the program will create it.
- You can skip this step. When you run the downloading, it suggests to set the output dir.
4. Set Timeout:
- Choose the value in milliseconds (ms).
- By default, it's 100 ms between every connection.
5. Start Download:
- The program will display real-time progress, including the number of files downloaded, elapsed time, and download speed.

### Example

```bash
Starting download...
Total files to download: 12118
Downloaded: 1823/12118 | Time elapsed: 00:08:50 | Speed: 3.44 photos/sec
```

6. Exit:
- Select the option to exit the program.

## Contributing

**Contributions are welcome!** 
Please open an issue or submit a pull request for any improvements or bug fixes.

## Authors

- **[serpo-dev](https://github.com/serpo-dev)**