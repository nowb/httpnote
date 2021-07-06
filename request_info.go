package httpnote

import (
	"bytes"
	"encoding/base64"
	"io"
	"mime/multipart"
	"net/http"
)

type RequestInfo struct {
	Method           string
	URL              *URLInfo
	Proto            string
	ProtoMajor       int
	ProtoMinor       int
	Header           map[string][]string
	BodyBytes        []byte
	BodyBase64       string
	ContentLength    int64
	TransferEncoding []string
	Close            bool
	Host             string
	Form             map[string][]string
	FormEncoded      string
	PostForm         map[string][]string
	PostFormEncoded  string
	MultipartForm    *FormInfo
	Trailer          map[string][]string
	RemoteAddr       string
	RequestURI       string
}

type URLInfo struct {
	Scheme      string
	Opaque      string
	User        *UserInfo
	Host        string
	Path        string
	RawPath     string
	ForceQuery  bool
	RawQuery    string
	Fragment    string
	RawFragment string
}

type UserInfo struct {
	Username string
	Password string
	Encoded  string
}

type FormInfo struct {
	Value map[string][]string
	File  map[string][]*FileHeaderInfo
}

type FileHeaderInfo struct {
	Filename   string
	MimeHeader map[string][]string
	Size       int64
	FileBytes  []byte
	FileBase64 string
}

func MapRequest(r *http.Request) *RequestInfo {
	if r == nil {
		return nil
	}

	ri := &RequestInfo{
		Method:           r.Method,
		Proto:            r.Proto,
		ProtoMajor:       r.ProtoMajor,
		ProtoMinor:       r.ProtoMinor,
		ContentLength:    r.ContentLength,
		TransferEncoding: r.TransferEncoding,
		Close:            r.Close,
		Host:             r.Host,
		RemoteAddr:       r.RemoteAddr,
		RequestURI:       r.RequestURI,
	}

	// set Header
	if len(r.Header) > 0 {
		ri.Header = r.Header.Clone()
	}

	// set Trailer
	if len(r.Trailer) > 0 {
		ri.Trailer = r.Trailer.Clone()
	}

	// set URL
	if r.URL != nil {
		ri.URL = &URLInfo{
			Scheme:      r.URL.Scheme,
			Opaque:      r.URL.Opaque,
			Host:        r.URL.Host,
			Path:        r.URL.Path,
			RawPath:     r.URL.RawPath,
			ForceQuery:  r.URL.ForceQuery,
			RawQuery:    r.URL.RawQuery,
			Fragment:    r.URL.Fragment,
			RawFragment: r.URL.RawFragment,
		}

		if user := r.URL.User; user != nil {
			ri.URL.User = &UserInfo{
				Username: user.Username(),
				Encoded:  user.String(),
			}

			if password, ok := user.Password(); ok {
				ri.URL.User.Password = password
			}
		}
	}

	// set BodyBytes and BodyBase64
	bodyBytes, err := io.ReadAll(r.Body)
	if err == nil && len(bodyBytes) > 0 {
		ri.BodyBytes = bodyBytes
		ri.BodyBase64 = base64.StdEncoding.EncodeToString(bodyBytes)

		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	_ = r.Body.Close()

	// set Form
	if len(r.Form) > 0 {
		ri.Form = r.Form
		ri.FormEncoded = r.Form.Encode()
	}

	// set PostForm
	if len(r.PostForm) > 0 {
		ri.PostForm = r.PostForm
		ri.PostFormEncoded = r.PostForm.Encode()
	}

	// set MultipartForm
	if r.MultipartForm != nil {
		ri.MultipartForm = MapForm(r.MultipartForm)
	}

	return ri
}

func MapFileHeader(fh *multipart.FileHeader) *FileHeaderInfo {
	if fh == nil {
		return nil
	}

	fhi := &FileHeaderInfo{
		Filename:   fh.Filename,
		MimeHeader: fh.Header,
		Size:       fh.Size,
	}

	// read file
	file, err := fh.Open()
	if err == nil {
		fileBytes, err := io.ReadAll(file)
		if err == nil {
			fhi.FileBytes = fileBytes
			fhi.FileBase64 = base64.StdEncoding.EncodeToString(fileBytes)
		}
	}

	return fhi
}

func MapForm(f *multipart.Form) *FormInfo {
	if f == nil {
		return nil
	}

	var fi FormInfo

	if len(f.Value) > 0 {
		fi.Value = f.Value
	}

	if len(f.File) > 0 {
		fileMap := make(map[string][]*FileHeaderInfo, len(f.File))

		for k, v := range f.File {
			var fhiSlice []*FileHeaderInfo

			for _, fh := range v {
				if fh != nil {
					fhiSlice = append(fhiSlice, MapFileHeader(fh))
				}
			}

			fileMap[k] = fhiSlice
		}

		fi.File = fileMap
	}

	return &fi
}
