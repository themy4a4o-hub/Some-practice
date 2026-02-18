package service

import "time"

func (m *Service) MakeMaskGreatAgain(s string) string {
	time.Sleep(50 * time.Millisecond)
	input := []byte(s)                    //Вводимая строка
	buffer := make([]byte, 0, len(input)) //то,что нужно заблюрить
	target := []byte("http://")           //одновременно и начало,после которого надо блюрить и то,что я ищу в вводимой строке
	isMatch := true
	for i := 0; i < len(input); i++ {
		if i+len(target) <= len(input) { //цикл на проверку соответствия количеству байтов
			for k := 0; k < len(target); k++ { //проверка на соответствие таргету (http://)
				if input[i+k] != target[k] {
					isMatch = false //если не соответствует,то false
					break           //прерываем цикл
				}
			}
			if isMatch { //если соответствует,то всё ок
				buffer = append(buffer, target...) //сходу добавляем http:// в буфер
				j := i + len(target)
				for j < len(input) && input[j] != ' ' { //цикл ,который ищет конец ссылки
					j++
				} //из -за этой фигурной скобки я мучался 3 дня
				stars := j - (i + len(target)) //переменная,которая будет блюрить
				for k := 0; k < stars; k++ {   // цикл,который добавляет звездочки в равном колличеству байтов до конца ссылки
					buffer = append(buffer, '*') // добавляю в буффер звездочки
				}
				i = j - 1 // и потом возвращаем j на один назад
				continue  // продолжаем

			}
		}
		buffer = append(buffer, input[i])
	}

	return string(buffer) //возвращаем заблюренную строку
}
