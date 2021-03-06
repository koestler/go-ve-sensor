// Package sample is a sample server driver
package ftpServer

import (
	"crypto/tls"
	"errors"
	"fmt"
	"os"
	"log"

	"github.com/fclairamb/ftpserver/server"
	"strings"
	"time"
	"sort"
	"path"
	"github.com/koestler/go-ve-sensor/config"
	"github.com/koestler/go-ve-sensor/storage"
	"strconv"
)

type VirtualFileSystem struct {
	directories *DirectoryList
	files       *FileList
}

// MainDriver defines a very basic ftpServer driver
type MainDriver struct {
	vfs        VirtualFileSystem
	listenHost string
	listenPort int
	cameras    []*config.FtpCameraConfig
}

// ClientDriver defines a very basic client driver
type ClientDriver struct {
	vfs    VirtualFileSystem
	device *storage.Device
}

// NewSampleDriver creates a sample driver
func NewDriver(listenHost string, listenPort int, cameras []*config.FtpCameraConfig) (*MainDriver, error) {
	// create new virtual in-memory filesystem
	driver := &MainDriver{
		vfs: VirtualFileSystem{
			directories: NewDirectoryList(),
			files:       NewFileList(),
		},
		listenHost: listenHost,
		listenPort: listenPort,
		cameras:    cameras,
	}

	return driver, nil
}

// GetSettings returns some general settings around the server setup
func (driver *MainDriver) GetSettings() (*server.Settings, error) {
	var settings server.Settings

	settings.ListenAddr = driver.listenHost + ":" + strconv.Itoa(driver.listenPort)
	settings.PublicHost = "::1"
	settings.DataPortRange = &server.PortRange{Start: 2122, End: 2200}
	settings.DisableMLSD = true
	settings.NonStandardActiveDataPort = false

	return &settings, nil
}

// GetTLSConfig returns a TLS Certificate to use
func (driver *MainDriver) GetTLSConfig() (*tls.Config, error) {
	return nil, errors.New("tls not supported")
}

// WelcomeUser is called to send the very first welcome message
func (driver *MainDriver) WelcomeUser(cc server.ClientContext) (string, error) {
	log.Printf("ftpcam-driver: WelcomeUser cc.ID=%v", cc.ID())

	cc.SetDebug(true)
	return fmt.Sprintf(
		"Welcome on go-ve-sensor ftpServer, your ID is %d, your IP:port is %s", cc.ID(), cc.RemoteAddr(),
	), nil
}

// AuthUser authenticates the user and selects an handling driver
func (driver *MainDriver) AuthUser(cc server.ClientContext, user, pass string) (server.ClientHandlingDriver, error) {
	log.Printf("ftpcam-driver: AuthUser cc.ID=%v", cc.ID())

	for _, camera := range driver.cameras {
		if camera.User == user && camera.Password == pass {
			device, err := storage.GetByName(camera.Name)
			if err != nil {
				return nil, err
			}

			return &ClientDriver{
				vfs:    driver.vfs,
				device: device,
			}, nil
		}
	}

	return nil, errors.New("bad username or password")
}

// UserLeft is called when the user disconnects, even if he never authenticated
func (driver *MainDriver) UserLeft(cc server.ClientContext) {
	log.Printf("ftpcam-driver: UserLeft cc.ID=%v", cc.ID())
}

// ChangeDirectory changes the current working directory
func (driver *ClientDriver) ChangeDirectory(cc server.ClientContext, directory string) error {
	log.Printf("ftpcam-driver: ChangeDirectory cc.ID=%v directory=%v", cc.ID(), directory)

	// create directories on the fly
	driver.vfs.directories.Create(path.Clean(directory))
	return nil
}

// MakeDirectory creates a directory
func (driver *ClientDriver) MakeDirectory(cc server.ClientContext, directory string) error {
	log.Printf("ftpcam-driver: MakeDirectory, cc.ID=%v directory=%v", cc.ID(), directory)
	driver.vfs.directories.Create(path.Clean(directory))
	return nil;
}

// ListFiles lists the files of a directory
func (driver *ClientDriver) ListFiles(cc server.ClientContext) ([]os.FileInfo, error) {
	log.Printf("ftpcam-driver: ListFiles cc.ID=%v cc.Path=%v", cc.ID(), cc.Path())

	dirPath := getDirPath(cc.Path())

	files := make([]os.FileInfo, 0)
	for directory := range driver.vfs.directories.Iterate() {
		if !strings.HasPrefix(directory, dirPath) {
			continue
		}

		reminder := directory[len(dirPath):]
		if len(reminder) < 1 || strings.Contains(reminder, "/") {
			// subdir -> ignore
			continue
		}

		files = append(files,
			VirtualFileInfo{
				name:     reminder,
				mode:     os.FileMode(0666) | os.ModeDir,
				size:     4096,
				modified: time.Now(),
			},
		)
	}

	fileList := driver.vfs.files.getFilesInsidePath(dirPath)
	for _, file := range fileList {
		files = append(files, file.getFileInfo(file.filePath[len(dirPath):]))
	}

	return files, nil
}

