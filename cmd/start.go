/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/spf13/cobra"
)

type LanguageConfig struct {
	Language string
	DockerImage string
	// ExecutionCommand func()
	Extension string
}

var LanguageConfigs = map[string]LanguageConfig{
	"cpp": {
		Language: "cpp",
		DockerImage: "gcc:trixie",
		// ExecutionCommand: handleDockerfileExecution,
		Extension: "cpp",
	},
	"go": {
		Language: "golang",
		DockerImage: "golang:alpine",
		// ExecutionCommand: handleDocker`fileExecution,
		Extension: "go",
	},
	"java": {
		Language: "java",
		DockerImage: "openjdk:26-trixie",
		// ExecutionCommand: handleDockerfileExecution,
		Extension: "java",
	},
	"js": {
		Language: "javascript",
		DockerImage: "node:alpine",
		// ExecutionCommand: handleDockerfileExecution,
		Extension: "js",
	},
	"py": {
		Language: "python",
		DockerImage: "python:alpine",
		// ExecutionCommand: handleDockerfileExecution,
		Extension: "py",
	},
	"rs": {
		Language: "rust",
		DockerImage: "rust:alpine",
		// ExecutionCommand: handleDockerfileExecution,
		Extension: "rs",
	},
	"node": {
		Language: "node",
		DockerImage: "node:alpine",
		// ExecutionCommand: handleDockerfileExecution,
		Extension: "js",
	},
}

type Cursor struct{}

func (cursor *Cursor) hide() {
    fmt.Printf("\033[?25l")
}

func (cursor *Cursor) show() {
    fmt.Printf("\033[?25h")
}

func (cursor *Cursor) moveUp(rows int) {
    fmt.Printf("\033[%dF", rows)
}

func (cursor *Cursor) moveDown(rows int) {
    fmt.Printf("\033[%dE", rows)
}

func (cursor *Cursor) clearLine() {
    fmt.Printf("\033[2K")
}

type pullEvent struct {
    ID             string `json:"id"`
    Status         string `json:"status"`
    Error          string `json:"error,omitempty"`
    Progress       string `json:"progress,omitempty"`
    ProgressDetail struct {
        Current int `json:"current"`
        Total   int `json:"total"`
    } `json:"progressDetail"`
}

func handleDockerfileExecution( filename string, image string) string {

	if image == "node:alpine" {
		return fmt.Sprintf(`
			if [ -f /app/%s ]; then
				echo "✓ File exists"
				node /app/%s
			else
				echo "✗ File not found!"
				ls -la /app
				exit 1
			fi
		`, filename, filename)
	}

	if image == "gcc:trixie" {
		cppTrimmedFilename := strings.TrimSuffix(filename, ".cpp")
		return fmt.Sprintf(`
			if [ -f /app/%s ]; then
				echo "✓ File exists"
				g++ /app/%s -o /app/%s
				/app/%s
			else
				echo "✗ File not found!"
				ls -la /app
				exit 1
			fi
		`, filename, filename, cppTrimmedFilename, cppTrimmedFilename)
	}

	if image == "golang:alpine" {
		return fmt.Sprintf(`
			if [ -f /app/%s ]; then
				echo "✓ File exists"
				go run /app/%s
			else
				echo "✗ File not found!"
				ls -la /app
				exit 1
			fi
		`, filename, filename)
	}

	if image == "openjdk:26-trixie" {
		return fmt.Sprintf(`
			if [ -f /app/%s ]; then
				echo "✓ File exists"
				java /app/%s
			else
				echo "✗ File not found!"
				ls -la /app
				exit 1
			fi
		`, filename, filename)
	}

	if image == "python:alpine" {
		return fmt.Sprintf(`
			if [ -f /app/%s ]; then
				echo "✓ File exists"
				python3 /app/%s
			else
				echo "✗ File not found!"
				ls -la /app
				exit 1
			fi
		`, filename, filename)
	}

	if image == "rust:alpine" {
		rustTrimmedFilename := strings.TrimSuffix(filename, ".rs")	
		return fmt.Sprintf(`
			if [ -f /app/%s ]; then
				echo "✓ File exists"
				rustc /app/%s
				/app/%s
			else
				echo "✗ File not found!"
				ls -la /app
				exit 1
			fi
		`, filename, filename, rustTrimmedFilename)
	}

	return ""
	
}

