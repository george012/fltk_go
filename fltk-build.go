//go:build ignore

package main

import (
	"errors"
	"fmt"
	"github.com/george012/fltk_go/config"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	commit      = "tags/release-1.4.2"
	fltkRepoURL = "https://github.com/fltk/fltk.git"
	patchFile   = "../../lib/fltk-1.4.patch"
	buildDir    = "fltk_build"
	sourceDir   = "fltk"
)

type BuildConfig struct {
	OS           string
	Arch         string
	LibDir       string
	IncludeDir   string
	CurrentDir   string
	CMakeOptions []string
}

func main() {
	if runtime.GOOS == "darwin" {
		// 在macOS上自动执行双架构编译
		buildForDarwinUniversal()
	} else {
		// 其他平台正常编译
		buildForPlatform(runtime.GOOS, runtime.GOARCH)
	}
}

func buildForDarwinUniversal() {
	fmt.Println("Building universal binary for macOS (arm64 + amd64)")

	// 先构建arm64
	buildForPlatform("darwin", "arm64")

	// 再构建amd64
	buildForPlatform("darwin", "amd64")

	// 合并为通用二进制
	createUniversalBinary()
}

func buildForPlatform(os, arch string) {
	fmt.Printf("\nBuilding for OS: %s, architecture: %s\n", os, arch)

	cfg := prepareBuildConfig(os, arch)
	checkRequiredTools()
	createDirectories(cfg)
	setupFLTKSource(cfg)
	runCMake(cfg)
	buildAndInstall(cfg)
	generateCgoFile(cfg)
}

func createUniversalBinary() {
	fmt.Println("\nCreating universal binary...")

	arm64LibDir := filepath.Join("lib", "darwin", "arm64")
	amd64LibDir := filepath.Join("lib", "darwin", "amd64")
	universalLibDir := filepath.Join("lib", "darwin", "universal")
	arm64IncludeDir := filepath.Join("include", "darwin", "arm64")
	universalIncludeDir := filepath.Join("include", "darwin", "universal")

	if err := os.MkdirAll(universalLibDir, 0750); err != nil {
		fail(fmt.Sprintf("Could not create universal lib directory: %v", err))
	}
	if err := os.MkdirAll(universalIncludeDir, 0750); err != nil {
		fail(fmt.Sprintf("Could not create universal include directory: %v", err))
	}

	// 合并所有.a文件
	files, err := os.ReadDir(arm64LibDir)
	if err != nil {
		fail(fmt.Sprintf("Error reading arm64 lib directory: %v", err))
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".a") {
			continue
		}

		arm64Lib := filepath.Join(arm64LibDir, file.Name())
		amd64Lib := filepath.Join(amd64LibDir, file.Name())
		universalLib := filepath.Join(universalLibDir, file.Name())

		// 使用lipo工具合并
		cmd := exec.Command("lipo", "-create", "-output", universalLib, arm64Lib, amd64Lib)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err = cmd.Run(); err != nil {
			fail(fmt.Sprintf("Error creating universal binary for %s: %v", file.Name(), err))
		}
	}

	// 复制arm64的头文件到universal目录
	copyDir(arm64IncludeDir, universalIncludeDir)

	// 生成universal的cgo文件
	generateUniversalCgoFile()

	fmt.Println("Successfully created universal binaries for macOS")
}

func copyDir(src, dst string) {
	srcInfo, err := os.Stat(src)
	if err != nil {
		fail(fmt.Sprintf("Error accessing source directory %s: %v", src, err))
	}
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		fail(fmt.Sprintf("Error creating destination directory %s: %v", dst, err))
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		fail(fmt.Sprintf("Error reading source directory %s: %v", src, err))
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			copyDir(srcPath, dstPath)
		} else {
			srcFile, err := os.Open(srcPath)
			if err != nil {
				fail(fmt.Sprintf("Error opening source file %s: %v", srcPath, err))
			}
			defer srcFile.Close()

			dstFile, err := os.Create(dstPath)
			if err != nil {
				fail(fmt.Sprintf("Error creating destination file %s: %v", dstPath, err))
			}
			defer dstFile.Close()

			if _, err := io.Copy(dstFile, srcFile); err != nil {
				fail(fmt.Sprintf("Error copying file %s to %s: %v", srcPath, dstPath, err))
			}
		}
	}
}

