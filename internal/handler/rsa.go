package handler

import (
	"bytes"
	"encoding/base64"
	"io"
	"log/slog"
	"math/big"
	"mime/multipart"
	"net/http"

	"github.com/PritOriginal/cryptolabs-back/internal/services/crypto"
	"github.com/PritOriginal/problem-map-server/pkg/handlers"
	"github.com/PritOriginal/problem-map-server/pkg/responses"
)

type RsaHandler struct {
	handlers.BaseHandler
	s crypto.RSA
}

func NewRsaHandler(log *slog.Logger, s crypto.RSA) *RsaHandler {
	return &RsaHandler{handlers.BaseHandler{Log: log}, s}
}

func (h *RsaHandler) GenerateKeys() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pub, priv, err := h.s.GenerateKeys(2048)
		if err != nil {
			h.RenderInternalError(w, r, handlers.HandlerError{Msg: "failed generate keys", Err: err})
			return
		}

		type RsaKeys struct {
			Public  string `json:"public"`
			Private string `json:"private"`
		}

		keys := RsaKeys{}

		publicKey := base64.StdEncoding.EncodeToString(append(pub.N.Bytes(), pub.E.Bytes()...))
		keys.Public = publicKey

		privateKey := base64.StdEncoding.EncodeToString(append(priv.N.Bytes(), priv.D.Bytes()...))
		keys.Private = privateKey

		h.Render(w, r, responses.SucceededRenderer(keys))
	}
}
func (h *RsaHandler) Encrypt() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(32 << 10) // 32 MB
		if err != nil {
			h.RenderInternalError(w, r, handlers.HandlerError{Msg: "error parse multipart form", Err: err})
			return
		}

		data, key, err := h.readDataAndKey(w, r)
		if err != nil {
			h.RenderError(w, r, handlers.HandlerError{Msg: "failed read data and key", Err: err}, responses.ErrBadRequest)
			return
		}

		pub := crypto.PublicKey{N: new(big.Int).SetBytes(key[:256]), E: new(big.Int).SetBytes(key[256:])}
		ciphertext, err := h.s.Encrypt(&pub, data)
		if err != nil {
			h.RenderInternalError(w, r, handlers.HandlerError{Msg: "failed encrypt data", Err: err})
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Disposition", "attachment; filename=ciphertext.txt")
		w.Write(ciphertext)
	}
}
func (h *RsaHandler) Decrypt() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(32 << 10) // 32 MB
		if err != nil {
			h.RenderInternalError(w, r, handlers.HandlerError{Msg: "error parse multipart form", Err: err})
			return
		}

		ciphertext, key, err := h.readDataAndKey(w, r)
		if err != nil {
			h.RenderError(w, r, handlers.HandlerError{Msg: "failed read data and key", Err: err}, responses.ErrBadRequest)
			return
		}

		priv := crypto.PrivateKey{N: new(big.Int).SetBytes(key[:256]), D: new(big.Int).SetBytes(key[256:])}
		data, err := h.s.Decrypt(&priv, ciphertext)
		if err != nil {
			h.RenderInternalError(w, r, handlers.HandlerError{Msg: "failed decrypt data", Err: err})
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Disposition", "attachment; filename=test.txt")
		w.Write(data)
	}
}

func (h *RsaHandler) readDataAndKey(w http.ResponseWriter, r *http.Request) ([]byte, []byte, error) {
	data, err := readFile(r.MultipartForm.File["data"][0])
	if err != nil {
		return nil, nil, err
	}
	keyBase64, err := readFile(r.MultipartForm.File["key"][0])
	if err != nil {
		return nil, nil, err
	}

	key, err := base64.StdEncoding.DecodeString(string(keyBase64))
	if err != nil {
		return nil, nil, err
	}

	return data, key, nil
}

func readFile(header *multipart.FileHeader) ([]byte, error) {
	file, err := header.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
