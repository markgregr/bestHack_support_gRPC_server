package csvsaver

import (
	"encoding/csv"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
)

// Структура для элементов JSON
type ClusterData struct {
	ClusterIndex int `json:"cluster_index"`
	DurationTime int `json:"duration_time"`
	ReactionTime int `json:"reaction_time"`
}

// Структура для хранения статистики
type ClusterStats struct {
	TotalDuration int
	TotalReaction int
	Count         int
}

func AddDataToJSON(jsonFile string, data ClusterData, log *logrus.Logger) error {
	const op = "utils.CsvSaver.AddDataToJSON"
	log.WithField("method", op)

	file, err := os.ReadFile("/app/data/input.json")
	if err != nil {
		log.WithError(err).Error("failed to open file")
		return err
	}

	// Парсим JSON данные
	var jsonData []ClusterData
	err = json.Unmarshal(file, &jsonData)
	if err != nil {
		log.WithError(err).Error("failed to unmarshal JSON")
		return err
	}

	// Добавляем новые данные
	jsonData = append(jsonData, data)

	// Сохраняем обновленные данные
	newData, err := json.Marshal(jsonData)
	if err != nil {
		log.WithError(err).Error("failed to marshal JSON")
		return err
	}

	err = os.WriteFile("/app/data/input.json", newData, 0666)
	if err != nil {
		log.WithError(err).Error("failed to write file")
		return err
	}
	log.Info("data added to JSON")
	return nil
}

func AvgCsv(inputFile, outputFile string, log *logrus.Logger) error {
	const op = "utils.CsvSaver.AvgCsv"
	log.WithField("method", op)

	file, err := os.ReadFile("/app/data/input.json")
	if err != nil {
		log.WithError(err).Error("failed to open file")
		return err
	}

	// Парсим JSON данные
	var data []ClusterData
	err = json.Unmarshal(file, &data)
	if err != nil {
		log.WithError(err).Error("failed to unmarshal JSON")
		return err
	}
	// Словарь для сбора статистики
	stats := make(map[int]*ClusterStats)

	// Собираем статистику по каждому кластеру
	for _, item := range data {
		if _, exists := stats[item.ClusterIndex]; !exists {
			stats[item.ClusterIndex] = &ClusterStats{}
		}
		cluster := stats[item.ClusterIndex]
		cluster.TotalDuration += item.DurationTime
		cluster.TotalReaction += item.ReactionTime
		cluster.Count++
	}

	// Открываем CSV файл для добавления данных
	outFile, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.WithError(err).Error("failed to open CSV file")
		return err
	}
	defer outFile.Close()

	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	// Write the header if the file is new or empty
	info, err := outFile.Stat()
	if err != nil {
		log.WithError(err).Error("failed to stat CSV file")
		return err
	}
	if info.Size() == 0 {
		if err := writer.Write([]string{"ClusterIndex", "AvgDuration", "AvgReaction"}); err != nil {
			log.WithError(err).Error("failed to write header")
			return err
		}
	}

	// Write data to CSV
	for cluster, stats := range stats {
		avgDuration := float64(stats.TotalDuration) / float64(stats.Count)
		avgReaction := float64(stats.TotalReaction) / float64(stats.Count)
		clusterStr := strconv.Itoa(cluster)
		avgDurationStr := strconv.FormatFloat(avgDuration, 'f', -1, 64)
		avgReactionStr := strconv.FormatFloat(avgReaction, 'f', -1, 64)

		if err := writer.Write([]string{clusterStr, avgDurationStr, avgReactionStr}); err != nil {
			log.WithError(err).Error("failed to write CSV row")
			return err
		}
	}
	log.Info("data saved to CSV")
	return nil
}