func generateUniversalCgoFile() {
	cgoFilename := "cgo_darwin_universal.go"
	cgoFile, err := os.Create(cgoFilename)
	if err != nil {
		fail(fmt.Sprintf("Error opening cgo file %s for writing: %v", cgoFilename, err))
	}
	defer cgoFile.Close()

	fmt.Fprintln(cgoFile, "//go:build darwin\n")
	fmt.Fprintln(cgoFile, fmt.Sprintf("package %s\n", config.ProjectName))
	fmt.Fprintln(cgoFile, "// #cgo darwin,arm64,darwin,amd64 CXXFLAGS: -std=c++11")

	// 使用arm64的配置作为基础，修改库路径为universal
	cfg := prepareBuildConfig("darwin", "arm64")
	fltkConfigPath := filepath.Join(buildDir, "build_arm64", "bin", "fltk-config")
	if err := makeExecutable(fltkConfigPath); err != nil {
		fail(fmt.Sprintf("Error making fltk-config executable: %v", err))
	}

	// Get CXX flags
	cxxFlags := getCommandOutput(fltkConfigPath, []string{"--use-gl", "--use-images", "--use-forms", "--cxxflags"})
	cxxFlags = strings.ReplaceAll(cxxFlags, cfg.CurrentDir, "${SRCDIR}")
	cxxFlags = strings.ReplaceAll(cxxFlags, cfg.IncludeDir, filepath.Join("include", cfg.OS, "universal"))
	if cfg.OS == "openbsd" {
		cxxFlags = "-I/usr/X11R6/include " + cxxFlags
	}
	fmt.Fprintf(cgoFile, "// #cgo darwin,arm64,darwin,amd64 CPPFLAGS: %s", cxxFlags)
	if cxxFlags[len(cxxFlags)-1] != '\n' {
		fmt.Fprintln(cgoFile, "")
	}

	// Get LD flags and modify for universal
	ldFlags := getCommandOutput(fltkConfigPath, []string{"--use-gl", "--use-images", "--use-forms", "--ldstaticflags"})
	ldFlags = strings.ReplaceAll(ldFlags, cfg.CurrentDir, "${SRCDIR}")
	ldFlags = strings.ReplaceAll(ldFlags, " -weak_framework", "")
	ldFlags = strings.ReplaceAll(ldFlags, cfg.LibDir, filepath.Join("lib", cfg.OS, "universal"))
	if cfg.OS == "openbsd" {
		ldFlags = "-L/usr/X11R6/lib " + ldFlags
	}
	fmt.Fprintf(cgoFile, "// #cgo darwin,arm64,darwin,amd64 LDFLAGS: %s", ldFlags)
	if ldFlags[len(ldFlags)-1] != '\n' {
		fmt.Fprintln(cgoFile, "")
	}

	fmt.Fprintln(cgoFile, "import \"C\"")
}

func validateEnvironment() {
	if runtime.GOOS == "" {
		fail("GOOS environment variable is empty")
	}
	if runtime.GOARCH == "" {
		fail("GOARCH environment variable is empty")
	}
}

func prepareBuildConfig(buildOS, buildArch string) BuildConfig {
	currentDir, err := os.Getwd()
	if err != nil {
		fail(fmt.Sprintf("Cannot get current directory: %v", err))
	}

	cfg := BuildConfig{
		OS:         buildOS,
		Arch:       buildArch,
		LibDir:     filepath.Join("lib", buildOS, buildArch),
		IncludeDir: filepath.Join("include", buildOS, buildArch),
		CurrentDir: currentDir,
		CMakeOptions: []string{
			"-DCMAKE_BUILD_TYPE=Release",
			"-DFLTK_BUILD_TEST=OFF",
			"-DFLTK_BUILD_EXAMPLES=OFF",
			"-DFLTK_BUILD_FLUID=OFF",
			"-DFLTK_BUILD_HTML_DOCS=OFF",
			"-DFLTK_BUILD_PDF_DOCS=OFF",
			"-DFLTK_BUILD_FLTK_OPTIONS=OFF",
			"-DFLTK_USE_SYSTEM_LIBJPEG=OFF",
			"-DFLTK_USE_SYSTEM_LIBPNG=OFF",
			"-DFLTK_USE_SYSTEM_ZLIB=OFF",
		},
	}

	// Add platform-specific options
	if cfg.OS == "darwin" {
		cfg.CMakeOptions = append(cfg.CMakeOptions, "-DCMAKE_OSX_DEPLOYMENT_TARGET=12.0")
		cfg.CMakeOptions = append(cfg.CMakeOptions, "-DCMAKE_OSX_SYSROOT=/Library/Developer/CommandLineTools/SDKs/MacOSX.sdk")
		if cfg.Arch == "amd64" {
			cfg.CMakeOptions = append(cfg.CMakeOptions, "-DCMAKE_OSX_ARCHITECTURES=x86_64")
		} else if runtime.GOARCH == "arm64" {
			cfg.CMakeOptions = append(cfg.CMakeOptions, "-DCMAKE_OSX_ARCHITECTURES=arm64")
		} else {
			fmt.Printf("Unsupported MacOS architecture, %s\n", cfg.Arch)
			os.Exit(1)
		}
	}

	return cfg
}

