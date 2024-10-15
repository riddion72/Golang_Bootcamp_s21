package model

import (
    "time"
)

// Entry структура для хранения данных записи
type Anomalies struct {
    ID        uint      `gorm:"primaryKey"`    // ID записи
    SessionID string    `gorm:"not null"`      // Идентификатор сессии
    Frequency float64   `gorm:"not null"`      // Частота
    Timestamp time.Time `gorm:"not null"`      // Время записи
}
