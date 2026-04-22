package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/auth/middlewares"
	"github.com/nutrixpos/pos/modules/auth/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	Config     config.Config
	Logger     logger.ILogger
	Collection *mongo.Collection
	JWTUtil    *middlewares.JWTUtil
}

func NewAuthHandler(conf config.Config, logger logger.ILogger, collection *mongo.Collection) *AuthHandler {
	return &AuthHandler{
		Config:     conf,
		Logger:     logger,
		Collection: collection,
		JWTUtil:    middlewares.NewJWTUtil(conf.Auth.JWTSecret, conf.Auth.JWTExpireHrs),
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	var user models.User
	err := h.Collection.FindOne(ctx, bson.M{"username": req.Username}).Decode(&user)
	if err != nil {
		h.Logger.Error("user not found", "username", req.Username)
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if !middlewares.CheckPassword(req.Password, user.PasswordHash) {
		h.Logger.Error("invalid password", "username", req.Username)
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := h.JWTUtil.GenerateToken(user)
	if err != nil {
		h.Logger.Error("failed to generate token", "error", err)
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.LoginResponse{
		Token: token,
		User:  user.ToResponse(),
	})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	count, err := h.Collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		h.Logger.Error("failed to count users", "error", err)
		http.Error(w, "failed to register", http.StatusInternalServerError)
		return
	}

	defaultRoles := []string{}
	if count == 0 {
		defaultRoles = []string{"superuser"}
	} else {
		defaultRoles = []string{"cashier"}
	}

	var userReq models.User
	if len(req.Roles) > 0 {
		userReq.Roles = req.Roles
		for _, role := range userReq.Roles {
			if role == "superuser" {
				http.Error(w, "cannot create superuser", http.StatusForbidden)
				return
			}
		}
	} else {
		userReq.Roles = defaultRoles
	}

	hashedPassword, err := middlewares.HashPassword(req.Password)
	if err != nil {
		h.Logger.Error("failed to hash password", "error", err)
		http.Error(w, "failed to register", http.StatusInternalServerError)
		return
	}

	user := models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Roles:        userReq.Roles,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	result, err := h.Collection.InsertOne(ctx, user)
	if err != nil {
		h.Logger.Error("failed to insert user", "error", err)
		http.Error(w, "failed to register", http.StatusInternalServerError)
		return
	}

	user.ID = result.InsertedID.(primitive.ObjectID)

	token, err := h.JWTUtil.GenerateToken(user)
	if err != nil {
		h.Logger.Error("failed to generate token", "error", err)
		http.Error(w, "failed to register", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.LoginResponse{
		Token: token,
		User:  user.ToResponse(),
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	authCtx := r.Context().Value("auth_ctx")
	if authCtx == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	claims, ok := authCtx.(*middlewares.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	var user models.User
	err = h.Collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user.ToResponse())
}

func (h *AuthHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	cursor, err := h.Collection.Find(ctx, bson.M{})
	if err != nil {
		h.Logger.Error("failed to find users", "error", err)
		http.Error(w, "failed to get users", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var users []models.UserResponse
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			h.Logger.Error("failed to decode user", "error", err)
			continue
		}
		users = append(users, user.ToResponse())
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *AuthHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	authCtx := r.Context().Value("auth_ctx")
	if authCtx == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if _, ok := authCtx.(*middlewares.Claims); !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userID := r.URL.Query().Get("id")
	if userID == "" {
		http.Error(w, "user id required", http.StatusBadRequest)
		return
	}

	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	var targetUser models.User
	err = h.Collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&targetUser)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	isTargetSuperuser := false
	for _, role := range targetUser.Roles {
		if role == "superuser" {
			isTargetSuperuser = true
			break
		}
	}

	if isTargetSuperuser {
		http.Error(w, "cannot delete superuser", http.StatusForbidden)
		return
	}

	_, err = h.Collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		h.Logger.Error("failed to delete user", "error", err)
		http.Error(w, "failed to delete user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "user deleted"})
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	authCtx := r.Context().Value("auth_ctx")
	if authCtx == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	claims, ok := authCtx.(*middlewares.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	isSuperuser := false
	for _, role := range claims.Roles {
		if role == "superuser" {
			isSuperuser = true
			break
		}
	}

	if !isSuperuser {
		http.Error(w, "only superuser can change user passwords", http.StatusForbidden)
		return
	}

	var req models.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.UserID == "" || req.Password == "" {
		http.Error(w, "user_id and password are required", http.StatusBadRequest)
		return
	}

	oid, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	hashedPassword, err := middlewares.HashPassword(req.Password)
	if err != nil {
		h.Logger.Error("failed to hash password", "error", err)
		http.Error(w, "failed to change password", http.StatusInternalServerError)
		return
	}

	ctx := context.Background()

	_, err = h.Collection.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{"password_hash": hashedPassword, "updated_at": time.Now()}})
	if err != nil {
		h.Logger.Error("failed to change password", "error", err)
		http.Error(w, "failed to change password", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "password changed"})
}

func (h *AuthHandler) ChangeMyPassword(w http.ResponseWriter, r *http.Request) {
	authCtx := r.Context().Value("auth_ctx")
	if authCtx == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	claims, ok := authCtx.(*middlewares.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		CurrentPassword string `json:"current_password"`
		Password        string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.CurrentPassword == "" || req.Password == "" {
		http.Error(w, "current_password and password are required", http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	var user models.User

	user_id_obj, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	err = h.Collection.FindOne(ctx, bson.M{"_id": user_id_obj}).Decode(&user)
	if err != nil {
		h.Logger.Error("user not found", "error", err)
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	if !middlewares.CheckPassword(req.CurrentPassword, user.PasswordHash) {
		http.Error(w, "invalid current password", http.StatusUnauthorized)
		return
	}

	hashedPassword, err := middlewares.HashPassword(req.Password)
	if err != nil {
		h.Logger.Error("failed to hash password", "error", err)
		http.Error(w, "failed to change password", http.StatusInternalServerError)
		return
	}

	_, err = h.Collection.UpdateOne(ctx, bson.M{"_id": user_id_obj}, bson.M{"$set": bson.M{"password_hash": hashedPassword, "updated_at": time.Now()}})
	if err != nil {
		h.Logger.Error("failed to change password", "error", err)
		http.Error(w, "failed to change password", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "password changed"})
}
