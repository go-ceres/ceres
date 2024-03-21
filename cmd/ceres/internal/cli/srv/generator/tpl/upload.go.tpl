package service

import (
	"github.com/go-ceres/ceres/pkg/common/config"
    "github.com/go-ceres/ceres/pkg/common/errors"
    "github.com/go-ceres/ceres/pkg/transport/http"
	"fmt"
    "crypto/md5"
    "github.com/google/uuid"
    "path"
	"os"
	"strings"
	"path/filepath"
)

type Options struct {
    UploadPath string `json:"uploadPath"` // 上传路径
}

type UploadService struct {
    options *Options
}


// NewUploadService 创建服务
func NewUploadService() *UploadService {
	opts := &Options{
		UploadPath: "./static",
    }
    if err := config.Get("application.upload").Scan(opts); err != nil {
        panic(err)
    }
    // 如果目录不存在，则创建
    exists, _ := pathExists(opts.UploadPath)
    if !exists {
        err := createPath(opts.UploadPath)
        if err != nil {
            panic(err)
        }
    }
	return &UploadService{
		options: opts,
    }
}

// RegisterService 注册路由
func (u *UploadService) RegisterServer(srv *http.Server)  {
    srv.POST("/upload", u.Upload)
    srv.Static("/upload/*filepath", u.options.UploadPath)
}

// Upload 上传方法
func (u *UploadService) Upload(ctx *http.Context) error {
    form, err := ctx.Request().MultipartForm()
    if err != nil {
        return err
    }
    // 设置正常数据响应
    filenames := make([]string, 0)
    ip := ctx.GetRequestHeader("X-Real-IP")
    uErrors := make([]error, 0)
    if fileheaders, ok := form.File["file"]; ok {
        for _, fileheader := range fileheaders {
            filename := strings.ToLower(fileheader.Filename)
            f, err := fileheader.Open()
            if err != nil {
                uErrors = append(uErrors, err)
                continue
            }
            filebuf := make([]byte, fileheader.Size)
            _, err = f.Read(filebuf)
            if err != nil {
                uErrors = append(uErrors, err)
                continue
            }
            ext := filepath.Ext(filename)
            filename = encodeMd5(ip+"_"+filename)+uuid.NewString() + ext
            file, err := os.Create(path.Join(u.options.UploadPath, filename))
            if err != nil {
                uErrors = append(uErrors, err)
                continue
            }
            defer file.Close()
            _, err = file.Write(filebuf)
            if err != nil {
                uErrors = append(uErrors, err)
                continue
            }
            filenames = append(filenames, filename)
        }
    }
    if len(uErrors) > 0 {
        return uErrors[0]
    } else {
        if len(filenames) == 1 {
            _ = ctx.Result(http.StatusOK, filenames[0])
        } else if len(filenames) > 1 {
            _ = ctx.Result(http.StatusOK, filenames)
        } else {
            return errors.New(500,"UPLOAD_FILE","文件上传失败")
        }
    }

    return nil
}

// pathExists 检查指定文件夹是否存在
func pathExists(path string) (bool, error) {
    abs, err := filepath.Abs(path)
    if err != nil {
        return false, err
    }
    ext := filepath.Ext(path)
    if len(ext) != 0 {
        return false, fmt.Errorf(`the path "%s" is not file path`, path)
    }
    _, err = os.Stat(abs)
    if os.IsNotExist(err) {
        return false, nil
    } else {
        return true, nil
    }
}

// createPath 创建文件夹
func createPath(paths ...string) error {
    for _, v := range paths {
        exist, err := pathExists(v)
        if err != nil {
            return err
        }
        if !exist {
            if err = os.MkdirAll(v, os.ModePerm); err != nil {
                return err
            }
        }
    }
    return nil
}

// encodeMd5 md5加密
func encodeMd5(str string) string {
    return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

