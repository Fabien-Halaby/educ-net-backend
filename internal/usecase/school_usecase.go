package usecase

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"educnet/internal/domain"
	"educnet/internal/repository"
	"educnet/internal/utils"
)

// ! CreateSchoolInput DTO d'entrée
type CreateSchoolInput struct {
	SchoolName    string `json:"school_name"`
	AdminEmail    string `json:"admin_email"`
	AdminPassword string `json:"admin_password"`
	AdminName     string `json:"admin_name"`
	Phone         string `json:"phone"`
	Address       string `json:"address"`
}

// ! CreateSchoolOutput DTO de sortie
type CreateSchoolOutput struct {
	School *domain.School
	Admin  *domain.User
	Token  string
}

// ! SchoolUseCase interface
type SchoolUseCase interface {
	CreateSchool(input CreateSchoolInput) (*CreateSchoolOutput, error)
}

// ! schoolUseCase implémentation
type schoolUseCase struct {
	db         *sql.DB
	schoolRepo repository.SchoolRepository
	userRepo   repository.UserRepository
	jwtSecret  string
}

// ! NewSchoolUseCase crée un nouveau use case
func NewSchoolUseCase(
	db *sql.DB,
	schoolRepo repository.SchoolRepository,
	userRepo repository.UserRepository,
	jwtSecret string,
) SchoolUseCase {
	return &schoolUseCase{
		db:         db,
		schoolRepo: schoolRepo,
		userRepo:   userRepo,
		jwtSecret:  jwtSecret,
	}
}

// ! CreateSchool crée une nouvelle école avec son admin
func (uc *schoolUseCase) CreateSchool(input CreateSchoolInput) (*CreateSchoolOutput, error) {
	//! 1. Validation basique
	if err := uc.validateInput(input); err != nil {
		return nil, err
	}

	//! 2. Créer slug
	slug := utils.CreateSlug(input.SchoolName)

	//! 3. Vérifier si slug existe
	exists, err := uc.schoolRepo.ExistsBySlug(slug)
	if err != nil {
		return nil, fmt.Errorf("failed to check school slug: %w", err)
	}
	if exists {
		return nil, domain.ErrSchoolAlreadyExists
	}

	//! 4. Vérifier si email existe
	exists, err = uc.userRepo.ExistsByEmail(input.AdminEmail)
	if err != nil {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if exists {
		return nil, domain.ErrEmailAlreadyExists
	}

	//! 5. Séparer prénom et nom
	firstName, lastName := utils.SplitFullName(input.AdminName)
	fmt.Println("SCHOOL_NAME: ", input)
	//! 6. Créer entités domain (avec validation métier)
	school, err := domain.NewSchool(
		input.SchoolName,
		slug,
		input.Address,
		input.AdminEmail,
		input.Phone,
	)
	if err != nil {
		return nil, err
	}

	admin, err := domain.NewAdminUser(
		0, //! school_id sera set après création école
		input.AdminEmail,
		input.AdminPassword,
		firstName,
		lastName,
		input.Phone,
	)
	if err != nil {
		return nil, err
	}

	//! 7. Transaction database
	tx, err := uc.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	//! 8. Créer école
	if err := uc.schoolRepo.Create(school); err != nil {
		return nil, fmt.Errorf("failed to create school: %w", err)
	}

	//! 9. Créer admin
	admin.SchoolID = school.ID
	if err := uc.userRepo.Create(admin); err != nil {
		return nil, fmt.Errorf("failed to create admin: %w", err)
	}

	//! 10. Mettre à jour école avec admin_user_id
	school.SetAdmin(admin.ID)
	if err := uc.schoolRepo.Update(school); err != nil {
		return nil, fmt.Errorf("failed to update school: %w", err)
	}

	//! 11. Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	//! 12. Générer JWT token
	token, err := uc.generateToken(admin)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	//! 13. Retourner résultat
	return &CreateSchoolOutput{
		School: school,
		Admin:  admin,
		Token:  token,
	}, nil
}

// ! validateInput valide l'input
func (uc *schoolUseCase) validateInput(input CreateSchoolInput) error {
	if input.SchoolName == "" {
		return domain.ErrSchoolNameRequired
	}
	if input.AdminEmail == "" {
		return domain.ErrEmailRequired
	}
	if input.AdminPassword == "" || len(input.AdminPassword) < 6 {
		return domain.ErrPasswordTooShort
	}
	if input.AdminName == "" {
		return domain.ErrNameRequired
	}
	return nil
}

// ! generateToken génère un JWT token
func (uc *schoolUseCase) generateToken(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   user.ID,
		"school_id": user.SchoolID,
		"email":     user.Email,
		"role":      user.Role,
		"exp":       time.Now().Add(time.Hour * 24 * 7).Unix(), //! 7 jours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.jwtSecret))
}
