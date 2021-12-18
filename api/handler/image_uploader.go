package handler

import (
	"bytes"
	"github.com/google/uuid"
	"github.com/xxarupakaxx/grpc-api/api/gen/api"
	"io"
	"net/http"
	"sync"
)

type ImageUploadHandler struct {
	sync.Mutex
	files map[string][]byte
}

func (i *ImageUploadHandler) Upload(stream api.ImageUploadService_UploadServer) error {
	req,err := stream.Recv()
	if err != nil {
		return err
	}

	meta := req.GetFileMeta()
	filename := meta.Filename

	u,err := uuid.NewRandom()
	if err != nil {
		return err
	}
	uuid := u.String()

	buf := &bytes.Buffer{}

	for true {
		r,err := stream.Recv()
		if err == io.EOF {
			break
		}else if err != nil {
			return err
		}

		chunk := r.GetData()
		_,err = buf.Write(chunk)
		if err != nil {
			return err
		}
	}

	data := buf.Bytes()
	mimeType := http.DetectContentType(data)

	i.files[filename] = data

	err = stream.SendAndClose(&api.ImageUploadResponse{
		Uuid:        uuid,
		Size:        int32(len(data)),
		ContentType: mimeType,
		Filename:    filename,
	})

	return err
}

func NewImageUploadHandler() *ImageUploadHandler {
	return &ImageUploadHandler{files: make(map[string][]byte)}
}