func checkRequiredTools() {
	requiredTools := []string{"git", "cmake"}
	for _, tool := range requiredTools {
		if _, err := exec.LookPath(tool); err != nil {
			fail(fmt.Sprintf("Cannot find %s binary: %v", tool, err))
		}
	}
}

func createDirectories(cfg BuildConfig) {
	dirs := []string{cfg.LibDir, cfg.IncludeDir, buildDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0750); err != nil {
			fail(fmt.Sprintf("Could not create directory %s: %v", dir, err))
		}
	}
}

func setupFLTKSource(cfg BuildConfig) {
	fltkSourceDir := filepath.Join(buildDir, sourceDir)

	if _, err := os.Stat(fltkSourceDir); errors.Is(err, fs.ErrNotExist) {
		fmt.Println("Cloning FLTK repository")
		runCommand("git", []string{"clone", fltkRepoURL}, buildDir)
	} else if stat, err := os.Stat(fltkSourceDir); err == nil && stat.IsDir() {
		fmt.Println("Found existing FLTK directory")

		if cfg.OS == "windows" {
			runCommand("git", []string{"checkout", "src/Fl_win32.cxx"}, fltkSourceDir)
		}

		runCommand("git", []string{"fetch"}, fltkSourceDir)
	} else {
		fail(fmt.Sprintf("Location for FLTK source code, %s, is not directory", fltkSourceDir))
	}

	runCommand("git", []string{"checkout", commit}, fltkSourceDir)

	if cfg.OS == "windows" {
		runCommand("git", []string{"apply", patchFile}, fltkSourceDir)
	}
}

func runCMake(cfg BuildConfig) {
	generator := "Unix Makefiles"
	if cfg.OS == "windows" {
		generator = "MinGW Makefiles"
	}
	cmakeBuildDir := "build_" + cfg.Arch
	args := []string{
		"-G", generator,
		"-S", sourceDir,
		"-B", cmakeBuildDir,
	}
	args = append(args, cfg.CMakeOptions...)
	args = append(args,
		"-DCMAKE_INSTALL_PREFIX="+cfg.CurrentDir,
		"-DCMAKE_INSTALL_INCLUDEDIR="+cfg.IncludeDir,
		"-DCMAKE_INSTALL_LIBDIR="+cfg.LibDir,
		"-DFLTK_INCLUDEDIR="+filepath.Join(cfg.CurrentDir, "include", cfg.OS, cfg.Arch),
		"-DFLTK_LIBDIR="+filepath.Join(cfg.CurrentDir, "lib", cfg.OS, cfg.Arch),
	)

	runCommand("cmake", args, buildDir)
}

func buildAndInstall(cfg BuildConfig) {
	cmakeBuildDir := "build_" + cfg.Arch
	buildArgs := []string{"--build", cmakeBuildDir, "--parallel"}
	if cfg.OS == "openbsd" {
		buildArgs = []string{"--build", cmakeBuildDir}
	}

	runCommand("cmake", buildArgs, buildDir)
	runCommand("cmake", []string{"--install", cmakeBuildDir}, buildDir)
}

func generateCgoFile(cfg BuildConfig) {
	cgoFilename := fmt.Sprintf("cgo_%s_%s.go", cfg.OS, cfg.Arch)
	cgoFile, err := os.Create(cgoFilename)
	if err != nil {
		fail(fmt.Sprintf("Error opening cgo file %s for writing: %v", cgoFilename, err))
	}
	defer cgoFile.Close()

	fmt.Fprintf(cgoFile, "//go:build %s && %s\n\n", cfg.OS, cfg.Arch)
	fmt.Fprintln(cgoFile, fmt.Sprintf("package %s\n", config.ProjectName))
	fmt.Fprintf(cgoFile, "// #cgo %s,%s CXXFLAGS: -std=c++11\n", cfg.OS, cfg.Arch)

	if cfg.OS != "windows" {
		generateUnixCgoFlags(cgoFile, cfg)
	} else {
		generateWindowsCgoFlags(cgoFile, cfg)
	}

	fmt.Fprintln(cgoFile, "import \"C\"")
}

