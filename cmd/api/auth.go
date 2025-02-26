package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/MohammadBohluli/social-app-go/internal/mailer"
	"github.com/MohammadBohluli/social-app-go/internal/store"
	"github.com/MohammadBohluli/social-app-go/pkg"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type RegisterUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserWithActivateToken struct {
	*store.User
	Token string `json:"token"`
}

type CreateUserTokenPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// registerUserHandler godoc
//
//	@Summary		Registers a user
//	@Description	Registers a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserRequest		true	"User credentials"
//	@Success		201		{object}	UserWithActivateToken	"User registered"
//	@Router			/auth/user [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserRequest
	if err := pkg.ReadJson(w, r, &payload); err != nil {
		pkg.BadRequestError(w, r, err)
		return
	}

	hashPassword, err := pkg.Hash(payload.Password)
	if err != nil {
		pkg.InternalServerError(w, r, err)
		return
	}

	user := store.User{
		Username: payload.Username,
		Email:    payload.Email,
		Password: hashPassword,
	}

	palinToken := uuid.New().String()
	hash := sha256.Sum256([]byte(palinToken))
	hashToken := hex.EncodeToString(hash[:])

	ctx := r.Context()
	err = app.store.Users.CreateAndInvite(ctx, &user, hashToken, app.config.mail.exp)
	if err != nil {

		switch err {
		case store.ErrDuplicateEmail:
			pkg.BadRequestError(w, r, err)
		case store.ErrDuplicateUsername:
			pkg.BadRequestError(w, r, err)
		default:
			pkg.InternalServerError(w, r, err)
		}
		return

	}

	u := UserWithActivateToken{
		User:  &user,
		Token: palinToken,
	}

	activationURL := fmt.Sprintf("%s/confirm/%s", app.config.frontEndURL, hashToken)
	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.Username,
		ActivationURL: activationURL,
	}

	isProdEnv := app.config.evn == "development"
	status, err := app.mailer.Send(mailer.UserWellcomeTemplate, user.Username, user.Email, vars, !isProdEnv)

	if err != nil {
		app.logger.Errorw("error sending wellcome email", "error", err)

		// roleback user creation if email fails(SAGA pattern)
		if err := app.store.Users.Delete(ctx, user.ID); err != nil {
			app.logger.Errorw("error sending wellcome email", "error", err)

		}

		pkg.InternalServerError(w, r, err)
		return
	}

	app.logger.Infow("Email sent", "status code", status)

	if err := pkg.JsonResponse(w, http.StatusCreated, u); err != nil {
		pkg.InternalServerError(w, r, err)
		return
	}
}

// createTokenHandler godoc

// @Summary		Creates a token
// @Description	Creates a token for a user
// @Tags			authentication
// @Accept			json
// @Produce		json
// @Param			payload	body		CreateUserTokenPayload	true	"User credentials"
// @Success		200		{string}	string					"Token"
// @Router			/auth/token [post]
func (app *application) createTokenHandler(w http.ResponseWriter, r *http.Request) {

	// parse payload credential
	var payload CreateUserTokenPayload
	if err := pkg.ReadJson(w, r, &payload); err != nil {
		pkg.BadRequestError(w, r, err)
		return
	}
	// fetch the user(check if suer exist) from the payload

	user, err := app.store.Users.GetByEmail(r.Context(), payload.Email)
	if err != nil {
		switch err {
		case store.ErrorNotFound:
			pkg.UnAuthorizedBasicErrorResponse(w, r, err)
		default:
			pkg.InternalServerError(w, r, err)
		}
		return
	}
	// generate the token -> add claims
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(app.config.auth.token.exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": "gopherSocial",
	}

	token, err := app.authenticator.GenerateToken(claims)
	if err != nil {
		pkg.InternalServerError(w, r, err)
		return
	}
	// send it to the client

	if err := pkg.JsonResponse(w, http.StatusCreated, token); err != nil {
		pkg.InternalServerError(w, r, err)
		return
	}
}
