package service

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

type Service struct {
	Prod Produser
	Pres Presenter
}

func (m *Service) Run(ctx context.Context) error {
	slog.Info("Service.Run Стартанул")
	startTime := time.Now()
	inputchan := make(chan string)
	outputchan := make(chan string)
	input, err := m.Prod.Produce()
	slog.Info("Продюсер вернулся", "Колличество строк", len(input), "error", err)
	if err != nil {
		return err
	}
	if len(input) == 0 {
		slog.Warn("Продюсер вернул пустой входной файл")
		return nil
	}
	var wg sync.WaitGroup
	var wgWriter sync.WaitGroup
	go func() {
		defer func() {
			close(inputchan)
			slog.Info("Горутины чтения отработали", "Отправлено строк :", len(input))
		}()
		sentCount := 0
		for _, line := range input {
			select {
			case <-ctx.Done():
				slog.Warn("Сработал контекст на чтении файла", "Прочитано", sentCount, "Всего", len(input))
				return
			case inputchan <- line:
				sentCount++
				if sentCount%5000 == 0 {
					slog.Info("Читает из файла", "Прочитано", sentCount, "Всего", len(input))
				}
			}
		}
	}()
	for i := 0; i < 10; i++ {
		workerID := i
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
				slog.Debug("Воркер закончил", "id", workerID)
			}()
			processed := 0

			for {
				select {
				case <-ctx.Done():
					slog.Debug("Воркер завершил по контексту", "id", workerID, "Из входного канала записано", processed)
					return
				case line, ok := <-inputchan:
					if !ok {
						slog.Debug("Входной канал закрыт", "id", workerID, "Из входного канала записано", processed)
						return
					}
					maskedline := m.MakeMaskGreatAgain(line)
					select {
					case <-ctx.Done():
						slog.Debug("Отработал контекст на запись в выходной канал", "id", workerID)
						return
					case outputchan <- maskedline:
						processed++
						if processed%1000 == 0 {
							slog.Debug("Прогресс записи в выходной канал", "id", workerID, "Записано строк", processed)
						}
					}
				}
			}
		}()
	}
	wgWriter.Add(1)
	var masked []string
	go func() {
		defer wgWriter.Done()
		writed := 0
		for {
			select {
			case <-ctx.Done():
				slog.Warn("Закончил по контексту писать", "Записал", writed)
				return
			case maskedline, ok := <-outputchan:
				if !ok {
					slog.Info("Выходной канал закрыт", "Всего записано", writed)
					return
				}
				masked = append(masked, maskedline)
				writed++
				if writed%5000 == 0 {
					slog.Info("Прогресс написания", "Написано", writed)
				}
			}
		}
	}()
	slog.Info("Ждем всех воркеров")
	wg.Wait()
	slog.Info("Все воркеры закончили,закрываем канал")
	close(outputchan)
	slog.Info("Ждем писателя")
	wgWriter.Wait()
	slog.Info("Процесс завершился", "Конечные результаты", len(masked), "Количество входных строк", len(input), "Ожидания", len(input), "Фактически", len(masked))
	slog.Info("Отправка презентору", "Количество строк", len(masked))
	if err := m.Pres.Present(masked); err != nil {
		slog.ErrorContext(ctx, "Ошибка презика", "Error", err)
		return err
	}
	slog.Info("Service.Run отработал успешно", "Время", time.Since(startTime))
	return nil
}
func NewService(Prod Produser, Pres Presenter) Service {
	return Service{Prod: Prod, Pres: Pres}

}
