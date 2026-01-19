package public

import (
	"context"
	"database/sql"
	"errors"

	"github.com/outlet-sh/outlet/internal/db"
	"github.com/outlet-sh/outlet/internal/svc"
	"github.com/outlet-sh/outlet/internal/types"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type CreateInitialAdminLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Create initial admin account (only works when no admin exists)
func NewCreateInitialAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateInitialAdminLogic {
	return &CreateInitialAdminLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateInitialAdminLogic) CreateInitialAdmin(req *types.CreateInitialAdminRequest) (resp *types.CreateInitialAdminResponse, err error) {
	// Check if any admin already exists
	adminCount, err := l.svcCtx.DB.CountUsers(l.ctx, db.CountUsersParams{
		FilterRole:   "super_admin",
		FilterStatus: nil,
	})
	if err != nil {
		l.Error("failed to count admin users:", err)
		return nil, err
	}

	if adminCount > 0 {
		return nil, errors.New("admin account already exists")
	}

	// Validate passwords match
	if req.Password != req.ConfirmPassword {
		return nil, errors.New("passwords do not match")
	}

	// Validate password strength (minimum 8 characters)
	if len(req.Password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		l.Error("failed to hash password:", err)
		return nil, err
	}

	// Set default name if not provided
	name := req.Name
	if name == "" {
		name = "Admin"
	}

	// Create the admin user
	_, err = l.svcCtx.DB.CreateUser(l.ctx, db.CreateUserParams{
		ID:            uuid.New().String(),
		Email:         req.Email,
		PasswordHash:  string(hashedPassword),
		Name:          name,
		Role:          "super_admin",
		Status:        "active",
		EmailVerified: 1,
		Phone:         sql.NullString{Valid: false},
	})
	if err != nil {
		l.Error("failed to create admin user:", err)
		return nil, err
	}

	// Save company name to platform settings if provided
	if req.Company != "" {
		_, err = l.svcCtx.DB.UpsertPlatformSetting(l.ctx, db.UpsertPlatformSettingParams{
			Key:         "platform.company",
			ValueText:   sql.NullString{String: req.Company, Valid: true},
			Description: sql.NullString{String: "Company name", Valid: true},
			Category:    "platform",
			IsSensitive: sql.NullInt64{Int64: 0, Valid: true},
		})
		if err != nil {
			l.Error("failed to save company setting:", err)
			// Don't fail the entire operation, just log the error
		}
	}

	// Save timezone to platform settings if provided
	if req.Timezone != "" {
		_, err = l.svcCtx.DB.UpsertPlatformSetting(l.ctx, db.UpsertPlatformSettingParams{
			Key:         "platform.timezone",
			ValueText:   sql.NullString{String: req.Timezone, Valid: true},
			Description: sql.NullString{String: "Platform timezone for scheduling", Valid: true},
			Category:    "platform",
			IsSensitive: sql.NullInt64{Int64: 0, Valid: true},
		})
		if err != nil {
			l.Error("failed to save timezone setting:", err)
			// Don't fail the entire operation, just log the error
		}
	}

	l.Infof("Initial admin account created for %s", req.Email)

	return &types.CreateInitialAdminResponse{
		Success: true,
		Message: "Admin account created successfully",
	}, nil
}
