package verify

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"net/smtp"
	"regexp"

	"purpleschool/3-validation-api/configs"
	"purpleschool/3-validation-api/internal/payload"
	"purpleschool/3-validation-api/internal/response"

	"github.com/jordan-wright/email"
)

type VerifyHandler struct {
	config *configs.Config
}

const smtpName string = "smtp.mail.ru"

var Result payload.Email

func NewVerifyHandler(router *http.ServeMux, config *configs.Config) {
	handler := VerifyHandler{
		config: config,
	}

	router.HandleFunc("GET /verify/{hash}", handler.verify())
	router.HandleFunc("POST /send", handler.send())
}

func (h *VerifyHandler) send() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var el payload.EmailRequest

		if err := json.NewDecoder(r.Body).Decode(&el); err != nil {
			log.Printf("ошибка при декодировании данных: %v", err)
			response.JSON(w, nil, http.StatusBadRequest)
			return
		}

		if err := validateEmail(el.Email); err != nil {
			log.Printf("ошибка при валидации email: %v", err)
			response.JSON(w, err.Error(), http.StatusBadRequest)
			return
		}

		hash := generateShortKey(8)

		Result.Email = el.Email
		Result.Hash = hash

		if err := h.createEmail(Result); err != nil {
			log.Printf("ошибка при отправке письма: %v", err)
			response.JSON(w, "ошибка при отправке письма", http.StatusInternalServerError)
			return
		}

		response.JSON(w, Result, http.StatusOK)
	}
}

func (h *VerifyHandler) verify() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.PathValue("hash")

		if hash == "" {
			response.JSON(w, "незаполнен хэш", http.StatusBadRequest)
			return
		}

		var resp bool

		if Result.Hash != hash {
			response.JSON(w, resp, http.StatusBadRequest)
			return
		}

		Result = payload.Email{}

		resp = true
		response.JSON(w, resp, http.StatusOK)
	}
}

func (h *VerifyHandler) createEmail(value payload.Email) error {

	e := email.NewEmail()
	e.From = fmt.Sprintf("Test <%s>", h.config.Email)
	e.To = []string{value.Email}
	e.Subject = "Subject"

	verificationLink := fmt.Sprintf("http://localhost:8081/verify/%s", value.Hash)

	e.Text = []byte(fmt.Sprintf(
		`Подтвердите ваш email:
	1. Нажмите на ссылку (может потребоваться копирование):
	%s
	2. Или скопируйте и вставьте в браузер:
	%s`, verificationLink, verificationLink))

	err := e.SendWithStartTLS(h.config.Address,
		smtp.PlainAuth("", h.config.Email, h.config.Password, smtpName), &tls.Config{
			ServerName:         smtpName,
			InsecureSkipVerify: true, // Только для отладки!
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func generateShortKey(length int) string {
	const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	key := make([]byte, length)

	for i := range key {
		n, _ := rand.Int(rand.Reader, big.NewInt(62))
		key[i] = base62Chars[n.Int64()]
	}

	return string(key)
}

func validateEmail(email string) error {
	if email == "" {
		return errors.New("незаполнен email")
	}

	match, err := regexp.MatchString(`[A-Za-z0-9\._%+\-]+@[A-Za-z0-9\.\-]+\.[A-Za-z]{2,}`, email)
	if err != nil {
		return errors.New("некорректный email")
	}

	if !match {
		return errors.New("некорректный email")
	}

	return nil
}