func PullImage(dockerImageName string) bool {
	// https://riptutorial.com/docker/example/31980/image-pulling-with-progress-bars--written-in-go
	// https://docs.docker.com/reference/api/engine/sdk/examples/#run-a-container -> // TODO: try this one too, seems simpler
    client, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

    if err != nil {
        panic(err)
    }

	defer client.Close()

	imageList, err := client.ImageList(context.Background(), image.ListOptions{All: true})
	if err != nil {
		panic(err)
	}
	for _, image := range imageList {
		fmt.Println(image.RepoTags)
	}
	
	imageExists := false
	for _, image := range imageList {
		if image.RepoTags[0] == dockerImageName {
			imageExists = true
		}
	}
	if imageExists {
		fmt.Println("Image already exists")
		return true
	}

    resp, err := client.ImagePull(context.Background(), dockerImageName, image.PullOptions{})

    if err != nil {
        panic(err)
    }

    cursor := Cursor{}
    layers := make([]string, 0)
    oldIndex := len(layers)

    var event *pullEvent
    decoder := json.NewDecoder(resp)

    fmt.Printf("\n")
    cursor.hide()

    for {
        if err := decoder.Decode(&event); err != nil {
            if err == io.EOF {
                break
            }

            panic(err)
        }

        imageID := event.ID

        // Check if the line is one of the final two ones
        if strings.HasPrefix(event.Status, "Digest:") || strings.HasPrefix(event.Status, "Status:") {
            fmt.Printf("%s\n", event.Status)
            continue
        }

        // Check if ID has already passed once
        index := 0
        for i, v := range layers {
            if v == imageID {
                index = i + 1
                break
            }
        }

        // Move the cursor
        if index > 0 {
            diff := index - oldIndex

            if diff > 1 {
                down := diff - 1
                cursor.moveDown(down)
            } else if diff < 1 {
                up := diff*(-1) + 1
                cursor.moveUp(up)
            }

            oldIndex = index
        } else {
            layers = append(layers, event.ID)
            diff := len(layers) - oldIndex

            if diff > 1 {
                cursor.moveDown(diff) // Return to the last row
            }

            oldIndex = len(layers)
        }

        cursor.clearLine()

        if event.Status == "Pull complete" {
            fmt.Printf("%s: %s\n", event.ID, event.Status)
        } else {
            fmt.Printf("%s: %s %s\n", event.ID, event.Status, event.Progress)
        }

    }

    cursor.show()

    if strings.Contains(event.Status, fmt.Sprintf("Downloaded newer image for %s", dockerImageName)) {
        return true
    }

    return false
}

func CreateContainer(dockerImageName string, filePath string) {
	
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		panic(err)
	}
	for _, ctr := range containers {
		fmt.Printf("%s %s (status: %s)\n", ctr.ID, ctr.Image, ctr.Status)
	}
	containerExists := false
	containerId := ""
	for _, ctr := range containers {
		if ctr.Image == dockerImageName {
			containerId = ctr.ID
			containerExists = true
		}
	}
	if containerExists {
		fmt.Println("Container already exists, removing...")
		cli.ContainerRemove(context.Background(), containerId, container.RemoveOptions{Force: true})
		containerExists = false
	}	

	filename := strings.Split(filePath, "/")[len(strings.Split(filePath, "/"))-1]

	if !containerExists {
		fmt.Println("Container does not exist, creating...", filename)

		createdContainer, err := cli.ContainerCreate(context.Background(), &container.Config{
			Image: dockerImageName,
			Cmd:   []string{"/bin/sh", "-c", handleDockerfileExecution(filename, dockerImageName)},
			Tty:   false,
			WorkingDir: "/app",
		}, nil, nil, nil, "")
		if err != nil {
			panic(err)
		}

		fmt.Println("Container created", filename)
		containerId = createdContainer.ID
	}

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	fileContent, _ := os.ReadFile(filePath)
	tw.WriteHeader(&tar.Header{
        Name: filename,
        Mode: 0644,
        Size: int64(len(fileContent)),
    })
    tw.Write(fileContent)
    tw.Close()

	fmt.Println("Copying file to container")

	if err := cli.CopyToContainer(context.Background(), containerId, "/app/", &buf, container.CopyToContainerOptions{}); 
	err != nil {
		panic(err)
	}

	

	fmt.Println("Starting container")

	if err := cli.ContainerStart(context.Background(), containerId, container.StartOptions{}); 
	err != nil {
		panic(err)
	}

	fmt.Println("Waiting for container to finish")

	statusCh, errCh := cli.ContainerWait(context.Background(), containerId, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	fmt.Println("Getting logs")

	out, err := cli.ContainerLogs(context.Background(), containerId, container.LogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	fmt.Println("Streaming logs")
	
	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
}



var startCmd = &cobra.Command{
	Use:   "start",
	Short: "probably could have started at jump, just wanted to call it jump start",
	Long: `Probably could have started at jump, just wanted to call it jump start.
			You can use it to run your code locally and test it against the test cases.
			It's a great way to practice your coding skills and get better at solving problems.
			
			Usage:
			jump start [file]`,
	Run: func(cmd *cobra.Command, args []string) {
		fileName := strings.Split(args[0], "/")[len(strings.Split(args[0], "/"))-1]
		fileExtension := strings.Split(fileName, ".")[len(strings.Split(fileName, "."))-1]

		// TODO: handle language not found
		// TODO: validate input file path and file extension
		// TODO: add language specific Cmd on CreateContainer

		PullImage(LanguageConfigs[fileExtension].DockerImage)
		CreateContainer(LanguageConfigs[fileExtension].DockerImage, args[0])

	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
