package file

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func FileOperations() {
	fmt.Println("=== File Operations ===")

	filename := "test.txt"
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	content := "Hello, World!\nThis is a test file.\nGo is awesome!"
	_, err = file.WriteString(content)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return
	}
	fmt.Printf("✅ Created and wrote to %s\n", filename)

	file, err = os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	fmt.Println("📖 File content:")
	for scanner.Scan() {
		fmt.Printf("  %s\n", scanner.Text())
	}

	fileInfo, err := os.Stat(filename)
	if err != nil {
		fmt.Printf("Error getting file info: %v\n", err)
		return
	}

	fmt.Printf("📊 File info:\n")
	fmt.Printf("  Name: %s\n", fileInfo.Name())
	fmt.Printf("  Size: %d bytes\n", fileInfo.Size())
	fmt.Printf("  Mode: %s\n", fileInfo.Mode())
	fmt.Printf("  Modified: %s\n", fileInfo.ModTime())
	fmt.Printf("  Is Directory: %t\n", fileInfo.IsDir())
}

func FolderOperations() {
	fmt.Println("\n=== Folder Operations ===")

	dirName := "test_folder"
	err := os.Mkdir(dirName, 0755)
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
	} else {
		fmt.Printf("✅ Created directory: %s\n", dirName)
	}

	nestedDir := "test_folder/nested/deep"
	err = os.MkdirAll(nestedDir, 0755)
	if err != nil {
		fmt.Printf("Error creating nested directories: %v\n", err)
	} else {
		fmt.Printf("✅ Created nested directories: %s\n", nestedDir)
	}

	fmt.Println("📁 Directory contents:")
	entries, err := os.ReadDir(".")
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		entryType := "📄"
		if entry.IsDir() {
			entryType = "📁"
		}
		fmt.Printf("  %s %s (%d bytes, %s)\n",
			entryType, entry.Name(), info.Size(), info.ModTime().Format("15:04:05"))
	}

	fmt.Println("\n🌳 Directory tree walk:")
	err = filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		level := strings.Count(path, string(os.PathSeparator))
		indent := strings.Repeat("  ", level)
		entryType := "📄"
		if d.IsDir() {
			entryType = "📁"
		}
		fmt.Printf("%s%s %s\n", indent, entryType, d.Name())
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
	}
}

func OSOperations() {
	fmt.Println("\n=== OS Operations ===")

	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting working directory: %v\n", err)
	} else {
		fmt.Printf("📍 Current directory: %s\n", wd)
	}

	originalDir, _ := os.Getwd()
	err = os.Chdir("test_folder")
	if err != nil {
		fmt.Printf("Error changing directory: %v\n", err)
	} else {
		fmt.Println("✅ Changed to test_folder")
		newDir, _ := os.Getwd()
		fmt.Printf("📍 New directory: %s\n", newDir)

		os.Chdir(originalDir)
		fmt.Println("✅ Changed back to original directory")
	}

	fmt.Println("\n🌍 Environment Variables:")
	fmt.Printf("  HOME: %s\n", os.Getenv("HOME"))
	fmt.Printf("  USER: %s\n", os.Getenv("USER"))
	fmt.Printf("  PATH: %s\n", os.Getenv("PATH"))
	fmt.Printf("  GOPATH: %s\n", os.Getenv("GOPATH"))

	fmt.Println("\n💻 System Information:")
	fmt.Printf("  OS: %s\n", runtime.GOOS)
	fmt.Printf("  Architecture: %s\n", runtime.GOARCH)
	fmt.Printf("  Go Version: %s\n", runtime.Version())
	fmt.Printf("  Num CPU: %d\n", runtime.NumCPU())

	fmt.Printf("  Process ID: %d\n", os.Getpid())
	fmt.Printf("  Parent Process ID: %d\n", os.Getppid())
}

