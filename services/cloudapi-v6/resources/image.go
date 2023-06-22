package resources

import (
	"bufio"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"path/filepath"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"

	"github.com/fatih/structs"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/kardianos/ftps"
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
	Username          string
	Password          string
}

// ImagesService is a wrapper around ionoscloud.Image
type ImagesService interface {
	Upload(ctx context.Context, properties UploadProperties) error
	List(params ListQueryParams) (Images, *Response, error)
	Get(imageId string, params QueryParams) (*Image, *Response, error)
	Update(imageId string, imgProp ImageProperties, params QueryParams) (*Image, *Response, error)
	Delete(imageId string, params QueryParams) (*Response, error)
}

type imagesService struct {
	client  *ionoscloud.APIClient
	context context.Context
}

var _ ImagesService = &imagesService{}

func NewImageService(client *client.Client, ctx context.Context) ImagesService {
	return &imagesService{
		client:  client.CloudClient,
		context: ctx,
	}
}

func (s *imagesService) Upload(ctx context.Context, p UploadProperties) error {
	tlsConfig := tls.Config{
		InsecureSkipVerify: p.SkipVerify,
		ServerName:         p.Url,
		RootCAs:            p.ServerCertificate,
		MaxVersion:         tls.VersionTLS12,
	}

	dialOptions := ftps.DialOptions{
		Host:        p.Url,
		Port:        p.Port,
		Username:    p.Username,
		Passowrd:    p.Password,
		ExplicitTLS: true,
		TLSConfig:   &tlsConfig,
	}

	c, err := ftps.Dial(ctx, dialOptions)
	if err != nil {
		return fmt.Errorf("dialing FTP server failed. Are you using the correct username & password?: %w", err)
	}

	err = c.Chdir(filepath.Dir(p.Path))
	if err != nil {
		return fmt.Errorf("failed while changing directory within FTP server: %w", err)
	}

	files, err := c.List(ctx)
	if err != nil {
		return fmt.Errorf("failed while listing files within FTP server: %w", err)
	}

	// Check if there already exists an image with the given name at the location
	desiredFileName := filepath.Base(p.Path)
	var errExists error
	for _, f := range files {
		if f.Name == desiredFileName {
			errExists = fmt.Errorf("%s already exists at %s", desiredFileName, p.Url)
		}
	}

	err = c.Upload(ctx, desiredFileName, p.DataBuffer)
	if err != nil {
		err = fmt.Errorf("failed while uploading %s to FTP server: %w", desiredFileName, err)
		if errExists != nil {
			err = fmt.Errorf("%w: Note, your error could be because of: %w", err, errExists)
		}
		return err
	}

	return c.Close()
}

func (s *imagesService) List(params ListQueryParams) (Images, *Response, error) {
	req := s.client.ImagesApi.ImagesGet(s.context)
	if !structs.IsZero(params) {
		if params.Filters != nil {
			for k, v := range *params.Filters {
				for _, val := range v {
					req = req.Filter(k, val)
				}
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
