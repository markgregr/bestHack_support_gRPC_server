package dataprocessing

import (
	"encoding/csv"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"os"
	"sort"
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

// Helper functions
func median(numbers []int) float64 {
	sort.Ints(numbers)
	n := len(numbers)
	if n%2 == 0 {
		return float64(numbers[n/2-1]+numbers[n/2]) / 2.0
	}
	return float64(numbers[n/2])
}

func mean(numbers []int) float64 {
	total := 0
	for _, number := range numbers {
		total += number
	}
	return float64(total) / float64(len(numbers))
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
	// Aggregate data by cluster
	clusterDurations := make(map[int][]int)
	clusterReactions := make(map[int][]int)
	for _, item := range data {
		clusterDurations[item.ClusterIndex] = append(clusterDurations[item.ClusterIndex], item.DurationTime)
		clusterReactions[item.ClusterIndex] = append(clusterReactions[item.ClusterIndex], item.ReactionTime)
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

	header := []string{"ClusterIndex", "AvgDuration", "MedianDuration", "StdDevDuration", "AvgReaction", "MedianReaction", "StdDevReaction"}
	if err := writer.Write(header); err != nil {
		log.WithError(err).Error("failed to write header")
		return err
	}

	// Write data
	for cluster, durations := range clusterDurations {
		reactions := clusterReactions[cluster]
		avgDuration := mean(durations)
		medianDuration := median(durations)
		avgReaction := mean(reactions)
		medianReaction := median(reactions)
		record := []string{
			strconv.Itoa(cluster),
			strconv.FormatFloat(avgDuration, 'f', 2, 64),
			strconv.FormatFloat(medianDuration, 'f', 2, 64),
			strconv.FormatFloat(avgReaction, 'f', 2, 64),
			strconv.FormatFloat(medianReaction, 'f', 2, 64),
		}
		if err := writer.Write(record); err != nil {
			log.WithError(err).Error("failed to write record")
			return err
		}
	}

	log.Info("Statistical data saved to CSV for analysis")
	return nil
}
