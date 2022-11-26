package web

import (
	lru "github.com/hashicorp/golang-lru"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type FileUploader struct {
	FileField   string
	DstPathFunc func(*multipart.FileHeader) string
}

func (f FileUploader) Handle() HandleFunc {

	if f.FileField == "" {
		f.FileField = "file"
	}

	if f.DstPathFunc == nil {
		// 设置默认值
	}

	return func(ctx *Context) {
		// 1.读到文件内容
		// 2.计算出目标路径
		// 3.保存文件
		// 4.返回响应
		file, fileHeader, err := ctx.Req.FormFile(f.FileField)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("上传失败" + err.Error())
			return
		}
		defer file.Close()

		// 计算目标路径
		dst := f.DstPathFunc(fileHeader)
		err = os.MkdirAll(filepath.Dir(dst), 0o666)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("上传失败" + err.Error())
			return
		}
		dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0o666)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("上传失败" + err.Error())
			return
		}
		defer dstFile.Close()

		// buf 会影响性能
		_, err = io.CopyBuffer(dstFile, file, nil)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("上传失败" + err.Error())
			return
		}
		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte("上传成功")
	}
}

//type FileUploaderOption func(uploader *FileUploader)
//
//func NewFileUploader(opts ...FileUploaderOption) *FileUploader {
//	res := &FileUploader{
//		FileField: "file",
//		DstPathFunc: func(header *multipart.FileHeader) string {
//			return filepath.Join("testdata", "upload", uuid.New().String())
//		},
//	}
//	for _, opt := range opts {
//		opt(res)
//	}
//	return res
//}

type FileDownloader struct {
	Dir string
}

func (f FileDownloader) Handle() HandleFunc {
	return func(ctx *Context) {
		// 用的是 xxx?file=xxx
		sv := ctx.QueryValue("file")
		req, err := sv.val, sv.err
		if err != nil {
			ctx.RespStatusCode = http.StatusBadRequest
			ctx.RespData = []byte("找不到目标文件")
			return
		}
		req = filepath.Clean(req)
		dst := filepath.Join(f.Dir, req)
		// 做一个校验，防止相对路径引起攻击者下载了你的系统文件
		//dst, err = filepath.Abs(dst)
		//if strings.Contains(f.Dir, req) {
		//
		//}
		fn := filepath.Base(dst)
		header := ctx.Resp.Header()
		header.Set("Content-Disposition", "attachment;filename="+fn)
		header.Set("Content-Description", "File Transfer")
		header.Set("Content-Type", "application/octet-stream")
		header.Set("Content-Transfer-Encoding", "binary")
		header.Set("Expires", "0")
		header.Set("Cache-Control", "must-revalidate")
		header.Set("Pragma", "public")

		http.ServeFile(ctx.Resp, ctx.Req, dst)
	}
}

type StaticResourceHandlerOption func(handler *StaticResourceHandler)

// StaticResourceHandler 两个层面上
// 1. 大文件不魂村
// 2. 控制住了缓存的文件的数量
// 所以，最多消耗多少内存？ size(cache) * maxSize
type StaticResourceHandler struct {
	dir                     string
	cache                   *lru.Cache
	extensionContentTypeMap map[string]string
	// 大文件不缓存
	maxSize int
}

func NewStaticResourceHandler(dir string, opts ...StaticResourceHandlerOption) (*StaticResourceHandler, error) {
	// 总共缓存 key-value
	c, err := lru.New(100 * 1024 * 1024)
	if err != nil {
		return nil, err
	}
	res := &StaticResourceHandler{
		dir:   dir,
		cache: c,
		// 10 MB，文件大小超过这个值，就不会缓存
		maxSize: 1024 * 1024 * 10,
		extensionContentTypeMap: map[string]string{
			"jpeg": "image/jpeg",
			"jpe":  "image/jpeg",
			"jpg":  "image/jpeg",
			"png":  "image/png",
			"pdf":  "image/pdf",
		},
	}

	for _, opt := range opts {
		opt(res)
	}

	return res, nil
}

func StaticWithMaxFileSize(maxSize int) StaticResourceHandlerOption {
	return func(handler *StaticResourceHandler) {
		handler.maxSize = maxSize
	}
}

func StaticWithMaxCache(c *lru.Cache) StaticResourceHandlerOption {
	return func(handler *StaticResourceHandler) {
		handler.cache = c
	}
}

func StaticWithMoreExtension(extMap map[string]string) StaticResourceHandlerOption {
	return func(handler *StaticResourceHandler) {
		for ext, contentType := range extMap {
			handler.extensionContentTypeMap[ext] = contentType
		}
	}
}

func (s *StaticResourceHandler) Handle(ctx *Context) {
	// 1.拿到目标文件名
	// 2.定位到目标文件，并且读出来
	// 3.返回给前端

	sv := ctx.PathValue("file")
	file, err := sv.val, sv.err
	if err != nil {
		ctx.RespStatusCode = http.StatusBadRequest
		ctx.RespData = []byte("请求路径不对")
		return
	}
	dst := filepath.Join(s.dir, file)
	ext := filepath.Ext(dst)[1:]
	header := ctx.Resp.Header()
	contentType := s.extensionContentTypeMap[ext][1:]

	if data, ok := s.cache.Get(file); ok {
		header.Set("Content-Type", contentType)
		header.Set("Content-Length", strconv.Itoa(len(data.([]byte))))
		ctx.RespData = data.([]byte)
		ctx.RespStatusCode = http.StatusOK
	}

	data, err := os.ReadFile(dst)
	if err != nil {
		ctx.RespStatusCode = http.StatusInternalServerError
		ctx.RespData = []byte("服务器错误")
		return
	}
	if len(data) <= s.maxSize {
		s.cache.Add(file, data)
	}
	header.Set("Content-Type", contentType)
	header.Set("Content-Length", strconv.Itoa(len(data)))
	ctx.RespData = data
	ctx.RespStatusCode = http.StatusOK
}
