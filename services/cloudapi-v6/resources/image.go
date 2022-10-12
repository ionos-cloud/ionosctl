package resources

import (
	"bufio"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/kardianos/ftps"
	"github.com/spf13/viper"
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
	Path       string // File name, server path (not local) and file extension included
	DataBuffer *bufio.Reader
}
type FTPServerProperties struct {
	Url               string // Server URL without any directory path. Example: ftp-fkb.ionos.com
	Port              int
	SkipVerify        bool           // Skip FTP server certificate verification. WARNING man-in-the-middle attack possible
	ServerCertificate *x509.CertPool // If FTP server uses self signed certificates, put this in tlsConfig. IONOS FTP Servers in prod DON'T need this
	Timeout           int            // Timeout in seconds
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
	tlsConfig := tls.Config{
		InsecureSkipVerify: p.SkipVerify,
		ServerName:         p.Url,
		RootCAs:            p.ServerCertificate,
		MaxVersion:         tls.VersionTLS12,
	}

	dialOptions := ftps.DialOptions{
		Host:        p.Url,
		Port:        p.Port,
		Username:    viper.GetString(config.Username),
		Passowrd:    viper.GetString(config.Password),
		ExplicitTLS: true,
		TLSConfig:   &tlsConfig,
	}

	ctx := context.Background()
	if s.context != nil {
		ctx = s.context
	}
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Duration(p.Timeout)*time.Second))
	defer cancel()

	c, err := ftps.Dial(ctx, dialOptions)
	if err != nil {
		return err
	}
	fmt.Printf("Connected to %s\n", p.Url)

	err = c.Chdir(filepath.Dir(p.Path))
	if err != nil {
		fmt.Printf("Failed to change to %s\n", filepath.Dir(p.Path))
		return err
	}

	files, err := c.List(ctx)
	if err != nil {
		fmt.Println("Failed to list")
		return err
	}

	// Check if there already exists an image with the given name at the location
	desiredFileName := filepath.Base(p.Path)
	for _, f := range files {
		if f.Name == desiredFileName {
			//err := c.RemoveFile(desiredFileName)
			//if err != nil {
			//	return err
			//}
			return fmt.Errorf("%s already exists at %s", desiredFileName, p.Url)
		}
	}

	err = c.Upload(ctx, desiredFileName, p.DataBuffer)
	if err != nil {
		fmt.Printf("Failed uploading %s to %s!\n", filepath.Base(p.Path), p.Url)
		return err
	}
	fmt.Printf("Uploaded %s to %s!\n", filepath.Base(p.Path), p.Url)

	return c.Close()
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
