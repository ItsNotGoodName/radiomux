//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var Default = Build

const (
	packageName = "github.com/ItsNotGoodName/radiomux"
	distDir     = "dist"
	name        = "radiomux"
)

var goos []string = []string{
	"linux",
	"windows",
	"darwin",
}
var goarch [][]string = [][]string{
	{"amd64", "arm", "arm64"},
	{"amd64"},
	{"arm64"},
}

func ldflags() string {
	buildPackageName := packageName + "/internal/build/build"
	version := "nightly" // TODO: figure out version automatically
	commit, _ := sh.Output("git", "rev-parse", "HEAD")
	date := time.Now().UTC().Format(time.RFC3339)
	repoURL := fmt.Sprintf("https://%s", packageName)

	return fmt.Sprintf("-s -w -X %s.Version=%s -X %s.Commit=%s -X %s.Date=%s -X %s.RepoURL=%s",
		buildPackageName, version,
		buildPackageName, commit,
		buildPackageName, date,
		buildPackageName, repoURL,
	)
}

func fileName(os, arch string) string {
	return fmt.Sprintf("%s_%s_%s", name, os, arch)
}

func binaryFileName(os, arch string) string {
	s := fileName(os, arch)
	if os == "windows" {
		return s + ".exe"
	}
	return s
}

func mkdirDistDir() error {
	return sh.Run("mkdir", "-p", distDir)
}

func modTidy() error {
	fmt.Println("Downloading...")
	return sh.RunV("go", "mod", "tidy")
}

func generate() error {
	mg.Deps(modTidy)
	fmt.Println("Generating...")
	return sh.RunV("go", "generate", "./...")
}

func BuildAPK() error {
	mg.Deps(mkdirDistDir)

	apkName := "radiomuxplayer"
	outputDebugAPKPath := path.Join(distDir, apkName+"-debug.apk")
	fmt.Println("Building", outputDebugAPKPath)

	if err := sh.RunV("chmod", "+x", "./android/gradlew"); err != nil {
		return err
	}

	if err := sh.RunV("sh", "-c", "cd android && ./gradlew build"); err != nil {
		return err
	}

	if err := sh.RunV("sh", "-c", "cd android && ./gradlew assembleDebug"); err != nil {
		return err
	}

	return sh.Run("cp", "./android/app/build/outputs/apk/debug/app-debug.apk", outputDebugAPKPath)
}

func BuildServer() error {
	mg.Deps(modTidy, generate, mkdirDistDir)

	main := packageName + "/cmd/radiomux"
	ldflags := ldflags()

	for osIndex, os := range goos {
		for _, arch := range goarch[osIndex] {
			binaryFilePath := path.Join(distDir, binaryFileName(os, arch))
			fmt.Println("Building", binaryFilePath)

			if err := sh.RunWith(map[string]string{
				"CGO_ENABLED": "0",
				"GOOS":        os,
				"GOARCH":      arch,
			}, "go", "build", "-ldflags", ldflags, "-o", binaryFilePath, main); err != nil {
				return err
			}
		}
	}

	return nil
}

func Build() error {
	mg.Deps(mkdirDistDir, BuildServer, BuildAPK)

	artifactsDir := path.Join(distDir, "artifacts")
	if err := sh.Run("mkdir", "-p", artifactsDir); err != nil {
		return err
	}

	for osIndex, os := range goos {
		for _, arch := range goarch[osIndex] {
			binaryFilePath := path.Join(distDir, binaryFileName(os, arch))
			files := []string{binaryFilePath, "README.md", "LICENSE"}

			if os == "windows" {
				zipFilePath := path.Join(artifactsDir, fileName(os, arch)+".zip")
				fmt.Println("Compressing", zipFilePath)

				args := append([]string{zipFilePath}, files...)

				sh.Run("zip", args...)
			} else {
				tarFilePath := path.Join(artifactsDir, fileName(os, arch)+".tar.gz")
				fmt.Println("Compressing", tarFilePath)

				args := append([]string{"-czf", tarFilePath}, files...)

				sh.Run("tar", args...)
			}
		}
	}

	return nil
}

func Clean() error {
	return os.RemoveAll("dist")
}
