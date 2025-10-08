package utils

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
	bybit "github.com/wuhewuhe/bybit.go.api"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"
)

func NewBybitHttpClient(apiKey string, APISecret string, baseURL string) *bybit.Client {
	c := &bybit.Client{
		APIKey:     apiKey,
		APISecret:  APISecret,
		BaseURL:    baseURL,
		HTTPClient: http.DefaultClient,
		Logger:     log.New(os.Stderr, bybit.Name, log.LstdFlags),
	}

	return c
}

func NewPgxPool(dsn string, maxConns int32) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to parse config: %v", err)
	}

	poolConfig.MinConns = 50
	poolConfig.MaxConns = maxConns

	dbpool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}
	return dbpool, nil
}

func NewRedisClient(host, port, password string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	})
}

func NewClickHouseConn(host, port, database, username, password string, maxOpenConnection int) (driver.Conn, error) {
	var (
		ctx       = context.Background()
		conn, err = clickhouse.Open(&clickhouse.Options{
			Addr: []string{fmt.Sprintf("%s:%s", host, port)},
			Auth: clickhouse.Auth{
				Database: database,
				Username: username,
				Password: password,
			},

			MaxIdleConns:     maxOpenConnection,
			ConnMaxLifetime:  time.Hour,
			ConnOpenStrategy: clickhouse.ConnOpenInOrder,

			ClientInfo: clickhouse.ClientInfo{
				Products: []struct {
					Name    string
					Version string
				}{
					{Name: "an-example-go-client", Version: "0.1"},
				},
			},

			Debugf: func(format string, v ...interface{}) {
				fmt.Printf(format, v)
			},
			//TLS: &tls.Config{
			//	InsecureSkipVerify: true,
			//},
		})
	)

	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("Exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		return nil, err
	}
	return conn, nil
}

func RequestServer(url string) (string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second, // Устанавливаем таймаут для запроса
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-OK status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	return string(body), nil
}

// AddTime добавляет указанное количество времени к заданному моменту времени и возвращает
// новый момент времени в формате Unix миллисекунд.
//
// Параметры:
// - tf (string): Тип времени, который определяет, как именно добавлять время. Допустимые значения:
//   - "M": Добавляет месяцы. Если последний месяц (декабрь), результат будет первым днем следующего года.
//   - "W": Добавляет недели. Каждая неделя составляет 168 часов.
//   - "D": Добавляет дни. Каждый день составляет 24 часа.
//   - "240": Добавляет 4 часа (240 минут).
//   - "60": Добавляет 1 час (60 минут).
//   - "30": Добавляет 30 минут.
//   - "15": Добавляет 15 минут.
//   - "5": Добавляет 5 минут.
//   - "1": Добавляет 1 минуту.
//
// - lastTime (time.Time): Исходный момент времени, к которому будет добавлено время.
// - k (time.Duration): Количество единиц времени, которое следует добавить.
//
// Возвращаемое значение:
// - int64: Новый момент времени в формате Unix миллисекунд после добавления указанного количества времени.
//
// Примеры использования:
// - AddTime("D", time.Now(), 5) - Добавляет 5 дней к текущему времени и возвращает результат в Unix миллисекундах.
// - AddTime("W", time.Now(), 2) - Добавляет 2 недели к текущему времени и возвращает результат в Unix миллисекундах.
// - AddTime("M", time.Now(), 1) - Добавляет 1 месяц к текущему времени и возвращает результат в Unix миллисекундах.
//
// Если тип времени (tf) не распознан, функция возвращает исходный момент времени в формате Unix миллисекунд.
func AddTime(tf string, lastTime time.Time, k time.Duration) int64 {
	switch tf {
	case "M":
		if lastTime.Month() == 12 {
			return time.Date(lastTime.Year()+1, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()
		} else {
			return time.Date(lastTime.Year(), lastTime.Month()+1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()
		}
	case "W":
		return lastTime.Add(time.Hour * 168 * k).UnixMilli()
	case "D":
		return lastTime.Add(time.Hour * 24 * k).UnixMilli()
	case "240":
		return lastTime.Add(time.Hour * 4 * k).UnixMilli()
	case "60":
		return lastTime.Add(time.Hour * k).UnixMilli()
	case "30":
		return lastTime.Add(time.Minute * 30 * k).UnixMilli()
	case "15":
		return lastTime.Add(time.Minute * 15 * k).UnixMilli()
	case "5":
		return lastTime.Add(time.Minute * 5 * k).UnixMilli()
	case "1":
		return lastTime.Add(time.Minute * k).UnixMilli()
	default:
		return lastTime.UnixMilli()
	}
}

func MinuteLength(timeframe string) float64 {
	switch timeframe {
	case "5":
		return 5
	case "15":
		return 15
	case "30":
		return 30
	case "60":
		return 60
	case "240":
		return 240
	}

	return 0
}

func WriteCsv(filepath string, headers []string, records [][]string) error {
	file, err := os.Create(fmt.Sprintf("%s.csv", filepath))
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write(headers)
	if err != nil {
		return err
	}

	for _, record := range records {
		err = writer.Write(record)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetMinuteInUnixMilli(ms int64) int64 {
	return (ms / (1000 * 60)) % 60
}

func MakeAbsSlice(s []float64) {
	for i := 0; i < len(s); i++ {
		s[i] = math.Abs(s[i])
	}
}

func ReverseSlice(s []float64) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// ParseFloat converts string to float64 with error handling
func ParseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func TimeFormat(date int) string {
	return time.UnixMilli(int64(date)).UTC().Format(time.DateTime)
}
