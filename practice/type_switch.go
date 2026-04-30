package practice

import "fmt"

func Describe(v any) string {
	switch val := v.(type) {
	case int:
		return fmt.Sprintf("целое число: %d", val)
	case float64:
		return fmt.Sprintf("число с плавающей точкой: %.2f", val)
	case string:
		return fmt.Sprintf("строка длиной %d символов: %s", len(val), val)
	case bool:
		if val {
			return "логическое значение: истина"
		}
		return "логическое значение: ложь"
	case []int:
		sum := 0
		for _, n := range val {
			sum += n
		}
		return fmt.Sprintf("срез целых чисел: %v, сумма элементов: %d", val, sum)
	default:
		return fmt.Sprintf("неизвестный тип: %T", val)
	}
}
