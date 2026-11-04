package reason

// mutant names one deliberately injected defect (M1–M8) used only by the
// reproducible mutation test in mutation_test.go. The test binary sets it from
// the REASON_MUTANT environment variable in TestMain and re-executes itself
// once per mutant; in every other build the value is empty and no defect is
// active.
var mutant string
