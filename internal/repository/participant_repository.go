package repository

import (
	"errors"
	"slido-clone-backend/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ParticipantRepository struct {
	Repository[entity.Participant]
	Log *logrus.Logger
}

func NewParticipantRepository(log *logrus.Logger) *ParticipantRepository {
	return &ParticipantRepository{
		Log: log,
	}
}

func (r *ParticipantRepository) FindByRoomIDAndUserID(db *gorm.DB, roomID uint, userID uint) (*entity.Participant, error) {
	var participant entity.Participant
	err := db.Where("room_id = ? AND user_id = ?", roomID, userID).First(&participant).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &participant, err
}

// List mencari participant berdasarkan room ID dengan pagination
func (r *ParticipantRepository) List(db *gorm.DB, roomID uint, page int, size int) ([]entity.Participant, int64, error) {
	var participants []entity.Participant
	err := db.Where("room_id = ?", roomID).Offset((page - 1) * size).Limit(size).Find(&participants).Error
	if err != nil {
		return nil, 0, err
	}

	var total int64
	err = db.Model(&entity.Participant{}).Where("room_id = ?", roomID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	return participants, total, nil
}

// FindParticipantInRoom mencari participant berdasarkan room ID dan participant ID
func (r *ParticipantRepository) FindParticipantInRoom(db *gorm.DB, roomID uint, participantID uint) (*entity.Participant, error) {
	var participant entity.Participant
	err := db.Where("room_id = ? AND id = ?", roomID, participantID).First(&participant).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &participant, err
}

// ListLeaderboard mengambil daftar xp berdasarkan room id, diurutkan berdasarkan poin tertinggi dan dibatasi hingga 10 entri
func (r *ParticipantRepository) ListLeaderboard(db *gorm.DB, roomID uint) ([]entity.Participant, error) {
	var leaderboard []entity.Participant
	query := db.Where("room_id = ?", roomID).Order("xp_score DESC").Limit(10)

	err := query.Find(&leaderboard).Error
	if err != nil {
		return nil, err
	}
	return leaderboard, nil
}

// GetRankAndScoreByParticipantID mengambil peringkat dan skor peserta dalam sebuah room
func (r *ParticipantRepository) GetRankAndScoreByParticipantID(db *gorm.DB, roomID, participantID uint) (int64, int64, error) {
	// get xp score of the participant
	var xpScore int64
	err := db.Model(&entity.Participant{}).Where("id = ? AND room_id = ?", participantID, roomID).Select("xp_score").Take(&xpScore).Error
	if err != nil {
		return 0, 0, err
	}

	// count participants with higher xp score
	var count int64
	err = db.Model(&entity.Participant{}).Where("room_id = ? AND xp_score > ?", roomID, xpScore).Count(&count).Error
	if err != nil {
		return 0, 0, err
	}

	// rank
	rank := count + 1

	return rank, xpScore, nil
}

// CountByRoomID menghitung jumlah participant dalam sebuah room
func (r *ParticipantRepository) CountByRoomID(db *gorm.DB, roomID uint) (int64, error) {
	var count int64
	err := db.Model(entity.Participant{}).Where("room_id = ?", roomID).Count(&count).Error
	return count, err
}