func generateUnixCgoFlags(cgoFile *os.File, cfg BuildConfig) {
	fltkConfigPath := filepath.Join(buildDir, "build_"+cfg.Arch, "bin", "fltk-config")
	if err := makeExecutable(fltkConfigPath); err != nil {
		fail(fmt.Sprintf("Error making fltk-config executable: %v", err))
	}

	// Get CXX flags
	cxxFlags := getCommandOutput(fltkConfigPath, []string{"--use-gl", "--use-images", "--use-forms", "--cxxflags"})
	cxxFlags = strings.ReplaceAll(cxxFlags, cfg.CurrentDir, "${SRCDIR}")
	if cfg.OS == "openbsd" {
		cxxFlags = "-I/usr/X11R6/include " + cxxFlags
	}
	fmt.Fprintf(cgoFile, "// #cgo %s,%s CPPFLAGS: %s", cfg.OS, cfg.Arch, cxxFlags)
	if cxxFlags[len(cxxFlags)-1] != '\n' {
		fmt.Fprintln(cgoFile, "")
	}

	// Get LD flags
	ldFlags := getCommandOutput(fltkConfigPath, []string{"--use-gl", "--use-images", "--use-forms", "--ldstaticflags"})
	ldFlags = strings.ReplaceAll(ldFlags, cfg.CurrentDir, "${SRCDIR}")
	ldFlags = strings.ReplaceAll(ldFlags, " -weak_framework", "")
	if cfg.OS == "openbsd" {
		ldFlags = "-L/usr/X11R6/lib " + ldFlags
	}
	fmt.Fprintf(cgoFile, "// #cgo %s,%s LDFLAGS: %s", cfg.OS, cfg.Arch, ldFlags)
	if ldFlags[len(ldFlags)-1] != '\n' {
		fmt.Fprintln(cgoFile, "")
	}
}

func generateWindowsCgoFlags(cgoFile *os.File, cfg BuildConfig) {
	libDirWithSlashes := filepath.ToSlash(cfg.LibDir)
	includeDirWithSlashes := filepath.ToSlash(cfg.IncludeDir)

	fmt.Fprintf(cgoFile, "// #cgo %s,%s CPPFLAGS: -I${SRCDIR}/%s -I${SRCDIR}/%s/FL/images -D_LARGEFILE_SOURCE -D_LARGEFILE64_SOURCE -D_FILE_OFFSET_BITS=64\n",
		cfg.OS, cfg.Arch, includeDirWithSlashes, includeDirWithSlashes)
	fmt.Fprintf(cgoFile, `// #cgo %s,%s LDFLAGS: -mwindows ${SRCDIR}/%s/libfltk_images.a ${SRCDIR}/%s/libfltk_jpeg.a ${SRCDIR}/%s/libfltk_png.a ${SRCDIR}/%s/libfltk_z.a ${SRCDIR}/%s/libfltk_gl.a -lglu32 -lopengl32 ${SRCDIR}/%s/libfltk_forms.a ${SRCDIR}/%s/libfltk.a -lgdiplus -lole32 -luuid -lcomctl32 -lws2_32 -lwinspool -lmsvcrt
`, cfg.OS, cfg.Arch, libDirWithSlashes, libDirWithSlashes, libDirWithSlashes, libDirWithSlashes, libDirWithSlashes, libDirWithSlashes, libDirWithSlashes)
}

func runCommand(name string, args []string, dir string) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fail(fmt.Sprintf("Error running %s %v: %v", name, args, err))
	}
}

func getCommandOutput(name string, args []string) string {
	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	if err != nil {
		fail(fmt.Sprintf("Error running %s %v: %v", name, args, err))
	}
	return string(output)
}

func makeExecutable(path string) error {
	stat, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("checking for file: %w", err)
	}
	if !stat.Mode().IsRegular() {
		return fmt.Errorf("file is not a regular file")
	}
	return os.Chmod(path, stat.Mode().Perm()|0111)
}

func fail(message string) {
	fmt.Println(message)
	os.Exit(1)
}
