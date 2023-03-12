package goserge

/*
*

	    InterpolateFloat64 returns an interpolated value between the currentFloat and targetFloat based on the given interpolateFactor.
	    The function calculates the interpolated float value between the currentFloat and the targetFloat using the given interpolateFactor.
	    The currentFloat and targetFloat parameters are of type float64, which means they are floating-point numbers with a precision of 64 bits.
	    The interpolateFactor parameter is also of type float64, and it represents the percentage of interpolation between the current and target values.
	    The calculation is done using a linear interpolation formula, where the function adds the difference between targetFloat and currentFloat,
	    multiplied by the interpolateFactor, to the currentFloat.
		@param currentFloat the current float value to be interpolated
		@param targetFloat the target float value to be interpolated
		@param interpolateFactor the interpolation factor used to calculate the interpolated value between the currentFloat and targetFloat
		@return the interpolated float value between the currentFloat and targetFloat based on the given interpolateFactor
*/
func InterpolateFloat64(currentFloat, targetFloat, interpolateFactor float64) float64 {
	return currentFloat + (targetFloat-currentFloat)*interpolateFactor
}