func EnvironmentVariables() {
	fmt.Println("\n=== Environment Variables ===")

	err := os.Setenv("MY_CUSTOM_VAR", "Hello from Go!")
	if err != nil {
		fmt.Printf("Error setting environment variable: %v\n", err)
	} else {
		fmt.Println("✅ Set MY_CUSTOM_VAR")
	}

	value := os.Getenv("MY_CUSTOM_VAR")
	fmt.Printf("📝 MY_CUSTOM_VAR = %s\n", value)

	fmt.Println("\n🌍 All Environment Variables:")
	envVars := os.Environ()
	for i, env := range envVars {
		if i >= 10 {
			fmt.Printf("  ... and %d more\n", len(envVars)-10)
			break
		}
		fmt.Printf("  %s\n", env)
	}

	value, exists := os.LookupEnv("PATH")
	if exists {
		fmt.Printf("\n🔍 PATH exists: %s\n", value[:50]+"...")
	} else {
		fmt.Println("🔍 PATH does not exist")
	}
}

func FileManipulation() {
	fmt.Println("\n=== File Manipulation ===")

	srcFile := "test.txt"
	dstFile := "test_copy.txt"

	err := copyFile(srcFile, dstFile)
	if err != nil {
		fmt.Printf("Error copying file: %v\n", err)
	} else {
		fmt.Printf("✅ Copied %s to %s\n", srcFile, dstFile)
	}

	newName := "test_renamed.txt"
	err = os.Rename(dstFile, newName)
	if err != nil {
		fmt.Printf("Error renaming file: %v\n", err)
	} else {
		fmt.Printf("✅ Renamed %s to %s\n", dstFile, newName)
	}

	if _, err := os.Stat(newName); os.IsNotExist(err) {
		fmt.Printf("❌ File %s does not exist\n", newName)
	} else {
		fmt.Printf("✅ File %s exists\n", newName)
	}

	fileInfo, err := os.Stat(newName)
	if err != nil {
		fmt.Printf("Error getting file info: %v\n", err)
	} else {
		fmt.Printf("📋 File permissions: %s\n", fileInfo.Mode())
	}
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return destFile.Sync()
}

func TemporaryFiles() {
	fmt.Println("\n=== Temporary Files ===")

	tempFile, err := os.CreateTemp("", "go_temp_*.txt")
	if err != nil {
		fmt.Printf("Error creating temp file: %v\n", err)
		return
	}
	defer os.Remove(tempFile.Name())

	fmt.Printf("✅ Created temporary file: %s\n", tempFile.Name())

	_, err = tempFile.WriteString("This is a temporary file!\n")
	if err != nil {
		fmt.Printf("Error writing to temp file: %v\n", err)
		return
	}

	tempFile.Seek(0, 0)
	content, err := io.ReadAll(tempFile)
	if err != nil {
		fmt.Printf("Error reading temp file: %v\n", err)
		return
	}
	fmt.Printf("📖 Temp file content: %s", string(content))

	tempDir, err := os.MkdirTemp("", "go_temp_dir_*")
	if err != nil {
		fmt.Printf("Error creating temp directory: %v\n", err)
		return
	}
	defer os.RemoveAll(tempDir)

	fmt.Printf("✅ Created temporary directory: %s\n", tempDir)
}

func FilePermissions() {
	fmt.Println("\n=== File Permissions ===")

	filename := "permission_test.txt"
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()
	defer os.Remove(filename)

	file.WriteString("Permission test file")

	fileInfo, err := os.Stat(filename)
	if err != nil {
		fmt.Printf("Error getting file info: %v\n", err)
		return
	}

	fmt.Printf("📋 File permissions: %s\n", fileInfo.Mode())
	fmt.Printf("📋 Octal permissions: %o\n", fileInfo.Mode().Perm())

	err = os.Chmod(filename, 0755)
	if err != nil {
		fmt.Printf("Error changing permissions: %v\n", err)
		return
	}

	fileInfo, _ = os.Stat(filename)
	fmt.Printf("📋 New permissions: %s\n", fileInfo.Mode())
}

func RunAllFileExamples() {
	fmt.Println("🎯 Go File & OS Operations Examples")
	fmt.Println("====================================")

	FileOperations()
	FolderOperations()
	OSOperations()
	EnvironmentVariables()
	FileManipulation()
	TemporaryFiles()
	FilePermissions()

	fmt.Println("\n✅ All file operations completed!")
}
