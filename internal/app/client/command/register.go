package command

import (
	"bufio"
	"context"
	"fmt"
	"io"

	"github.com/bbquite/go-pass-keeper/internal/models"
	"github.com/bbquite/go-pass-keeper/internal/service/client"
)

type RegisterCommand struct {
	authService *client.ClientAuthService
	reader      io.Reader
	writer      io.Writer
}

func NewRegisterCommand(authService *client.ClientAuthService, reader io.Reader, writer io.Writer) *RegisterCommand {
	return &RegisterCommand{
		authService: authService,
		reader:      reader,
		writer:      writer,
	}
}

func (c *RegisterCommand) Name() string {
	return "REGISTER"
}

func (c *RegisterCommand) Desc() string {
	return "To access other commands, you need to log in with your login and password."
}

func (c *RegisterCommand) Usage() string {
	return "TO BE ADDED"
}

func (c *RegisterCommand) Execute() error {

	var data models.UserRegisterData
	scanner := bufio.NewScanner(c.reader)

	_, err := fmt.Fprint(c.writer, "Enter login: ")
	if err != nil {
		return err
	}

	if scanner.Scan() {
		data.Username = scanner.Text()
	} else {
		return scanner.Err()
	}

	_, err = fmt.Fprint(c.writer, "Enter password: ")
	if err != nil {
		return err
	}

	if scanner.Scan() {
		data.Password = scanner.Text()
	} else {
		return scanner.Err()
	}

	data.Email = "temp@mail.ru"

	err = c.authService.RegisterUser(context.Background(), &data)
	if err != nil {
		return err
	}

	return nil
}
