package httpnote

import (
	"bytes"
	"encoding/base64"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/rs/zerolog"
)

type RequestInfo struct {
	Method           string
	URL              *URLInfo
	Proto            string
	ProtoMajor       int
	ProtoMinor       int
	Header           map[string][]string
	Body             []byte
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

func (ri *RequestInfo) MarshalZerologObject(e *zerolog.Event) {
	e.Str("method", ri.Method)

	if ri.URL != nil {
		e.Object("url", ri.URL)
	}

	e.
		Str("proto", ri.Proto).
		Int("proto_major", ri.ProtoMajor).
		Int("proto_minor", ri.ProtoMinor)

	if len(ri.Header) > 0 {
		header := zerolog.Dict()
		for k, v := range ri.Header {
			header.Strs(k, v)
		}

		e.Dict("header", header)
	}

	if len(ri.Body) > 0 {
		e.Bytes("body", ri.Body)
	}

	if ri.BodyBase64 != "" {
		e.Str("body_base64", ri.BodyBase64)
	}

	if ri.ContentLength > 0 {
		e.Int64("content_length", ri.ContentLength)
	}

	if len(ri.TransferEncoding) > 0 {
		e.Strs("transfer_encoding", ri.TransferEncoding)
	}

	if ri.Close {
		e.Bool("close", ri.Close)
	}

	e.Str("host", ri.Host)

	if len(ri.Form) > 0 {
		form := zerolog.Dict()
		for k, v := range ri.Form {
			form.Strs(k, v)
		}

		e.Dict("form", form)
	}

	if ri.FormEncoded != "" {
		e.Str("form_encoded", ri.FormEncoded)
	}

	if len(ri.PostForm) > 0 {
		postForm := zerolog.Dict()
		for k, v := range ri.PostForm {
			postForm.Strs(k, v)
		}

		e.Dict("post_form", postForm)
	}

	if ri.PostFormEncoded != "" {
		e.Str("post_form_encoded", ri.PostFormEncoded)
	}

	if ri.MultipartForm != nil {
		e.Object("multipart_form", ri.MultipartForm)
	}

	if len(ri.Trailer) > 0 {
		trailer := zerolog.Dict()
		for k, v := range ri.Trailer {
			trailer.Strs(k, v)
		}

		e.Dict("trailer", trailer)
	}

	e.
		Str("remote_addr", ri.RemoteAddr).
		Str("request_uri", ri.RequestURI)
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

func (ui *URLInfo) MarshalZerologObject(e *zerolog.Event) {
	if ui.Scheme != "" {
		e.Str("scheme", ui.Scheme)
	}

	if ui.Opaque != "" {
		e.Str("opaque", ui.Opaque)
	}

	if ui.User != nil {
		e.Object("user", ui.User)
	}

	if ui.Host != "" {
		e.Str("host", ui.Host)
	}

	if ui.Path != "" {
		e.Str("path", ui.Path)
	}

	if ui.RawPath != "" {
		e.Str("raw_path", ui.RawPath)
	}

	if ui.ForceQuery {
		e.Bool("force_query", ui.ForceQuery)
	}

	if ui.RawQuery != "" {
		e.Str("raw_query", ui.RawQuery)
	}

	if ui.Fragment != "" {
		e.Str("fragment", ui.Fragment)
	}

	if ui.RawFragment != "" {
		e.Str("raw_fragment", ui.RawFragment)
	}
}

type UserInfo struct {
	Username string
	Password string
	Encoded  string
}

func (ui *UserInfo) MarshalZerologObject(e *zerolog.Event) {
	if ui.Username != "" {
		e.Str("username", ui.Username)
	}

	if ui.Password != "" {
		e.Str("password", ui.Password)
	}

	if ui.Encoded != "" {
		e.Str("encoded", ui.Encoded)
	}
}

type FormInfo struct {
	Value map[string][]string
	File  map[string]FileHeaderInfos
}

type FileHeaderInfos []*FileHeaderInfo

func (fhis FileHeaderInfos) MarshalZerologArray(a *zerolog.Array) {
	for _, fhi := range fhis {
		a.Object(fhi)
	}
}

func (fi *FormInfo) MarshalZerologObject(e *zerolog.Event) {
	if len(fi.Value) > 0 {
		value := zerolog.Dict()
		for k, v := range fi.Value {
			value.Strs(k, v)
		}

		e.Dict("value", value)
	}

	if len(fi.File) > 0 {
		file := zerolog.Dict()
		for k, v := range fi.File {
			file.Array(k, v)
		}

		e.Dict("file", file)
	}
}

type FileHeaderInfo struct {
	Filename   string
	MimeHeader map[string][]string
	Size       int64
	File       []byte
	FileBase64 string
}

func (fhi *FileHeaderInfo) MarshalZerologObject(e *zerolog.Event) {
	if fhi.Filename != "" {
		e.Str("filename", fhi.Filename)
	}

	if len(fhi.MimeHeader) > 0 {
		mimeHeader := zerolog.Dict()
		for k, v := range fhi.MimeHeader {
			mimeHeader.Strs(k, v)
		}

		e.Dict("mime_header", mimeHeader)
	}

	if fhi.Size != 0 {
		e.Int64("size", fhi.Size)
	}

	if len(fhi.File) > 0 {
		e.Bytes("file", fhi.File)
	}

	if fhi.FileBase64 != "" {
		e.Str("file_base64", fhi.FileBase64)
	}
}

func MapRequest(r *http.Request, encodeBytes bool) *RequestInfo {
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

	// set Body and BodyBase64
	bodyBytes, err := io.ReadAll(r.Body)
	if err == nil && len(bodyBytes) > 0 {
		if !encodeBytes {
			ri.Body = bodyBytes
		} else {
			ri.BodyBase64 = base64.StdEncoding.EncodeToString(bodyBytes)
		}

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
		ri.MultipartForm = MapForm(r.MultipartForm, encodeBytes)
	}

	return ri
}

func MapFileHeader(fh *multipart.FileHeader, encodeBytes bool) *FileHeaderInfo {
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
			if !encodeBytes {
				fhi.File = fileBytes
			} else {
				fhi.FileBase64 = base64.StdEncoding.EncodeToString(fileBytes)
			}
		}
	}

	return fhi
}

func MapForm(f *multipart.Form, encodeBytes bool) *FormInfo {
	if f == nil {
		return nil
	}

	var fi FormInfo

	if len(f.Value) > 0 {
		fi.Value = f.Value
	}

	if len(f.File) > 0 {
		fileMap := make(map[string]FileHeaderInfos, len(f.File))

		for k, v := range f.File {
			var fhiSlice []*FileHeaderInfo

			for _, fh := range v {
				if fh != nil {
					fhiSlice = append(fhiSlice, MapFileHeader(fh, encodeBytes))
				}
			}

			fileMap[k] = fhiSlice
		}

		fi.File = fileMap
	}

	return &fi
}
