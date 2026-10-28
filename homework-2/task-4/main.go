package main

import "fmt"

type Profile struct {
	Name       string
	Age        int
	Email      string
	Newsletter bool
}

type ProfilePatch struct {
	Name       *string
	Age        *int
	Email      *string
	Newsletter *bool
}

func newInt(v int) *int          { return &v }
func newBool(v bool) *bool       { return &v }
func newString(v string) *string { return &v }

func ApplyPatch(p *Profile, patch ProfilePatch) {
	if p == nil {
		return
	}

	if patch.Name != nil {
		p.Name = *patch.Name
	}
	if patch.Age != nil {
		p.Age = *patch.Age
	}
	if patch.Email != nil {
		p.Email = *patch.Email
	}
	if patch.Newsletter != nil {
		p.Newsletter = *patch.Newsletter
	}
}

func main() {
	profile := Profile{
		Name:       "Анна",
		Age:        30,
		Email:      "anna@example.com",
		Newsletter: true,
	}

	fmt.Printf("До патча: %+v\n", profile)

	patch := ProfilePatch{
		Age:        new(31),
		Newsletter: newBool(false),
	}

	ApplyPatch(&profile, patch)
	fmt.Printf("После патча: %+v\n", profile)
	fmt.Println("\n--- Проверка на nil ---")
	ApplyPatch(&profile, ProfilePatch{})
	fmt.Printf("Пустой патч: %+v\n", profile)
	ApplyPatch(nil, patch)
	fmt.Println("Передача nil-профиля прошла безопасно.")
}

// Вопрос: почему в патче нельзя просто объявить Newsletter bool?
// Ответ:
// В Go у типа bool есть "нулевое значение" (zero value), которое равно false.
// Если мы объявим в патче просто `Newsletter bool`, мы не сможем отличить две ситуации:
// 1. Пользователь ЯВНО передал false (хочет отписаться от рассылки).
// 2. Пользователь ВООБЩЕ НЕ ПЕРЕДАЛ это поле (и оно автоматически стало false).
//
// Использование указателя *bool решает эту проблему:
// - nil означает «поле не передано, ничего не меняем».
// - &false означает «явно передан false, нужно обновить профиль».
