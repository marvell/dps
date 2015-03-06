package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fsouza/go-dockerclient"
	"github.com/marvell/tablewriter"
	"github.com/mgutz/ansi"
)

var (
	dockerClient *docker.Client

	all bool
)

func init() {
	flag.BoolVar(&all, "a", false, "Show all containers.")
}

func main() {
	var err error

	flag.Parse()

	endpoint := "unix:///var/run/docker.sock"
	if docker_host := os.Getenv("DOCKER_HOST"); docker_host != "" {
		endpoint = docker_host
	}
	dockerClient, err = docker.NewClient(endpoint)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	containers, err := containers()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	tableHeader := []string{"ID", "Name", "Image", "IP", "Status"}
	tableData := make([][]string, 0)
	for _, container := range containers {
		tableRow := make([]string, 0)

		ip, err := containerIP(container.ID)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		tableRow = append(tableRow, container.ID[0:12])
		tableRow = append(tableRow, container.Names[0][1:])
		tableRow = append(tableRow, container.Image)
		tableRow = append(tableRow, ip)
		tableRow = append(tableRow, colorStatus(container.Status))

		tableData = append(tableData, tableRow)
	}

	draw(tableHeader, tableData)
}

func colorStatus(status string) string {
	if len(status) == 0 {
		return status
	}

	if status[0:2] == "Up" {
		return ansi.Color(status, "green")
	} else if status[0:6] == "Exited" {
		return ansi.Color(status, "red")
	} else {
		return status
	}
}

func draw(header []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)

	for _, row := range data {
		table.Append(row)
	}

	table.Render()
}

func containers() ([]docker.APIContainers, error) {
	return dockerClient.ListContainers(docker.ListContainersOptions{All: all})
}

func containerIP(id string) (string, error) {
	container, err := dockerClient.InspectContainer(id)
	if err != nil {
		return "", err
	}

	return container.NetworkSettings.IPAddress, nil
}