func getDirPath(dirPath string) string {
	dirPath = path.Clean(dirPath)
	if len(dirPath) > 1 {
		dirPath += "/"
	}
	return dirPath
}

func (fl *FileList) getFilesInsidePath(dirPath string) (ret []*VirtualFile) {
	ret = make([]*VirtualFile, 0, fl.Length())

	for item := range fl.Iterate() {
		filePath := item.Path
		file := item.Value

		if !strings.HasPrefix(filePath, dirPath) {
			continue
		}

		reminder := filePath[len(dirPath):]
		if len(reminder) < 1 || strings.Contains(reminder, "/") {
			// subdir -> ignore
			continue
		}
		ret = append(ret, file)
	}
	sort.Sort(VirtualFileByCreated(ret))
	return
}

func (vf *VirtualFile) getFileInfo(name string) VirtualFileInfo {
	return VirtualFileInfo{
		name:     name,
		mode:     os.FileMode(0666),
		size:     vf.Size(),
		modified: vf.modified,
	}
}

// OpenFile opens a file in 3 possible modes: read, write, appending write (use appropriate flags)
func (driver *ClientDriver) OpenFile(cc server.ClientContext, filePath string, flag int) (server.FileStream, error) {
	//log.Printf("ftpcam-driver: OpenFile cc.ID=%v filePath=%v flag=%v", cc.ID(), filePath, flag)

	// cleanup filesystem
	driver.vfs.pathRetention(getDirPath(path.Dir(filePath)))

	// If we are writing and we are not in append mode, we should remove the file
	if (flag & os.O_WRONLY) != 0 {
		flag |= os.O_CREATE
		if (flag & os.O_APPEND) == 0 {
			driver.vfs.files.Delete(filePath)
		}
	}

	if (flag & os.O_CREATE) != 0 {
		driver.vfs.files.Create(
			filePath,
			&VirtualFile{
				device:   driver.device,
				filePath: filePath,
				modified: time.Now(),
			},
		)
	}

	file, ok := driver.vfs.files.Get(filePath)
	if !ok {
		return nil, os.ErrNotExist
	}

	return file, nil
}

// GetFileInfo gets some info around a file or a directory
func (driver *ClientDriver) GetFileInfo(cc server.ClientContext, path string) (os.FileInfo, error) {
	log.Printf("ftpcam-driver: GetFileInfo cc.ID=%v path=%v", cc.ID(), path)

	if file, ok := driver.vfs.files.Get(path); ok {
		return file.getFileInfo(path), nil
	} else if ok := driver.vfs.directories.Exists(path); !ok {
		return &VirtualFileInfo{
			name:     path,
			mode:     os.FileMode(0666) | os.ModeDir,
			size:     4096,
			modified: time.Now(),
		}, nil
	}

	return nil, os.ErrNotExist
}

// CanAllocate gives the approval to allocate some data
func (driver *ClientDriver) CanAllocate(cc server.ClientContext, size int) (bool, error) {
	log.Printf("ftpcam-driver: CanAllocate cc.ID=%v size=%v", cc.ID(), size)
	return true, nil
}

// ChmodFile changes the attributes of the file
func (driver *ClientDriver) ChmodFile(cc server.ClientContext, path string, mode os.FileMode) error {
	log.Printf("ftpcam-driver: ChmodFile cc.ID=%v path=%v, mode=%v", cc.ID(), path, mode)
	return os.ErrPermission
}

// DeleteFile deletes a file or a directory
func (driver *ClientDriver) DeleteFile(cc server.ClientContext, path string) error {
	log.Printf("ftpcam-driver: DeleteFile cc.ID=%v path=%v", cc.ID(), path)
	return os.ErrPermission
}

// RenameFile renames a file or a directory
func (driver *ClientDriver) RenameFile(cc server.ClientContext, from, to string) error {
	log.Printf("ftpcam-driver: RenameFile cc.ID=%v from=%v to=%v", cc.ID(), from, to)
	return os.ErrPermission
}

func (vfs *VirtualFileSystem) pathRetention(dirPath string) {
	// get file list ordered by modified asc
	fileList := vfs.files.getFilesInsidePath(dirPath)

	// delete all but last 5 files
	for i := 0; i <= len(fileList)-5; i++ {
		//log.Printf("ftpcam-driver: virtualFileSystem cleanup filePath=%v", fileList[i].filePath)
		vfs.files.Delete(fileList[i].filePath)
	}
}
