package helper

func CalculateShippingCost(distance float64) float64 {
	const BASE_COST = 50000
	const COST_PER_KM = 15000

	var cost float64

	if distance < 1 {
		return BASE_COST
	}

	if distance > 1 {
		cost = BASE_COST + (COST_PER_KM * (distance - 1))
	}
	return cost
}