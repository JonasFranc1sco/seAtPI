package upload

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)
// Recebe o caminho do arquivo e a URL do endpoint
func UploadCSV(path, url string) error {
	// Abrir o arquivo
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Criar o corpo da requisição em memória
	var body bytes.Buffer
	w := multipart.NewWriter(&body)

	// Criar o campo do formulário
	fw, err := w.CreateFormFile("file", filepath.Base(path))

	// Copiar o conteúdo do arquivo
	if _, err := io.Copy(fw, f); err != nil {
		return err
	}

	// Fechar o multipart
	w.Close()

	// Montar requisição
	req, err := http.NewRequest(http.MethodPost, url, &body)
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Enviar requisição
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Ler a resposta
	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println(resp.Status, string(respBody))
	return nil
}