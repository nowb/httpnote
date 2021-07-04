package httpnote

type RequestInfo struct {
	Method           string              `json:"string"`
	URL              *URLInfo            `json:"url"`
	Proto            string              `json:"proto"`
	ProtoMajor       int                 `json:"proto_major"`
	ProtoMinor       int                 `json:"proto_minor"`
	Header           map[string][]string `json:"header,omitempty"`
	BodyBytes        []byte              `json:"body_bytes,omitempty"`
	BodyBase64       string              `json:"body_base64,omitempty"`
	ContentLength    int64               `json:"content_length,omitempty"`
	TransferEncoding []string            `json:"transfer_encoding,omitempty"`
	Close            bool                `json:"close,omitempty"`
	Host             string              `json:"host"`
	Form             map[string][]string `json:"form,omitempty"`
	FormEncoded      string              `json:"form_encoded,omitempty"`
	PostForm         map[string][]string `json:"post_form,omitempty"`
	PostFormEncoded  string              `json:"post_form_encoded,omitempty"`
	Trailer          map[string][]string `json:"trailer,omitempty"`
	RemoteAddr       string              `json:"remote_addr"`
	RequestURI       string              `json:"request_uri"`
}

type URLInfo struct {
	Scheme      string    `json:"scheme,omitempty"`
	Opaque      string    `json:"opaque,omitempty"`
	User        *UserInfo `json:"user,omitempty"`
	Host        string    `json:"host,omitempty"`
	Path        string    `json:"path,omitempty"`
	RawPath     string    `json:"raw_path,omitempty"`
	ForceQuery  bool      `json:"force_query,omitempty"`
	RawQuery    string    `json:"raw_query,omitempty"`
	Fragment    string    `json:"fragment,omitempty"`
	RawFragment string    `json:"raw_fragment,omitempty"`
}

type UserInfo struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Encoded  string `json:"encoded,omitempty"`
}
