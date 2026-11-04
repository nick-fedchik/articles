package reason

import (
	"os"
	"os/exec"
	"testing"
)

func TestMain(m *testing.M) {
	mutant = os.Getenv("REASON_MUTANT")
	os.Exit(m.Run())
}

// mutants lists the injected defects; each must make the rest of the test
// suite fail, otherwise the suite does not detect that class of error.
var mutants = []struct{ id, defect string }{
	{"M1", "усі типові правила діють як строгі"},
	{"M2a", "корінь свідчення замінено на ретранслятор"},
	{"M3a", "рівність кількості підтримки не фіксується як нічия"},
	{"M3b", "висновки правил не об'єднуються зі спостереженнями (немає JOIN)"},
	{"M4", "конфлікт виведеного факту не враховується у зіставленні"},
	{"M5", "ключ атома без екранування (колізія p(a,b) і p(\"a,b\"))"},
	{"M6", "прямий прогін пропускається для спостережуваної цілі"},
	{"M8", "перевірка циклу правил вимкнена"},
}

func TestMutationsAreDetected(t *testing.T) {
	if mutant != "" {
		t.Skip("запуск усередині мутанта")
	}
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	detected := 0
	for _, mu := range mutants {
		cmd := exec.Command(exe, "-test.count=1")
		cmd.Env = append(os.Environ(), "REASON_MUTANT="+mu.id)
		out, runErr := cmd.CombinedOutput()
		if runErr != nil {
			detected++
			t.Logf("%s виявлено: %s", mu.id, mu.defect)
			continue
		}
		t.Errorf("%s НЕ виявлено (%s)\n%s", mu.id, mu.defect, out)
	}
	t.Logf("виявлено %d із %d мутантів", detected, len(mutants))
}
