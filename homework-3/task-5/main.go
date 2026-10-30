package main

import (
	"errors"
	"fmt"
)

var ErrIndexOutOfRange = errors.New("индекс вне диапазона")

func RemoveElement(nums []int, index int) ([]int, error) {
	if index < 0 || index >= len(nums) {
		return nil, fmt.Errorf("%w: %d [0, %d)", ErrIndexOutOfRange, index, len(nums))
	}

	result := make([]int, 0, len(nums)-1)
	result = append(result, nums[:index]...)
	result = append(result, nums[index+1:]...)

	return result, nil
}

func main() {
	original1 := []int{1, 2, 2, 3, 1}
	res1, err1 := RemoveElement(original1, 3)
	fmt.Printf("Результат: %v, %v\n", res1, err1)
	fmt.Printf("Исходный: %v\n", original1)
	original2 := []int{1, 2, 3}
	_, err2 := RemoveElement(original2, 5)
	fmt.Printf("%v\n", err2)

	if errors.Is(err2, ErrIndexOutOfRange) {
		fmt.Println("Мы точно знаем, что это ошибка выхода за границы!")
	}

}

// Вопрос: почему решение через append(nums[:index], nums[index+1:]...) испортит исходный слайс?

// В Go слайс - это лишь "окно" (указатель) на скрытый базовый массив (backing array).
// Выражение nums[:index] создает новый слайс, но он смотрит в ту же самую память, что и nums.
//
// Функция append, если ей хватает capacity (а здесь ее хватает, так как это кусок nums),
// не создает новый массив, а начинает перезаписывать элементы в старом.
// Она берет элементы из nums[index+1:] и записывает их поверх элементов, начиная с индекса index.
// Поскольку память общая, исходный слайс nums тоже изменится.
//
// Чтобы этого избежать, мы создаем новый слайс через make().
// Это выделяет абсолютно новый участок памяти, и мы копируем туда данные, не трогая оригинал.
