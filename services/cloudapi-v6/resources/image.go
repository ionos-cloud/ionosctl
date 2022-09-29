package resources

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/jlaffaye/ftp"
	"github.com/spf13/viper"
	"io"
	"path/filepath"
	"time"
)

type Image struct {
	ionoscloud.Image
}

type Images struct {
	ionoscloud.Images
}

type ImageProperties struct {
	ionoscloud.ImageProperties
}

// UploadProperties contains info needed to initialize an FTP connection to IONOS server and upload an image.
type UploadProperties struct {
	ImageFileProperties
	FTPServerProperties
}

type ImageFileProperties struct {
	Path   string // File name, server path (not local) and file extension included
	DataIO io.Reader
}
type FTPServerProperties struct {
	Url  string // Server URL without any directory path. Example: ftp-fkb.ionos.com
	Port int
}

// ImagesService is a wrapper around ionoscloud.Image
type ImagesService interface {
	Upload(properties UploadProperties) error
	List(params ListQueryParams) (Images, *Response, error)
	Get(imageId string, params QueryParams) (*Image, *Response, error)
	Update(imageId string, imgProp ImageProperties, params QueryParams) (*Image, *Response, error)
	Delete(imageId string, params QueryParams) (*Response, error)
}

type imagesService struct {
	client  *Client
	context context.Context
}

var _ ImagesService = &imagesService{}

func NewImageService(client *Client, ctx context.Context) ImagesService {
	return &imagesService{
		client:  client,
		context: ctx,
	}
}

func (s *imagesService) Upload(p UploadProperties) error {
	// Uncomment for FTPS lib
	//conn := new(ftps.FTPS)
	//
	//conn.TLSConfig.InsecureSkipVerify = true // often necessary in shared hosting environments
	//conn.Debug = false
	//
	//err := conn.Connect(p.Url, p.Port)
	//if err != nil {
	//	return err
	//}
	//
	//err = conn.Login(viper.GetString(config.Username), viper.GetString(config.Password))
	//if err != nil {
	//	return err
	//}
	//
	//err = conn.ChangeWorkingDirectory(filepath.Dir(p.Path))
	//if err != nil {
	//	return err
	//}
	//
	//// TODO: Large uploads fail. Try buffering data, or changing timeout somehow: StoreFile -> net.Write -> net.SetDeadline
	//err = conn.StoreFile(filepath.Base(p.Path), p.Data)
	//if err != nil {
	//	return err
	//}
	//
	//err = conn.Quit()
	//if err != nil {
	//	return err
	//}

	tlsConfig := tls.Config{
		InsecureSkipVerify: true, // TODO: INSECURE. Change this before prod! Client susceptible to "Man-in-the-middle" attacks.
	}

	c, err := ftp.Dial(fmt.Sprintf("%s:%d", p.Url, p.Port), ftp.DialWithTimeout(30*time.Second), ftp.DialWithExplicitTLS(&tlsConfig))
	if err != nil {
		return err
	}
	fmt.Printf("Connected to %s\n", p.Url)

	err = c.Login(viper.GetString(config.Username), viper.GetString(config.Password))
	if err != nil {
		return err
	}
	fmt.Printf("Logged in;\n")

	// Do something with the FTP conn
	err = c.ChangeDir(filepath.Dir(p.Path))
	if err != nil {
		fmt.Printf("Failed to change to %s\n", filepath.Dir(p.Path))
		return err
	}
	fmt.Printf("Dir changed to %s\n", filepath.Dir(p.Path))

	err = c.Stor(filepath.Base(p.Path), p.DataIO)
	if err != nil {
		fmt.Printf("Failed uploading %s to %s!\n", p.Path, p.Url)
		return err
	}
	fmt.Printf("Uploaded %s to %s!\n", p.Path, p.Url)

	if err := c.Quit(); err != nil {
		return err
	}

	return nil
}

func (s *imagesService) List(params ListQueryParams) (Images, *Response, error) {
	req := s.client.ImagesApi.ImagesGet(s.context)
	if !structs.IsZero(params) {
		if params.Filters != nil {
			for k, v := range *params.Filters {
				req = req.Filter(k, v)
			}
		}
		if params.OrderBy != nil {
			req = req.OrderBy(*params.OrderBy)
		}
		if params.MaxResults != nil {
			req = req.MaxResults(*params.MaxResults)
		}
		if !structs.IsZero(params.QueryParams) {
			if params.QueryParams.Depth != nil {
				req = req.Depth(*params.QueryParams.Depth)
			}
			if params.QueryParams.Pretty != nil {
				// Currently not implemented
				req = req.Pretty(*params.QueryParams.Pretty)
			}
		}
	}
	images, resp, err := s.client.ImagesApi.ImagesGetExecute(req)
	return Images{images}, &Response{*resp}, err
}

func (s *imagesService) Get(imageId string, params QueryParams) (*Image, *Response, error) {
	req := s.client.ImagesApi.ImagesFindById(s.context, imageId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	image, resp, err := s.client.ImagesApi.ImagesFindByIdExecute(req)
	return &Image{image}, &Response{*resp}, err
}

func (s *imagesService) Update(imageId string, imgProp ImageProperties, params QueryParams) (*Image, *Response, error) {
	req := s.client.ImagesApi.ImagesPatch(s.context, imageId).Image(imgProp.ImageProperties)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	image, resp, err := s.client.ImagesApi.ImagesPatchExecute(req)
	return &Image{image}, &Response{*resp}, err
}

func (s *imagesService) Delete(imageId string, params QueryParams) (*Response, error) {
	req := s.client.ImagesApi.ImagesDelete(s.context, imageId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	resp, err := s.client.ImagesApi.ImagesDeleteExecute(req)
	return &Response{*resp}, err
}
