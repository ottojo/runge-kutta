package rungekutta

type RungeKutta struct {
	Y                map[float64]float64
	t_0              float64
	y_0              float64
	phi              Phi
	deltaT           float64
	calculatedNumber int
}

type Phi interface {
	phi(t, y_t float64) float64
}

func (r *RungeKutta) calcUntil(t float64) {
	for r.t_0+float64(r.calculatedNumber)*r.deltaT < t {
		r.CalcY(r.calculatedNumber + 1)
		r.calculatedNumber++
	}
}

func (r *RungeKutta) CalcY(i int) {
	k1, k2, k3, k4 := r.k(i - 1)
	r.Y[r.t(i)] = r.Y[r.t(i-1)] + (k1/6+k2/3+k3/3+k4/6)*r.deltaT
}

// Berechnet k_1 bis k_4 für Zeitschritt i
func (r *RungeKutta) k(i int) (k1 float64, k2 float64, k3 float64, k4 float64) {
	k1 = r.phi.phi(r.t(i), r.Y[r.t(i)])
	k2 = r.phi.phi(r.t(i)+r.deltaT/2, r.Y[r.t(i)]+r.deltaT*k1/2)
	k3 = r.phi.phi(r.t(i)+r.deltaT/2, r.Y[r.t(i)]+r.deltaT*k2/2)
	k3 = r.phi.phi(r.t(i)+r.deltaT, r.Y[r.t(i)]+r.deltaT*k3)
	return
}

// Zeitpunkt nach i Zeitschritten
func (r *RungeKutta) t(i int) float64 {
	return r.t_0 + r.deltaT*float64(i)
}
